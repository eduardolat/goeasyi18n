package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Get the lang query param, you are responsible in your app to get the lang
	// In this example is just passed as a query param but you can use a cookie
	// or an url segment or whatever you want
	lang := "en" // Default language
	if r.URL.Query().Get("lang") != "" {
		lang = r.URL.Query().Get("lang")
	}

	// Check if language is supported (you can return a 404 error if you want)
	hasLanguage := i18n.HasLanguage(lang)

	// Execute the template using the detected language
	templateData := map[string]any{
		"Lang":        lang,
		"HasLanguage": hasLanguage,
	}

	// Instead of using the {{Translate ...}} inside the template, you can also
	// translate the text here and pass it to the template for example:
	// someTranslation := i18n.Translate(lang, "some_translation_key")

	responseText, err := ExecuteTemplate("./views/index.html", templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Write([]byte(responseText))
}

func main() {
	InitializeI18n()

	http.HandleFunc("/", handler)

	log.Println("Listening on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
