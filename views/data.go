package views

import (
	"log"

	"../models"
)

const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"
	AlertMsgGeneric = "Something went wrong please try again"
)

type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

type Alert struct {
	Level   string
	Message string
}

type PublicError interface {
	error
	Public() string
}

func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}

	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}

}
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
