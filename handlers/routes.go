package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/casperin/go-zero/handlers/cookie"
	"github.com/casperin/go-zero/models"
)

type unauthedConfig struct {
	db         *sqlx.DB
	isLoggedIn bool
	tplArgs    map[string]interface{}
}

type authedConfig struct {
	db         *sqlx.DB
	isLoggedIn bool
	me         models.User
	tplArgs    map[string]interface{}
}

type unauthedHandler struct {
	path    string
	handler func(http.ResponseWriter, *http.Request, unauthedConfig) error
}

type authedHandler struct {
	path    string
	handler func(http.ResponseWriter, *http.Request, authedConfig) error
}

var unauthedHandlers = []unauthedHandler{
	{"/", index},
	{"/users", users},
	{"/users/new", usersNew},
	{"/users/create", usersCreate},
	{"/users/edit", usersEdit},
	{"/users/update", usersUpdate},
	{"/users/authenticate", usersAuthenticate},
}

var authedHandlers = []authedHandler{
	{"/home", home},
	{"/logout", usersLogout},
}

func Routes(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, email, e := cookie.GetUser(r)
		isLoggedIn := e == nil
		user := models.User{Id: id, Email: email}

		// Unauthed handlers
		for _, h := range unauthedHandlers {
			if h.path != r.URL.Path {
				continue
			}
			e := h.handler(w, r, unauthedConfig{
				db:         db,
				isLoggedIn: isLoggedIn,
				tplArgs: map[string]interface{}{
					"IsLoggedIn": isLoggedIn,
					"Me":         user, // may be empty user
					"Info":       cookie.Consume(w, r, "info"),
					"Alert":      cookie.Consume(w, r, "alert"),
				},
			})
			if e != nil {
				log.Println(e)
			}
			return
		}

		// Authed handlers
		for _, h := range authedHandlers {
			if h.path != r.URL.Path {
				continue
			}
			if !isLoggedIn {
				cookie.SetInfo(w, "You are not logged in")
				http.Redirect(w, r, "/", 302)
				return
			}
			e = h.handler(w, r, authedConfig{
				db:         db,
				isLoggedIn: true,
				me:         user,
				tplArgs: map[string]interface{}{
					"IsLoggedIn": true,
					"Me":         user,
					"Info":       cookie.Consume(w, r, "info"),
					"Alert":      cookie.Consume(w, r, "alert"),
				},
			})
			if e != nil {
				log.Println(e)
			}
			return
		}

		w.WriteHeader(404)
		e = renderTemplate(w, r, nil, "404.html")
		if e != nil {
			log.Println(e)
		}
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tplArgs map[string]interface{}, paths ...string) error {
	tplPaths := make([]string, len(paths)+1)
	tplPaths[0] = "handlers/templates/layout.html"
	for i, p := range paths {
		tplPaths[i+1] = "handlers/templates/" + p
	}
	tpl, e := template.ParseFiles(tplPaths...)
	if e != nil {
		return e
	}
	return tpl.Execute(w, tplArgs)
}

func renderError(w http.ResponseWriter, e error, statusCode ...int) error {
	s := 500
	if len(statusCode) > 0 {
		s = statusCode[0]
	}
	tpl, e := template.ParseFiles("handlers/templates/layout.html", "handlers/templates/error.html")
	if e != nil {
		return e
	}
	return tpl.Execute(w, map[string]interface{}{"StatusCode": s, "Error": e})
}
