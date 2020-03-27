package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yanzay/tbot/v2"
)

type application struct {
	client  *tbot.Client
	votings map[string]*voting
}

type voting struct {
	ups   int
	downs int
}

func mainButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         "🛵 משלוח",
		CallbackData: "new_delivery",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         "📦 הזמנות",
		CallbackData: "delivery_info",
	}
	button3 := tbot.InlineKeyboardButton{
		Text:         "📝 הירשם",
		CallbackData: "new_member",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2, button3},
		},
	}
}

func (a *application) stat(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		start := time.Now()
		h(u)
		log.Printf("Handle time: %v, update: %+v", time.Now().Sub(start), u.Message)
	}
}

func main() {
	token := "1072843930:AAHrQ4ieqegwF9yyXO31DYHyw4Dn3Cza0A8"
	bot := tbot.New(token)
	app := &application{
		votings: make(map[string]*voting),
	}
	app.client = bot.Client()
	bot.Use(app.stat) // add stat middleware to bot
	bot.HandleMessage("/oz", app.ozHandler)
	bot.HandleMessage("/end", app.endHandler)
	bot.HandleCallback(app.callbackHandler)

	bot.Start()
}

func (a *application) endHandler(m *tbot.Message) {
	msg, _ := a.client.SendMessage(m.Chat.ID, "הזמנה נמצאת בטיפול", tbot.OptInlineKeyboardMarkup(mainButtons()))
	votingID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)
	a.votings[votingID] = &voting{}
}

func (a *application) ozHandler(m *tbot.Message) {
	msg, _ := a.client.SendMessage(m.Chat.ID, ":פעולות", tbot.OptInlineKeyboardMarkup(mainButtons()))
	votingID := fmt.Sprintf("%s:%d", m.Chat.ID, msg.MessageID)
	a.votings[votingID] = &voting{}
}

func (a *application) callbackHandler(cq *tbot.CallbackQuery) {
	votingID := fmt.Sprintf("%s:%d", cq.Message.Chat.ID, cq.Message.MessageID)
	switch cq.Data {
	case "new_delivery":
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("הזן רשימה חדש"))
		msg, _ := a.client.SendMessage(cq.Message.Chat.ID, "רשימה חדשה")
		fmt.Println("get new delivery", votingID, "messageID:", msg.MessageID)
		return
	case "delivery_info":
		fmt.Println("get delivery info", votingID)
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("הזמנה נמצאת בטיפול"))
		_, _ = a.client.SendMessage(cq.Message.Chat.ID, `קיימת הזמנה 1
🟢 בטיפול
⏱ זמן הגעה משואר 10:20`)
		return
	case "new_member":
		a.client.AnswerCallbackQuery(cq.ID, tbot.OptText("משתמש חדש , תהליך רישום החל"))
		_, _ = a.client.SendMessage(cq.Message.Chat.ID, "הזן כתובת")
		return
	}

}
