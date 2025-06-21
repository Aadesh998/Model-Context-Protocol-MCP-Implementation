package chat

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Conversation struct {
	Messages []Message `json:"messages"`
}

func SaveConversation(messages []Message) error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("chat_history/conversation_%s.json", timestamp)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(Conversation{Messages: messages})
}
