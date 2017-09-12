package db

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
)

func (context *APIContext) LoadImage(id string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(context.Config.API.ImageFilePath, id))
}

func (context *APIContext) StoreImage(logo []byte) (string, error) {
	sha := sha1.New()
	sha.Write([]byte(logo))
	shasum := hex.EncodeToString(sha.Sum(nil))

	err := ioutil.WriteFile(filepath.Join(context.Config.API.ImageFilePath, shasum), logo, 0644)
	if err != nil {
		return "", err
	}
	return shasum, nil
}

func (context *APIContext) BuildImageURL(id string) string {
	if id == "" {
		return ""
	}

	u, _ := url.Parse(context.Config.Web.BaseURL)
	u.Path = path.Join(u.Path, "images", id)
	return u.String()
}