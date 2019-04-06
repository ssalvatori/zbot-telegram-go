package zbot

import (
	"encoding/json"
	"testing"

	command "github.com/ssalvatori/zbot-telegram-go/commands"
	"github.com/ssalvatori/zbot-telegram-go/db"
	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestProcessingIsCommandDisabled(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	command.DisabledCommands = []string{
		"learn",
		"version",
	}

	botMsg := tb.Message{Text: "!learn", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "", result, "command disabled")

}

func Test_ProcessingVersion(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	buildTime = "2017-05-06 09:59:21.318841424 +0300 EEST"
	command.DisabledCommands = nil

	botMsg := tb.Message{
		Text: "!version",
		Sender: &tb.User{
			Username: "zbot_test",
		},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "zbot golang version ["+version+"] commit [undefined] build-time ["+buildTime+"]", result, "!version default")
}

func TestProcessingStats(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{Text: "!stats", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "Count: 666", "!stats")
}

func TestProcessingPing(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{Text: "!ping", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "pong!!", "!ping")
}

func TestProcessingRand(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Rand_def: db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
	}

	botMsg := tb.Message{Text: "!rand", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "[hola] - [gatolinux]", result, "!rand")
}

func TestProcessingGet(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:   "666",
		File:    "hola.db",
		Term:    "hola",
		Meaning: "foo bar!",
	}

	botMsg := tb.Message{Text: "? hola", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "[hola] - [foo bar!]", "? def fail")

}

func TestProcessingFind(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:   "666",
		File:    "hola.db",
		Term:    "hola",
		Meaning: "foo bar!",
	}

	botMsg := tb.Message{Text: "!find hola", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, result, "hola", "!find fail")
}

func TestProcessingSearch(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
	}

	botMsg := tb.Message{Text: "!search hola", Sender: &tb.User{Username: "zbot_test"}}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "hola chao foobar", result, "!rand")
}

func TestProcessingUserLevel(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
	}

	botMsg := tb.Message{
		Text:   "!level",
		Sender: &tb.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "ssalvato level 666", result, "!level self user")
}

func TestProcessingUserIgnoreList(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
		User_ignored: []db.UserIgnore{
			{Username: "ssalvato", Since: "1231", Until: "4564"},
		},
	}

	botMsg := tb.Message{
		Text:   "!ignore list",
		Sender: &tb.User{FirstName: "ssalvato", Username: "ssalvato"},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "[ @ssalvato ] since [01-01-1970 00:20:31 UTC] until [01-01-1970 01:16:04 UTC]", result, "!ignore list")
}

func TestProcessingUserIgnoreInsert(t *testing.T) {

	dbMock := &db.MockZbotDatabase{
		Level:        "666",
		File:         "hola.db",
		Term:         "hola",
		Meaning:      "foo bar!",
		Find_terms:   []string{"hola", "chao", "foo_bar"},
		Rand_def:     db.DefinitionItem{Term: "hola", Meaning: "gatolinux"},
		Search_terms: []string{"hola", "chao", "foobar"},
		User_ignored: []db.UserIgnore{{Username: "ssalvatori", Since: "1231", Until: "4564"}},
	}

	botMsg := tb.Message{
		Text:   "!ignore add rigo",
		Sender: &tb.User{FirstName: "ssalvatori", Username: "ssalvatori"},
	}
	result := processing(dbMock, botMsg)
	assert.Equal(t, "User [rigo] ignored for 10 minutes", result, "!ignore add OK")

	botMsg = tb.Message{
		Text:   "!ignore add ssalvatori",
		Sender: &tb.User{FirstName: "ssalvatori", Username: "ssalvatori"},
	}
	result = processing(dbMock, botMsg)
	assert.Equal(t, "You can't ignore yourself", result, "!ignore add myself")

}

func TestProcessingLearnReplyTo(t *testing.T) {
	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{Text: "!learn arg1",
		Sender: &tb.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tb.Message{
			Text: "message in reply-to",
			Sender: &tb.User{
				Username: "otheruser",
			},
		},
	}
	result := processing(dbMock, botMsg)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestMessageProcessing(t *testing.T) {
	dbMock := &db.MockZbotDatabase{
		Level: "666",
		File:  "hola.db",
	}

	botMsg := tb.Message{Text: "!learn arg1",
		Sender: &tb.User{
			Username:  "ssalvatori",
			FirstName: "stefano",
		},
		ReplyTo: &tb.Message{
			Text: "message in reply-to",
			Sender: &tb.User{
				Username: "otheruser",
			},
		},
	}

	result := messagesProcessing(dbMock, &botMsg)

	assert.Equal(t, "[arg1] - [otheruser message in reply-to]", result, "!learn with replayto")
}

func TestGetDisabledCommands(t *testing.T) {

	commands := `["level","ignore"]`
	jsonRaw := json.RawMessage(commands)
	binary, _ := jsonRaw.MarshalJSON()

	SetDisabledCommands(binary)

	assert.Equal(t, []string{"ignore", "level"}, GetDisabledCommands(), "Get Disabled Commands")

	commands = `["level]`
	jsonRaw = json.RawMessage(commands)
	binary, _ = jsonRaw.MarshalJSON()
	SetDisabledCommands(binary)

	assert.Equal(t, []string(nil), GetDisabledCommands(), "Get Disabled Commands")

}

/*
func TestMessagesProcessing(t *testing.T) {
	dbMock := &db.MockZbotDatabase{
		Ignore_User: true,
	}
	msgChan := make(chan tb.Message)
	bot := &tb.Bot{Messages: msgChan}

	msgObj := tb.Message{
		Text:   "!hola",
		Sender: tb.User{FirstName: "Stefano", Username: "Ssalvato"},
	}
	bot.Messages <- msgObj
	go messagesProcessing(dbMock, bot)
}
*/
