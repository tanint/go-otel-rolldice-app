package events

type RollEvent struct {
	RollID    string `json:"roll_id"`
	Result    int    `json:"result"`
	Timestamp string `json:"timestamp"`
}
