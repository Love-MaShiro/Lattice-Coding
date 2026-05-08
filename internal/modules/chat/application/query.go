package application

type SessionPageQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type MessageQuery struct {
	SessionID uint64 `json:"session_id"`
	Limit     int    `json:"limit"`
}
