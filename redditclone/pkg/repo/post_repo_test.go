package repo

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/repo/mocks"
	"redditclone/pkg/user"
	"redditclone/pkg/vote"
	"testing"
)

var (
	postEx = &post.Post{
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
	}

	postEx2 = &post.Post{
		Author:           user.User{ID: 3, Login: "arin0"},
		AuthorID:         "arin0",
		Category:         "music",
		Comments:         []comment.Comment{{postEx.Author, "textCom", "89789", 1}},
		Created:          "2022-05-10T13:31:10+03:00",
		ID:               "2",
		Score:            1,
		Text:             "Post Text ",
		Title:            "Post Title ex",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            1,
		Votes:            []vote.Vote{{User: 3, Vote: 1}},
	}
	postExCom = &post.Post{
		Author:           user.User{ID: 3, Login: "arin0"},
		AuthorID:         "arin0",
		Category:         "music",
		Comments:         []comment.Comment{{postEx2.Author, "textCom", "89789", 1}},
		Created:          "2022-05-10T13:31:10+03:00",
		ID:               "3",
		Score:            1,
		Text:             "Post Text ",
		Title:            "Post Title ex",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            1,
		Votes:            []vote.Vote{{User: 3, Vote: 1}},
	}
)

func TestAdd(t *testing.T) {
	db := InitMock()
	db.(*mocks.PostDataFunctional).
		On("Len").
		Return(int64(2))
	db.(*mocks.PostDataFunctional).
		On("Add", postEx).
		Return(postEx, nil)

	DB := NewPostDB(db)
	res, err := DB.Add(postEx)

	assert.NoError(t, err)
	assert.Equal(t, postEx, res)
	db.(*mocks.PostDataFunctional).
		On("Add", postEx2).
		Return(nil, errors.New("mocked-user-db-error"))

	res0, err0 := DB.Add(postEx2)
	assert.Empty(t, res0)
	assert.EqualError(t, err0, "mocked-user-db-error")
}

func InitMock() post.PostDataFunctional {
	db := &mocks.PostDataFunctional{}
	db.
		On("Len").
		Return(int64(2))
	return db
}

func TestGets(t *testing.T) {
	db := InitMock()
	id := postEx.ID
	db.(*mocks.PostDataFunctional).
		On("Get", id).
		Return(postEx, nil)

	DB := NewPostDB(db)
	res, err := DB.Get(id)

	assert.NoError(t, err)
	assert.Equal(t, postEx, res)
	db.(*mocks.PostDataFunctional).
		On("GetFilter", context.TODO(), bson.M{}).
		Return([]*post.Post{postEx, postEx2}, nil)

	res0, err := DB.GetAll()
	assert.NoError(t, err)
	assert.Equal(t, []*post.Post{postEx, postEx2}, res0)

	db.(*mocks.PostDataFunctional).
		On("GetFilter", context.TODO(), bson.M{"category": "music"}).
		Return([]*post.Post{postEx2}, nil)

	res1, err := DB.GetInCategory("music")
	assert.NoError(t, err)
	assert.Equal(t, []*post.Post{postEx2}, res1)

	userName := postEx2.AuthorID
	db.(*mocks.PostDataFunctional).
		On("GetFilter", context.TODO(), bson.M{"authorID": userName}).
		Return([]*post.Post{postEx2}, nil)

	res2, err := DB.GetFromUser(userName)
	assert.NoError(t, err)
	assert.Equal(t, []*post.Post{postEx2}, res2)
}

func TestAddComment(t *testing.T) {
	db := InitMock()
	postId := postEx.ID
	postans := postEx
	postans.Comments = []comment.Comment{
		{postEx2.Author, "textCom", "89789", 0},
	}
	db.(*mocks.PostDataFunctional).
		On("Get", postId).
		Return(postEx, nil)
	db.(*mocks.PostDataFunctional).
		On("AddComm", postEx).
		Return(postans, nil)

	DB := NewPostDB(db)

	res, err := DB.AddComment(postId, "textCom", &postEx2.Author)

	assert.NoError(t, err)
	assert.Equal(t, postEx, res)

	db.(*mocks.PostDataFunctional).
		On("Get", "89").
		Return(nil, errors.New("some db err"))
	db.(*mocks.PostDataFunctional).
		On("AddComm", postEx).
		Return(postans, nil)

	res0, err0 := DB.AddComment("89", "textCom0", &postEx2.Author)

	assert.Empty(t, res0)
	assert.EqualError(t, err0, "some db err")

	db.(*mocks.PostDataFunctional).
		On("Get", "93").
		Return(postEx2, nil)
	db.(*mocks.PostDataFunctional).
		On("AddComm", postEx2).
		Return(nil, errors.New("some db err"))

	res1, err1 := DB.AddComment("93", "textCom0", &postEx2.Author)

	assert.Empty(t, res1)
	assert.EqualError(t, err1, "some db err")

}

func TestDeleteComment(t *testing.T) {
	db := InitMock()
	postId := postExCom.ID
	postans := *postExCom
	postans.Comments = []comment.Comment{}
	DB := NewPostDB(db)

	db.(*mocks.PostDataFunctional).
		On("Get", postId).
		Return(postExCom, nil)
	db.(*mocks.PostDataFunctional).
		On("DeleteComm", &postans).
		Return(&postans, nil)

	res, err := DB.DeleteComment(postId, int64(1))

	assert.NoError(t, err)
	assert.Equal(t, &postans, res)

	db.(*mocks.PostDataFunctional).
		On("Get", "89").
		Return(nil, errors.New("some db err"))
	db.(*mocks.PostDataFunctional).
		On("DeleteComm", postExCom).
		Return(&postans, nil)

	res0, err0 := DB.DeleteComment("89", 0)

	assert.Empty(t, res0)
	assert.EqualError(t, err0, "some db err")

	postans2 := *postEx2
	postans2.Comments = []comment.Comment{}
	db.(*mocks.PostDataFunctional).
		On("Get", "93").
		Return(postEx2, nil)
	db.(*mocks.PostDataFunctional).
		On("DeleteComm", &postans2).
		Return(nil, errors.New("some db err"))

	res1, err1 := DB.DeleteComment("93", 0)

	assert.Empty(t, res1)
	assert.EqualError(t, err1, "some db err")
}

func TestUpdateVote(t *testing.T) {
	db := InitMock()
	DB := NewPostDB(db)
	postEx = &post.Post{Author: user.User{ID: 3, Login: "arin0"}, AuthorID: "arin0", Category: "programming", Comments: []comment.Comment{}, Created: "2022-05-10T13:31:10+03:00", ID: "1", Score: 1, Text: "Post Text exemple", Title: "Post Title ex", Type: "text", UpvotePercentage: 100, Views: 1, Votes: []vote.Vote{{User: 3, Vote: 1}}}
	postId := postEx.ID
	postans := *postEx
	postans.Score = 0
	postans.UpvotePercentage = 0
	postans.Votes = []vote.Vote{}

	db.(*mocks.PostDataFunctional).
		On("Get", postId).
		Return(postEx, nil)
	db.(*mocks.PostDataFunctional).
		On("UpVote", &postans).
		Return(&postans, nil)

	res, err := DB.UpdateVote(1, postId, &user.User{ID: 3, Login: "arin0"})

	assert.NoError(t, err)
	assert.Equal(t, &postans, res)

	postEx4 := &post.Post{Author: user.User{ID: 3, Login: "arin0"}, AuthorID: "arin0", Category: "programming", Comments: []comment.Comment{}, Created: "2022-05-10T13:31:10+03:00", ID: "4", Score: 1, Text: "Post Text exemple", Title: "Post Title ex", Type: "text", UpvotePercentage: 100, Views: 1, Votes: []vote.Vote{{User: 3, Vote: 1}}}
	postId = postEx4.ID
	postans3 := *postEx4
	postans3.Score = 2
	postans3.UpvotePercentage = 100
	postans3.Votes = []vote.Vote{{User: 3, Vote: 1}, {User: 1, Vote: 1}}

	db.(*mocks.PostDataFunctional).
		On("Get", postId).
		Return(postEx4, nil)
	db.(*mocks.PostDataFunctional).
		On("UpVote", &postans3).
		Return(&postans3, nil)

	res3, err := DB.UpdateVote(1, postId, &user.User{ID: 1, Login: "arin0"})

	assert.NoError(t, err)
	assert.Equal(t, &postans3, res3)

	res3, err3 := DB.UpdateVote(6, postId, &user.User{ID: 3, Login: "arin0"})
	assert.Empty(t, res3)
	assert.EqualError(t, err3, vote.ErrBadVote.Error())
	db.(*mocks.PostDataFunctional).
		On("Get", "89").
		Return(nil, errors.New("some db err"))
	db.(*mocks.PostDataFunctional).
		On("UpVote", &postans).
		Return(&postans, nil)

	res0, err0 := DB.UpdateVote(-1, "89", &user.User{ID: 3, Login: "arin0"})

	assert.Empty(t, res0)
	assert.EqualError(t, err0, "some db err")

	postans2 := *postEx
	postans2.Score = -1
	postans2.UpvotePercentage = 0
	postans2.Votes = []vote.Vote{{User: 3, Vote: -1}}
	db.(*mocks.PostDataFunctional).
		On("Get", "93").
		Return(postEx, nil)
	db.(*mocks.PostDataFunctional).
		On("UpVote", &postans2).
		Return(nil, errors.New("some db err"))

	res1, err1 := DB.UpdateVote(-1, "93", &user.User{ID: 3, Login: "arin0"})

	assert.Empty(t, res1)
	assert.EqualError(t, err1, "some db err")

}
func TestDelete(t *testing.T) {
	db := InitMock()
	id := postEx.ID
	db.(*mocks.PostDataFunctional).
		On("Delete", id).
		Return(true, nil)

	DB := NewPostDB(db)
	res, err := DB.Delete(id)

	assert.NoError(t, err)
	assert.Equal(t, true, res)
}
