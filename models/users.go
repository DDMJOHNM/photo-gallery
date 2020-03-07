package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var userPepper = "secret-random-string"

var (
	ErrInvalidID = errors.New("models: id provided was invalid")
)

var (
	// ErrNotFound is returned when a resource cannot be found // in the database.
	ErrNotFound = errors.New("models: resource not found")
)

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

type UserService struct {
	db *gorm.DB
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrNotFound
		}
		return err
	}
	return err
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	db.AutoMigrate(&User{})
	return &UserService{
		db: db,
	}, nil
}

func (us *UserService) Close() error {
	return us.db.Close()
}

func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email=?", email)
	err := first(db, user)
	return &user, err
}

func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
