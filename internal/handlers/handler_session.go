package handlers

import (
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
)

func (apiCfg *ApiConfig) HandlerSession(w http.ResponseWriter, r *http.Request, userID string) {
	helpers.RespondWithJSON(w, 200, userID)
}
