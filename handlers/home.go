package handlers

import (
	"net/http"

	"github.com/casperin/go-zero/handlers/cookie"
)

func index(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	_, _, e := cookie.GetUser(r)
	if e == nil {
		// User is logged in, so we forward to home
		http.Redirect(w, r, "/home", 302)
		return nil
	}
	return renderTemplate(w, r, conf.tplArgs, "home/index.html")
}

func home(w http.ResponseWriter, r *http.Request, conf authedConfig) error {
	return renderTemplate(w, r, conf.tplArgs, "home/show.html")
}
