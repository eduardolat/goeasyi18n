package goeasyi18n

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestTranslate(t *testing.T) {
	t.Run("simple general tests", func(t *testing.T) {
		i18n := NewI18n(Config{FallbackLanguageName: "en"})

		// Add English translations
		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome",
				Male:    "Welcome, sir",
				Female:  "Welcome, ma'am",
			},
			TranslateString{
				Key:  "emails",
				One:  "You have one email",
				Many: "You have many emails",
			},
			TranslateString{
				Key:     "greetings",
				Default: "Hello {{.Name}}",
			},
		})

		// Add Spanish translations
		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Bienvenido",
				Male:    "Bienvenido, señor",
				Female:  "Bienvenida, señora",
			},
			TranslateString{
				Key:  "emails",
				One:  "Tienes un correo",
				Many: "Tienes muchos correos",
			},
			TranslateString{
				Key:     "greetings",
				Default: "Hola {{.Name}}",
			},
		})

		tests := []struct {
			lang     string
			key      string
			options  Options
			expected string
		}{
			{"en", "welcome", Options{}, "Welcome"},
			{"en", "welcome", Options{Gender: createPtr("male")}, "Welcome, sir"},
			{"en", "welcome", Options{Gender: createPtr("female")}, "Welcome, ma'am"},
			{"en", "emails", Options{Count: createPtr(1)}, "You have one email"},
			{"en", "emails", Options{Count: createPtr(5)}, "You have many emails"},
			{"es", "welcome", Options{}, "Bienvenido"},
			{"es", "welcome", Options{Gender: createPtr("Male")}, "Bienvenido, señor"},
			{"es", "welcome", Options{Gender: createPtr("Female")}, "Bienvenida, señora"},
			{"es", "emails", Options{Count: createPtr(1)}, "Tienes un correo"},
			{"es", "emails", Options{Count: createPtr(5)}, "Tienes muchos correos"},
			// Test fallback language
			{"xxx", "welcome", Options{}, "Welcome"},
			// Test fallback key
			{"en", "xxx", Options{}, ""},
			// Test data interpolation
			{"en", "greetings", Options{Data: map[string]string{"Name": "John"}}, "Hello John"},
			{"es", "greetings", Options{Data: map[string]string{"Name": "John"}}, "Hola John"},
		}

		for _, test := range tests {
			t.Run(test.key, func(t *testing.T) {
				got := i18n.Translate(test.lang, test.key, test.options)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("method HasLanguage should work", func(t *testing.T) {
		i18n := NewI18n(Config{})
		i18n.AddLanguage("en", TranslateStrings{})
		i18n.AddLanguage("es", TranslateStrings{})

		if !i18n.HasLanguage("en") {
			t.Errorf("expected language en to exist")
		}

		if !i18n.HasLanguage("es") {
			t.Errorf("expected language es to exist")
		}

		if i18n.HasLanguage("xxx") {
			t.Errorf("expected language xxx to not exist")
		}
	})

	t.Run("english should be the default fallback language even if multiple langs are added", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Bienvenido",
			},
		})

		got := i18n.Translate("xxx", "welcome", Options{})
		expected := "Welcome"

		if got != expected {
			t.Errorf("expected %s; got %s", expected, got)
		}
	})

	t.Run("the custom fallback language should be used if set", func(t *testing.T) {
		i18n := NewI18n(Config{FallbackLanguageName: "es"})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Bienvenido",
			},
		})

		got := i18n.Translate("xxx", "welcome", Options{})
		expected := "Bienvenido"

		if got != expected {
			t.Errorf("expected %s; got %s", expected, got)
		}
	})

	t.Run("should fallback complex use case", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:      "welcomefallbacked",
				ManyMale: "Welcome, you have {{.EmailQty}} emails, sir {{.Name}}",
			},
		})

		got := i18n.Translate("xxx", "welcomefallbacked", Options{
			Count:  createPtr(5),
			Gender: createPtr("male"),
			Data: map[string]string{
				"EmailQty": "5",
				"Name":     "John",
			},
		})
		expected := "Welcome, you have 5 emails, sir John"

		if got != expected {
			t.Errorf("expected %s; got %s", expected, got)
		}
	})

	t.Run("the pluralization should work with default options (only one and many)", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "print_emails",
				Default: "You have emails",
				One:     "You have one email",
				Many:    "You have {{.EmailQty}} emails",
			},
		})

		tests := []struct {
			lang     string
			key      string
			options  Options
			expected string
		}{
			{"en", "print_emails", Options{}, "You have emails"},
			{"en", "print_emails", Options{Count: createPtr(0), Data: Data{"EmailQty": 0}}, "You have 0 emails"},
			{"en", "print_emails", Options{Count: createPtr(1), Data: Data{"EmailQty": 1}}, "You have one email"},
			{"en", "print_emails", Options{Count: createPtr(5), Data: Data{"EmailQty": 5}}, "You have 5 emails"},
		}

		for _, test := range tests {
			t.Run(test.key, func(t *testing.T) {
				got := i18n.Translate(test.lang, test.key, test.options)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("the pluralization should work with custom pluralization function", func(t *testing.T) {
		i18n := NewI18n(Config{})

		myPluralizationFn := func(count int) string {
			if count == 0 {
				return "Zero"
			}
			if count == 1 {
				return "One"
			}
			if count == 2 {
				return "Two"
			}
			if count == 3 {
				return "Few"
			}
			return "Many"
		}

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "print_emails",
				Default: "You have emails",
				Zero:    "You have no emails",
				One:     "You have one email",
				Two:     "You have two emails",
				Few:     "You have three emails",
				Many:    "You have {{.EmailQty}} emails",
			},
		})

		i18n.SetPluralizationFunc("en", myPluralizationFn)

		tests := []struct {
			lang     string
			key      string
			options  Options
			expected string
		}{
			{"en", "print_emails", Options{}, "You have emails"},
			{"en", "print_emails", Options{Count: createPtr(0), Data: Data{"EmailQty": 0}}, "You have no emails"},
			{"en", "print_emails", Options{Count: createPtr(1), Data: Data{"EmailQty": 1}}, "You have one email"},
			{"en", "print_emails", Options{Count: createPtr(2), Data: Data{"EmailQty": 2}}, "You have two emails"},
			{"en", "print_emails", Options{Count: createPtr(3), Data: Data{"EmailQty": 3}}, "You have three emails"},
			{"en", "print_emails", Options{Count: createPtr(100), Data: Data{"EmailQty": 100}}, "You have 100 emails"},
		}

		for _, test := range tests {
			t.Run(test.key, func(t *testing.T) {
				got := i18n.Translate(test.lang, test.key, test.options)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("should handle gendered translations in multiple languages", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:       "hello_message",
				Default:   "Hello",
				Male:      "Hello sir",
				Female:    "Hello ma'am",
				NonBinary: "Hello friend",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:       "hello_message",
				Default:   "Hola",
				Male:      "Hola amigo",
				Female:    "Hola amiga",
				NonBinary: "Hola amigue",
			},
		})

		tests := []struct {
			lang     string
			key      string
			options  Options
			expected string
		}{
			{"en", "hello_message", Options{}, "Hello"},
			{"en", "hello_message", Options{Gender: createPtr("somethingelse")}, "Hello"},
			{"en", "hello_message", Options{Gender: createPtr("male")}, "Hello sir"},
			{"en", "hello_message", Options{Gender: createPtr("female")}, "Hello ma'am"},
			{"en", "hello_message", Options{Gender: createPtr("nonbinary")}, "Hello friend"},
			{"en", "hello_message", Options{Gender: createPtr("non-binary")}, "Hello friend"},
			{"es", "hello_message", Options{}, "Hola"},
			{"es", "hello_message", Options{Gender: createPtr("somethingelse")}, "Hola"},
			{"es", "hello_message", Options{Gender: createPtr("male")}, "Hola amigo"},
			{"es", "hello_message", Options{Gender: createPtr("female")}, "Hola amiga"},
			{"es", "hello_message", Options{Gender: createPtr("nonbinary")}, "Hola amigue"},
			{"es", "hello_message", Options{Gender: createPtr("non-binary")}, "Hola amigue"},
		}

		for _, test := range tests {
			t.Run(test.key, func(t *testing.T) {
				got := i18n.Translate(test.lang, test.key, test.options)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("should handle complex pluralized-gendered translations in multiple languages", func(t *testing.T) {
		i18n := NewI18n(Config{})

		myPluralizationFn := func(count int) string {
			if count == 0 {
				return "Zero"
			}
			if count == 1 {
				return "One"
			}
			if count == 2 {
				return "Two"
			}
			if count == 3 {
				return "Few"
			}
			return "Many"
		}

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:           "hello_emails",
				Default:       "Hello, you have emails",
				ZeroMale:      "Hello, you have no emails, sir",
				OneMale:       "Hello, you have one email, sir",
				TwoMale:       "Hello, you have two emails, sir",
				FewMale:       "Hello, you have three emails, sir",
				ManyMale:      "Hello, you have {{.EmailQty}} emails, sir",
				ZeroFemale:    "Hello, you have no emails, ma'am",
				OneFemale:     "Hello, you have one email, ma'am",
				TwoFemale:     "Hello, you have two emails, ma'am",
				FewFemale:     "Hello, you have three emails, ma'am",
				ManyFemale:    "Hello, you have {{.EmailQty}} emails, ma'am",
				ZeroNonBinary: "Hello, you have no emails, friend",
				OneNonBinary:  "Hello, you have one email, friend",
				TwoNonBinary:  "Hello, you have two emails, friend",
				FewNonBinary:  "Hello, you have three emails, friend",
				ManyNonBinary: "Hello, you have {{.EmailQty}} emails, friend",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:           "hello_emails",
				Default:       "Hola, tienes correos",
				ZeroMale:      "Hola, no tienes correos, amigo",
				OneMale:       "Hola, tienes un correo, amigo",
				TwoMale:       "Hola, tienes dos correos, amigo",
				FewMale:       "Hola, tienes tres correos, amigo",
				ManyMale:      "Hola, tienes {{.EmailQty}} correos, amigo",
				ZeroFemale:    "Hola, no tienes correos, amiga",
				OneFemale:     "Hola, tienes un correo, amiga",
				TwoFemale:     "Hola, tienes dos correos, amiga",
				FewFemale:     "Hola, tienes tres correos, amiga",
				ManyFemale:    "Hola, tienes {{.EmailQty}} correos, amiga",
				ZeroNonBinary: "Hola, no tienes correos, amigue",
				OneNonBinary:  "Hola, tienes un correo, amigue",
				TwoNonBinary:  "Hola, tienes dos correos, amigue",
				FewNonBinary:  "Hola, tienes tres correos, amigue",
				ManyNonBinary: "Hola, tienes {{.EmailQty}} correos, amigue",
			},
		})

		i18n.SetPluralizationFunc("en", myPluralizationFn)
		i18n.SetPluralizationFunc("es", myPluralizationFn)

		tests := []struct {
			lang     string
			key      string
			options  Options
			expected string
		}{
			{"en", "hello_emails", Options{}, "Hello, you have emails"},
			{"en", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(0)}, "Hello, you have no emails, sir"},
			{"en", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(1)}, "Hello, you have one email, sir"},
			{"en", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(2)}, "Hello, you have two emails, sir"},
			{"en", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(3)}, "Hello, you have three emails, sir"},
			{"en", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(100), Data: Data{"EmailQty": 100}}, "Hello, you have 100 emails, sir"},
			{"en", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(0)}, "Hello, you have no emails, ma'am"},
			{"en", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(1)}, "Hello, you have one email, ma'am"},
			{"en", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(2)}, "Hello, you have two emails, ma'am"},
			{"en", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(3)}, "Hello, you have three emails, ma'am"},
			{"en", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(100), Data: Data{"EmailQty": 100}}, "Hello, you have 100 emails, ma'am"},
			{"en", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(0)}, "Hello, you have no emails, friend"},
			{"en", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(1)}, "Hello, you have one email, friend"},
			{"en", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(2)}, "Hello, you have two emails, friend"},
			{"en", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(3)}, "Hello, you have three emails, friend"},
			{"en", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(100), Data: Data{"EmailQty": 100}}, "Hello, you have 100 emails, friend"},
			{"es", "hello_emails", Options{}, "Hola, tienes correos"},
			{"es", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(0)}, "Hola, no tienes correos, amigo"},
			{"es", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(1)}, "Hola, tienes un correo, amigo"},
			{"es", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(2)}, "Hola, tienes dos correos, amigo"},
			{"es", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(3)}, "Hola, tienes tres correos, amigo"},
			{"es", "hello_emails", Options{Gender: createPtr("male"), Count: createPtr(100), Data: Data{"EmailQty": 100}}, "Hola, tienes 100 correos, amigo"},
			{"es", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(0)}, "Hola, no tienes correos, amiga"},
			{"es", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(1)}, "Hola, tienes un correo, amiga"},
			{"es", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(2)}, "Hola, tienes dos correos, amiga"},
			{"es", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(3)}, "Hola, tienes tres correos, amiga"},
			{"es", "hello_emails", Options{Gender: createPtr("female"), Count: createPtr(100), Data: Data{"EmailQty": 100}}, "Hola, tienes 100 correos, amiga"},
			{"es", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(0)}, "Hola, no tienes correos, amigue"},
			{"es", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(1)}, "Hola, tienes un correo, amigue"},
			{"es", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(2)}, "Hola, tienes dos correos, amigue"},
			{"es", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(3)}, "Hola, tienes tres correos, amigue"},
			{"es", "hello_emails", Options{Gender: createPtr("nonbinary"), Count: createPtr(100), Data: Data{"EmailQty": 100}}, "Hola, tienes 100 correos, amigue"},
		}

		for _, test := range tests {
			t.Run(test.key, func(t *testing.T) {
				got := i18n.Translate(test.lang, test.key, test.options)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("should return empty strings on edge incorrect cases", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "hello_message",
				Default: "Hello",
				Male:    "Hello sir",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:     "hello_message",
				Default: "Hola",
				Male:    "Hola amigo",
			},
		})

		tests := []struct {
			lang     string
			key      string
			options  Options
			expected string
		}{
			{"xxx", "xxx", Options{}, ""},
			{"en", "xxx", Options{}, ""},
			{"es", "xxx", Options{}, ""},
		}

		for _, test := range tests {
			t.Run(test.key, func(t *testing.T) {
				got := i18n.Translate(test.lang, test.key, test.options)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("T should be alias for Translate", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome",
			},
		})

		got := i18n.Translate("en", "welcome", Options{})
		gotT := i18n.T("en", "welcome", Options{})
		expected := "Welcome"

		if got != expected {
			t.Errorf("expected %s; got %s", expected, got)
		}

		if gotT != expected {
			t.Errorf("expected %s; got %s", expected, gotT)
		}
	})

	t.Run("test basic templating function", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome",
			},
		})

		templateFunc := i18n.NewTemplatingTranslateFunc()

		result := execI18nTemplate(templateFunc, `{{Translate "lang" "en" "key" "welcome"}}`)

		if result != "Welcome" {
			t.Errorf("expected %s; got %s", "Welcome", result)
		}
	})

	t.Run("test templating function with multiple interpolations", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome {{.Name}} {{.SurName}}",
			},
		})

		templateFunc := i18n.NewTemplatingTranslateFunc()

		result := execI18nTemplate(templateFunc, `{{Translate "lang" "en" "key" "welcome" "Name" "John" "SurName" "Doe" }}`)
		expected := "Welcome John Doe"

		if result != expected {
			t.Errorf("expected %s; got %s", expected, result)
		}
	})

	t.Run("test complex templating function", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:        "welcome_emails",
				Default:    "Welcome, you have emails",
				One:        "Welcome, you have one email",
				Many:       "Welcome, you have {{.EmailQty}} emails",
				OneMale:    "Welcome, you have one email, sir",
				ManyMale:   "Welcome, you have {{.EmailQty}} emails, sir",
				OneFemale:  "Welcome, you have one email, madam",
				ManyFemale: "Welcome, you have {{.EmailQty}} emails, madam",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:        "welcome_emails",
				Default:    "Bienvenido, tienes correos",
				One:        "Bienvenido, tienes un correo",
				Many:       "Bienvenido, tienes {{.EmailQty}} correos",
				OneMale:    "Bienvenido, tienes un correo, amigo",
				ManyMale:   "Bienvenido, tienes {{.EmailQty}} correos, amigo",
				OneFemale:  "Bienvenida, tienes un correo, amiga",
				ManyFemale: "Bienvenida, tienes {{.EmailQty}} correos, amiga",
			},
		})

		templateFunc := i18n.NewTemplatingTranslateFunc()

		tests := []struct {
			templateText string
			expected     string
		}{
			{`{{Translate "lang" "en" "key" "welcome_emails"}}`, "Welcome, you have emails"},
			{`{{Translate "lang" "en" "key" "welcome_emails" "count" "1"}}`, "Welcome, you have one email"},
			{`{{Translate "lang" "en" "key" "welcome_emails" "count" "5" "EmailQty" "5"}}`, "Welcome, you have 5 emails"},
			{`{{Translate "lang" "en" "key" "welcome_emails" "count" "1" "gender" "male"}}`, "Welcome, you have one email, sir"},
			{`{{Translate "lang" "en" "key" "welcome_emails" "count" "5" "gender" "male" "EmailQty" "5"}}`, "Welcome, you have 5 emails, sir"},
			{`{{Translate "lang" "en" "key" "welcome_emails" "count" "1" "gender" "female"}}`, "Welcome, you have one email, madam"},
			{`{{Translate "lang" "en" "key" "welcome_emails" "count" "5" "gender" "female" "EmailQty" "5"}}`, "Welcome, you have 5 emails, madam"},
			{`{{Translate "lang" "es" "key" "welcome_emails"}}`, "Bienvenido, tienes correos"},
			{`{{Translate "lang" "es" "key" "welcome_emails" "count" "1"}}`, "Bienvenido, tienes un correo"},
			{`{{Translate "lang" "es" "key" "welcome_emails" "count" "5" "EmailQty" "5"}}`, "Bienvenido, tienes 5 correos"},
			{`{{Translate "lang" "es" "key" "welcome_emails" "count" "1" "gender" "male"}}`, "Bienvenido, tienes un correo, amigo"},
			{`{{Translate "lang" "es" "key" "welcome_emails" "count" "5" "gender" "male" "EmailQty" "5"}}`, "Bienvenido, tienes 5 correos, amigo"},
			{`{{Translate "lang" "es" "key" "welcome_emails" "count" "1" "gender" "female"}}`, "Bienvenida, tienes un correo, amiga"},
			{`{{Translate "lang" "es" "key" "welcome_emails" "count" "5" "gender" "female" "EmailQty" "5"}}`, "Bienvenida, tienes 5 correos, amiga"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("complex template test - %v", i), func(t *testing.T) {
				got := execI18nTemplate(templateFunc, test.templateText)
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})

	t.Run("test translate function creation for specific lang", func(t *testing.T) {
		i18n := NewI18n(Config{})

		i18n.AddLanguage("en", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Welcome",
			},
		})

		i18n.AddLanguage("es", TranslateStrings{
			TranslateString{
				Key:     "welcome",
				Default: "Bienvenido",
			},
		})

		enTranslateFunc := i18n.NewLangTranslateFunc("en")
		esTranslateFunc := i18n.NewLangTranslateFunc("es")

		tests := []struct {
			translateFunc func(translateKey string, options Options) string
			expected      string
		}{
			{enTranslateFunc, "Welcome"},
			{esTranslateFunc, "Bienvenido"},
		}

		for i, test := range tests {
			t.Run(fmt.Sprintf("translate function creation test - %v", i), func(t *testing.T) {
				got := test.translateFunc("welcome", Options{})
				if got != test.expected {
					t.Errorf("expected %s; got %s", test.expected, got)
				}
			})
		}
	})
}

// Function to create pointer to a value
func createPtr[T string | int](s T) *T {
	return &s
}

// Function to create template and execute it
func execI18nTemplate(fn func(...any) string, templateText string) string {
	tmpl, err := template.New("main").Funcs(
		template.FuncMap{
			"Translate": fn,
		},
	).Parse(templateText)

	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, nil)
	if err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}

	return result.String()
}
