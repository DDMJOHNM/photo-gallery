package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


type Services struct {
	Gallery GalleryService
	User    UserService
	Image   ImageService
	db      *gorm.DB
}

type ServicesConfig func(*Services) error

func NewServices(cfgs ...ServicesConfig) (*Services, error) {


	for _, cfgs := range cfgs{
		if err := cfg(&s); err != nil{
			return nil , err
		}
		
	}
	
	return &s, nil

	//632

	// db, err := gorm.Open(dialect, connectionInfo)
	// if err != nil {
	// 	return nil, err
	// }
	// db.LogMode(true)

	// return &Services{
	// 	User: NewUserService(db),
	// 	//Gallery: &galleryGorm{},
	// 	Gallery: NewGalleryService(db),
	// 	Image:   NewImageService(),
	// 	db:      db,
	// }, nil

}


func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
