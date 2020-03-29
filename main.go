package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yanzay/tbot/v2"
)

type application struct {
	client              *tbot.Client
	customers           map[string]*customer
	onBoardingCustomers map[string]*onBoardingCustomer
}

type onBoardingCustomer struct {
	username    string
	homeAddress string
	cellPhone   string
	description string
	done        bool
}

type customer struct {
	username    string
	homeAddress string
	cellPhone   string
	description string
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

func (a *application) messageListener(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		start := time.Now()
		h(u)
		log.Printf("Handle time: %v, update: %+v", time.Now().Sub(start), u.Message)
	}
}

func main() {
	bot := tbot.New(os.Getenv("TELEGRAM_TOKEN"))
	app := &application{
		customers:           make(map[string]*customer),
		onBoardingCustomers: make(map[string]*onBoardingCustomer),
	}
	app.client = bot.Client()
	bot.Use(app.messageListener) // add messageListener middleware to bot
	bot.HandleMessage("/oz", app.ozHandler)
	bot.HandleMessage("/end", app.endHandler)
	bot.HandleMessage("/manager", app.mangerHandler)
	bot.HandleCallback(app.callbackHandler)

	bot.Start()
}

func (a *application) mangerHandler(m *tbot.Message) {
	_, _ = a.client.SendMessage(m.Chat.ID, "good morning manager")
}

func (a *application) endHandler(m *tbot.Message) {
	_, _ = a.client.SendMessage(m.Chat.ID, "הזמנה נמצאת בטיפול", tbot.OptInlineKeyboardMarkup(mainButtons()))

}

func (a *application) ozHandler(m *tbot.Message) {
	// username := m.Chat.Username
	// if username == "" {
	// 	a.client.SendMessage(m.Chat.ID, "הגדר שם משתמש (username)")
	// 	return
	// }
	// _, ok := a.customers[username]
	// if !ok {
	// 	a.client.SendMessage(m.Chat.ID, "הזן כתובת הבית")
	// 	a.onBoardingCustomers[username] = &onBoardingCustomer{
	// 		username: username,
	// 	}
	// 	return
	// }
	a.client.SendMessage(m.Chat.ID, ":פעולות", tbot.OptInlineKeyboardMarkup(mainButtons()))
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

		_, _ = a.client.SendMessage("@GolanDeliveryBot", "לקוח חדש , פתיחת כרטיס")

		return
	}

}
