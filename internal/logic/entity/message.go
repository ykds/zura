package entity

type Message struct {
	BaseModel
	SessionId int64  `json:"session_id"`
	Body      string `json:"body"`
	Read      bool   `json:"read"`
}
