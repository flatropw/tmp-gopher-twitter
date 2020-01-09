package models

import (
	"errors"
	"github.com/flatropw/gopher-twitter/internal/app/db"
	"time"
)

type Subscriber struct {
	User *User
}

type Subscription struct {
	Id           uint `json:"id"`
	SubscriberId uint `json:"subscriber_id"`
	SubscribedId uint `json:"subscribed_id"`
	Status       bool `json:"status"`
}

func (sub *Subscriber) SubscribeTo(Id uint) (*Subscription, error) {
	subscription, err := sub.GetSubscriptionOn(Id)
	if err != nil {
		return &Subscription{}, nil
	}

	if sub.User.Id == Id {
		return &Subscription{}, errors.New("you cannot subscribe to yourself")
	}

	if subscription.Id > 0 {
		row := db.Instance.Db.QueryRow(db.SubUpdateStatusQuery, !subscription.Status, time.Now().Unix(), subscription.Id)
		err = row.Scan(&subscription.Status)
	} else {
		row := db.Instance.Db.QueryRow(db.SubInsertQuery, sub.User.Id, Id, time.Now().Unix(), time.Now().Unix())
		err = row.Scan(&subscription.Id, &subscription.SubscriberId, &subscription.SubscribedId, &subscription.Status)
	}

	return subscription, err
}

func (sub *Subscriber) GetSubscriptionOn(subscribedId uint) (*Subscription, error) {
	var subscription = &Subscription{}
	row := db.Instance.Db.QueryRow(db.SubHasAlreadySubscribedQuery, sub.User.Id, subscribedId)
	err := row.Scan(&subscription.Id, &subscription.SubscriberId, &subscription.SubscribedId, &subscription.Status)
	return subscription, err
}

func (sub *Subscriber) GetSubscriptions() (subs []Subscription, err error) {
	rows, err := db.Instance.Db.Query(db.SubGetActiveQuery, sub.User.Id)
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		return
	}

	for rows.Next() {
		var s Subscription
		err = rows.Scan(&s.Id, &s.SubscriberId, &s.SubscribedId, &s.Status)
		if err != nil {
			return
		}
		subs = append(subs, s)
	}

	return
}
