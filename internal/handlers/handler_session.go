package handlers

import (
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
)

func (apiCfg *ApiConfig) HandlerSession(w http.ResponseWriter, r *http.Request, user mappers.User) {
	helpers.RespondWithJSON(w, 200, user)
}
