package users_handler

import (
	"net/http"

	"github.com/lucidfy/lucid/app/models/users"
	"github.com/lucidfy/lucid/pkg/engines"
	"github.com/lucidfy/lucid/pkg/errors"
	"github.com/lucidfy/lucid/pkg/facade/session"
)

func delete(T engines.EngineContract) *errors.AppError {
	engine := T.(engines.NetHttpEngine)
	w := engine.ResponseWriter
	r := engine.HttpRequest
	ses := session.File(w, r)
	req := engine.Request
	res := engine.Response
	url := engine.URL

	//> prepare message and status
	message := "Successfully Deleted!"
	status := http.StatusOK

	//> validate "id" if exists
	id := req.Input("id", nil).(string)
	if app_err := users.Exists("id", &id); app_err != nil {
		return app_err
	}

	//> now get the data
	data, app_err := users.Find(&id, nil)
	if app_err != nil {
		return app_err
	}

	//> and delete the data
	data.Delete()

	//> response: for api based
	if req.WantsJson() {
		return res.Json(map[string]interface{}{
			"success": message,
		}, status)
	}

	//> response: for form based, just redirect
	ses.PutFlash("success", message)
	return url.RedirectPrevious()
}