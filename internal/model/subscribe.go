package model

type Subscribe struct {
	UserID         int `db:"user_id"`
	SubscribedToId int `db:"subscribed_to_id"`
}
