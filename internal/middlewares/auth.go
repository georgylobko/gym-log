package middlewares

import (
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
	"github.com/georgylobko/gym-log/internal/mappers"
)

type autherHandler func(http.ResponseWriter, *http.Request, mappers.User)

func MiddlewareAuth(handler autherHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := helpers.GetParsedToken(w, r)
		if err != nil {
			helpers.RespondWithError(w, 401, fmt.Sprintf("Unauthorized %s", err))
			return
		}

		handler(w, r, user)
	}
}
