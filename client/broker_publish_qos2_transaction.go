// brokerPublishQOS2 just saves the original PUBLISH message for later handlers firing.
//
// [the sender] MUST NOT re-send the PUBLISH once it has sent the corresponding
// PUBREL message.
// MQTT specification v. 5.0, chapter 4.3.3 QoS 2: Exactly once delivery
//
// Hence, we fire subscription handlers as soon as the corresponding PUBREL is
// received.

package client

import (
	"fmt"

	snMsgs "github.com/energomonitor/bisquitt/messages"
	"github.com/energomonitor/bisquitt/transactions"
)

type brokerPublishQOS2Transaction struct {
	*transactions.TransactionBase
	client  *Client
	publish *snMsgs.PublishMessage
}

func newBrokerPublishQOS2Transaction(client *Client, msgID uint16) *brokerPublishQOS2Transaction {
	tLog := client.log.WithTag(fmt.Sprintf("PUBLISH2b(%d)", msgID))
	tLog.Debug("Created.")
	return &brokerPublishQOS2Transaction{
		TransactionBase: transactions.NewTransactionBase(
			func() {
				client.transactions.Delete(msgID)
				tLog.Debug("Deleted.")
			},
		),
		client: client,
	}
}

func (t *brokerPublishQOS2Transaction) Publish(publish *snMsgs.PublishMessage) error {
	t.publish = publish
	pubrec := snMsgs.NewPubrecMessage()
	pubrec.CopyMessageID(publish)
	return t.client.send(pubrec)
}

func (t *brokerPublishQOS2Transaction) Pubrel(pubrel *snMsgs.PubrelMessage) error {
	pubcomp := snMsgs.NewPubcompMessage()
	pubcomp.CopyMessageID(pubrel)
	topic, err := t.client.topicForPublish(t.publish)
	if err != nil {
		return err
	}
	t.client.messageHandlers.handle(t.client, topic, t.publish)
	err = t.client.send(pubcomp)
	if err != nil {
		return err
	}
	t.Success()
	return nil
}
