package comment

import (
	"redditclone/pkg/user"
)

type Comment struct {
	Author  user.User `json:"author" bson:"author"`
	Body    string    `json:"body" bson:"body"`
	Created string    `json:"created" bson:"created"`
	ID      int64     `json:"id" bson:"id"`
}

type CommentRepo interface {
	GetAll() []*Comment
	Get(id int64) (*Comment, error)
	Create(text string, author *user.User) (*Comment, error)
	Delete(id int64) []*Comment
}
