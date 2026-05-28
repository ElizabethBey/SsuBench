package model

import "time"

type Payment struct {
	ID         int       `json:"id"`
	TaskID     int       `json:"task_id"`
	FromUserID int       `json:"from_user_id"`
	ToUserID   int       `json:"to_user_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
