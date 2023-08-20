package middlewares

import (
	"fmt"
	"net/http"

	"github.com/georgylobko/gym-log/internal/helpers"
)

type autherHandler func(http.ResponseWriter, *http.Request, string)

func MiddlewareAuth(handler autherHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.GetParsedToken(w, r)
		if err != nil {
			helpers.RespondWithError(w, 401, fmt.Sprintf("Unauthorized %s", err))
			return
		}

		handler(w, r, userId)
	}
}
