package models

import (
	"fmt"
	"github.com/Shyp/go-dberror"
	"github.com/flatropw/gopher-twitter/internal/app/db"
	"time"
)

type Subscriber struct {
	user *User
}

type Subscription struct {
	Id uint
	Status bool
}

func (sub *Subscriber) SubscribeTo(Id uint) (subscription *Subscription, err error) {
	subscription, err = sub.GetSubscriptionOf(Id)
	if err != nil {
		return
	}

	if subscription.Id > 0 {
		_, err = db.Instance.Db.Exec(db.SubUpdateStatusQuery, !subscription.Status, time.Now().Unix(), subscription.Id)
	} else {
		row := db.Instance.Db.QueryRow(db.SubInsertQuery, sub.user.Id, Id, time.Now().Unix(), time.Now().Unix())
		err = row.Scan(&subscription.Id, &subscription.Status)
	}
	fmt.Println(err)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &Subscription{}, fmt.Errorf(e.Error())
	default:
		return
	}
}

func (sub *Subscriber) GetSubscriptionOf(subscribedId uint) (*Subscription, error) {
	var subscription = Subscription{}
	row := db.Instance.Db.QueryRow(db.SubHasAlreadySubscribedQuery, sub.user.Id, subscribedId)
	err := row.Scan(&subscription.Id, &subscription.Status)
	dbErr := dberror.GetError(err)
	switch e := dbErr.(type) {
	case *dberror.Error:
		return &Subscription{}, fmt.Errorf(e.Error())
	default:
		return &subscription, nil
	}
}
