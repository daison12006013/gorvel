package registrar

import (
	"github.com/daison12006013/gorvel/app/handlers"
	"github.com/daison12006013/gorvel/app/handlers/authhandler"
	"github.com/daison12006013/gorvel/app/handlers/samplehandler"
	"github.com/daison12006013/gorvel/app/handlers/usershandler"
	r "github.com/daison12006013/gorvel/pkg/facade/routes"
)

var Routes = &[]r.Routing{
	{
		Path:    "/",
		Name:    "welcome",
		Method:  r.Method{"GET"},
		Handler: handlers.Welcome,
	},
	{
		Path: "/users",
		Name: "users",
		Resources: r.Resources{
			"index":   usershandler.Lists,  //  GET    /users
			"create":  usershandler.Create, //  GET    /users/create
			"store":   usershandler.Store,  //  POST   /users
			"show":    usershandler.Show,   //  GET    /users/{id}
			"edit":    usershandler.Show,   //  GET    /users/{id}/edit
			"update":  usershandler.Update, //  PUT    /users/{id}, POST /users/{id}/update
			"destroy": usershandler.Delete, //  DELETE /users/{id}, POST /users/{id}/delete
		},
		Middlewares: r.Middlewares{"auth"},
	},
	{
		Path:    "/samples/requests",
		Name:    "",
		Method:  r.Method{"GET", "POST"},
		Handler: samplehandler.Requests,
	},
	{
		Path:    "/samples/storage",
		Name:    "",
		Method:  r.Method{"POST"},
		Handler: samplehandler.FileStorage,
	},
	{
		Path:   "/static",
		Name:   "static",
		Static: "./resources/static",
	},
	{
		Path:    "/docs",
		Prefix:  true,
		Name:    "docs",
		Method:  r.Method{"GET"},
		Handler: handlers.Docs,
	},
	{
		Path:    "/auth/via-cookie",
		Name:    "auth-via-cookie",
		Method:  r.Method{"POST"},
		Handler: authhandler.ViaCookie,
	},
}
