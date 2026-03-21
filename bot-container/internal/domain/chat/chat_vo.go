package chatdomain

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

type Role string

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func NewMessage(role Role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}
