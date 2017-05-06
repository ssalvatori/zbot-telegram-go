package command

import (
	"testing"

	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
)

var statsCommand = StatsCommand{}

func TestStatsCommandOK(t *testing.T) {
	statsCommand.Db = &db.MockZbotDatabase{
		Level:    "7",
		Rand_def: db.DefinitionItem{Term: "foo", Meaning: "bar"},
	}
	assert.Equal(t, "Count: 7", statsCommand.ProcessText("!stats", userTest), "Stats Command")
	assert.Equal(t, "", statsCommand.ProcessText("!stats6", userTest), "Stats no next command")
	statsCommand.Next = &FakeCommand{}
	assert.Equal(t, "Fake OK", statsCommand.ProcessText("!stats6", userTest), "Stats next command")
}
