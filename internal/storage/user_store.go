package storage

// UserStore covers User CRUD and lookups. Membership of users on a Ticker
// is owned by TickerStore.
type UserStore interface {
	FindUsers(filter UserFilter, opts ...QueryOpt) ([]User, error)
	FindUserByID(id int, opts ...QueryOpt) (User, error)
	FindUsersByIDs(ids []int, opts ...QueryOpt) ([]User, error)
	FindUserByEmail(email string, opts ...QueryOpt) (User, error)
	SaveUser(user *User) error
	DeleteUser(user User) error
}
