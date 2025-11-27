package registry

import (
	"encoding/json"
	"fmt"

	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
)

type Constructor func(metadata map[string]interface{}) (message.EmailMessage, error)

type Registry struct {
	constructors map[string]Constructor
}

func New() *Registry {
	return &Registry{
		constructors: make(map[string]Constructor),
	}
}

func (r *Registry) Register(templateName string, constructor Constructor) {
	r.constructors[templateName] = constructor
}

func (r *Registry) Get(sqsMessageBody string) (message.EmailMessage, error) {
	var payload struct {
		Template string                 `json:"template"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := json.Unmarshal([]byte(sqsMessageBody), &payload); err != nil {
		return nil, fmt.Errorf("failed to parse SQS message: %w", err)
	}

	if payload.Template == "" {
		return nil, fmt.Errorf("template field is required")
	}

	constructor, exists := r.constructors[payload.Template]
	if !exists {
		return nil, fmt.Errorf("unknown template: %s", payload.Template)
	}

	return constructor(payload.Metadata)
}
