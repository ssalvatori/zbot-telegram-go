package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/tucnak/telebot"

	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ssalvatori/zbot-telegram-go/commands"
	"github.com/ssalvatori/zbot-telegram-go/db"
)

const version string = "1.0"
const dbFile string = "./sample.db"

var levelsConfig = command.Levels{
	Ignore: 100,
	Lock:   1000,
}

func main() {

	log.Info("Loading zbot-telegram")
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	bot, err := telebot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	database := &db.SqlLite{
		File: dbFile,
	}

	err = database.Init()
	defer database.Close()

	if err != nil {
		log.Fatal(err)
	}

	go database.UserCleanIgnore()
	bot.Messages = make(chan telebot.Message, 1000)
	go messagesProcessing(database, bot)

	bot.Start(1 * time.Second)
}

func messagesProcessing(db db.ZbotDatabase, bot *telebot.Bot) {
	output := make(chan string)
	for message := range bot.Messages {

		//we're going to process only the message starting with ! or ?
		processingMsg := regexp.MustCompilePOSIX(`^[!|?].*`)

		//check if the user isn't on the ignore_list
		log.Debug(fmt.Sprintf("Checking user [%s] ", strings.ToLower(message.Sender.Username)))
		ignore, err := db.UserCheckIgnore(strings.ToLower(message.Sender.Username))
		if err != nil {
			log.Error(err)
		}
		if !ignore {
			if processingMsg.MatchString(message.Text) {
				log.Debug(fmt.Sprintf("Received a message from %s with the text: %s", message.Sender.Username, message.Text))
				go sendResponse(bot, db, message, output)
			}
		} else {
			log.Debug(fmt.Sprintf("User [%s] ignored", strings.ToLower(message.Sender.Username)))
		}
	}
}

func sendResponse(bot *telebot.Bot, db db.ZbotDatabase, msg telebot.Message, output chan string) {
	response := processing(db, msg)
	bot.SendMessage(msg.Chat, response, nil)
}

func buildUser(sender telebot.User) command.User {
	user := command.User{}
	user.Ident = strings.ToLower(sender.FirstName)
	user.Username = strings.ToLower(sender.FirstName)

	if sender.Username != "" {
		user.Username = strings.ToLower(sender.Username)
	}

	return user
}

func processing(db db.ZbotDatabase, msg telebot.Message) string {

	user := buildUser(msg.Sender)

	// TODO: how to clean this code
	commands := &command.PingCommand{}
	versionCommand := &command.VersionCommand{Version: version}
	statsCommand := &command.StatsCommand{Db: db}
	randCommand := &command.RandCommand{Db: db}
	topCommand := &command.TopCommand{Db: db}
	lastCommand := &command.LastCommand{Db: db}
	getCommand := &command.GetCommand{Db: db}
	findCommand := &command.FindCommand{Db: db}
	searchCommand := &command.SearchCommand{Db: db}
	learnCommand := &command.LearnCommand{Db: db}
	levelCommand := &command.LevelCommand{Db: db}
	ignoreCommand := &command.IgnoreCommand{Db: db, Levels: levelsConfig}
	lockCommand := &command.LockCommand{Db: db, Levels: levelsConfig}
	externalCommand := &command.ExternalCommand{
		PathModules: "./modules/",
	}

	commands.Next = versionCommand
	versionCommand.Next = statsCommand
	statsCommand.Next = randCommand
	randCommand.Next = topCommand
	topCommand.Next = lastCommand
	lastCommand.Next = getCommand
	getCommand.Next = findCommand
	findCommand.Next = searchCommand
	searchCommand.Next = learnCommand
	learnCommand.Next = levelCommand
	levelCommand.Next = ignoreCommand
	lockCommand.Next = lockCommand

	ignoreCommand.Next = externalCommand

	outputMsg := commands.ProcessText(msg.Text, user)

	return outputMsg
}
