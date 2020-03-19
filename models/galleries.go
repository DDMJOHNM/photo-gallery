package models

import "github.com/jinzhu/gorm"

var _ GalleryDB = &galleryGorm{}

type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

type galleryService struct {
	GalleryDB
}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

//
type galleryGorm struct {
	db *gorm.DB
}

type galleryValFn func(*Gallery) error

func (gg *galleryGorm) Create(gallery *Gallery) error {
	//return nil
	return gg.db.Create(gallery).Error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

func runGalleryValFns(gallery *Gallery, fns ...galleryValFn) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}
