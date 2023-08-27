package goeasyi18n

import (
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
			options  TranslateOptions
			expected string
		}{
			{"en", "welcome", TranslateOptions{}, "Welcome"},
			{"en", "welcome", TranslateOptions{Gender: CreatePtr("male")}, "Welcome, sir"},
			{"en", "welcome", TranslateOptions{Gender: CreatePtr("female")}, "Welcome, ma'am"},
			{"en", "emails", TranslateOptions{Count: CreatePtr(1)}, "You have one email"},
			{"en", "emails", TranslateOptions{Count: CreatePtr(5)}, "You have many emails"},
			{"es", "welcome", TranslateOptions{}, "Bienvenido"},
			{"es", "welcome", TranslateOptions{Gender: CreatePtr("Male")}, "Bienvenido, señor"},
			{"es", "welcome", TranslateOptions{Gender: CreatePtr("Female")}, "Bienvenida, señora"},
			{"es", "emails", TranslateOptions{Count: CreatePtr(1)}, "Tienes un correo"},
			{"es", "emails", TranslateOptions{Count: CreatePtr(5)}, "Tienes muchos correos"},
			// Test fallback language
			{"xxx", "welcome", TranslateOptions{}, "Welcome"},
			// Test fallback key
			{"en", "xxx", TranslateOptions{}, ""},
			// Test data interpolation
			{"en", "greetings", TranslateOptions{Data: map[string]string{"Name": "John"}}, "Hello John"},
			{"es", "greetings", TranslateOptions{Data: map[string]string{"Name": "John"}}, "Hola John"},
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

		got := i18n.Translate("xxx", "welcome", TranslateOptions{})
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

		got := i18n.Translate("xxx", "welcome", TranslateOptions{})
		expected := "Bienvenido"

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
			options  TranslateOptions
			expected string
		}{
			{"en", "print_emails", TranslateOptions{}, "You have emails"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(0), Data: Data{"EmailQty": 0}}, "You have 0 emails"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(1), Data: Data{"EmailQty": 1}}, "You have one email"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(5), Data: Data{"EmailQty": 5}}, "You have 5 emails"},
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
			options  TranslateOptions
			expected string
		}{
			{"en", "print_emails", TranslateOptions{}, "You have emails"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(0), Data: Data{"EmailQty": 0}}, "You have no emails"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(1), Data: Data{"EmailQty": 1}}, "You have one email"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(2), Data: Data{"EmailQty": 2}}, "You have two emails"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(3), Data: Data{"EmailQty": 3}}, "You have three emails"},
			{"en", "print_emails", TranslateOptions{Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "You have 100 emails"},
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
			options  TranslateOptions
			expected string
		}{
			{"en", "hello_message", TranslateOptions{}, "Hello"},
			{"en", "hello_message", TranslateOptions{Gender: CreatePtr("somethingelse")}, "Hello"},
			{"en", "hello_message", TranslateOptions{Gender: CreatePtr("male")}, "Hello sir"},
			{"en", "hello_message", TranslateOptions{Gender: CreatePtr("female")}, "Hello ma'am"},
			{"en", "hello_message", TranslateOptions{Gender: CreatePtr("nonbinary")}, "Hello friend"},
			{"en", "hello_message", TranslateOptions{Gender: CreatePtr("non-binary")}, "Hello friend"},
			{"es", "hello_message", TranslateOptions{}, "Hola"},
			{"es", "hello_message", TranslateOptions{Gender: CreatePtr("somethingelse")}, "Hola"},
			{"es", "hello_message", TranslateOptions{Gender: CreatePtr("male")}, "Hola amigo"},
			{"es", "hello_message", TranslateOptions{Gender: CreatePtr("female")}, "Hola amiga"},
			{"es", "hello_message", TranslateOptions{Gender: CreatePtr("nonbinary")}, "Hola amigue"},
			{"es", "hello_message", TranslateOptions{Gender: CreatePtr("non-binary")}, "Hola amigue"},
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
			options  TranslateOptions
			expected string
		}{
			{"en", "hello_emails", TranslateOptions{}, "Hello, you have emails"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(0)}, "Hello, you have no emails, sir"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(1)}, "Hello, you have one email, sir"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(2)}, "Hello, you have two emails, sir"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(3)}, "Hello, you have three emails, sir"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "Hello, you have 100 emails, sir"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(0)}, "Hello, you have no emails, ma'am"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(1)}, "Hello, you have one email, ma'am"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(2)}, "Hello, you have two emails, ma'am"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(3)}, "Hello, you have three emails, ma'am"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "Hello, you have 100 emails, ma'am"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(0)}, "Hello, you have no emails, friend"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(1)}, "Hello, you have one email, friend"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(2)}, "Hello, you have two emails, friend"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(3)}, "Hello, you have three emails, friend"},
			{"en", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "Hello, you have 100 emails, friend"},
			{"es", "hello_emails", TranslateOptions{}, "Hola, tienes correos"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(0)}, "Hola, no tienes correos, amigo"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(1)}, "Hola, tienes un correo, amigo"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(2)}, "Hola, tienes dos correos, amigo"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(3)}, "Hola, tienes tres correos, amigo"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("male"), Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "Hola, tienes 100 correos, amigo"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(0)}, "Hola, no tienes correos, amiga"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(1)}, "Hola, tienes un correo, amiga"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(2)}, "Hola, tienes dos correos, amiga"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(3)}, "Hola, tienes tres correos, amiga"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("female"), Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "Hola, tienes 100 correos, amiga"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(0)}, "Hola, no tienes correos, amigue"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(1)}, "Hola, tienes un correo, amigue"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(2)}, "Hola, tienes dos correos, amigue"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(3)}, "Hola, tienes tres correos, amigue"},
			{"es", "hello_emails", TranslateOptions{Gender: CreatePtr("nonbinary"), Count: CreatePtr(100), Data: Data{"EmailQty": 100}}, "Hola, tienes 100 correos, amigue"},
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
			options  TranslateOptions
			expected string
		}{
			{"xxx", "xxx", TranslateOptions{}, ""},
			{"en", "xxx", TranslateOptions{}, ""},
			{"es", "xxx", TranslateOptions{}, ""},
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

		got := i18n.Translate("en", "welcome", TranslateOptions{})
		gotT := i18n.T("en", "welcome", TranslateOptions{})
		expected := "Welcome"

		if got != expected {
			t.Errorf("expected %s; got %s", expected, got)
		}

		if gotT != expected {
			t.Errorf("expected %s; got %s", expected, gotT)
		}
	})
}

// Function to create pointer to a value
func CreatePtr[T string | int](s T) *T {
	return &s
}
