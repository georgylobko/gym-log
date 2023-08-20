package handlers

import (
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
)

func HandlerErr(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithError(w, 400, "Something went wrong")
}
