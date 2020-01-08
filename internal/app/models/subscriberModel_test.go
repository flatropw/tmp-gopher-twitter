package models

import (
	"fmt"
	"github.com/flatropw/gopher-twitter/internal/app/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

var _ = db.Instance.Connect("postgres://postgres:postgres@localhost:5432/twitter?sslmode=disable")

func TestUser_GetById(t *testing.T) {
	u := User{}
	us, err := u.GetById(15)
	assert.Nil(t, err)
	expected := uint(15)
	assert.Equal(t, expected, us.Id)
}

func TestSubscriber_GetSubscriptionOn(t *testing.T) {
	u := User{}
	us, err := u.GetById(14)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(us)

	sub := Subscriber{us}
	res, err := sub.SubscribeTo(15)
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
	}
	res, err = sub.GetSubscriptionOn(15)
	assert.Nil(t, err)
	assert.Equal(t, &Subscription{}, res)
}

func TestSubscriber_GetSubscriptions(t *testing.T) {
	var _ = db.Instance.Connect("postgres://postgres:postgres@localhost:5432/twitter?sslmode=disable")
	u := User{}
	us, err := u.GetById(14)
	if err != nil {
		fmt.Println(err)
	}

	sub := Subscriber{us}
	res, err := sub.GetSubscriptions()

	assert.Equal(t, []Subscription {}, res)
}

func TestTweet_GetByUserIds(t *testing.T) {
	var _ = db.Instance.Connect("postgres://postgres:postgres@localhost:5432/twitter?sslmode=disable")
	tw := Tweet{}
	ids := []uint {4, 15}
	res, err := tw.GetByUserIds(ids, 10)
	assert.Nil(t, err)
	assert.Equal(t, 1, res)
}
