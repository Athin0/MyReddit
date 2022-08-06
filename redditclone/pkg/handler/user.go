package handler

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"time"
)

type UserHandler struct {
	UserRepo user.UserRepo
	Logger   *zap.SugaredLogger
	Sessions session.SessRepo
}

var (
	ExampleTokenSecret = []byte("супер секретный ключ")
)

type LoginForm struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

func jsonError(w io.Writer, status int, msg string) {
	resp, errMarshal := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	if errMarshal != nil {
		log.Println("Error in Marshaling response", errMarshal)
		return
	}
	_, err := w.Write(resp)
	if err != nil {
		log.Println("Error of write", err)
		return
	}
}

func (h *UserHandler) Re(w http.ResponseWriter, r *http.Request) {
	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		h.Logger.Infow("err in Re", errRead)
		return
	}
	err1 := r.Body.Close()
	if err1 != nil {
		h.Logger.Infow("Error of close req body", err1)
		return
	}
	fd := &LoginForm{}
	err := json.Unmarshal(body, fd)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	us, exist := h.UserRepo.Authorize(fd.Login, fd.Password)
	if exist != nil {
		h.Logger.Infow(exist.Error())
		http.Error(w, "Authorize error", http.StatusUnauthorized)
		jsonError(w, http.StatusUnauthorized, "bad login or password")
		return
	}
	_, errCreate := h.Sessions.Create(w, us.ID)
	if errCreate != nil {
		http.Error(w, "Authorize error", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": us.Login,
			"id":       us.ID,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 1200,
	})
	tokenString, err := token.SignedString(ExampleTokenSecret)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp, errMrsh := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if errMrsh != nil {
		h.Logger.Infow("Err of Marshal", errMrsh)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func (h *UserHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	body, errRead := ioutil.ReadAll(r.Body)
	if errRead != nil {
		h.Logger.Infow("err in RegisterPage", errRead)
		return
	}
	err1 := r.Body.Close()
	if err1 != nil {
		h.Logger.Infow("Error of close req body", err1)
		return
	}

	fd := &LoginForm{}
	err := json.Unmarshal(body, fd)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	t, errUser := h.UserRepo.Authorize(fd.Login, fd.Password)
	if errUser != user.ErrNoUser {
		if t != nil && errUser != user.ErrBadPass {
			Arr, errM := json.Marshal(map[string][]ErrForm{
				"errors": {
					{
						Location: "body",
						Msg:      "already exists",
						Param:    "username",
						Value:    fd.Login,
					},
				}})
			if errM != nil {
				h.Logger.Infow("Error in Marshaling response", errM)
			}
			http.Error(w, "Empty comment", http.StatusUnprocessableEntity)
			_, err = w.Write(Arr)
			if err != nil {
				h.Logger.Infow("Error of write", err)
				return
			}
			return
		} else {
			http.Error(w, "err in Autorise", http.StatusInternalServerError)
			return
		}
	}
	us, exist := h.UserRepo.AddUserInRepo(fd.Login, fd.Password)
	if exist != nil {
		h.Logger.Infow(exist.Error())
		http.Error(w, "bad login or password", http.StatusUnauthorized)
		return
	}
	_, errCreate := h.Sessions.Create(w, us.ID)
	if errCreate != nil {
		h.Logger.Infow("errin Create:", errCreate)
		http.Error(w, "err in Create", http.StatusInternalServerError)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": map[string]interface{}{
			"username": us.Login,
			"id":       us.ID,
		},
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + 1200,
	})
	tokenString, err := token.SignedString(ExampleTokenSecret)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp, errMrsh := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if errMrsh != nil {
		h.Logger.Infow("Err of Marshal", errMrsh)
		return
	}

	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}
