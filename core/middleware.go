package core

import (
	"errors"
	"github.com/firat.atmaca/go-mongo/modules/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieData := sessions.Default(ctx)
		tokenString := cookieData.Get("token").(string)

		if tokenString == "" {
			_ = ctx.AbortWithError(http.StatusNoContent, errors.New("no value for token"))
			return
		}

		parse, err := auth.Parse(tokenString)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, gin.Error{Err: err})
		}

		ctx.Set("pass", tokenString)
		ctx.Set("id", parse.Id)
		ctx.Set("email", parse.Email)
		ctx.Next()
	}
}
