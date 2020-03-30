package main

import (
	"log"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yanzay/tbot/v2"
	"github.com/yorikya/telegramcl/app"
)

var (
	telgramToken = ""
	workManager = ""
	messageLang = ""
)

func init() {
	telgramToken = os.Getenv("TELEGRAM_TOKEN")
	if telgramToken == "" {
		panic("TELEGRAM_TOKEN environment varaible is missing")
	}

	workManager = os.Getenv("WORK_MANAGER")
	if workManager == "" {				
		panic("WORK_MANAGER environment varaible is missing")
	}

	messageLang = os.Getenv("MESSAGE_LANG")
}

func main() {
	users := cache.New(24*time.Hour, 24*time.Hour)
	onBoardUsers := cache.New(24*time.Hour, 24*time.Hour)
	orders := cache.New(24*time.Hour, 24*time.Hour)
	
	bot := tbot.New(telgramToken)
	app := app.New(users, onBoardUsers, orders, bot, workManager, messageLang)

	log.Panic(app.Start())
}
