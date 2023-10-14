package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

func main() {
	// 1. Create a new i18n instance
	i18n := goeasyi18n.NewI18n()

	// 2. Load your translations from a database, api, etc.
	enTranslationsString := `[
		{
			"Key": "hello_world",
			"Default": "Hello World"
		}
	]`
	esTranslationsBytes := []byte(`[
		{
			"Key": "hello_world",
			"Default": "Hola Mundo"
		}
	]`)

	// 3. Then you can load the translations from your source
	// In this case we are using JSON but you can also use YAML
	enTranslations, err := goeasyi18n.LoadFromJsonString(enTranslationsString)
	if err != nil {
		panic(err)
	}

	// You can also load from bytes instead of strings
	esTranslations, err := goeasyi18n.LoadFromJsonBytes(esTranslationsBytes)
	if err != nil {
		panic(err)
	}

	// 4. Add your languages with their translations
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 5. Get the translations
	ten := i18n.T("en", "hello_world")
	tes := i18n.T("es", "hello_world")

	fmt.Println(ten)
	fmt.Println(tes)

	/*
		Prints:
		Hello World
		Hola Mundo
	*/
}
