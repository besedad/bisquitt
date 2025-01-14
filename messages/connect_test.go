package messages

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectMessage(t *testing.T) {
	clientID := []byte("client-id")
	cleanSession := true
	will := true
	duration := uint16(90)
	msg := NewConnectMessage(clientID, cleanSession, will, duration)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*messages.ConnectMessage", reflect.TypeOf(msg).String(), "Type should be ConnectMessage")
		assert.Equal(t, will, msg.Will, "Bad Will value")
		assert.Equal(t, cleanSession, msg.CleanSession, "Bad CleanSession value")
		assert.Equal(t, uint8(1), msg.ProtocolID, "Default ProtocolID should be 1")
		assert.Equal(t, duration, msg.Duration, "Bad Duration value")
		assert.Equal(t, clientID, msg.ClientID, "Bad ClientID value")
	}
}

func TestConnectMarshal(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBuffer(nil)

	msg1 := NewConnectMessage([]byte("test-client"), true, true, 75)
	err := msg1.Write(buf)
	if err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(buf.Bytes())
	msg2, err := ReadPacket(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(msg1, msg2.(*ConnectMessage))
}
