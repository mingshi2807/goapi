package middleware

import (
	"fmt"
	"net/http"

	"github.com/mingshi2807/goapi/api"
	"github.com/mingshi2807/goapi/internal/tools"
	log "github.com/sirupsen/logrus"
)

var UnAuthorizedError = fmt.Errorf("Invalid username or token.")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usrname := r.URL.Query().Get("usrname")
		token := r.Header.Get("Authorization")
		var err error

		if usrname == "" || token == "" {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		var database *tools.DatabaseInterface
		database, err = tools.NewDatabase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		loginDetails := (*database).GetUserLoginDetails(usrname)

		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(UnAuthorizedError)
			api.RequestErrorHandler(w, UnAuthorizedError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
