package command

import (
	"errors"
	"regexp"

	"github.com/ssalvatori/zbot-telegram-go/user"
)

// PingCommand command definition
type PingCommand struct {
	//	Next   HandlerCommand
	//	Levels Levels
}

// ProcessText run command
func (handler *PingCommand) ProcessText(text string, user user.User) (string, error) {

	commandPattern := regexp.MustCompile(`^!ping$`)

	if commandPattern.MatchString(text) {
		return "pong!!", nil
	}

	return "", errors.New("text doesn't match")
}
