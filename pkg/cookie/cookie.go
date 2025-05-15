package cookie

import "net/http"

func SetCookie(w http.ResponseWriter,
	name, value, path, domain string,
	secure, httpOnly bool,
	sameSite http.SameSite) {

	http.SetCookie(w, &http.Cookie{
		Name: name,
		Value: value,
		Path: path,
		Domain: domain,
		Secure: secure,
		HttpOnly: httpOnly,
		SameSite: sameSite,
	})
}