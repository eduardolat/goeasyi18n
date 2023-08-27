package goeasyi18n

import (
	"reflect"
	"strings"
)

type i18n struct {
	languages            map[string]TranslateStrings
	pluralizationFuncs   map[string]PluralizationFunc
	fallbackLanguageName string
}

// Config for the i18n object
type Config struct {
	// Default: "en"
	FallbackLanguageName string
}

// Create a new i18n object
func NewI18n(config Config) *i18n {
	instance := i18n{
		languages:          make(map[string]TranslateStrings),
		pluralizationFuncs: make(map[string]PluralizationFunc),
	}

	if config.FallbackLanguageName != "" {
		instance.fallbackLanguageName = config.FallbackLanguageName
	} else {
		instance.fallbackLanguageName = "en"
	}

	return &instance
}

// Default pluralization function that is used if no pluralization function is set for the language
func DefaultPluralizationFunc(count int) string {
	if count == 1 {
		return "One"
	}
	return "Many"
}

// Function to determinate the gender and sanitize it
func createGenderForm(input string) string {

	input = strings.ToLower(input)

	if input == "male" {
		return "Male"
	}

	if input == "female" {
		return "Female"
	}

	if input == "nonbinary" || input == "non-binary" {
		return "NonBinary"
	}

	return ""

}

// Add a language to the i18n object with its translations
func (t *i18n) AddLanguage(languageName string, translateStrings TranslateStrings) {
	t.languages[languageName] = translateStrings
	t.SetPluralizationFunc(languageName, DefaultPluralizationFunc)
}

// Set the pluralization function for a language
func (t *i18n) SetPluralizationFunc(languageName string, fn PluralizationFunc) {
	t.pluralizationFuncs[languageName] = fn
}

// Additional options for the Translate function
type TranslateOptions struct {
	Data   any
	Count  *int
	Gender *string // male, female, nonbinary, non-binary (case insensitive)
}

// Translate a string using its key and the language name
func (t *i18n) Translate(languageName string, translateKey string, options TranslateOptions) string {
	// Get lang and fallback if not found
	lang, okLang := t.languages[languageName]
	fallbackLang, okFallbackLang := t.languages[t.fallbackLanguageName]
	if !okLang && !okFallbackLang {
		return ""
	}
	if !okLang {
		lang = fallbackLang
	}

	// Get the translate string from key or fallback if not found
	var translateString TranslateString
	for _, ts := range lang {
		if ts.Key == translateKey {
			translateString = ts
			break
		}
	}
	if translateString.Key == "" {
		for _, ts := range fallbackLang {
			if ts.Key == translateKey {
				translateString = ts
				break
			}
		}
	}
	if translateString.Key == "" {
		return ""
	}

	// Get the string key to be used
	mode := "Default" // Default - Pluralized - Gendered - PluralizedGendered
	if options.Count != nil && options.Gender != nil {
		mode = "PluralizedGendered"
	}
	if mode == "Default" && options.Count != nil {
		mode = "Pluralized"
	}
	if mode == "Default" && options.Gender != nil {
		mode = "Gendered"
	}

	// Get the plural and gender forms to be used if needed
	var pluralForm, genderForm string
	if mode == "Pluralized" || mode == "PluralizedGendered" {
		pluralizationFunc := t.pluralizationFuncs[languageName]
		pluralForm = pluralizationFunc(*options.Count)
	}
	if mode == "Gendered" || mode == "PluralizedGendered" {
		genderForm = createGenderForm(*options.Gender)
	}

	// Get the string key to be used
	stringKey := "Default"
	if mode == "Pluralized" {
		stringKey = pluralForm
	}
	if mode == "Gendered" && genderForm != "" {
		stringKey = genderForm
	}
	if mode == "PluralizedGendered" && genderForm != "" {
		stringKey = pluralForm + genderForm
	}

	// Get the translation
	translation := translateString.Default
	reflected := reflect.ValueOf(translateString)
	field := reflected.FieldByName(stringKey)
	if field.IsValid() {
		translation = field.String()
	}

	// Execute the template
	if options.Data != nil {
		translation = ExecuteTemplate(translation, options.Data)
	}

	return translation
}

// T is a shortcut for Translate
func (t *i18n) T(languageName string, translateKey string, options TranslateOptions) string {
	return t.Translate(languageName, translateKey, options)
}
