package messages

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisconnectStruct(t *testing.T) {
	duration := uint16(123)
	msg := NewDisconnectMessage(duration)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*messages.DisconnectMessage", reflect.TypeOf(msg).String(), "Type should be DisconnectMessage")
		assert.Equal(t, duration, msg.Duration, "Bad Duration value")
	}
}

func TestDisconnectMarshal(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBuffer(nil)

	msg1 := NewDisconnectMessage(75)
	err := msg1.Write(buf)
	if err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(buf.Bytes())
	msg2, err := ReadPacket(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(msg1, msg2.(*DisconnectMessage))
}
