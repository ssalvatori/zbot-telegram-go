package command

import (
	"regexp"
	"github.com/ssalvatori/zbot-telegram-go/db"
	log "github.com/Sirupsen/logrus"
	"fmt"
)

type StatsCommand struct {
	Db db.ZbotDatabase
	Next     HandlerCommand
}

func (handler *StatsCommand) ProcessText(text string) string {

	commandPattern := regexp.MustCompile(`^!stats`)

	if(commandPattern.MatchString(text)) {
		statTotal, err := handler.Db.Statistics()
		if err != nil {
			log.Error(err)
			return "Error!"
		}
		return fmt.Sprintf("Count: %s",statTotal)
	} else {
		return handler.Next.ProcessText(text)
	}

}

