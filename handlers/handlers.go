package handlers

import (
	"errors"
	"fmt"
	"github.com/firat.atmaca/go-mongo/encrypt"
	"github.com/firat.atmaca/go-mongo/model"
	"github.com/firat.atmaca/go-mongo/modules/auth"
	"github.com/firat.atmaca/go-mongo/modules/config"
	"github.com/firat.atmaca/go-mongo/modules/database/storage"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"regexp"
	"time"
)

type UserHandler struct {
	App         *config.AppTools
	UserStorage storage.IUserStorage
}

func NewUserHandler(app *config.AppTools, userStorage *storage.UserStorage) *UserHandler {
	return &UserHandler{
		App:         app,
		UserStorage: userStorage,
	}
}

func (g *UserHandler) Home() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"resp": "Welcome"})
	}
}

func (g *UserHandler) SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *model.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.Password, _ = encrypt.Hash(user.Password)

		if err := g.App.Validator.Struct(&user); err != nil {
			if _, ok := err.(*validator.InvalidValidationError); !ok {
				_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
				g.App.InfoLogger.Println(err)
				return
			}
		}

		ok, status, err := g.UserStorage.InsertUser(user)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, errors.New("error while adding new user"))
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if !ok {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		switch status {
		case 1:
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Registered Successfully",
			})
			return

		case 2:
			ctx.JSON(http.StatusFound, gin.H{
				"message": "Existing Account, Go to the Login page",
			})
			return

		}
	}
}

func (g *UserHandler) SignIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user *model.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, gin.Error{Err: err})
		}

		regMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		ok := regMail.MatchString(user.Email)

		if ok {
			res, checkErr := g.UserStorage.VerifyUser(user.Email)

			if checkErr != nil {
				_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("unregistered user %v", checkErr))
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "unregistered user"})
				return
			}

			id := (res["_id"]).(primitive.ObjectID)
			password := (res["password"]).(string)

			verified, err := encrypt.Verify(user.Password, password)
			if err != nil {
				_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("cannot verify user details"))
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect login details"})
				return
			}

			switch {
			case verified:
				cookieData := sessions.Default(ctx)

				userInfo := map[string]interface{}{
					"ID":       id,
					"email":    user.Email,
					"password": user.Password,
				}
				cookieData.Set("data", userInfo)

				if err := cookieData.Save(); err != nil {
					log.Println("error from the session storage")
					_ = ctx.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
					return
				}
				// generate the jwt token
				t1, t2, err := auth.Generate(user.Email, id)
				if err != nil {
					_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("token no generated : %v ", err))
				}

				cookieData.Set("token", t1)

				if err := cookieData.Save(); err != nil {
					log.Println("error from the session storage")
					_ = ctx.AbortWithError(http.StatusNotFound, gin.Error{Err: err})
					return
				}

				tk := map[string]string{"t1": t1, "t2": t2}

				_, updateErr := g.UserStorage.UpdateInfo(id, tk)
				if updateErr != nil {
					_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("unregistered user %v", updateErr))
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect login details"})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"message":       "Successfully Logged in",
					"email":         user.Email,
					"id":            id,
					"session_token": t1,
				})
			case !verified:
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect login details"})
				return
			}
		}
	}
}

func (g *UserHandler) DashBoard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"resp": "Welcome to Go App Dashboard"})
	}
}
