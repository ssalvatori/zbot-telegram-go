package command

import (
	"fmt"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/ssalvatori/zbot-telegram-go/db"

	"strings"
)

type LevelCommand struct {
	Db     db.ZbotDatabase
	Next   HandlerCommand
	Levels Levels
}

func (handler *LevelCommand) AddUser(userToCheck string, user string) string {
	return "not ready"
}

func (handler *LevelCommand) DelUser(userToCheck string, user string) string {
	return "not ready"
}

func (handler *LevelCommand) GetLevel(userToCheck string, user string) string {
	result := ""
	if IsUserAllow(handler.Db, user, 0) {
		level, err := handler.Db.UserLevel(user)
		if err != nil {
			log.Error(err)
		}
		result = fmt.Sprintf("%s level %s", user, level)
	}
	return result
}

func (handler *LevelCommand) ProcessText(text string, user User) string {
	commandPattern := regexp.MustCompile(`^!level(\s|$)(\S*)\s?(\S+)?\s?(\d+)?`)
	result := ""

	if commandPattern.MatchString(text) {
		subcommand := commandPattern.FindStringSubmatch(text)
		log.Debug("level subcommand: ", subcommand[2])
		log.Debug(strings.Join(subcommand, "-"))
		switch subcommand[2] {
			case "add":
				result = handler.AddUser(subcommand[2], user.Username)
			case "del":
				result = handler.DelUser(subcommand[2], user.Username)
			default:
				result = handler.GetLevel(subcommand[2], user.Username)
		}
	} else {
		if handler.Next != nil {
			result = handler.Next.ProcessText(text, user)
		}
	}
	return result
}
