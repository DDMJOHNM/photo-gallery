package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"net/url"
)

type imageService struct {
}

type Image struct {
	GalleryID uint
	Filename  string
}

type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func (i *Image) Path() string { 
	temp := url.URL{
	Path: "/" + i.RelativePath(), 
	}
	return temp.String() 
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

// func (i *Image) Path() string {
// 	return "/" + i.RelativePath()
// }

func (i *Image) RelativePath() string {
	galleryID := fmt.Sprintf("%v", i.GalleryID)
	return filepath.ToSlash(filepath.Join("images", "galleries", galleryID, i.Filename))
}

func NewImageService() ImageService {
	return &imageService{}
}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil

}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries",
		fmt.Sprintf("%v", galleryID))
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}

	fmt.Println(filepath.FromSlash(galleryPath))
	fmt.Println(filepath.ToSlash(galleryPath))

	return galleryPath, nil
	//john's fix for windows file system
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			Filename:  filepath.Base(imgStr),
			GalleryID: galleryID,
		}
	}

	// for i := range strings {
	// 	strings[i] = "/" + filepath.ToSlash(strings[i])
	// }

	return ret, nil
}
