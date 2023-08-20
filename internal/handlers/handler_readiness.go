package handlers

import (
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, 200, struct{}{})
}
