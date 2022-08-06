package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http/httptest"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/user"
	"redditclone/pkg/vote"
	"strings"
	"testing"
)

func TestPostHandlerList(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи

	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "programming",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetAll().
		Return(resultPost, nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	service.List(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost)
	if !bytes.Contains(body, ans) {
		t.Errorf("no text found")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetAll().
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()

	service.List(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}
}

func TestPostHandlerCategory(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи

	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "programming",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetInCategory("programming").
		Return(resultPost, nil)

	req := httptest.NewRequest("GET", "/api/posts/programming", strings.NewReader(`{"CATEGORY_NAME":"programming"}`))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"CATEGORY_NAME": "programming",
	})
	service.Category(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost)

	if !bytes.Contains(body, ans) {
		t.Errorf("no text found:")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetInCategory("programming").
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("GET", "/api/posts/programming", nil)
	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"CATEGORY_NAME": "programming",
	})
	service.Category(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bea")
}

func TestPostHandlerGetFromUser(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи

	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "programming",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetFromUser("arin0").
		Return(resultPost, nil)

	req := httptest.NewRequest("GET", "/api/posts/arin0", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"USER_LOGIN": "arin0",
	})
	service.GetPostsOfUser(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost)

	if !bytes.Contains(body, ans) {
		t.Errorf("no text found:")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().GetFromUser("arin0").
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("GET", "/api/posts/arin0", nil)
	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"USER_LOGIN": "arin0",
	})
	service.GetPostsOfUser(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bea")
}

func TestPostHandlerGet(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи

	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "programming",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Get("1").
		Return(resultPost[0], nil)

	req := httptest.NewRequest("GET", "/api/post/1", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})
	service.Get(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[0])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}
	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Get("1").
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("GET", "/api/posts/1", nil)
	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})

	service.Get(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bea")

}

func TestPostHandlerAdd(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "music",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
		&post.Post{
			Author:   user.User{ID: 3, Login: "arin0"},
			Category: "music",
			Text:     "Post Text exemple",
			Title:    "Post Title ex",
			Type:     "text",
		},
	}

	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().Add(resultPost[1]).
		Return(resultPost[0], nil)
	st.EXPECT().UpdateVote(1, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(resultPost[0], nil)
	req := httptest.NewRequest("POST", "/api/posts", strings.NewReader(`{"category": "music", "type": "text", "title": "Post Title ex", "text": "Post Text exemple"}`))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.Add(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[0])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}
	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().Add(resultPost[1]).
		Return(nil, fmt.Errorf("no results"))
	//st.EXPECT().UpdateVote(1, resultPost[0].ID, &user.User{ID: 3, Login: "arin0"}).
	//	Return(resultPost[0], nil)
	req1 := httptest.NewRequest("POST", "/api/posts", strings.NewReader(`{"category": "music", "type": "text", "title": "Post Title ex", "text": "Post Text exemple"}`))

	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})

	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")

	service.Add(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibXl1c2VybmFtZTIiLCJpZCI6IjYyNjZlNTVmMTg1NzNiMDAwOWNjYWIwNSJ9LCJpYXQiOjE2NTIyMDE5NzYsImV4cCI6MTY1MjgwNjc3Nn0.MZTTa4TCyDF98L4X4PC7Kz9pxZ6DZzFlP0zdOuv6KjU")
}

func TestPostHandlerDeleteComment(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи

	}

	resultPost := []*post.Post{
		&post.Post{
			Author:   user.User{ID: int64(3), Login: "arin0"},
			AuthorID: "arin0",
			Category: "programming",
			Comments: []comment.Comment{
				{Author: user.User{ID: int64(3), Login: "arin0"},
					Body: "fgh", Created: "2022-05-10T20:45:17+03:00", ID: int64(1)},
			},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
		&post.Post{
			Author:           user.User{ID: int64(3), Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "programming",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: int64(3), Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().DeleteComment("1", int64(1)).
		Return(resultPost[1], nil)

	req := httptest.NewRequest("GET", "/api/post/1/1", nil)
	w := httptest.NewRecorder()

	req.Header.Add("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID":    "1",
		"COMMENT_ID": "1",
	})
	service.DeleteComment(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[1])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().DeleteComment("2", int64(1)).
		Return(nil, fmt.Errorf("no results"))

	req1 := httptest.NewRequest("GET", "/api/post/2/1", nil)
	w1 := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/json")

	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID":    "2",
		"COMMENT_ID": "1",
	})

	service.DeleteComment(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

}

func TestPostHandlerAddComment(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	resultPost := []*post.Post{
		&post.Post{
			Author:   user.User{ID: int64(3), Login: "arin0"},
			AuthorID: "arin0",
			Category: "programming",
			Comments: []comment.Comment{
				{Author: user.User{ID: int64(3), Login: "arin0"},
					Body: "fgh", Created: "2022-05-10T20:45:17+03:00", ID: int64(1)},
			},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
		&post.Post{
			Author:   user.User{ID: int64(3), Login: "arin0"},
			AuthorID: "arin0",
			Category: "programming",
			Comments: []comment.Comment{
				{Author: user.User{ID: int64(3), Login: "arin0"},
					Body: "fgh", Created: "2022-05-10T20:45:17+03:00", ID: int64(1)},
			},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().AddComment("1", "fgh", &(resultPost[0].Comments[0].Author)).
		Return(resultPost[0], nil)
	req := httptest.NewRequest("POST", "/api/post/1", strings.NewReader(`{"comment": "fgh"}`))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.AddComment(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[0])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}

	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат
	st.EXPECT().AddComment("1", "fgh", &(resultPost[0].Comments[0].Author)).
		Return(nil, fmt.Errorf("no results"))
	req1 := httptest.NewRequest("POST", "/api/post/1", strings.NewReader(`{"comment": "fgh"}`))

	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})
	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.AddComment(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//st.EXPECT().AddComment("1", "fgh", &(resultPost[0].Comments[0].Author)).
	//	Return(nil, fmt.Errorf("no results"))
	req2 := httptest.NewRequest("POST", "/api/post/1", strings.NewReader(`{"comment": ""}`))

	w2 := httptest.NewRecorder()
	req2 = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})
	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.AddComment(w2, req2)

	resp2 := w2.Result()
	if resp2.StatusCode != 422 {
		t.Errorf("expected resp status 422, got %d", resp2.StatusCode)
		return
	}

}

func TestPostHandlerUpdate(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "music",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().UpdateVote(1, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(resultPost[0], nil)
	req := httptest.NewRequest("GET", "/api/post/1/upvote", strings.NewReader(`{"category": "music", "type": "text", "title": "Post Title ex", "text": "Post Text exemple"}`))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.Upvote(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[0])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}
	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().UpdateVote(1, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(nil, fmt.Errorf("no results"))
	req1 := httptest.NewRequest("GET", "/api/post/1/upvote", nil)

	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})

	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")

	service.Upvote(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibXl1c2VybmFtZTIiLCJpZCI6IjYyNjZlNTVmMTg1NzNiMDAwOWNjYWIwNSJ9LCJpYXQiOjE2NTIyMDE5NzYsImV4cCI6MTY1MjgwNjc3Nn0.MZTTa4TCyDF98L4X4PC7Kz9pxZ6DZzFlP0zdOuv6KjU")
}

func TestPostHandlerDownvote(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "music",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: -1}},
		},
	}

	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().UpdateVote(-1, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(resultPost[0], nil)
	req := httptest.NewRequest("GET", "/api/post/1/downvote", strings.NewReader(`{"category": "music", "type": "text", "title": "Post Title ex", "text": "Post Text exemple"}`))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.Downvote(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[0])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}
	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().UpdateVote(-1, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(nil, fmt.Errorf("no results"))
	req1 := httptest.NewRequest("GET", "/api/post/1/downvote", nil)

	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})

	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")

	service.Downvote(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibXl1c2VybmFtZTIiLCJpZCI6IjYyNjZlNTVmMTg1NzNiMDAwOWNjYWIwNSJ9LCJpYXQiOjE2NTIyMDE5NzYsImV4cCI6MTY1MjgwNjc3Nn0.MZTTa4TCyDF98L4X4PC7Kz9pxZ6DZzFlP0zdOuv6KjU")
}

func TestPostHandlerUnvote(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	resultPost := []*post.Post{
		&post.Post{
			Author:           user.User{ID: 3, Login: "arin0"},
			AuthorID:         "arin0",
			Category:         "music",
			Comments:         []comment.Comment{},
			Created:          "2022-05-10T13:31:10+03:00",
			ID:               "1",
			Score:            1,
			Text:             "Post Text exemple",
			Title:            "Post Title ex",
			Type:             "text",
			UpvotePercentage: 100,
			Views:            1,
			Votes:            []vote.Vote{{User: 3, Vote: 1}},
		},
		&post.Post{
			Author:   user.User{ID: 3, Login: "arin0"},
			Category: "music",
			Text:     "Post Text exemple",
			Title:    "Post Title ex",
			Type:     "text",
		},
	}

	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().UpdateVote(0, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(resultPost[0], nil)
	req := httptest.NewRequest("GET", "/api/post/1/unvote", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.Unvote(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	ans, _ := json.Marshal(resultPost[0])

	if !bytes.Contains(body, ans) {
		t.Errorf("Bad ans")
		return
	}
	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().UpdateVote(0, "1", &user.User{ID: 3, Login: "arin0"}).
		Return(nil, fmt.Errorf("no results"))
	req1 := httptest.NewRequest("GET", "/api/post/1/unvote", nil)

	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})

	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")

	service.Unvote(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibXl1c2VybmFtZTIiLCJpZCI6IjYyNjZlNTVmMTg1NzNiMDAwOWNjYWIwNSJ9LCJpYXQiOjE2NTIyMDE5NzYsImV4cCI6MTY1MjgwNjc3Nn0.MZTTa4TCyDF98L4X4PC7Kz9pxZ6DZzFlP0zdOuv6KjU")
}

func TestPostHandlerDeletePost(t *testing.T) {

	// мы передаём t сюда, это надо чтобы получить корректное сообщение если тесты не пройдут
	ctrl := gomock.NewController(t)

	// Finish сравнит последовательсноть вызовов и выведет ошибку если последовательность другая
	defer ctrl.Finish()

	st := post.NewMockPostRepo(ctrl) //хз работает ли
	service := &PostHandler{
		PostRepo: st,
		Logger:   zap.NewNop().Sugar(), // не пишет логи
	}

	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().Delete("1").
		Return(true, nil)
	req := httptest.NewRequest("DELETE", "/api/post/1", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, map[string]string{
		"POST_ID": "1",
	})
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")
	service.DeletePost(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Errorf("expected resp status 200, got %d", resp.StatusCode)
		return
	}
	// GetPhotos error
	// тут мы записываем последовтаельность вызовов и результат

	st.EXPECT().Delete("1").
		Return(false, fmt.Errorf("no results"))
	req1 := httptest.NewRequest("DELETE", "/api/post/1", nil)

	w1 := httptest.NewRecorder()
	req1 = mux.SetURLVars(req1, map[string]string{
		"POST_ID": "1",
	})

	req1.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")

	service.DeletePost(w1, req1)

	resp1 := w1.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp1.StatusCode)
		return
	}

	st.EXPECT().Delete("1").
		Return(false, nil)
	req2 := httptest.NewRequest("DELETE", "/api/post/1", nil)

	w2 := httptest.NewRecorder()
	req2 = mux.SetURLVars(req2, map[string]string{
		"POST_ID": "1",
	})

	req2.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTkxNzk4MDMsImlhdCI6MTY1MjE3ODYwMywidXNlciI6eyJpZCI6MywidXNlcm5hbWUiOiJhcmluMCJ9fQ.v0kJrua2GXD_3931ZPCw3ydnerK333LUsdTFYS2aYAE")

	service.DeletePost(w2, req2)

	resp2 := w2.Result()
	if resp1.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp2.StatusCode)
		return
	}

	//req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibXl1c2VybmFtZTIiLCJpZCI6IjYyNjZlNTVmMTg1NzNiMDAwOWNjYWIwNSJ9LCJpYXQiOjE2NTIyMDE5NzYsImV4cCI6MTY1MjgwNjc3Nn0.MZTTa4TCyDF98L4X4PC7Kz9pxZ6DZzFlP0zdOuv6KjU")
}
