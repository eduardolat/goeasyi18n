package main

import (
	"github.com/eduardolat/goeasyi18n"
)

var i18n *goeasyi18n.I18n

func InitializeI18n() {
	i18n = goeasyi18n.NewI18n(goeasyi18n.Config{})

	enTranslations, err := goeasyi18n.LoadFromYaml("./translations/en.yaml")
	if err != nil {
		panic(err)
	}

	esTranslations, err := goeasyi18n.LoadFromYaml("./translations/es.yaml")
	if err != nil {
		panic(err)
	}

	ptTranslations, err := goeasyi18n.LoadFromYaml("./translations/pt.yaml")
	if err != nil {
		panic(err)
	}

	frTranslations, err := goeasyi18n.LoadFromYaml("./translations/fr.yaml")
	if err != nil {
		panic(err)
	}

	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)
	i18n.AddLanguage("pt", ptTranslations)
	i18n.AddLanguage("fr", frTranslations)
}
