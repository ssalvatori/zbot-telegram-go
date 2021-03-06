package command

import (
	"errors"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram/db"
	"github.com/ssalvatori/zbot-telegram/user"
	"github.com/ssalvatori/zbot-telegram/utils"
)

// WhoCommand definition
type WhoCommand struct {
	Db db.ZbotDatabase
}

//SetDb set db connection if the module need it
func (handler *WhoCommand) SetDb(db db.ZbotDatabase) {
	handler.Db = db
}

//ProcessText run command
func (handler *WhoCommand) ProcessText(text string, user user.User, chat string, private bool) (string, error) {

	if private {
		return "", ErrNextCommand
	}

	commandPattern := regexp.MustCompile(`^!who\s(\S*)$`)

	if commandPattern.MatchString(text) {
		if checkLearnCommandOnChannel(chat) {
			return "", ErrLearnDisabledChannel
		}

		term := commandPattern.FindStringSubmatch(text)
		item := db.Definition{
			Term: term[1],
		}
		Item, err := handler.Db.Get(item.Term, chat)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return "", fmt.Errorf("Definition [%s] not found", item.Term)
			}
			log.Error(err.Error())
			return "", ErrInternalError
		}

		err = handler.Db.IncreaseHits(Item.ID)
		if err != nil {
			log.Error(err.Error())
			return "", ErrInternalError
			// if !errors.Is(err, db.ErrInternalError) {
			// 	return "Internal Error, check logs", nil
			// }
		}

		return fmt.Sprintf("[%s] by [%s] on [%v] hits [%d]", Item.Term, Item.Author, utils.ConvertToDateToUTC(Item.UpdatedAt), Item.Hits), nil
	}

	return "", ErrNextCommand
}
