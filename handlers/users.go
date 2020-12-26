package handlers

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/casperin/go-zero/handlers/cookie"
	"github.com/casperin/go-zero/models"
)

func users(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	id := r.URL.Query().Get("id")

	// No id, means we index all users
	if id == "" {
		users := []models.User{}
		e := conf.db.Select(&users, "select * from users")
		if e != nil {
			return renderError(w, e)
		}
		conf.tplArgs["Users"] = users
		return renderTemplate(w, r, conf.tplArgs, "users/index.html")
	}

	// Single user
	user := models.User{}
	e := conf.db.Get(&user, "select * from users where id = $1", id)
	if e != nil {
		return renderError(w, e)
	}
	conf.tplArgs["User"] = user
	return renderTemplate(w, r, conf.tplArgs, "users/show.html")
}

func usersNew(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	return renderTemplate(w, r, conf.tplArgs, "users/new.html")
}

func usersCreate(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	if email == "" || password == "" {
		conf.tplArgs["Alert"] = "Both email and password are required"
		conf.tplArgs["Email"] = email
		return renderTemplate(w, r, conf.tplArgs, "users/new.html")
	}
	hash, e := bcrypt.GenerateFromPassword([]byte(password), 10)
	if e != nil {
		return renderError(w, e)
	}
	id := 0
	rows, e := conf.db.Query("insert into users (email, password) values ($1, $2) returning id", email, hash)
	if e != nil {
		return renderError(w, e)
	}
	defer rows.Close()
	rows.Next()
	e = rows.Scan(&id)
	if e != nil {
		return renderError(w, e)
	}
	path := fmt.Sprintf("/users?id=%d", id)
	http.Redirect(w, r, path, 302)
	return nil
}

func usersEdit(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	id := r.URL.Query().Get("id")

	user := models.User{}
	e := conf.db.Get(&user, "select * from users where id = $1", id)
	if e != nil {
		return renderError(w, e)
	}
	conf.tplArgs["User"] = user
	return renderTemplate(w, r, conf.tplArgs, "users/edit.html")
}

func usersUpdate(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	id := r.FormValue("id")
	email := r.FormValue("email")
	old_email := r.FormValue("old_email")
	password := r.FormValue("password")
	if email != old_email {
		_, e := conf.db.Query("update users set email=$1 where id=$2", email, id)
		if e != nil {
			return renderError(w, e)
		}
	}
	if password != "" {
		hash, e := bcrypt.GenerateFromPassword([]byte(password), 10)
		if e != nil {
			return renderError(w, e)
		}
		_, e = conf.db.Query("update users set password=$1 where id=$2", hash, id)
		if e != nil {
			return renderError(w, e)
		}
	}
	http.Redirect(w, r, "/users?id="+id, 302)
	return nil
}

func usersAuthenticate(w http.ResponseWriter, r *http.Request, conf unauthedConfig) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	persistent := r.FormValue("persistent")
	user := models.User{}
	e := conf.db.Get(&user, "select * from users where email=$1", email)
	if e == nil {
		e = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	}
	if e != nil {
		conf.tplArgs["Alert"] = "Email or password is not correct"
		conf.tplArgs["Email"] = email
		return renderTemplate(w, r, conf.tplArgs, "home/index.html")
	}
	var expires time.Time
	if persistent == "true" {
		expires = time.Now().Add(time.Hour * 24 * 90)
	}
	e = cookie.Login(w, user.Id, user.Email, expires)
	if e != nil {
		return renderError(w, e)
	}
	cookie.SetInfo(w, "You are logged in")
	http.Redirect(w, r, "/home", 302)
	return nil
}

func usersLogout(w http.ResponseWriter, r *http.Request, _ authedConfig) error {
	cookie.Logout(w)
	cookie.SetInfo(w, "You are logged out")
	http.Redirect(w, r, "/", 302)
	return nil
}
