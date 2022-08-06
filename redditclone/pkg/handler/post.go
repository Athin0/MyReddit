package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"redditclone/pkg/post"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"sort"
	"strconv"
	"strings"
)

type PostHandler struct {
	PostRepo post.PostRepo
	Logger   *zap.SugaredLogger
	Sessions session.SessRepo
}

type PostForm struct {
	Category string `json:"category"`
	Text     string `json:"text"`
	Title    string `json:"title"`
	Type     string `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
}
type CommentForm struct {
	Comment string `json:"comment"`
}

var (
	AplJSON = "application/json"
)

type ErrForm struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Msg      string `json:"msg"`
	Value    string `json:"value,omitempty"`
}

func (h *PostHandler) List(w http.ResponseWriter, _ *http.Request) {
	elems, err := h.PostRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	sort.Sort(PostSort(elems))
	resp, errMarshal := json.Marshal(elems)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func (h *PostHandler) Category(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category, err0 := vars["CATEGORY_NAME"]
	if !err0 {
		http.Error(w, `{"error": "bad category"}`, http.StatusBadGateway)
		return
	}
	elems, err := h.PostRepo.GetInCategory(category)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	sort.Sort(PostSort(elems))
	resp, err2 := json.Marshal(elems)
	if err2 != nil {
		h.Logger.Infow("Error of marshal", err2)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}

}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err0 := vars["POST_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	elem, err := h.PostRepo.Get(id)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	resp, errMrsh := json.Marshal(elem)
	if errMrsh != nil {
		h.Logger.Infow("Error of Marshal", errMrsh)
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}
func (h *PostHandler) Add(w http.ResponseWriter, r *http.Request) {

	body, err1 := ioutil.ReadAll(r.Body)
	if err1 != nil {
		h.Logger.Infow("Error of Read body", err1)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	err := r.Body.Close()
	if err != nil {
		h.Logger.Infow("Error of close req body", err)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	item := &post.Post{}
	err0 := json.Unmarshal(body, item)
	if err0 != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}

	inToken := r.Header.Get("authorization")
	u := GetUserFromToken(inToken)
	item.Author = *u
	ans, err1 := h.PostRepo.Add(item)
	if err1 != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	_, err4 := h.PostRepo.UpdateVote(int(1), ans.ID, u)
	if err4 != nil {
		h.Logger.Infow("Err in UpdateVote: ", err4)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	resp, err3 := json.Marshal(ans)
	if err3 != nil {
		h.Logger.Infow("Error of Marshal: ", err3)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	_, err2 := w.Write(resp)
	if err2 != nil {
		h.Logger.Infow("Error of write", err2)
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *PostHandler) AddComment(w http.ResponseWriter, r *http.Request) {

	body, err3 := ioutil.ReadAll(r.Body)
	if err3 != nil {
		h.Logger.Infow("Error of Reading body", err3)
	}
	err1 := r.Body.Close()
	if err1 != nil {
		h.Logger.Infow("Error of close req body", err1)
		return
	}

	vars := mux.Vars(r)
	id, err0 := vars["POST_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item := &CommentForm{}
	err := json.Unmarshal(body, item)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "cant unpack payload")
		return
	}
	if item.Comment == "" {
		Arr, errM := json.Marshal(map[string][]ErrForm{
			"errors": {
				{
					Location: "body",
					Param:    "comment",
					Msg:      "is required",
				},
			}})
		if errM != nil {
			h.Logger.Infow("Error in Marshaling response", errM)
		}
		http.Error(w, "", http.StatusUnprocessableEntity)
		_, err = w.Write(Arr)
		if err != nil {
			h.Logger.Infow("Error of write", err)
			return
		}
		return
	}
	inToken := r.Header.Get("authorization")
	u := GetUserFromToken(inToken)
	elem, err := h.PostRepo.AddComment(id, item.Comment, u)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	resp, errMarshal := json.Marshal(elem)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func (h *PostHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idPost, err6 := vars["POST_ID"]
	if !err6 {
		h.Logger.Infow("Error in Vars/Delate")
	}

	idComment, err0 := vars["COMMENT_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	id1, _ := strconv.Atoi(idComment)
	elem, err := h.PostRepo.DeleteComment(idPost, int64(id1))
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	resp, errMarshal := json.Marshal(elem)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func (h *PostHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPost, err0 := vars["POST_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	inToken := r.Header.Get("authorization")
	u := GetUserFromToken(inToken)
	elem, err := h.PostRepo.UpdateVote(1, idPost, u)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}

	resp, errMarshal := json.Marshal(elem)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func (h *PostHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPost, err0 := vars["POST_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	inToken := r.Header.Get("authorization")

	u := GetUserFromToken(inToken)
	elem, err := h.PostRepo.UpdateVote(-1, idPost, u)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	resp, errMarshal := json.Marshal(elem)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func (h *PostHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idPost, err0 := vars["POST_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	inToken := r.Header.Get("authorization")

	u := GetUserFromToken(inToken)
	elem, err := h.PostRepo.UpdateVote(0, idPost, u)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, `Template error`, http.StatusInternalServerError)
		return
	}
	resp, errMarshal := json.Marshal(elem)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

func GetUserFromToken(inToken string) *user.User {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return ExampleTokenSecret, nil
	}
	inToken = strings.Split(inToken, " ")[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(inToken, claims, hashSecretGetter)
	if err != nil {
		log.Println(err.Error())
	}

	op := claims["user"].(map[string]interface{})
	u := &user.User{
		Login: op["username"].(string),
		ID:    int64(op["id"].(float64)),
	}
	return u
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idPost, err0 := vars["POST_ID"]
	if !err0 {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	ok, err := h.PostRepo.Delete(idPost)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if ok {
		jsonError(w, http.StatusOK, "success")
	} else {
		jsonError(w, http.StatusInternalServerError, "error of delete")
	}
}

func (h *PostHandler) GetPostsOfUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err0 := vars["USER_LOGIN"]
	if !err0 {
		http.Error(w, `{"error": "bad category"}`, http.StatusBadGateway)
		return
	}

	elems, err := h.PostRepo.GetFromUser(userID)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	resp, errMarshal := json.Marshal(elems)
	if errMarshal != nil {
		h.Logger.Infow("Error in Marshaling response", errMarshal)
	}
	_, err = w.Write(resp)
	if err != nil {
		h.Logger.Infow("Error of write", err)
		return
	}
}

type PostSort []*post.Post

func (a PostSort) Len() int           { return len(a) }
func (a PostSort) Less(i, j int) bool { return a[i].Score > a[j].Score }
func (a PostSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
