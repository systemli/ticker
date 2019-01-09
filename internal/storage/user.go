package storage

import (
	"errors"
	"github.com/asdine/storm/q"
	. "github.com/systemli/ticker/internal/model"
)

//FindUserByID returns user if one exists with the given id.
func FindUserByID(id int) (*User, error) {
	var user User

	err := DB.One("ID", id, &user)

	return &user, err
}

//FindUsers returns all users.
func FindUsers() ([]User, error) {
	var users []User

	err := DB.Select().Reverse().Find(&users)
	if err != nil {
		return users, err
	}

	return users, nil
}

//FindUsersByTicker returns all users associated with given ticker.
func FindUsersByTicker(ticker Ticker) ([]User, error) {
	var users []User

	query := DB.Select()
	err := query.Each(new(User), func(record interface{}) error {
		u := record.(*User)

		for _, id := range u.Tickers {
			if id == ticker.ID {
				users = append(users, *u)
			}
		}

		return nil
	})

	if err != nil {
		return users, err
	}

	return users, nil
}

//UserAuthenticate returns User when authentication was successful.
func UserAuthenticate(email, password string) (*User, error) {
	var user User

	err := DB.One("Email", email, &user)
	if err != nil {
		return &user, err
	}

	if user.Authenticate(password) {
		return &user, nil
	}

	return &user, errors.New("authentication failed")
}

//AddUsersToTicker append Ticker to the given slice of users.
func AddUsersToTicker(ticker Ticker, ids []int) error {
	var users []User

	err := DB.Select(q.In("ID", ids)).Find(&users)
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.IsSuperAdmin {
			continue
		}
		user.AddTicker(ticker)
		err = DB.Save(&user)
	}

	return err
}

//RemoveTickerFromUser remove ticker from user.
func RemoveTickerFromUser(ticker Ticker, user User) error {
	user.RemoveTicker(ticker)

	return DB.Save(&user)
}
