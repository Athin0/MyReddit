package handler

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"net/http/httptest"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
	"strings"
	"testing"
)

func TestUserLoginPage(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	sess := session.NewMockSessRepo(ctrl)
	service := &UserHandler{
		UserRepo: st,
		Logger:   zap.NewNop().Sugar(),
		Sessions: sess,
	}

	arrUser := []*user.User{
		&user.User{0,
			"qwerty",
			"asdfghjk"},
	}
	w := httptest.NewRecorder()
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(arrUser[0], nil)
	sess.EXPECT().Create(w, arrUser[0].ID).Return(&session.Session{}, nil)

	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))

	service.Re(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("some trabls with loginpage")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))

	w1 := httptest.NewRecorder()

	service.Re(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 401 {
		t.Errorf("expected resp status 401, got %d", resp1.StatusCode)
		return
	}

	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(arrUser[0], nil)
	sess.EXPECT().Create(w, arrUser[0].ID).Return(nil, fmt.Errorf("bad sess create"))

	req2 := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))

	service.Re(w, req2)

	resp2 := w1.Result()
	if resp2.StatusCode != 401 {
		t.Errorf("expected resp status 401, got %d", resp2.StatusCode)
		return
	}
}

func TestUserRegisterPage(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := user.NewMockUserRepo(ctrl)
	sess := session.NewMockSessRepo(ctrl)
	service := &UserHandler{
		UserRepo: st,
		Logger:   zap.NewNop().Sugar(),
		Sessions: sess,
	}

	arrUser := []*user.User{
		&user.User{0,
			"qwerty",
			"asdfghjk"},
	}

	w := httptest.NewRecorder()
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(nil, user.ErrNoUser)
	st.EXPECT().AddUserInRepo(arrUser[0].Login, arrUser[0].Password).
		Return(arrUser[0], nil)
	sess.EXPECT().Create(w, arrUser[0].ID).
		Return(&session.Session{}, nil)
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))

	service.RegisterPage(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("some trabls with loginpage")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))

	w1 := httptest.NewRecorder()

	service.RegisterPage(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}
	w2 := httptest.NewRecorder()
	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(nil, user.ErrNoUser)

	st.EXPECT().AddUserInRepo(arrUser[0].Login, arrUser[0].Password).
		Return(arrUser[0], nil)
	sess.EXPECT().Create(w2, int64(arrUser[0].ID)).Return(nil, fmt.Errorf("bad sess create"))

	req2 := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))

	service.RegisterPage(w2, req2)

	resp2 := w2.Result()
	if resp2.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp2.StatusCode)
		return
	}

	st.EXPECT().Authorize(arrUser[0].Login, arrUser[0].Password).
		Return(arrUser[0], nil)
	w3 := httptest.NewRecorder()

	req3 := httptest.NewRequest("POST", "/api/register", strings.NewReader(`{"username": "qwerty", "password": "asdfghjk"}`))
	service.RegisterPage(w3, req3)
	resp3 := w1.Result()
	if resp3.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp3.StatusCode)
		return
	}

}
