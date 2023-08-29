package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

/*
Custom Pluralization: Why is it useful? 🤔

Imagine you're building an app that shows the number of unread messages.
In English, you'd say "1 unread message" and "2 unread messages" - note the "s" at the end.
But what about languages where plural rules aren't so simple?

With custom pluralization, you can easily handle these cases without writing complex if-else
statements. Just define your plural rules once, and let the library do the heavy lifting!

This makes your code cleaner and your app more linguistically accurate. Win-win!

Don't get scared by the amount of code, it's just an example and in a real app you'll
probably set up this only once and then use it everywhere.

The available pluralization keys are:
- Zero
- One
- Two
- Few
- Many

Later in the translation strings you can use the keys to make your translations different
depending on the count and the key returned by the custom pluralization function.

Let's get started! 🚀
*/

// The custom pluralization works in the same way as the default pluralization
// The only difference is that you can define your own pluralization rules
// Start creating a function that returns a string and receives a int
func MyCustomPluralization(count int) string {
	// The count variable is the same as the Count field in the Options
	// The string returned by this function will be used as the pluralization key
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

func main() {
	// 1. Create a new i18n instance
	i18n := goeasyi18n.NewI18n(goeasyi18n.Config{})

	// 2. Create your translations
	// If something goes wrong, the default value is used
	// The default pluralization only works for one and many keys
	// If the count (later in the Options) is 1, the one key is used
	// If the count (later in the Options) is greater than 1, the many key is used
	enTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_emails",
			Default: "Hello, you have emails",
			Zero:    "Hello, you have no emails",
			One:     "Hello, you have one email",
			Two:     "Hello, you have two emails",
			Few:     "Hello, you have three emails",
			Many:    "Hello, you have {{.EmailQty}} emails",
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_emails",
			Default: "Hola, tienes correos",
			Zero:    "Hola, no tienes correos",
			One:     "Hola, tienes un correo",
			Two:     "Hola, tienes 2 correos",
			Few:     "Hola, tienes 3 correos",
			Many:    "Hola, tienes {{.EmailQty}} correos",
		},
	}

	// 3. Add your languages with their translations
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 4. Add your custom pluralization function
	// This method sets the pluralization function for the given language
	// In this case, we are setting the pluralization function for the "en" language
	// The "es" language will still use the default pluralization to see the differences
	i18n.SetPluralizationFunc("en", MyCustomPluralization)

	// 5. Create the Options (we are using a helper that lives at end of this file)
	zeroEmailOptions := MakeOptions(0)
	oneEmailOptions := MakeOptions(1)
	twoEmailOptions := MakeOptions(2)
	fewEmailsOptions := MakeOptions(3)
	manyEmailsOptions := MakeOptions(100)

	// 6. You are done! 🎉 Just get that translations!
	ten0 := i18n.T("en", "hello_emails", zeroEmailOptions)
	ten1 := i18n.T("en", "hello_emails", oneEmailOptions)
	ten2 := i18n.T("en", "hello_emails", twoEmailOptions)
	tenf := i18n.T("en", "hello_emails", fewEmailsOptions)
	tenm := i18n.T("en", "hello_emails", manyEmailsOptions)

	tes0 := i18n.T("es", "hello_emails", zeroEmailOptions)
	tes1 := i18n.T("es", "hello_emails", oneEmailOptions)
	tes2 := i18n.T("es", "hello_emails", twoEmailOptions)
	tesf := i18n.T("es", "hello_emails", fewEmailsOptions)
	tesm := i18n.T("es", "hello_emails", manyEmailsOptions)

	fmt.Println(ten0)
	fmt.Println(ten1)
	fmt.Println(ten2)
	fmt.Println(tenf)
	fmt.Println(tenm)
	fmt.Println(tes0)
	fmt.Println(tes1)
	fmt.Println(tes2)
	fmt.Println(tesf)
	fmt.Println(tesm)

	/*
		Prints:
		Hello, you have no emails
		Hello, you have one email
		Hello, you have two emails
		Hello, you have three emails
		Hello, you have 100 emails
		Hola, tienes 0 correos
		Hola, tienes un correo
		Hola, tienes 2 correos
		Hola, tienes 3 correos
		Hola, tienes 100 correos
	*/

	/*
		Note how the "en" language uses the custom pluralization function and the "es" language
		uses the default pluralization function. This is because we only set the custom
		pluralization function for the "en" language.

		You can also set the custom pluralization function for all the languages you want.

		The default pluralization function only checks if the count is 1, then returns
		the "One" key, otherwise it returns the "Many" key. Is that simple!
	*/
}

// Helper function to create the Options for pluralization
func MakeOptions(count int) goeasyi18n.Options {
	options := goeasyi18n.Options{
		Count: &count,
		Data:  map[string]any{"EmailQty": count},
	}
	return options
}
