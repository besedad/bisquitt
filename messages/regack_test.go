package messages

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegackStruct(t *testing.T) {
	topicID := uint16(123)
	returnCode := RC_ACCEPTED
	msg := NewRegackMessage(topicID, returnCode)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*messages.RegackMessage", reflect.TypeOf(msg).String(), "Type should be RegackMessage")
		assert.Equal(t, uint16(7), msg.MessageLength(), "Default Length should be 7")
		assert.Equal(t, topicID, msg.TopicID, "Bad TopicID value")
		assert.Equal(t, uint16(0), msg.MessageID(), "Default MessageID should be 0")
		assert.Equal(t, returnCode, msg.ReturnCode, "ReturnCode should be RC_ACCEPTED")
	}
}

func TestRegackMarshal(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBuffer(nil)

	msg1 := NewRegackMessage(123, RC_CONGESTION)
	msg1.SetMessageID(12)
	if err := msg1.Write(buf); err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(buf.Bytes())
	msg2, err := ReadPacket(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(msg1, msg2.(*RegackMessage))
}
