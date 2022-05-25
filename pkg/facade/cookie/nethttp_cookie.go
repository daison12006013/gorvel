package cookie

import (
	"net/http"
	"os"

	"github.com/golang-module/carbon"
	"github.com/lucidfy/lucid/pkg/errors"
	"github.com/lucidfy/lucid/pkg/facade/crypt"
	"github.com/lucidfy/lucid/pkg/helpers"
)

type NetHttpCookie struct {
	ResponseWriter http.ResponseWriter
	HttpRequest    *http.Request
}

func New(w http.ResponseWriter, r *http.Request) *NetHttpCookie {
	s := NetHttpCookie{
		ResponseWriter: w,
		HttpRequest:    r,
	}

	return &s
}

func (s *NetHttpCookie) CreateSessionCookie() interface{} {
	sessionKey := crypt.GenerateRandomString(20)
	s.Set(helpers.Getenv("SESSION_NAME", "lucid_session"), sessionKey)
	return sessionKey
}

func (s *NetHttpCookie) Set(name string, value interface{}) (bool, *errors.AppError) {
	encoded, err := crypt.Encrypt(value)
	if err == nil {
		lifetime, err := helpers.StringToInt(os.Getenv("SESSION_LIFETIME"))
		if err != nil {
			return false, err
		}
		cookie := &http.Cookie{
			Name:    name,
			Value:   encoded,
			Path:    "/",
			MaxAge:  lifetime,
			Expires: carbon.Now().AddSeconds(lifetime).Carbon2Time(),
			Domain:  os.Getenv("SESSION_DOMAIN"),
		}
		http.SetCookie(s.ResponseWriter, cookie)
		return true, nil
	}

	return false, err
}

func (s *NetHttpCookie) Get(name string) (interface{}, *errors.AppError) {
	if s.HttpRequest == nil {
		return nil, nil
	}

	cookie, err := s.HttpRequest.Cookie(name)
	if err != nil {
		return nil, errors.InternalServerError("s.HttpRequest.Cookie() error", err)
	}

	decoded, app_err := crypt.Decrypt(cookie.Value)
	if app_err != nil {
		return nil, app_err
	}

	return decoded, nil
}

func (s *NetHttpCookie) Expire(name string) {
	cookie := &http.Cookie{Name: name, Value: "", Path: "/", MaxAge: -1}
	http.SetCookie(s.ResponseWriter, cookie)
}