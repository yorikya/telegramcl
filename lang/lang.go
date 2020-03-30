package lang

import (
	"log"
)

type Lang interface {
	MainMenu() string
	NewDelivery()string
	DeliveryInfo() string
	Market()string
	Pharm()string
	DIY()string
	Private()string
	EnterPhone()string
	EnterAddress()string
	EnterCity()string
	EnterDscription()string
	OnboardingComplete()string
	PleaseDefineUserName()string
	ItemAdded()string
	CompleteOrderError()string
	OrderComplited()string
	DontHaveOpenOrder()string
	EndOrderError()string
	NotAuthorize()string
	OrderItemListExample()string
	GeneralError()string
	DeliveryUnavaible()string
	CompleteRegistartion()string
}

type LangEN struct {}

func NewLangEN() LangEN {
	return LangEN{}
}

func (LangEN) MainMenu() string {
	return "Main Menu"
}
func (LangEN)NewDelivery()string {
	return "New delivery"
}
func (LangEN)DeliveryInfo() string {
	return "Delivery Info"
}
func (LangEN)Market()string {
	return "Market"
}
func (LangEN)Pharm()string {
	return "Pharm"
}
func (LangEN)DIY()string {
	return "DIY"
}
func (LangEN)Private()string{
	return "Private"
}
func (LangEN)EnterPhone()string {
	return "please enter phone number (ex: 0521234567)"
}
func (LangEN)EnterAddress()string {
	return "please enter adress (ex: Zavitan 9)"
}
func (LangEN)EnterCity()string{
	return "please enter city (ex: Kazerin)"
}
func (LangEN)EnterDscription()string{
	return "please enter descritpion for courier (ex: House behind the post office wiht big white gate)"
}
func (LangEN)OnboardingComplete()string{
	return "on boarding complete! now you can try your first delivery"
}
func (LangEN)PleaseDefineUserName()string{
	return "please define @Username in you telegram app"
}
func (LangEN)ItemAdded()string{
	return "items added, to finish order /end"
}

func (LangEN)CompleteOrderError()string{
	return "complete order, somthing went wrong"
}
func (LangEN)OrderComplited()string {
	return "order completed"
}
func (LangEN)DontHaveOpenOrder()string{
	return "you dont have an open order"
}
func (LangEN)EndOrderError()string{
	return "end order error"
}
func (LangEN)NotAuthorize()string{
	return "you are not authorize"
}
func (LangEN)OrderItemListExample()string{
	return "please send items list with new line seperator\nExample:\nMilk 3% - 1\nElite Black Coffe - 4"
}
func (LangEN)GeneralError()string{
	return "somthing went wrong"
}
func (LangEN)DeliveryUnavaible()string {
	return "delivery unavaible now"
}
func (LangEN)CompleteRegistartion()string {
	return "you need complete registration /start"
}

type LangHE struct {}

func NewLangHE() LangHE {
	return LangHE{}
}

func (LangHE) MainMenu() string {
	return "Main Menu"
}
func (LangHE)NewDelivery()string {
	return "New delivery"
}
func (LangHE)DeliveryInfo() string {
	return "Delivery Info"
}
func (LangHE)Market()string {
	return "Market"
}
func (LangHE)Pharm()string {
	return "Pharm"
}
func (LangHE)DIY()string {
	return "DIY"
}
func (LangHE)Private()string{
	return "Private"
}
func (LangHE)EnterPhone()string {
	return "please enter phone number (ex: 0521234567)"
}
func (LangHE)EnterAddress()string {
	return "please enter adress (ex: Zavitan 9)"
}
func (LangHE)EnterCity()string{
	return "please enter city (ex: Kazerin)"
}
func (LangHE)EnterDscription()string{
	return "please enter descritpion for courier (ex: House behind the post office wiht big white gate)"
}
func (LangHE)OnboardingComplete()string{
	return "on boarding complete! now you can try your first delivery"
}
func (LangHE)PleaseDefineUserName()string{
	return "please define @Username in you telegram app"
}
func (LangHE)ItemAdded()string{
	return "items added, to finish order /end"
}

func (LangHE)CompleteOrderError()string{
	return "complete order, somthing went wrong"
}
func (LangHE)OrderComplited()string {
	return "order completed"
}
func (LangHE)DontHaveOpenOrder()string{
	return "you dont have an open order"
}
func (LangHE)EndOrderError()string{
	return "end order error"
}
func (LangHE)NotAuthorize()string{
	return "you are not authorize"
}
func (LangHE)OrderItemListExample()string{
	return "please send items list with new line seperator\nExample:\nMilk 3% - 1\nElite Black Coffe - 4"
}
func (LangHE)GeneralError()string{
	return "somthing went wrong"
}
func (LangHE)DeliveryUnavaible()string {
	return "delivery unavaible now"
}
func (LangHE)CompleteRegistartion()string {
	return "you need complete registration /start"
}

func New(l string) Lang {
	switch l {
	case "HE":
		log.Println("Message languge is Hebrew")
		return NewLangHE()
	default:
		log.Println("[default] Message languge is English")
		return NewLangEN()
	}
}