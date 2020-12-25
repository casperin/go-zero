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
	db      *sqlx.DB
	tplArgs map[string]interface{}
}

type authedConfig struct {
	db      *sqlx.DB
	user    models.User
	tplArgs map[string]interface{}
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
		// Unauthed handlers
		for _, h := range unauthedHandlers {
			if h.path != r.URL.Path {
				continue
			}
			e := h.handler(w, r, unauthedConfig{
				db: db,
				tplArgs: map[string]interface{}{
					"Info":  cookie.Consume(w, r, "info"),
					"Alert": cookie.Consume(w, r, "alert"),
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
			id, email, e := cookie.GetUser(r)
			if e != nil {
				cookie.SetInfo(w, "You are not logged in")
				http.Redirect(w, r, "/", 302)
				return
			}
			user := models.User{Id: id, Email: email}
			e = h.handler(w, r, authedConfig{
				db:   db,
				user: user,
				tplArgs: map[string]interface{}{
					"User":  user,
					"Info":  cookie.Consume(w, r, "info"),
					"Alert": cookie.Consume(w, r, "alert"),
				},
			})
			if e != nil {
				log.Println(e)
			}
			return
		}

		w.WriteHeader(404)
		e := renderTemplate(w, r, nil, "404.html")
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
