package models

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type PushMessage struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}
