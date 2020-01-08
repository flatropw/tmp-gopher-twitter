package models

import (
	"fmt"
	"github.com/Shyp/go-dberror"
	"github.com/flatropw/gopher-twitter/internal/app/db"
	"time"
)

type Subscriber struct {
	User *User
}

type Subscription struct {
	Id uint
	SubscriberId uint
	SubscribedId uint
	Status bool
}

func (sub *Subscriber) SubscribeTo(Id uint) (subscription *Subscription, err error) {
	subscription, err = sub.GetSubscriptionOn(Id)
	if err != nil {
		return
	}

	if subscription.Id > 0 {
		row := db.Instance.Db.QueryRow(db.SubUpdateStatusQuery, !subscription.Status, time.Now().Unix(), subscription.Id)
		err = row.Scan(&subscription.Status)

	} else {
		row := db.Instance.Db.QueryRow(db.SubInsertQuery, sub.User.Id, Id, time.Now().Unix(), time.Now().Unix())
		err = row.Scan(&subscription.Id, &subscription.SubscriberId, &subscription.SubscribedId, &subscription.Status)
	}
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &Subscription{}, fmt.Errorf(e.Error())
	default:
		return
	}
}

func (sub *Subscriber) GetSubscriptionOn(subscribedId uint) (*Subscription, error) {
	var subscription = Subscription{}
	row := db.Instance.Db.QueryRow(db.SubHasAlreadySubscribedQuery, sub.User.Id, subscribedId)
	err := row.Scan(&subscription.Id, &subscription.SubscriberId, &subscription.SubscribedId, &subscription.Status)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &Subscription{}, fmt.Errorf(e.Error())
	default:
		return &subscription, nil
	}
}

func (sub *Subscriber) GetSubscriptions() (subs []Subscription, err error) {
	rows, err := db.Instance.Db.Query(db.SubGetActiveQuery, sub.User.Id)
	defer func() {
		_ = rows.Close()
	}()
	if err != nil {
		return
	}
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return subs, fmt.Errorf(e.Error())
	default:
		for rows.Next() {
			var s Subscription
			err = rows.Scan(&s.Id, &s.SubscriberId, &s.SubscribedId, &s.Status)
			if err != nil {
				return
			}
			subs = append(subs, s)
		}
	}

	return
}
