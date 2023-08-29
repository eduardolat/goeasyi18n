package main

import (
	"github.com/eduardolat/goeasyi18n"
)

var i18n *goeasyi18n.I18n

func InitializeI18n() {
	i18n = goeasyi18n.NewI18n(goeasyi18n.Config{})

	enTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_message",
			Default: "Hello {{.Name}}, welcome to Go Easy i18n!",
		},
		{
			Key:  "unread_messages",
			One:  "You have {{.Qty}} unread message.",
			Many: "You have {{.Qty}} unread messages.",
		},
		{
			Key:       "friend_update",
			Male:      "{{.Name}} updated his status.",
			Female:    "{{.Name}} updated her status.",
			NonBinary: "{{.Name}} updated their status.",
		},
		{
			Key:        "friend_request",
			OneMale:    "He sent you a friend request.",
			ManyFemale: "She sent you {{.Qty}} friend requests.",
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_message",
			Default: "¡Hola {{.Name}}, bienvenido a Go Easy i18n!",
		},
		{
			Key:  "unread_messages",
			One:  "Tienes {{.Qty}} mensaje sin leer.",
			Many: "Tienes {{.Qty}} mensajes sin leer.",
		},
		{
			Key:       "friend_update",
			Male:      "{{.Name}} actualizó su estado.",
			Female:    "{{.Name}} actualizó su estado.",
			NonBinary: "{{.Name}} actualizó su estado.",
		},
		{
			Key:        "friend_request",
			OneMale:    "Él te envió una solicitud de amistad.",
			ManyFemale: "Ella te envió {{.Qty}} solicitudes de amistad.",
		},
	}

	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)
}
