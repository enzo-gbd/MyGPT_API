package builders

import (
	"github.com/enzo-gbd/GBA/internal/models"
	"github.com/google/uuid"
	"time"
)

// UserBuilder is a builder used to construct a message model with various properties.
type MessageBuilder struct {
	m models.Message
}

// NewMessageBuilder initializes and returns a new MessageBuilder with default values for a message.
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		m: models.Message{
			ID:      uuid.New(),
			Content: "",
			Sender:  "USER",
			Date:    time.Now(),
		},
	}
}

// Build constructs and returns the message model configured by the builder.
func (mb *MessageBuilder) Build() models.Message { return mb.m }

// BuildMessageInput constructs and returns a MessageInput model based on the message model configured by the builder.
func (mb *MessageBuilder) BuildMessageInput() models.MessageInput {
	return models.MessageInput{
		Content: mb.m.Content,
		Sender:  mb.m.Sender,
	}
}

// WithBase sets the base message model and returns the UserBuilder.
func (mb *MessageBuilder) WithBase(message models.Message) *MessageBuilder {
	mb.m = message
	return mb
}

// WhereContent sets the Content of the message being built and returns the MessageBuilder.
func (mb *MessageBuilder) WhereContent(content string) *MessageBuilder {
	mb.m.Content = content
	return mb
}

// WhereSender sets the Sender of the message being built and returns the MessageBuilder.
func (mb *MessageBuilder) WhereSender(sender string) *MessageBuilder {
	mb.m.Sender = sender
	return mb
}

// WhereDate sets the Dqte of the message being built and returns the MessageBuilder.
func (mb *MessageBuilder) WhereDate(date time.Time) *MessageBuilder {
	mb.m.Date = date
	return mb
}
