package main

import (
	"os"
	"text/template"
)

/*
	The templating feature allows you to make translations inside your templates.
	It works with both text/template and html/template.

	To pass the language use: "lang" "xxx"
	To pass the key use:      "key" "xxx"
	To pass the count use:    "count" "xxx"
	To pass the gender use:   "gender" "male/female/nonbinary"

	All other "key" "value" pairs will be converted to strings and passed to
	the translation as variables.
*/

const templateText = `Welcome to Go Easy i18n!

{{Translate "lang" "en" "key" "hello_message" "Name" "John Doe"}}
{{Translate "lang" "es" "key" "hello_message" "Name" "John Doe"}}

{{Translate "lang" "en" "key" "unread_messages" "count" "1" "Qty" "1"}}
{{Translate "lang" "en" "key" "unread_messages" "count" "10" "Qty" "10"}}
{{Translate "lang" "es" "key" "unread_messages" "count" "1" "Qty" "1"}}
{{Translate "lang" "es" "key" "unread_messages" "count" "10" "Qty" "10"}}

{{Translate "lang" "en" "key" "friend_update" "gender" "male" "Name" "John"}}
{{Translate "lang" "en" "key" "friend_update" "gender" "female" "Name" "Jane"}}
{{Translate "lang" "en" "key" "friend_update" "gender" "nonbinary" "Name" "Jane"}}
{{Translate "lang" "es" "key" "friend_update" "gender" "male" "Name" "John"}}
{{Translate "lang" "es" "key" "friend_update" "gender" "female" "Name" "Jane"}}
{{Translate "lang" "es" "key" "friend_update" "gender" "nonbinary" "Name" "Jane"}}

{{Translate "lang" "en" "key" "friend_request" "gender" "male" "count" "1" "Qty" "1"}}
{{Translate "lang" "en" "key" "friend_request" "gender" "female" "count" "10" "Qty" "10"}}
{{Translate "lang" "es" "key" "friend_request" "gender" "male" "count" "1" "Qty" "1"}}
{{Translate "lang" "es" "key" "friend_request" "gender" "female" "count" "10" "Qty" "10"}}
`

func main() {
	// 1. Initialize the i18n instance
	InitializeI18n()

	// 2. Create function to pass to the template func map
	translateFunc := i18n.NewTemplatingTranslateFunc()

	// 3. Create the template and pass the created function to the func map
	tmpl := template.Must(template.New("test").Funcs(template.FuncMap{
		"Translate": translateFunc, // You can use any name you want, for example: "T"
	}).Parse(templateText))

	// 4. Execute the template
	tmpl.Execute(os.Stdout, nil)
}
