package db

import "gorm.io/gorm"

type User struct {
	gorm.Model            // swagger:ignore
	ID             int64  `gorm:"primaryKey"  json:"id"`
	Username       string `gorm:"unique" json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `gorm:"unique" json:"email"`
	SecurityKey    string `json:"security_key"`
	HashedPassword string `json:"hashed_password"`
}

// Create User Account
func (store *Store) CreateUser(ua *User) (err error) {
	if err := store.Create(ua).Error; err != nil {
		return err
	}

	return
}

func (store *Store) UpdateUser(ua *User) (err error) {
	if err := store.Save(ua).Error; err != nil {
		return err
	}
	return
}

// Get by ID
func (store *Store) GetUserByID(uid int64) (*User, error) {
	user := User{ID: uid}
	if err := store.Where("id = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Get by ID
func (store *Store) GetUserByEmail(email string) (*User, error) {
	user := User{}
	if err := store.Where("email = ? OR username = ?", email, email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
