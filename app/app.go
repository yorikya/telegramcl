package app

import (
	"fmt"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/yanzay/tbot/v2"
)

const (
	//Onboarding states
	phoneRequire = iota
	addressRequire
	cityRequire
	descriptionRequire
	onboardingComplete
)

func mainButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         "🛵 משלוח",
		CallbackData: "new_delivery",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         "📦 הזמנות",
		CallbackData: "delivery_info",
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2},
		},
	}
}

func deliveryTypeButtons() *tbot.InlineKeyboardMarkup {
	button1 := tbot.InlineKeyboardButton{
		Text:         "Market",
		CallbackData: "new_market",
	}
	button2 := tbot.InlineKeyboardButton{
		Text:         "Pharm",
		CallbackData: "new_pharm",
	}
	button3 := tbot.InlineKeyboardButton{
		Text:         "DIY",
		CallbackData: "new_diy",
	}
	button4 := tbot.InlineKeyboardButton{
		Text:         "Private",
		CallbackData: "new_private",
	}
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{button1, button2, button3, button4},
		},
	}
}

type User struct {
	Username,
	Phone,
	Address,
	City,
	Description string
	OnboardingState int
}

type Order struct {
	OrderType string
	Items     []string
	UserInfo  *User
	Complete  bool
	StartTime time.Time
	EndTime   time.Time
}

type Application struct {
	users        *cache.Cache
	onBoradUsers *cache.Cache
	orders       *cache.Cache
	bot          *tbot.Client
	botServer    *tbot.Server
	manager,
	managerKey,
	managerChatID  string 
}

func New(users, onBoardUsers, orders *cache.Cache, bot *tbot.Server, manager string) *Application {
	a := &Application{
		users:        users,
		onBoradUsers: onBoardUsers,
		orders:       orders,
		botServer:    bot,
		bot:          bot.Client(),
		manager:      manager,
		managerKey:   "--manager",  //Telegram username not accept '--' chars, the manager key can't be oviriten with another username
	}

	bot.Use(a.messageListener) // add messageListener middleware to bot
	bot.HandleMessage("/start", a.startHandler)
	bot.HandleMessage("/end", a.endHandler)
	bot.HandleMessage("/work", a.startWorkHandler)
	bot.HandleMessage("/finish", a.stopWorkHandler)
	bot.HandleCallback(a.callbackHandler)

	return a
}

func (a *Application) Start() error {
	log.Println("stating application...")
	return a.botServer.Start()
}

func (a *Application) validateUserInfo(c tbot.Chat) bool {
	return c.Username == ""
}

func (a *Application) validateRegestredUser(username string) bool {
	_, found := a.users.Get(username)
	return found
}

func (a *Application) sendMessage(chatID, username,  text string) {
	if _, err := a.bot.SendMessage(chatID, text); err != nil {
		log.Printf("get an erro when send a message: %s to user: %s, ChatID: %s, error: %s", text, username, chatID, err)
	}
	log.Printf("send message to user: %s, message: %s", username, text)
}

func (a *Application) sendMessageToManger(msg string) {
	if a.managerChatID != "" {
		a.sendMessage(a.managerChatID, "--Manager", msg)
		return
	}
	log.Println("Not has a working manager, discard message", msg)
}

func (a *Application) sendMainMenu(chatID string) {
	a.bot.SendMessage(chatID, "Main:", tbot.OptInlineKeyboardMarkup(mainButtons()))
}

func (a *Application) canOrder(_ string) bool {
	_, found := a.users.Get(a.managerKey) //Check if manger start working
	return found
}

func (a *Application) completeOrder(m *tbot.Message) {
	if o, found := a.orders.Get(m.Chat.Username); found {
		order := o.(*Order)
		order.Items = append(order.Items, m.Text)
		
		a.sendMessage(m.Chat.ID, m.Chat.Username, "items added, to finish order /end")
		return
	}
	
	a.sendMessage(m.Chat.ID, m.Chat.Username, "complete order, somthing went wrong")
}

func (a *Application) onBoardingProcess(username string) bool {
	_, found := a.onBoradUsers.Get(username)
	return found

}

func (a *Application) onOrderProcess(username string) bool {
	_, found := a.orders.Get(username)
	return found
}

func (a *Application) startOnBoarding(username, chatID string) {
	a.onBoradUsers.Set(username, &User{
		Username: username,
	}, cache.DefaultExpiration)
	a.sendMessage(chatID, username,"please enter phone number (ex: 0521234567)")
}

func (a *Application) completeOnboarding(m *tbot.Message) {
	u, found := a.onBoradUsers.Get(m.Chat.Username)
	if !found {
		log.Println("error onboarding, can not find user '%s'", m.Chat.Username)
		return
	}

	user := u.(*User)
	switch user.OnboardingState {
	case phoneRequire:
		user.Phone = m.Text //TODO: Validate phone number
		user.OnboardingState = addressRequire
		a.sendMessage(m.Chat.ID, m.Chat.Username,"please enter adress (ex: Zavitan 9)")
	case addressRequire:
		user.Address = m.Text
		user.OnboardingState = cityRequire
		a.sendMessage(m.Chat.ID, m.Chat.Username,"please enter city (ex: Kazerin)")
	case cityRequire:
		user.City = m.Text
		user.OnboardingState = descriptionRequire
		a.sendMessage(m.Chat.ID, m.Chat.Username, "please enter descritpion for courier (ex: House behind the post office wiht big white gate)")
	case descriptionRequire:
		user.Description = m.Text
		user.OnboardingState = onboardingComplete

		a.users.Set(m.Chat.Username, user, cache.DefaultExpiration)

		a.onBoradUsers.Delete(m.Chat.Username) //TODO: create API for delete onboarding users (using for stucking users)

		a.sendMessageToManger(fmt.Sprintf("craete new user %+v", user))
		a.sendMessage(m.Chat.ID, m.Chat.Username, fmt.Sprintf("on boarding complete! now you can try your first delivery"))
		a.sendMainMenu(m.Chat.ID)
	default:
		log.Printf("unknown user state, user: '%s', state: %d", m.Chat.Username, user.OnboardingState)
	}
}

func (a *Application) startHandler(m *tbot.Message) {
	if a.validateUserInfo(m.Chat) {
		a.sendMessage(m.Chat.ID, "<nil>","please define @Username in you telegram app")
		return
	}

	if !a.validateRegestredUser(m.Chat.Username) {
		a.startOnBoarding(m.Chat.Username, m.Chat.ID)
		return
	}

	a.sendMainMenu(m.Chat.ID)
}

func (a *Application) endHandler(m *tbot.Message) {
	if a.validateRegestredUser(m.Chat.Username) {
		if o, found := a.orders.Get(m.Chat.Username); found {
			order := o.(*Order)
			order.Complete = true
			order.EndTime = time.Now()
			a.orders.Delete(m.Chat.Username)
			a.sendMessage(m.Chat.ID, m.Chat.Username, "order completed")
			a.sendMainMenu(m.Chat.ID)
			
			a.sendMessageToManger(fmt.Sprintf("get an new order:\n %+v", *order)) //TODO: Send order email
			return
		}
		a.sendMessage(m.Chat.ID, m.Chat.Username, "you dont have an open order")
		return
	}
	a.sendMessage(m.Chat.ID, m.Chat.Username, "end order error")
}

func (a *Application) startWorkHandler(m *tbot.Message) {
	if m.Chat.Username == a.manager {
		a.users.Set(a.managerKey, &User{
			Username:        a.manager,
			Phone:           "0521234567",
			OnboardingState: onboardingComplete,
		}, cache.DefaultExpiration)
		a.managerChatID = m.Chat.ID
		a.sendMessageToManger("start working, to finish work /finish")
		return
	}
	a.sendMessage(m.Chat.ID, m.Chat.Username, "you are not authorize")
}

func (a *Application) stopWorkHandler(m *tbot.Message) {
	if m.Chat.Username == a.manager {
		a.users.Delete(a.managerKey)
		a.sendMessage(m.Chat.ID, m.Chat.Username, "finish working")
		return
	}
	
	a.sendMessage(m.Chat.ID, m.Chat.Username, "you are not authorize")
}

func (a *Application) callbackHandler(cq *tbot.CallbackQuery) {
	switch cq.Data {
	case "new_market", "new_pharm", "new_diy", "new_private":
		if u, found := a.users.Get(cq.Message.Chat.Username); found {
			a.orders.Set(cq.Message.Chat.Username, &Order{
				UserInfo:  u.(*User),
				StartTime: time.Now(),
				OrderType: cq.Data,
			}, cache.DefaultExpiration)
			a.sendMessage(cq.Message.Chat.ID, cq.Message.Chat.Username, `please send items list with new line seperator
Example:
Milk 3% - 1
Elite Black Coffe - 4`)
			return
		}
		a.sendMessage(cq.Message.Chat.ID, cq.Message.Chat.Username, "somthing went wrong")
	case "new_delivery":
		if !a.validateRegestredUser(cq.Message.Chat.Username) {
			a.sendMessage(cq.Message.Chat.ID, cq.Message.Chat.Username, "somthing went wrong")
			return
		}

		if !a.canOrder(cq.Message.Chat.Username) {
			a.sendMessage(cq.Message.Chat.ID, cq.Message.Chat.Username, "delivery unavaible now")
			return
		} 

		if _, err := a.bot.SendMessage(cq.Message.Chat.ID, "New delivery", tbot.OptInlineKeyboardMarkup(deliveryTypeButtons())); err != nil {
			log.Println("get an error when send new order message")
		}

	case "delivery_info":
		a.sendMessage(cq.Message.Chat.ID, cq.Message.Chat.Username,`קיימת הזמנה 1
🟢 בטיפול
⏱ זמן הגעה משואר 10:20`)
	}
}

func (a *Application) messageListener(h tbot.UpdateHandler) tbot.UpdateHandler {
	return func(u *tbot.Update) {
		if u.Message != nil {
			switch {
			case a.onBoardingProcess(u.Message.Chat.Username):
				a.completeOnboarding(u.Message)
				return
			case (u.Message.Chat.Username != a.manager) && !a.validateRegestredUser(u.Message.Chat.Username) && (u.Message.Text != "/start"):
				a.sendMessage(u.Message.Chat.ID, u.Message.Chat.Username, "you need complete registration /start")
				return
			case a.onOrderProcess(u.Message.Chat.Username) && (u.Message.Text != "/end"):
				a.completeOrder(u.Message)
				return
			}
		}
		h(u)
	}
}
