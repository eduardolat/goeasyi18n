package goeasyi18n

import (
	"bytes"
	"html/template"
)

func ExecuteTemplate(templateStr string, data any) string {
	tmpl := template.Must(template.New("template").Parse(templateStr))

	b := new(bytes.Buffer)

	err := tmpl.Execute(b, data)
	if err != nil {
		return ""
	}

	return b.String()
}
