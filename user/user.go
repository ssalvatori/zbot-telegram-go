package user

import (
	"github.com/ssalvatori/zbot-telegram-go/db"

	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"
)

type User struct {
	Username string
	Ident    string
	Host     string
	Level    int
}

//Create a user using telegram information
func BuildUser(sender telebot.User, db db.ZbotDatabase) User {
	user := User{}
	user.Ident = strings.ToLower(sender.FirstName)

	if sender.Username != "" {
		user.Username = strings.ToLower(sender.Username)
	} else {
		user.Username = strings.ToLower(sender.FirstName)
	}

	user.Level = GetUserLevel(db, sender.Username)

	return user
}

// Get the current level for a user using its username
func GetUserLevel(Db db.ZbotDatabase, username string) int {
	userLevel, err := Db.UserLevel(username)

	if err != nil {
		log.Error(err)
		return 0
	}

	userLevelInt, _ := strconv.Atoi(userLevel)

	return userLevelInt
}

// Check if username has level greater or equal that a level given
func (u User) IsAllow(level int) bool {
	result := false

	if u.Level >= level {
		result = true
	}

	return result
}