package messages

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchGwStruct(t *testing.T) {
	radius := uint8(123)
	msg := NewSearchGwMessage(radius)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*messages.SearchGwMessage", reflect.TypeOf(msg).String(), "Type should be SearchGwMessage")
		assert.Equal(t, uint16(3), msg.MessageLength(), "Default Length should be 3")
		assert.Equal(t, radius, msg.Radius, "Bad Radius value")
	}
}

func TestSearchGwMarshal(t *testing.T) {
	assert := assert.New(t)
	buf := bytes.NewBuffer(nil)

	msg1 := NewSearchGwMessage(123)
	if err := msg1.Write(buf); err != nil {
		t.Fatal(err)
	}

	r := bytes.NewReader(buf.Bytes())
	msg2, err := ReadPacket(r)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(msg1, msg2.(*SearchGwMessage))
}
