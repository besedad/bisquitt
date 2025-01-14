package client

import (
	"fmt"

	msgs "github.com/energomonitor/bisquitt/messages"
	"github.com/energomonitor/bisquitt/transactions"
)

type publishQOS1Transaction struct {
	*transaction
}

func newPublishQOS1Transaction(client *Client, msgID uint16) *publishQOS1Transaction {
	tLog := client.log.WithTag(fmt.Sprintf("PUBLISH1(%d)", msgID))
	tLog.Debug("Created.")
	return &publishQOS1Transaction{
		transaction: &transaction{
			RetryTransaction: transactions.NewRetryTransaction(
				client.groupCtx, client.cfg.RetryDelay, client.cfg.RetryCount,
				func(lastMsg interface{}) error {
					tLog.Debug("Resend.")
					return client.send(lastMsg.(msgs.Message))
				},
				func() {
					client.transactions.Delete(msgID)
					tLog.Debug("Deleted.")
				},
			),
			client: client,
			log:    tLog,
		},
	}
}

func (t *publishQOS1Transaction) Puback(puback *msgs.PubackMessage) {
	if t.State != awaitingPuback {
		t.log.Debug("Unexpected message in %d: %v", t.State, puback)
		return
	}
	t.Success()
}
