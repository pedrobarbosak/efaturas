package middlewares

import (
	"log"
	"net/http"
	"os"
	"strings"

	"efaturas-xtreme/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware gin.HandlerFunc

const (
	header   = "Authorization"
	UserID   = "x-user-id"
	Username = "x-username"
	Password = "x-password"
)

func Auth(s *auth.Service) gin.HandlerFunc {
	logger := log.New(os.Stdout, "<auth-middleware> ", log.Flags())

	return func(ctx *gin.Context) {
		token := strings.TrimPrefix(ctx.GetHeader(header), "Bearer ")
		if token == "" {
			logger.Println("empty token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		session, err := s.GetAndExtend(ctx, token)
		if err != nil {
			logger.Println("failed to get and extend:", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set(UserID, session.UserID)
		ctx.Set(Username, session.Username)
		ctx.Set(Password, session.Password)
	}
}
