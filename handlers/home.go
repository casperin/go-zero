package handlers

import "net/http"

func index(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	if conf.isLoggedIn {
		http.Redirect(w, r, "/home", 302)
		return nil
	}
	return renderTemplate(w, r, conf.tplArgs, "home/index.html")
}

func home(w http.ResponseWriter, r *http.Request, conf authedConfig) error {
	return renderTemplate(w, r, conf.tplArgs, "home/show.html")
}
