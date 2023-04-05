package core

import (
	"github.com/firat.atmaca/go-mongo/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine, g *handlers.UserHandler) {
	router := r.Use(gin.Logger(), gin.Recovery())

	router.GET("/", g.Home())

	cookieData := cookie.NewStore([]byte("go-app"))
	router.Use(sessions.Sessions("session", cookieData))

	router.POST("/sign-up", g.SignUp())
	router.POST("/sign-in", g.SignIn())

	authRouter := r.Group("/auth", Authorization())
	{
		authRouter.GET("/dashboard")
	}
}
