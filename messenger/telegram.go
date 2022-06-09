package messenger

// TeleMsg represents body for telegram message
type TeleMsg struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type TgMessenger interface {
	SendMessage()
}

type tgMessenger struct {
}

type messenger struct {
}
