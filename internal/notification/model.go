package notification

const (
	// supported message types
	NotificationTypeSubscribe = "subscribe"

	// supported channels
	NotificationValueOutcomes    = "outcomes"
	NotificationValueLeaderboard = "leaderboard"
)

// Message simple model for communication through websocket
type Message struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
