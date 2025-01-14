package messages

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillTopicReqStruct(t *testing.T) {
	msg := NewWillTopicReqMessage()

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*messages.WillTopicReqMessage", reflect.TypeOf(msg).String(), "Type should be WillTopicReqMessage")
	}
}

func TestWillTopicReqMarshal(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBuffer(nil)

	msg1 := NewWillTopicReqMessage()
	if err := msg1.Write(buf); err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(buf.Bytes())
	msg2, err := ReadPacket(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(msg1, msg2.(*WillTopicReqMessage))
}
