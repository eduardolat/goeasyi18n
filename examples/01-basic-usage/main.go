package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

func main() {
	// 1. Create a new i18n instance
	// You can skip the goeasyi18n.Config{} entirely if you are
	// ok with the default values.
	i18n := goeasyi18n.NewI18n(goeasyi18n.Config{
		// You can set the fallback language (optional)
		// The default value is "en"
		FallbackLanguageName: "en",

		// You can disable the consistency check (optional)
		// By default, if you add a translation for a language
		// that has not the same keys as the other languages,
		// the i18n instance will log warnings.
		DisableConsistencyCheck: false,
	})

	// 2. Create your translations
	enTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_message",
			Default: "Hello, welcome to Go Easy i18n!",
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_message",
			Default: "¡Hola, bienvenido a Go Easy i18n!",
		},
	}

	// 3. Add your languages with their translations
	// The name of the language can be anything you want
	// You can use simple strings like "en" or "es"
	// You you can use the full language name like "english" or "spanish"
	// You can even use the IETF Language Tag like "en-US" or "es-ES"
	i18n.AddLanguage("en", enTranslations)

	// If the language has the same keys as the other languages,
	// the i18n instance will log warnings and return a slice of
	// inconsistencies as strings. You can disable this behavior
	// in the i18n instance config.
	inconsistencies := i18n.AddLanguage("es", esTranslations)
	fmt.Printf("Inconsistencies: %v\n", inconsistencies)
	// (no inconsistencies)
	// Prints: Inconsistencies: []

	// 4. You are done! 🎉 Just get that translations!
	t1 := i18n.Translate("en", "hello_message", goeasyi18n.Options{})
	// Or you can use the T method (it's just an alias for Translate)
	// and you can skip the options if you don't need them
	t2 := i18n.T("es", "hello_message")

	fmt.Println(t1)
	fmt.Println(t2)

	/*
		Prints:
		Hello, welcome to Go Easy i18n!
		¡Hola, bienvenido a Go Easy i18n!
	*/

	// 5. (Extra) You can check if a language exists in the i18n instance
	enExists := i18n.HasLanguage("en")
	xxExists := i18n.HasLanguage("xx")

	fmt.Printf("en exists: %v\n", enExists)
	fmt.Printf("xx exists: %v\n", xxExists)

	/*
		Prints:
		en exists: true
		xx exists: false
	*/

	// 6. (Extra) You can create a translate function for a specific language
	// to prevent passing the language name every time you want to translate
	// something
	translateEn := i18n.NewLangTranslateFunc("en")
	translateEs := i18n.NewLangTranslateFunc("es")

	// You can skip the options if you don't need them
	t3 := translateEn("hello_message", goeasyi18n.Options{})
	t4 := translateEs("hello_message")

	fmt.Println(t3)
	fmt.Println(t4)

	/*
		Prints:
		Hello, welcome to Go Easy i18n!
		¡Hola, bienvenido a Go Easy i18n!
	*/
}
