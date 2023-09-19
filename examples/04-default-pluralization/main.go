package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

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
			One:     "Hello, you have one email",
			Many:    "Hello, you have {{.EmailQty}} emails",
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:     "hello_emails",
			Default: "Hola, tienes correos",
			One:     "Hola, tienes un correo",
			Many:    "Hola, tienes {{.EmailQty}} correos",
		},
	}

	// 3. Add your languages with their translations
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 4. Create the Options
	// The Count field is a *int that contains a number which is used to
	// select the correct pluralization key
	oneEmail := 1 // Get this value from your database or wherever you want
	oneEmailOptions := &goeasyi18n.Options{
		Count: &oneEmail,
		Data: map[string]any{
			"EmailQty": oneEmail,
		},
	}

	manyEmails := 5 // Get this value from your database or wherever you want
	manyEmailsOptions := &goeasyi18n.Options{
		Count: &manyEmails,
		Data: map[string]any{
			"EmailQty": manyEmails,
		},
	}

	// 5. You are done! 🎉 Just get that translations!
	ten1 := i18n.T("en", "hello_emails", oneEmailOptions)
	ten2 := i18n.T("en", "hello_emails", manyEmailsOptions)

	tes1 := i18n.T("es", "hello_emails", oneEmailOptions)
	tes2 := i18n.T("es", "hello_emails", manyEmailsOptions)

	fmt.Println(ten1)
	fmt.Println(ten2)
	fmt.Println(tes1)
	fmt.Println(tes2)

	/*
		Prints:
		Hello, you have one email
		Hello, you have 5 emails
		Hola, tienes un correo
		Hola, tienes 5 correos
	*/
}
