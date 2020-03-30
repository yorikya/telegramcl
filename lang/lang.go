package lang

import (
	"log"
)

type Lang interface {
	GetLang() string
	Translate(string) string
}

type LangEN struct {
	lang string
}

func NewLangEN() *LangEN {
	return &LangEN{
		lang: "EN",
	}
}

func (l LangEN) GetLang() string {
	return l.lang
}

func (l LangEN) Translate(sentense string) string {
	return sentense
}

type LangHE struct {
	lang string
}

func NewLangHE() *LangHE {
	return &LangHE{
		lang: "HEB",
	}
}

func (l LangHE) GetLang() string {
	return l.lang
}

func (l LangHE) Translate(sentense string) string {
	return sentense
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