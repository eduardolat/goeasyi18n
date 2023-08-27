package goeasyi18n

import "testing"

func TestExecuteTemplate(t *testing.T) {

	t.Run("should execute a template", func(t *testing.T) {
		executed := ExecuteTemplate("Hello {{.Name}}", struct{ Name string }{Name: "World"})
		expected := "Hello World"

		if executed != expected {
			t.Errorf("Expected %s, got %s", expected, executed)
		}
	})

	t.Run("should execute a template with two interpolations", func(t *testing.T) {
		data := struct{ FirstName, SurName string }{FirstName: "John", SurName: "Doe"}
		executed := ExecuteTemplate("Hello {{.FirstName}} {{.SurName}}", data)
		expected := "Hello John Doe"

		if executed != expected {
			t.Errorf("Expected %s, got %s", expected, executed)
		}
	})

}
