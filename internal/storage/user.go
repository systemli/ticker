package storage

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	UserOrderFields = []string{"id", "created_at", "updated_at", "email", "is_super_admin"}
)

type User struct {
	ID                int `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Email             string `gorm:"uniqueIndex;not null"`
	EncryptedPassword string
	IsSuperAdmin      bool
	Tickers           []Ticker `gorm:"many2many:ticker_users;"`
}

func NewUser(email, password string) (User, error) {
	user := User{
		IsSuperAdmin:      false,
		Email:             email,
		EncryptedPassword: "",
	}

	pw, err := hashPassword(password)
	if err != nil {
		return user, err
	}

	user.EncryptedPassword = pw

	return user, nil
}

func (u *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password))
	return err == nil
}

func (u *User) UpdatePassword(password string) {
	pw, err := hashPassword(password)
	if err != nil {
		return
	}

	u.EncryptedPassword = pw
}

func (u *User) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                 u.ID,
		"created_at":         u.CreatedAt,
		"updated_at":         u.UpdatedAt,
		"email":              u.Email,
		"encrypted_password": u.EncryptedPassword,
		"is_super_admin":     u.IsSuperAdmin,
	}
}

func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

type UserFilter struct {
	Email        *string
	IsSuperAdmin *bool
	OrderBy      string
	Sort         string
}

func NewUserFilter(req *http.Request) UserFilter {
	filter := UserFilter{
		OrderBy: "id",
		Sort:    "desc",
	}

	if req == nil {
		return filter
	}

	if req.URL.Query().Get("order_by") != "" {
		opts := []string{"id", "created_at", "updated_at", "email", "is_super_admin"}
		for _, opt := range opts {
			if req.URL.Query().Get("order_by") == opt {
				filter.OrderBy = req.URL.Query().Get("order_by")
				break
			}
		}
	}
	if req.URL.Query().Get("sort") != "" {
		filter.Sort = req.URL.Query().Get("sort")
	}

	email := req.URL.Query().Get("email")
	isSuperAdmin := req.URL.Query().Get("is_super_admin")
	if email != "" {
		filter.Email = &email
	}

	if isSuperAdmin != "" {
		isSuperAdminBool, _ := strconv.ParseBool(isSuperAdmin)
		filter.IsSuperAdmin = &isSuperAdminBool
	}

	return filter
}
