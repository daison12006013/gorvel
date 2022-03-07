package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/daison12006013/gorvel/pkg/errors"
	"github.com/daison12006013/gorvel/pkg/facade/session"
	"github.com/daison12006013/gorvel/pkg/facade/urls"
	"github.com/daison12006013/gorvel/pkg/rules"
	"github.com/daison12006013/gorvel/pkg/rules/must"
	"github.com/gorilla/mux"
)

type MuxRequest struct {
	ResponseWriter http.ResponseWriter
	HttpRequest    *http.Request
}

func Mux(w http.ResponseWriter, r *http.Request) MuxRequest {
	t := MuxRequest{
		ResponseWriter: w,
		HttpRequest:    r,
	}
	return t
}

// Url  ----------------------------------------------------

func (t MuxRequest) CurrentUrl() string {
	h := t.HttpRequest.URL.Host
	if len(h) > 0 {
		return h
	}
	// if ever we can't resolve the Host from net/http, we can still base
	// from our env config
	return urls.BaseUrl(nil)
}

func (t MuxRequest) FullUrl() string {
	return t.CurrentUrl() + t.HttpRequest.URL.RequestURI()
}

func (t MuxRequest) PreviousUrl() string {
	return t.HttpRequest.Referer()
}

func (t MuxRequest) RedirectPrevious() {
	http.Redirect(t.ResponseWriter, t.HttpRequest, t.PreviousUrl(), http.StatusFound)
}

func (t MuxRequest) SetFlash(name string, value string) {
	name = "flash-" + name
	s := session.Mux(t.ResponseWriter, t.HttpRequest)
	s.Set(name, value)
}

func (t MuxRequest) GetFlash(name string) *string {
	name = "flash-" + name

	s := session.Mux(t.ResponseWriter, t.HttpRequest)
	if s == nil {
		return nil
	}

	value, err := s.Get(name)
	if (err != nil && err == http.ErrNoCookie) || value == nil {
		return nil
	}

	// delete the cookie by expiring it immediately!
	deleteCookie := &http.Cookie{Name: name, MaxAge: -1, Expires: time.Unix(1, 0), Path: "/"}
	http.SetCookie(t.ResponseWriter, deleteCookie)
	return value
}

// SetFlashMap sets a session flash based on json format
// make sure the values you're providing is set as map[string]interface{}
// therefore, we can stringify it into json format
func (t MuxRequest) SetFlashMap(name string, values interface{}) {
	j, err := json.Marshal(values.(map[string]interface{}))
	if err != nil {
		panic(err)
	}
	t.SetFlash(name, string(j))
}

// GetFlashMap this pulls a session flash from SetFlashMap, in which
// it will reverse the json into a map
func (t MuxRequest) GetFlashMap(name string) *map[string]interface{} {
	ret := map[string]interface{}{}
	flash := t.GetFlash(name)
	if flash != nil {
		json.Unmarshal([]byte(*t.GetFlash(name)), &ret)
	}
	return &ret
}

// Request  -------------------------------------------------

// All returns available http queries
func (t MuxRequest) All() interface{} {
	params := map[string]interface{}{}

	// via form inputs
	for idx, val := range t.HttpRequest.Form {
		if len(val) > 0 {
			params[idx] = val[0]
		}
	}

	// via query params
	for idx, val := range t.HttpRequest.URL.Query() {
		if len(val) > 0 {
			params[idx] = val[0]
		}
	}

	// via route params
	for idx, val := range mux.Vars(t.HttpRequest) {
		params[idx] = val
	}

	return params
}

// Get returns the specific value from the provided key
func (t MuxRequest) Get(k string) interface{} {
	// check the queries if exists
	val, ok := t.All().(map[string]interface{})[k]
	if ok {
		return val
	}
	return nil
}

// GetFirst returns the specifc value provided with default value
func (t MuxRequest) GetFirst(k string, dfault interface{}) interface{} {
	val := t.Get(k)
	if val == nil {
		return dfault
	}
	return val
}

// Input ist meant as proxy for GetFirst(...)
func (t MuxRequest) Input(k string, dfault interface{}) interface{} {
	return t.GetFirst(k, dfault)
}

// HasContentType checks if the string exists in the header
func (t MuxRequest) HasContentType(substr string) bool {
	contentType := t.HttpRequest.Header.Get("Content-Type")
	return strings.Contains(contentType, substr)
}

func (t MuxRequest) HasAccept(substr string) bool {
	accept := t.HttpRequest.Header.Get("Accept")
	return strings.Contains(accept, substr)
}

// IsForm checks if the request is an http form
func (t MuxRequest) IsForm() bool {
	return t.HasContentType("application/x-www-form-urlencoded")
}

// IsJson checks if the content type contains json
func (t MuxRequest) IsJson() bool {
	return t.HasContentType("json")
}

// IsMultipart checks if the content type contains multipart
func (t MuxRequest) IsMultipart() bool {
	return t.HasContentType("multipart")
}

func (t MuxRequest) WantsJson() bool {
	return t.HasAccept("json")
}

// --- Validator

func (t MuxRequest) Validator(setOfRules *must.SetOfRules) *errors.AppError {
	var errsChan = make(chan map[string]string)

	var wg sync.WaitGroup

	for inputField, inputRules := range *setOfRules {
		for _, inputRule := range inputRules {
			inputValue := t.Get(inputField).(string)
			wg.Add(1)
			go rules.Validate(
				inputField,
				inputValue,
				inputRule,
				errsChan,
				&wg,
			)
		}
	}

	go func() {
		wg.Wait()
		close(errsChan)
	}()

	validationErrors := map[string]interface{}{}
	for val := range errsChan {
		for k, v := range val {
			validationErrors[k] = v
		}
	}

	if len(validationErrors) > 0 {
		return &errors.AppError{
			Error:           fmt.Errorf("pkg.facade.request.muxrequest@Validator: Request validation error"),
			Message:         "Request validation error",
			Code:            http.StatusUnprocessableEntity,
			ValidationError: validationErrors,
		}
	}

	return nil
}

// GetIp returns the client IP address
// it resolves first from "x-forwarded-for"
// or else it goes check if "x-real-ip" exists
// or else we pull based on the remoteaddr under net/http
func (t MuxRequest) GetIp() string {
	ip := t.HttpRequest.Header.Get("X-Forwarded-For")
	if len(ip) == 0 {
		ip = t.HttpRequest.Header.Get("X-Real-Ip")
	}
	if len(ip) == 0 {
		ip = t.HttpRequest.RemoteAddr
	}
	return ip
}