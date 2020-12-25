package cookie

import "net/http"

func SetInfo(w http.ResponseWriter, msg string) {
	SetCookie(w, "info", msg)
}

func SetAlert(w http.ResponseWriter, msg string) {
	SetCookie(w, "alert", msg)
}

func Consume(w http.ResponseWriter, r *http.Request, name string) string {
	cookie, e := r.Cookie(name)
	if e != nil {
		return ""
	}
	SetCookie(w, name, "")
	return cookie.Value
}

func SetCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Path:  "/",
		Name:  name,
		Value: value,
	})
}
