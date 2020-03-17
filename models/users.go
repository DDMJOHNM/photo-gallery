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

var _ UserDB = &userGorm{}
var _ UserService = &userService{}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
	UserDB
}

type userService struct {
	UserDB
}

type userGorm struct {
	UserDB
	db   *gorm.DB
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	Close() error

	AutoMigrate() error
	DestructiveReset() error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userValFn func(*User) error


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

func NewUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	db.AutoMigrate(&User{})
	return &userGorm{
		db:   db,
	}, nil
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := NewUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	hmac := hash.NewHMAC(hmacSecretKey)

	uv:= &userValidator{
		UserDB: ug,
		hmac:   hmac,
	}

	return &userService{
		UserDB: uv,
	}, nil
}

func (us *userService) Close() error {
	return us.UserDB.Close()
}

func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email=?", email)
	err := first(db, &user)
	return &user, err
}

func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {

	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}


func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

func (ug *userGorm) DestructiveReset() error {
	err := ug.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return ug.AutoMigrate()
}

func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {

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

func (uv *userValidator)ByRemember(token string)(*User, error){
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user,uv.hmacRemember); err != nil{
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}


func (uv *userValidator) bcryptPassword(user *User) error {

	if user.Password == "" {
		return nil
	}
		pwBytes := []byte(user.Password + userPepper)
		hashedBytes, err := bcrypt.GenerateFromPassword(
			pwBytes, bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedBytes)
		user.Password = ""
		return nil


	return nil

}

func (uv *userValidator) hmacRemember (user *User) error{
	if user.Remember == ""{
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}


func (uv *userValidator) Create(user *User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
		}
		user.Remember = token
	}
	err := runUserValFns(user,
		uv.bcryptPassword,
		uv.hmacRemember,
		uv.setRememberIfUnset)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error{
	err := runUserValFns(user,
		uv.bcryptPassword,
		uv.hmacRemember)
	if err != nil{
		return err
	}
	return uv.UserDB.Update(user)
}


func (uv *userValidator) setRememberIfUnset(user *User)error {
	if user.Remember != ""{
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil{
		return err
	}
	user.Remember = token
	return nil
}

func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns{
		if err := fn(user); err != nil{
			return err
		}
	}
	return nil
}

func (uv *userValidator) idGreaterThan(n uint) userValFn{
	return userValFn(func(user *User) error{
		if user.ID <= n{
			return ErrInvalidID
		}
		return nil
	})
}

func (uv *userValidator) Delete(id uint) error{

	var user User
	user.ID = id

	err := runUserValFns(&user,uv.idGreaterThan(0))
	if err != nil{
		return err
	}
	return uv.UserDB.Delete(id)
}

