package storage

import "gorm.io/gorm"

// UserStore covers User CRUD and lookups. Membership of users on a Ticker
// is owned by TickerStore.
type UserStore interface {
	FindUsers(filter UserFilter, opts ...QueryOpt) ([]User, error)
	FindUserByID(id int, opts ...QueryOpt) (User, error)
	FindUsersByIDs(ids []int, opts ...QueryOpt) ([]User, error)
	FindUserByEmail(email string, opts ...QueryOpt) (User, error)
	SaveUser(user *User) error
	DeleteUser(user User) error

	WithUserTx(tx *gorm.DB) UserStore
}

// WithUserTx returns a UserStore scoped to the given transaction.
func (s *SqlStorage) WithUserTx(tx *gorm.DB) UserStore {
	return &SqlStorage{DB: tx, uploadPath: s.uploadPath}
}
