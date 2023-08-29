package main

import (
	"bytes"
	"html/template"
	"os"
)

func ExecuteTemplate(templatePath string, data any) (string, error) {

	// Prepare func map
	translateFunc := i18n.NewTemplatingTranslateFunc()
	funcs := template.FuncMap{
		"Translate": translateFunc, // You can use any name you want, for example: "T"
	}

	// Read template file
	fileContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", err
	}

	// Execute template
	tmpl := template.Must(template.New("test").Funcs(funcs).Parse(string(fileContent)))
	result := new(bytes.Buffer)
	tmpl.Execute(result, data)

	if err != nil {
		return "", err
	}

	return result.String(), nil
}
