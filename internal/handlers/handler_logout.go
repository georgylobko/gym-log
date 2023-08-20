package handlers

import (
	"net/http"
	"time"

	"github.com/georgylobko/gym-log/internal/helpers"
)

func (apiCfg *ApiConfig) HandlerLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)

	helpers.RespondWithJSON(w, 200, struct{}{})
}
