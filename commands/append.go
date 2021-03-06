package command

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
)

// AppendCommand definition
type AppendCommand struct {
	Db db.ZbotDatabase
}

// ProcessText run command
func (handler *AppendCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`(?s)^!append\s(\S*)\s(.*)`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}
		term := commandPattern.FindStringSubmatch(text)
		def := db.Definition{
			Term:    term[1],
			Meaning: term[2],
			Author:  fmt.Sprintf("%s!%s@telegram.bot", user.Username, user.Ident),
			Chat:    chat,
		}
		err := handler.Db.Append(def, chat)
		if err != nil {
			log.Error(err.Error())
			return "", ErrInternalError
		}

		def, err = handler.Db.Get(def.Term, chat)
		if err != nil {
			log.Error(err.Error())
			return "", ErrInternalError
		}
		return fmt.Sprintf("[%s] = [%s]", def.Term, def.Meaning), nil
	}

	return "", ErrNextCommand
}
