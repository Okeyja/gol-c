package v1

import (
	"github.com/gin-gonic/gin"
	"gol-c/api/v1/util"
	"gol-c/database"
	"gol-c/model"
	s_util "gol-c/util"
	"net/http"
	"time"
)

type LoginJsonForm struct {
	Username string
	Password string
}

// CreateSession
// @Summary Login
// @Description Login as a user
// @Tags auth
// @Success 200 {string} string    "ok"
// @Param param body LoginJsonForm true "Login json form"
// @Router /v1/auth/sessions [POST]
func CreateSession(c *gin.Context) {
	loginJsonForm := &LoginJsonForm{}
	util.GetJsonForm(c, loginJsonForm)
	user := &model.User{}
	database.GetByField(&model.User{Username: loginJsonForm.Username}, user, []string{})
	pass := util.CheckEncrypt(user.Password, loginJsonForm.Password)
	if pass {
		if user.EmailVerified {
			session := model.CreateSession(user)
			err := database.Create(c, session, "session", util.ErrorMessageStatus)
			if err != nil {
				return
			}
			util.SuccessDataMessage(c, gin.H{"renewal_code": session.SessionKey, "token_body": session.SessionToken}, "Create session succeeded!")
		} else {
			util.ErrorMessageStatus(c, "Create session failed: email not verified!", http.StatusUnauthorized)
		}
	} else {
		util.ErrorMessageStatus(c, "Create session failed: wrong username or password!", http.StatusUnauthorized)
	}
}

type RenewSessionForm struct {
	RenewalKey string
}

// RenewSession
// @Summary Keep login status
// @Description Keep login status as a user.
// @Tags auth
// @Success 200 {string} string    "ok"
// @Param param body RenewSessionForm true "session renewal form"
// @Router /v1/auth/sessions/any:renew [PUT]
func RenewSession(c *gin.Context) {
	renewSessionForm := &RenewSessionForm{}
	util.GetJsonForm(c, renewSessionForm)
	session := &model.Session{}
	database.GetByField(&model.Session{SessionKey: renewSessionForm.RenewalKey}, session, []string{})
	if time.Now().Before(session.ExpiredAt) {
		if session.RenewalStock > 0 || session.RenewalStock == -1 {
			body, _ := s_util.GenToken(session.User.Username, session.SessionKey)
			session.SessionToken = body
			session.RenewalStock -= 1
			err := database.Update(c, session, "session", util.ErrorMessageStatus)
			if err != nil {
				return
			}
			util.SuccessDataMessage(c, gin.H{"renewal_code": session.SessionKey, "token_body": session.SessionToken}, "Renew session succeeded!")
		} else {
			util.ErrorMessageStatus(c, "Session renewing failed, reached maximum renewal times.", http.StatusForbidden)
		}
	} else {
		util.ErrorMessageStatus(c, "Session expired, renewing failed.", http.StatusForbidden)
	}
}
