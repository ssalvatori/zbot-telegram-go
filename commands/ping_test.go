package command

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var pingCommand = PingCommand{}

func TestPingCommandOK(t *testing.T) {
	assert.Equal(t, "pong!!", pingCommand.ProcessText("!ping"), "Ping Command")
}
func TestPingCommandNoNext(t *testing.T) {
	assert.Equal(t, "", pingCommand.ProcessText("!ping6"), "Ping no next command")
}

func TestPingCommandNext(t *testing.T) {
	pingCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", pingCommand.ProcessText("!ping6"), "Ping  next command")
}