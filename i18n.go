package goeasyi18n

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type I18n struct {
	languages               map[string]TranslateStrings
	pluralizationFuncs      map[string]PluralizationFunc
	fallbackLanguageName    string
	disableConsistencyCheck bool
}

// Config is used to configure the i18n object
type Config struct {
	// Default: "en"
	FallbackLanguageName string
	// Default: false
	DisableConsistencyCheck bool
}

// NewI18n creates and returns a new i18n object
func NewI18n(config ...Config) *I18n {
	var pickedConfig Config
	if len(config) > 0 {
		pickedConfig = config[0]
	} else {
		pickedConfig = Config{}
	}

	instance := I18n{
		languages:               make(map[string]TranslateStrings),
		pluralizationFuncs:      make(map[string]PluralizationFunc),
		disableConsistencyCheck: pickedConfig.DisableConsistencyCheck,
	}

	if pickedConfig.FallbackLanguageName != "" {
		instance.fallbackLanguageName = pickedConfig.FallbackLanguageName
	} else {
		instance.fallbackLanguageName = "en"
	}

	return &instance
}

// DefaultPluralizationFunc is the function that is used if no
// custom pluralization function is set for the language
func DefaultPluralizationFunc(count int) string {
	if count == 1 {
		return "One"
	}
	return "Many"
}

// createGenderForm function to determinate the gender
// and sanitize it
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

	// If the gender is not valid, return empty string
	// it causes the use of the Default form of the
	// translation
	return ""

}

// CheckLanguageConsistency checks if a language is consistent
// with the other languages, it checks if the translations keys
// are the same in all languages
func (t *I18n) CheckLanguageConsistency(
	langNameToCheck string,
) (bool, []string) {
	langToCheck, exists := t.languages[langNameToCheck]
	if !exists {
		return false, []string{
			"goeasyi18n: the language '" + langNameToCheck + "' doesn't exist",
		}
	}

	inconsistencies := []string{}

	for langName, lang := range t.languages {
		if langName == langNameToCheck {
			continue
		}

		// Check if the new language has more keys
		// than existing languages
		for _, translateStringToCheck := range langToCheck {
			found := false
			for _, translateString := range lang {
				if translateString.Key == translateStringToCheck.Key {
					found = true
					break
				}
			}
			if !found {
				inconsistencies = append(
					inconsistencies,
					fmt.Sprintf(
						"goeasyi18n: the language '%s' has the key '%s' that doesn't exist in '%s'",
						langNameToCheck,
						translateStringToCheck.Key,
						langName,
					),
				)
			}
		}

		// Check if the new language has less keys
		// than existing languages
		for _, translateString := range lang {
			found := false
			for _, translateStringToCheck := range langToCheck {
				if translateString.Key == translateStringToCheck.Key {
					found = true
					break
				}
			}
			if !found {
				inconsistencies = append(
					inconsistencies,
					fmt.Sprintf(
						"goeasyi18n: the language '%s' has the key '%s' that doesn't exist in '%s'",
						langName,
						translateString.Key,
						langNameToCheck,
					),
				)
			}
		}
	}

	isConsistent := len(inconsistencies) == 0
	return isConsistent, inconsistencies
}

// AddLanguage adds a language to the i18n object with its translations
// and after that it check if the language is consistent with
// the other languages (can be disabled with the config)
//
// It returns a slice of errors as strings if the language is
// not consistent with the other languages and the consistency
// check is enabled
func (t *I18n) AddLanguage(
	languageName string,
	translateStrings TranslateStrings,
) []string {
	t.languages[languageName] = translateStrings
	t.SetPluralizationFunc(languageName, DefaultPluralizationFunc)

	if t.disableConsistencyCheck == false {
		isConsistent, errors := t.CheckLanguageConsistency(languageName)
		if isConsistent == false {
			errorMsg := strings.Join(errors, "\n")
			fmt.Println(errorMsg)
		}
		return errors
	}

	return nil
}

// HasLanguage checks if a language is available (if is loaded)
func (t *I18n) HasLanguage(languageName string) bool {
	_, ok := t.languages[languageName]
	return ok
}

// SetPluralizationFunc sets the pluralization function for a language
func (t *I18n) SetPluralizationFunc(languageName string, fn PluralizationFunc) {
	t.pluralizationFuncs[languageName] = fn
}

// Options are the additional options for the Translate function
type Options struct {
	Data   any
	Count  *int
	Gender *string // male, female, nonbinary, non-binary (case insensitive)
}

// Translate translates a string in a specific language using a key
// and additional optional options
func (t *I18n) Translate(
	languageName string,
	translateKey string,
	options ...Options,
) string {
	// Initialize options if not provided
	var pickedOptions Options
	if len(options) > 0 {
		pickedOptions = options[0]
	} else {
		pickedOptions = Options{}
	}

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
	if pickedOptions.Count != nil && pickedOptions.Gender != nil {
		mode = "PluralizedGendered"
	}
	if mode == "Default" && pickedOptions.Count != nil {
		mode = "Pluralized"
	}
	if mode == "Default" && pickedOptions.Gender != nil {
		mode = "Gendered"
	}

	// Get the plural and gender forms to be used if needed
	var pluralForm, genderForm string
	if mode == "Pluralized" || mode == "PluralizedGendered" {
		pluralizationFunc := t.pluralizationFuncs[languageName]
		pluralForm = pluralizationFunc(*pickedOptions.Count)
	}
	if mode == "Gendered" || mode == "PluralizedGendered" {
		genderForm = createGenderForm(*pickedOptions.Gender)
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
	if pickedOptions.Data != nil {
		translation = ExecuteTemplate(translation, pickedOptions.Data)
	}

	return translation
}

// T is a shortcut for Translate
func (t *I18n) T(
	languageName string,
	translateKey string,
	options ...Options,
) string {
	return t.Translate(languageName, translateKey, options...)
}

// NewTemplatingTranslateFunc creates a function that can be used in text/template and html/template.
// Just pass the created function to template.FuncMap.
//
// Example:
//
//	tempTemplate := template.New("main").Funcs(
//
//		template.FuncMap{
//			// 👇 "Translate" could be just "T" (for simplicity) or any other name you want.
//			"Translate": i18n.NewTemplatingTranslateFunc(),
//		},
//
//	)
//
// Then you can use it in your template like this:
//
//	{{ Translate "lang" "en" "key" "hello_emails" "gender" "nonbinary" "count" "100" "SomeData" "Anything" }}
//
// The format is key-value based and the order doesn't matter.
// This is the format:
//
//	{{ Translate "key1" "value1" "key2" "value2" ... }}
//
// Arguments:
//
// - "lang" "en": Language code (e.g., "en", "es").
//
// - "key" "hello_emails": Translation key.
//
// - "gender" "nonbinary": Gender for the translation (optional).
//
// - "count" "100": Count for pluralization (optional).
//
// - Additional key-value pairs will be added to the Data map.
//
// Arguments are passed in pairs. The first item in each pair is the key, and the second is the value.
//
// Key-Value Explanation:
//
// - Each argument is processed as a pair: the first string is considered the key and the second string is the value.
//
// - For example, in "lang" "en", "lang" is the key and "en" is the value.
//
// As you can imagine, "lang", "key", "gender", and "count" are reserved keys.
// You can use any other key you want to pass data to translation.
//
// Note: All arguments are strings. The function will attempt to convert "count" to an integer.
func (t *I18n) NewTemplatingTranslateFunc() func(args ...interface{}) string {
	return func(args ...interface{}) string {
		var lang, key string
		var gender *string
		var count *int
		data := make(Data)

		for i := 0; i < len(args); i += 2 {
			if i+1 >= len(args) {
				break
			}

			keyStr, ok1 := args[i].(string)
			valueStr, ok2 := args[i+1].(string)

			if !ok1 || !ok2 {
				continue
			}

			switch keyStr {
			case "lang":
				lang = valueStr
			case "key":
				key = valueStr
			case "count":
				intVal, err := strconv.Atoi(valueStr)
				if err == nil {
					count = &intVal
				}
			case "gender":
				gender = &valueStr
			default:
				data[keyStr] = valueStr
			}
		}

		options := Options{
			Count:  count,
			Gender: gender,
			Data:   data,
		}

		return t.Translate(lang, key, options)
	}
}

// NewLangTranslateFunc creates a function to translate a string in a specific language
// without the need to pass the language name every time.
func (t *I18n) NewLangTranslateFunc(
	languageName string,
) func(
	translateKey string,
	options ...Options,
) string {
	return func(translateKey string, options ...Options) string {
		return t.Translate(languageName, translateKey, options...)
	}
}
