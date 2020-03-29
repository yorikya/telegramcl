package main

import (
	"log"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yanzay/tbot/v2"
	"github.com/yorikya/telegramcl/app"
)

func main() {
	users := cache.New(24*time.Hour, 24*time.Hour)
	onBoardUsers := cache.New(24*time.Hour, 24*time.Hour)
	orders := cache.New(24*time.Hour, 24*time.Hour)
	manager := os.Getenv("WORK_MANAGER")
	if manager != "" {
		bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
		app := app.New(users, onBoardUsers, orders, bot, manager)

		log.Panic(app.Start())
	}
	log.Println("missing WORK_MANAGER environment varaible")
}
