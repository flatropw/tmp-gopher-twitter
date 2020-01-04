package middlewares

import (
	"context"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/flatropw/gopher-twitter/internal/app/models"
	u "github.com/flatropw/gopher-twitter/internal/app/utils"
	"net/http"
	"os"
	"strings"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u.IsAuthorizedRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		response := make(map[string] interface{})
		tokenHeader := r.Header.Get("Authorization") //Получение токена

		if len(tokenHeader) == 0 {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusUnauthorized)
			u.Response(w, response)
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			u.Response(w, response)
			return
		}

		tokenPart := splitted[1] //Получаем вторую часть токена
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_secret")), nil
		})

		if err != nil { //Неправильный токен, как правило, возвращает 403 http-код
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			u.Response(w, response)
			return
		}

		if !token.Valid { //токен недействителен, возможно, не подписан на этом сервере
			response = u.Message(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			u.Response(w, response)
			return
		}

		fmt.Println("auth")
		ctx := context.WithValue(r.Context(), "user", tk.Id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
