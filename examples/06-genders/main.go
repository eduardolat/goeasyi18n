package main

import (
	"fmt"

	"github.com/eduardolat/goeasyi18n"
)

/*
Gender Handling: Why is it useful?

Let's say your app has a feature that says, "John liked your post" or "Emily liked your post."
In some languages, the verb "liked" might change based on the gender of the person who
liked the post.

Example:
- English: "He liked your post" vs "She liked your post"
- Spanish: "A él le gustó tu publicación" vs "A ella le gustó tu publicación"

With gender-specific translations, you can easily adapt the sentence structure to fit the
gender, making your app more linguistically accurate and inclusive.

No need for messy if-else statements to handle gender variations. Just set it up once, and
the library takes care of the rest!
*/

func main() {
	// 1. Create a new i18n instance
	i18n := goeasyi18n.NewI18n(goeasyi18n.Config{})

	// 2. Create your translations
	// If something goes wrong, the default value is used
	// The gender keys are Male, Female and NonBinary
	enTranslations := goeasyi18n.TranslateStrings{
		{
			Key:       "friend_emails",
			Default:   "Hello, your friend have emails",
			Male:      "Hello, he has emails",
			Female:    "Hello, she has emails",
			NonBinary: "Hello, your friend have emails",
		},
	}

	esTranslations := goeasyi18n.TranslateStrings{
		{
			Key:       "friend_emails",
			Default:   "Hola, tu amigo tiene correos",
			Male:      "Hola, él tiene correos",
			Female:    "Hola, ella tiene correos",
			NonBinary: "Hola, tu amigue tiene correos",
		},
	}

	// 3. Add your languages with their translations
	i18n.AddLanguage("en", enTranslations)
	i18n.AddLanguage("es", esTranslations)

	// 4. Create the Options
	// The Gender field is a *string that contains the gender to use
	// Here you can use male, female, nonbinary or non-binary
	maleText := "male"
	femaleText := "female"
	nonbinaryText := "nonbinary"

	maleOptions := goeasyi18n.Options{
		Gender: &maleText,
	}

	femaleOptions := goeasyi18n.Options{
		Gender: &femaleText,
	}

	nonbinaryOptions := goeasyi18n.Options{
		Gender: &nonbinaryText,
	}

	// 5. You are done! 🎉 Just get that translations!
	ten1 := i18n.T("en", "friend_emails", maleOptions)
	ten2 := i18n.T("en", "friend_emails", femaleOptions)
	ten3 := i18n.T("en", "friend_emails", nonbinaryOptions)

	tes1 := i18n.T("es", "friend_emails", maleOptions)
	tes2 := i18n.T("es", "friend_emails", femaleOptions)
	tes3 := i18n.T("es", "friend_emails", nonbinaryOptions)

	fmt.Println(ten1)
	fmt.Println(ten2)
	fmt.Println(ten3)
	fmt.Println(tes1)
	fmt.Println(tes2)
	fmt.Println(tes3)

	/*
		Prints:
		Hello, he has emails
		Hello, she has emails
		Hello, your friend have emails
		Hola, él tiene correos
		Hola, ella tiene correos
		Hola, tu amigue tiene correos
	*/
}
