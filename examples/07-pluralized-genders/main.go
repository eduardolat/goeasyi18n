package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

/*
You have come very far!! 👏👏👏
Now that you know how to pluralize and manage gender in your translations,
let's do something interesting... Let's combine both features!

Combining both features you will gain acces to all these keys in your translations:

// For only pluralization
Zero, One, Two, Few, Many

// For only genders management
Male, Female, NonBinary

// For pluralization with male gender
ZeroMale, OneMale, TwoMale, FewMale, ManyMale

// For pluralization with female gender
ZeroFemale, OneFemale, TwoFemale, FewFemale, ManyFemale

// For pluralization with non binary gender
ZeroNonBinary, OneNonBinary, TwoNonBinary, FewNonBinary, ManyNonBinary
*/

func main() {
	// 1. Create a new i18n instance
	i18n := goeasyi18n.NewI18n(goeasyi18n.Config{})

	// 2. Create your translations
	// If something goes wrong, the default value is used
	enTranslations := goeasyi18n.TranslateStrings{
		{
			Key:        "friend_emails",
			Default:    "Hello, your friend have emails",
			OneMale:    "Hello, he one email",
			ManyFemale: "Hello, she has {{.EmailQty}} emails",
			// You can add as many combinations as you want
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:        "friend_emails",
			Default:    "Hola, tu amigo tiene correos",
			OneMale:    "Hola, él tiene un correo",
			ManyFemale: "Hola, ella tiene {{.EmailQty}} correos",
		},
	}

	// 3. Add your languages with their translations
	// Yoy can also add custom pluralization rules but for this
	// example we will use the default ones to keep it simple
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 4. Create the Options
	// The Gender field is a *string that contains the gender to use
	// Here you can use male, female, nonbinary or non-binary
	oneInt := 1
	manyInt := 10
	maleText := "male"
	femaleText := "female"

	oneMaleOptions := &goeasyi18n.Options{
		Gender: &maleText,
		Count:  &oneInt,
		Data: map[string]any{
			"EmailQty": oneInt,
		},
	}

	manyFemaleOptions := &goeasyi18n.Options{
		Gender: &femaleText,
		Count:  &manyInt,
		Data: map[string]any{
			"EmailQty": manyInt,
		},
	}

	// 5. You are done! 🎉 Just get that translations!
	ten1 := i18n.T("en", "friend_emails", oneMaleOptions)
	ten2 := i18n.T("en", "friend_emails", manyFemaleOptions)

	tes1 := i18n.T("es", "friend_emails", oneMaleOptions)
	tes2 := i18n.T("es", "friend_emails", manyFemaleOptions)

	fmt.Println(ten1)
	fmt.Println(ten2)
	fmt.Println(tes1)
	fmt.Println(tes2)

	/*
		Prints:
		Hello, he one email
		Hello, she has 10 emails
		Hola, él tiene un correo
		Hola, ella tiene 10 correos
	*/
}
