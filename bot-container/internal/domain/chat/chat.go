package chatdomain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrDuplicatedKey  = errors.New("username already exists")
	ErrRecordNotFount = errors.New("record not found")
	ErrEmptyContent   = errors.New("empty content")
)

type Chat struct {
	id      int64
	botID   uuid.UUID
	history *History
}

type History struct {
	messages []Message
}

func NewChat(id int64, botID uuid.UUID, systemPrompt string) *Chat {
	msg := NewMessage(RoleSystem, systemPrompt)
	history := NewHistory([]Message{msg})

	return &Chat{
		id:      id,
		botID:   botID,
		history: history,
	}
}

// USE ONLY FOR CREATING CHAT FROM SERVICE STORAGE!!!
func NewChatFromStorage(id int64, botID uuid.UUID, history *History) *Chat {
	return &Chat{
		id:      id,
		botID:   botID,
		history: history,
	}
}

func NewHistory(msgs []Message) *History {
	messages := make([]Message, 0, len(msgs))
	for _, msg := range msgs {
		messages = append(messages, msg)
	}

	return &History{
		messages: messages,
	}
}

func (c *Chat) ID() int64 {
	return c.id
}

func (c *Chat) BotID() uuid.UUID {
	return c.botID
}

func (c *Chat) History() History {
	return *c.history
}

func (c *Chat) AppendHistory(role Role, content string) error {
	if content == "" {
		return ErrEmptyContent
	}

	msg := NewMessage(role, content)
	c.history.messages = append(c.history.messages, msg)

	return nil
}

func (c *Chat) Messages() []Message {
	return c.history.messages
}

func (h *History) Messages() []Message {
	return h.messages
}
