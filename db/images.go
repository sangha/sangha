package db

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
)

// LoadImage reads an image from disk
func (context *APIContext) LoadImage(id string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(context.Config.API.ImageFilePath, id))
}

// StoreImage writes an image to disk and returns the hashsum
func (context *APIContext) StoreImage(logo []byte) (string, error) {
	if len(logo) == 0 {
		return "", errors.New("empty image data")
	}

	sha := sha1.New()
	sha.Write([]byte(logo))
	shasum := hex.EncodeToString(sha.Sum(nil))

	err := ioutil.WriteFile(filepath.Join(context.Config.API.ImageFilePath, shasum), logo, 0644)
	if err != nil {
		return "", err
	}
	return shasum, nil
}

// BuildImageURL returns the canonical URL for an image
func (context *APIContext) BuildImageURL(id string, placeholder string) string {
	u, _ := url.Parse(context.Config.Web.ImageURL)

	if id == "" {
		u.Path = path.Join(u.Path, placeholder+".png")
	} else {
		u.Path = path.Join(u.Path, "images", id)
	}

	return u.String()
}
