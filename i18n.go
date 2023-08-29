package goeasyi18n

import (
	"reflect"
	"strconv"
	"strings"
)

type I18n struct {
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
func NewI18n(config Config) *I18n {
	instance := I18n{
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
func (t *I18n) AddLanguage(languageName string, translateStrings TranslateStrings) {
	t.languages[languageName] = translateStrings
	t.SetPluralizationFunc(languageName, DefaultPluralizationFunc)
}

// Check if a language is available (if is loaded)
func (t *I18n) HasLanguage(languageName string) bool {
	_, ok := t.languages[languageName]
	return ok
}

// Set the pluralization function for a language
func (t *I18n) SetPluralizationFunc(languageName string, fn PluralizationFunc) {
	t.pluralizationFuncs[languageName] = fn
}

// Additional options for the Translate function
type TranslateOptions struct {
	Data   any
	Count  *int
	Gender *string // male, female, nonbinary, non-binary (case insensitive)
}

// Translate a string using its key and the language name
func (t *I18n) Translate(languageName string, translateKey string, options TranslateOptions) string {
	// Get lang and fallback if not found
	lang, okLang := t.languages[languageName]
	fallbackLang, okFallbackLang := t.languages[t.fallbackLanguageName]
	if !okLang && !okFallbackLang {
		return ""
	}
	if !okLang {
		copy(lang, fallbackLang)
		languageName = t.fallbackLanguageName
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
func (t *I18n) T(languageName string, translateKey string, options TranslateOptions) string {
	return t.Translate(languageName, translateKey, options)
}

/*
NewTemplatingTranslateFunc creates a function that can be used in text/template and html/template.
Just pass the created function to template.FuncMap.

Example:

tempTemplate := template.New("main").Funcs(

	template.FuncMap{
		// 👇 "Translate" could be just "T" (for simplicity) or any other name you want.
		"Translate": i18n.NewTemplatingTranslateFunc(),
	},

)

Then you can use it in your template like this:

{{Translate "lang:en" "key:hello_emails" "gender:nonbinary" "count:100" "SomeData:Anything"}}

Arguments:

- "lang:en": Language code (e.g., "en", "es").
- "key:hello_emails": Translation key.
- "gender:nonbinary": Gender for the translation (optional).
- "count:100": Count for pluralization (optional).
- Additional key-value pairs will be added to the Data map.

As you can imagine, lang, key, gender and count are reserved keys.
You can use any other key you want to pass data to translation.

Note: All arguments are strings. The function will attempt to convert "count" to an integer.
*/
func (t *I18n) NewTemplatingTranslateFunc() func(args ...any) string {
	return func(args ...any) string {
		var lang, key string
		var gender *string
		var count *int
		data := make(Data)

		for _, arg := range args {
			strArg, ok := arg.(string)
			if !ok {
				continue
			}

			parts := strings.SplitN(strArg, ":", 2)
			if len(parts) != 2 {
				continue
			}

			switch parts[0] {
			case "lang":
				lang = parts[1]
			case "key":
				key = parts[1]
			case "count":
				intVal, err := strconv.Atoi(parts[1])
				if err == nil {
					count = &intVal
				}
			case "gender":
				gender = &parts[1]
			default:
				data[parts[0]] = parts[1]
			}
		}

		options := TranslateOptions{
			Count:  count,
			Gender: gender,
			Data:   data,
		}

		return t.Translate(lang, key, options)
	}
}
