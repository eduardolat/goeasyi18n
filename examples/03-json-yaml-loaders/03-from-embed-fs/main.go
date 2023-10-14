package main

import (
	"embed"
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

//go:embed en/*.json
var enFS embed.FS

//go:embed es/*.yaml
var esFS embed.FS

func main() {
	// With the embed.FS feature you can load the translations from
	// any fs.FS implementation, like embed.FS.
	//
	// This allows you to bundle your translations with your app
	// in the same binary file.

	// 1. Create a new i18n instance
	i18n := goeasyi18n.NewI18n()

	// 2. Load your translations from JSON or YAML files inside the embed.FS
	// You can load one or more files like goeasyi18n.LoadFromJsonFS(fs, "en/t1.json", "en/t2.json")
	// You can use glob patterns like goeasyi18n.LoadFromJsonFS(fs, "en/*.json")
	// All the translation files get merged

	// Load english translations from JSON files
	enTranslations, err := goeasyi18n.LoadFromJsonFS(enFS, "en/*.json")
	if err != nil {
		panic(err)
	}

	// Load spanish translations from YAML files
	esTranslations, err := goeasyi18n.LoadFromYamlFS(esFS, "es/*.yaml")
	if err != nil {
		panic(err)
	}

	// 3. Add your languages with their translations
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 4. Crete the options for the translations with/without interpolations
	options := goeasyi18n.Options{}
	optionsWithName := goeasyi18n.Options{
		Data: map[string]string{
			"Name": "John Doe",
		},
	}

	// 5. Get the translations using the options (with the variables)
	ten1 := i18n.T("en", "hello_world", options)
	ten2 := i18n.T("en", "hello_user", optionsWithName)
	ten3 := i18n.T("en", "hello_admin", optionsWithName)

	tes1 := i18n.T("es", "hello_world", options)
	tes2 := i18n.T("es", "hello_user", optionsWithName)
	tes3 := i18n.T("es", "hello_admin", optionsWithName)

	fmt.Println(ten1)
	fmt.Println(ten2)
	fmt.Println(ten3)
	fmt.Println(tes1)
	fmt.Println(tes2)
	fmt.Println(tes3)

	/*
		Prints:
		Hello World from embed
		Hello John Doe from embed
		Hello John Doe, you are an admin, from embed
		Hola Mundo desde embed
		Hola John Doe desde embed
		Hola John Doe, eres un administrador, desde embed
	*/
}
