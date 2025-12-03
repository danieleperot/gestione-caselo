package emails

import (
	"bytes"
	"fmt"
	"text/template"
)

// RenderTemplate parses and executes a template with the given data
func RenderTemplate(templateName, templateContent string, data interface{}) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}
