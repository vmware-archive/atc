package auth

import (
	"context"
	"encoding/base64"
	"net/http"
)

const AuthCookieName = "ATC-Authorization"
const CSRFRequiredKey = "CSRFRequired"

type CookieSetHandler struct {
	Handler http.Handler
}

func (handler CookieSetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AuthCookieName)
	if err == nil {
		ctx := context.WithValue(r.Context(), CSRFRequiredKey, true)
		r = r.WithContext(ctx)
		if r.Header.Get("Authorization") == "" {
			data, err := base64.StdEncoding.DecodeString(cookie.Value);
			if err == nil {
				r.Header.Set("Authorization", string(data[:]))
			}
		}
	}
	handler.Handler.ServeHTTP(w, r)
}
