package messages

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillMsgUpdateStruct(t *testing.T) {
	willMsg := []byte("test-msg")
	msg := NewWillMsgUpdateMessage(willMsg)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*messages.WillMsgUpdateMessage", reflect.TypeOf(msg).String(), "Type should be WillMsgUpdateMessage")
		assert.Equal(t, willMsg, msg.WillMsg, "Bad WillMsg value")
	}
}

func TestWillMsgUpdateMarshal(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBuffer(nil)

	msg1 := NewWillMsgUpdateMessage([]byte("test-message"))
	if err := msg1.Write(buf); err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(buf.Bytes())
	msg2, err := ReadPacket(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(msg1, msg2.(*WillMsgUpdateMessage))
}
