package main

import (
	"embed"

	"github.com/eduardolat/goeasyi18n"
)

var i18n *goeasyi18n.I18n

//go:embed translations
var translationsFS embed.FS

func InitializeI18n() {
	i18n = goeasyi18n.NewI18n()

	enTranslations, err := goeasyi18n.LoadFromYamlFS(
		translationsFS,
		"translations/en.yaml",
	)
	if err != nil {
		panic(err)
	}

	esTranslations, err := goeasyi18n.LoadFromYamlFS(
		translationsFS,
		"translations/es.yaml",
	)
	if err != nil {
		panic(err)
	}

	ptTranslations, err := goeasyi18n.LoadFromYamlFS(
		translationsFS,
		"translations/pt.yaml",
	)
	if err != nil {
		panic(err)
	}

	frTranslations, err := goeasyi18n.LoadFromYamlFS(
		translationsFS,
		"translations/fr.yaml",
	)
	if err != nil {
		panic(err)
	}

	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)
	i18n.AddLanguage("pt", ptTranslations)
	i18n.AddLanguage("fr", frTranslations)
}
