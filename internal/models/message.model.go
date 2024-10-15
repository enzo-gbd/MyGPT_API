package models

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Message represents each session messages
// @Description Message hold the content, the sender (USER or GPT), and the date of the message
type Message struct {
	ID      uuid.UUID `gorm:"type:char(36);primary_key"`  // Unique identifier for the message
	Sender  string    `gorm:"type:varchar(255);not null"` // The Sender of the message
	Content string    `gorm:"type:text;not null"`         // The Content of the message
	Date    time.Time `gorm:"type:datetime;not null"`     // The sending Date of the message
}

// Validate performs validation on Message fields using ozzo-validation package.
// It ensures all necessary fields meet their expected format.
func (m Message) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Content, validation.Required, validation.Length(1, 10000)),
		validation.Field(&m.Sender, validation.Required, validation.In("USER", "GPT")),
		validation.Field(&m.Date, validation.Required),
	)
}

// BeforeCreate is a GORM hook that is called before a new Message record is created.
// It assigns a new UUID to the Message's ID.
func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

// String provides a string representation of the Message which includes
// personal details and contact information.
func (m Message) String() string {
	return fmt.Sprintf(`\"%s\", send by %s the %s`, m.Content, m.Sender, m.Date)
}

// IsEqualTo compares this Message instance with another to determine if they are identical.
func (m Message) IsEqualTo(other Message) bool {
	return m.Content == other.Content &&
		m.Sender == other.Sender &&
		m.Date == other.Date
}

// MessageInput represents the required fields for a Message creation.
// @Description Fields required to send a new Message.
type MessageInput struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

// Validate performs validation on MessageInput fields using ozzo-validation package.
// It ensures all necessary fields meet their expected format.
func (m MessageInput) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Sender, validation.Required, validation.Required, validation.In("USER", "GPT")),
		validation.Field(&m.Content, validation.Required, validation.Length(1, 10000)),
	)
}
