package model

import "time"

type BidStatus string

const (
	BidStatusPending  BidStatus = "pending"
	BidStatusAccepted BidStatus = "accepted"
	BidStatusRejected BidStatus = "rejected"
)

type Bid struct {
	ID         int       `json:"id"`
	TaskID     int       `json:"task_id"`
	ExecutorID int       `json:"executor_id"`
	Amount     float64   `json:"amount"`
	Status     BidStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateBidRequest struct {
	TaskID int     `json:"task_id"`
	Amount float64 `json:"amount"`
}
