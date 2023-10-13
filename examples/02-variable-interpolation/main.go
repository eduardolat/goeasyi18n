package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

func main() {
	// 1. Create a new i18n instance
	i18n := goeasyi18n.NewI18n()

	// 2. Create your translations
	// You can add any variables to your translations
	// Use the syntax {{.VariableName}}
	enTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_message",
			Default: "Hello {{.Name}} {{.SurName}}, welcome to Go Easy i18n!",
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_message",
			Default: "¡Hola {{.Name}} {{.SurName}}, bienvenido a Go Easy i18n!",
		},
	}

	// 3. Add your languages with their translations
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 4. Crete the options for the translation with the variables
	// The Data field is a map[string]any that contains the variables to be replaced
	options := goeasyi18n.Options{
		Data: map[string]any{
			"Name":    "John",
			"SurName": "Doe",
		},
	}

	// 5. Get the translations using the options (with the variables)
	t1 := i18n.T("en", "hello_message", options)
	t2 := i18n.T("es", "hello_message", options)

	fmt.Println(t1)
	fmt.Println(t2)

	/*
		Prints:
		Hello John Doe, welcome to Go Easy i18n!
		¡Hola John Doe, bienvenido a Go Easy i18n!
	*/
}
