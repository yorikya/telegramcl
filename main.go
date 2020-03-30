package main

import (
	"log"
	"os"
	"time"
	"fmt"

	"github.com/patrickmn/go-cache"
	"github.com/yanzay/tbot/v2"
	"github.com/yorikya/telegramcl/app"
)

var (
	telgramToken = ""
	workManager = ""
	messageLang = ""
)

func getEnvVarOrDefault(key, defaultVal string) (string, error) {
	if val := os.Getenv(key); val != "" {
		return val, nil
	}

	if defaultVal == ""{
		return "", fmt.Errorf("environment varaible %s is missing", key)
	}

	return defaultVal, nil
}

func init() {
	var err error
	telgramToken, err = getEnvVarOrDefault("TELEGRAM_TOKEN", "")
	if err != nil {
		panic("TELEGRAM_TOKEN environment varaible is missing")
	}

	workManager, err = getEnvVarOrDefault("WORK_MANAGER", "")
	if err != nil {				
		panic("WORK_MANAGER environment varaible is missing")
	}

	messageLang, _ = getEnvVarOrDefault("MESSAGE_LANG", "EN")
}

func main() {
	users := cache.New(24*time.Hour, 24*time.Hour)
	onBoardUsers := cache.New(24*time.Hour, 24*time.Hour)
	orders := cache.New(24*time.Hour, 24*time.Hour)
	
	bot := tbot.New(telgramToken)
	app := app.New(users, onBoardUsers, orders, bot, workManager, messageLang)

	log.Panic(app.Start())
}
