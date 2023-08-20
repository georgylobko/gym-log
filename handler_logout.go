package main

import (
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	respondWithJSON(w, 200, struct{}{})
}
