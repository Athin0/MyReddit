package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"math"
	"redditclone/pkg/comment"
	"redditclone/pkg/post"
	"redditclone/pkg/user"
	"redditclone/pkg/vote"
	"strconv"
	"time"
)

type PostDB struct {
	data      post.PostDataFunctional
	lastIndex int64
}

func NewPostDB(db post.PostDataFunctional) *PostDB {
	return &PostDB{data: db,
		lastIndex: db.Len(),
	}
}

func (m *PostDB) Add(c *post.Post) (*post.Post, error) {
	c.Created = time.Now().Format(time.RFC3339)
	c.UpvotePercentage = 0
	c.Views = 0
	c.Score = 0
	c.AuthorID = c.Author.Login
	c.Comments = []comment.Comment{}
	c.Votes = []vote.Vote{}
	c.ID = strconv.Itoa(int(m.lastIndex + 1))
	_, err := m.data.Add(c)
	if err != nil {
		log.Println("err in Add PostDB:", err)
		return nil, err
	}
	return c, nil
}
func (m *PostDB) Get(id string) (*post.Post, error) {
	return m.data.Get(id)
}
func (m *PostDB) GetAll() ([]*post.Post, error) {
	return m.data.GetFilter(context.TODO(), bson.M{})
}
func (m *PostDB) GetInCategory(c string) ([]*post.Post, error) {
	return m.data.GetFilter(context.TODO(), bson.M{"category": c})
}
func (m *PostDB) GetFromUser(userName string) ([]*post.Post, error) {
	return m.data.GetFilter(context.TODO(), bson.M{"authorID": userName})
}
func (m *PostDB) AddComment(id string, text string, author *user.User) (*post.Post, error) {
	post, err := m.Get(id)
	if err != nil {
		log.Println("err in AddComment:", err)
		return nil, err
	}
	post.Comments, err = comment.Create(post.Comments, text, author)
	if err != nil {
		return nil, err
	}
	return m.data.AddComm(post)
}
func (m *PostDB) DeleteComment(idPost string, idComment int64) (*post.Post, error) {
	post, err := m.Get(idPost)
	if err != nil {
		log.Println(err, "err in AddComment")
		return nil, err
	}
	post.Comments = comment.Delete(post.Comments, idComment)
	return m.data.DeleteComm(post)
}
func (m *PostDB) UpdateVote(coin int, idPost string, author *user.User) (*post.Post, error) {
	ans, err := m.Get(idPost)
	if err != nil {
		log.Println(err, "err in UpdateVote:", err)
		return nil, err
	}
	c, err := vote.MakeVoteArr(ans.Votes, coin, author)
	if err != nil {
		log.Println("err in UpdVote, aft makeVote:", err)
		return nil, err
	}
	ans.Votes = c
	UpdateScore(ans)
	return m.data.UpVote(ans)
}
func (m *PostDB) Delete(id string) (bool, error) {
	return m.data.Delete(id)
}

func UpdateScore(post *post.Post) {
	score := 0
	upvotes := 0
	votes := len(post.Votes)
	for _, item := range post.Votes {
		score += item.Vote
		if item.Vote == 1 {
			upvotes++
		}
	}
	post.Score = score
	if votes == 0 {
		post.UpvotePercentage = 0
		return
	}
	post.UpvotePercentage = int(math.Abs(float64(upvotes) / float64(votes) * 100))
}
