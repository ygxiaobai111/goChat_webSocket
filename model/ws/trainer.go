package ws

type Trainer struct {
	Content   string `json:"content" bson:"content"`
	StartTime int64  `json:"startTime" bson:"startTime"` // 创建时间
	EndTime   int64  `json:"endTime" bson:"endTime"`     // 过期时间
	Read      uint   `json:"read" bson:"read"`
}

type Result struct {
	StartTime int64
	Msg       string
	Content   interface{}
	From      string
}
