package models

import (
	"errors"

	"../hash"
	"../rand"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var userPepper = "secret-random-string"

const hmacSecretKey = "secret-hmac-key"

var (
	ErrInvalidID = errors.New("models: id provided was invalid")
)

var (
	ErrInvalidPassword = errors.New(
		"models : incorrect password provided")
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
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

//page 310

type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
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
	hmac := hash.NewHMAC(hmacSecretKey)
	db.AutoMigrate(&User{})
	return &UserService{
		db:   db,
		hmac: hmac,
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
	err := first(db, &user)
	return &user, err
}

func (us *UserService) ByRemember(token string)(*User, error){

	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?",rememberHash),&user)
	if err != nil{
		return nil, err
	}
	return &user, nil

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

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	user.RememberHash = us.hmac.Hash(user.Remember)

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

	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}
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

func (us *UserService) Authenticate(email, password string) (*User, error) {

	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPepper))

	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
	return foundUser, nil
}
