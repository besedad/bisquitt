// The MQTT broker watches the connection state using keepalive (PING* messages)
// _after_ the connection is established (after a CONNECT message is received by
// the broker).  Before the connection is established, the MQTT-SN gateway
// must watch the connection itself because a malicious client could leave the
// connection half-established (=> possible DoS attack vulnerability).
// Hence, we must use time-limited connectTransaction.

package gateway

import (
	"context"
	"errors"
	"fmt"

	mqttPackets "github.com/eclipse/paho.mqtt.golang/packets"
	snMsgs "github.com/energomonitor/bisquitt/messages"
	"github.com/energomonitor/bisquitt/transactions"
	"github.com/energomonitor/bisquitt/util"
)

var Cancelled = errors.New("transaction cancelled")

type connectTransaction struct {
	*transactions.TimedTransaction
	handler       *handler
	log           util.Logger
	authEnabled   bool
	mqConnect     *mqttPackets.ConnectPacket
	authenticated bool
}

func newConnectTransaction(ctx context.Context, h *handler, authEnabled bool, mqConnect *mqttPackets.ConnectPacket) *connectTransaction {
	tLog := h.log.WithTag("CONNECT")
	tLog.Debug("Created.")
	return &connectTransaction{
		TimedTransaction: transactions.NewTimedTransaction(
			ctx, connectTransactionTimeout,
			func() {
				h.transactions.DeleteByType(snMsgs.CONNECT)
				tLog.Debug("Deleted.")
			},
		),
		handler:     h,
		log:         tLog,
		authEnabled: authEnabled,
		mqConnect:   mqConnect,
	}
}

func (t *connectTransaction) Start(ctx context.Context) error {
	t.handler.group.Go(func() error {
		select {
		case <-t.Done():
			if err := t.Err(); err != nil {
				if err == Cancelled {
					return nil
				}
				return fmt.Errorf("CONNECT: %s", err)
			}
			t.log.Debug("CONNECT transaction finished successfully.")
			return nil
		case <-ctx.Done():
			t.log.Debug("CONNECT transaction cancelled.")
			return nil
		}
	})

	if t.authEnabled {
		t.log.Debug("Waiting for AUTH message.")
		return nil
	}

	if t.mqConnect.WillFlag {
		// Continue with WILLTOPICREQ.
		return t.handler.snSend(snMsgs.NewWillTopicReqMessage())
	}

	return t.handler.mqttSend(t.mqConnect)
}

func (t *connectTransaction) Auth(snMsg *snMsgs.AuthMessage) error {
	// Extract username and password from PLAIN data.
	if snMsg.Method == snMsgs.AUTH_PLAIN {
		user, password, err := snMsgs.DecodePlain(snMsg)
		if err != nil {
			t.Fail(err)
			return err
		}
		t.mqConnect.UsernameFlag = true
		t.mqConnect.Username = user
		t.mqConnect.PasswordFlag = true
		t.mqConnect.Password = password
	} else {
		if err := t.SendConnack(snMsgs.RC_NOT_SUPPORTED); err != nil {
			return err
		}
		err := fmt.Errorf("Unknown auth method: %#v.", snMsg.Method)
		t.Fail(err)
		return err
	}

	if t.mqConnect.WillFlag {
		// Continue with WILLTOPICREQ.
		return t.handler.snSend(snMsgs.NewWillTopicReqMessage())
	}

	// All information successfully gathered - send MQTT connect.
	return t.handler.mqttSend(t.mqConnect)
}

func (t *connectTransaction) WillTopic(snWillTopic *snMsgs.WillTopicMessage) error {
	t.mqConnect.WillQos = snWillTopic.QOS
	t.mqConnect.WillRetain = snWillTopic.Retain
	t.mqConnect.WillTopic = snWillTopic.WillTopic

	// Continue with WILLMSGREQ.
	return t.handler.snSend(snMsgs.NewWillMsgReqMessage())
}

func (t *connectTransaction) WillMsg(snWillMsg *snMsgs.WillMsgMessage) error {
	t.mqConnect.WillMessage = snWillMsg.WillMsg

	// All information successfully gathered - send MQTT connect.
	return t.handler.mqttSend(t.mqConnect)
}

func (t *connectTransaction) Connack(mqConnack *mqttPackets.ConnackPacket) error {
	if mqConnack.ReturnCode != mqttPackets.Accepted {
		// We misuse RC_CONGESTION here because MQTT-SN spec v. 1.2 does not define
		// any suitable return code.
		if err := t.SendConnack(snMsgs.RC_CONGESTION); err != nil {
			return err
		}
		returnCodeStr, ok := mqttPackets.ConnackReturnCodes[mqConnack.ReturnCode]
		if !ok {
			returnCodeStr = "unknown code!"
		}
		err := fmt.Errorf(
			"CONNECT refused by MQTT broker with return code %d (%s).",
			mqConnack.ReturnCode, returnCodeStr)
		t.Fail(err)
		return err
	}

	// Must be set before snSend to avoid race condition in tests.
	t.handler.setState(util.StateActive)
	if err := t.SendConnack(snMsgs.RC_ACCEPTED); err != nil {
		t.Fail(err)
		return err
	}
	t.Success()
	return nil
}

// Inform client that the CONNECT request was refused.
func (t *connectTransaction) SendConnack(code snMsgs.ReturnCode) error {
	snConnack := snMsgs.NewConnackMessage(code)
	if err := t.handler.snSend(snConnack); err != nil {
		t.Fail(err)
		return err
	}
	return nil
}
