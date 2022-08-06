package post

import (
	"context"
	"redditclone/pkg/comment"
	"redditclone/pkg/user"
	"redditclone/pkg/vote"
)

type Post struct {
	Author           user.User         `json:"author" bson:"author"`
	AuthorID         string            `bson:"authorID" json:"authorID"`
	Category         string            `json:"category" bson:"category"`
	Comments         []comment.Comment `json:"comments" bson:"comments" `
	Created          string            `json:"created" bson:"created"`
	ID               string            `json:"id" bson:"id"`
	Score            int               `json:"score" bson:"score"`
	Text             string            `json:"text,omitempty" bson:"text"`
	URL              string            `json:"url,omitempty" bson:"url"`
	Title            string            `json:"title" bson:"title"`
	Type             string            `json:"type" bson:"type"`
	UpvotePercentage int               `json:"upvotePercentage" bson:"upvotePercentage"`
	Views            int               `json:"views" bson:"views"`
	Votes            []vote.Vote       `json:"votes" bson:"votes"`
}

type PostDataFunctional interface {
	Len() int64
	Add(c *Post) (*Post, error)
	Get(id string) (*Post, error)
	GetFilter(ctx context.Context, filter interface{}) ([]*Post, error)
	AddComm(post *Post) (*Post, error)
	DeleteComm(post *Post) (*Post, error)
	UpVote(ans *Post) (*Post, error)
	Delete(id string) (bool, error)
}

//go:generate mockgen -source=post.go -destination=repo_mock.go -package=post PostRepo
type PostRepo interface {
	GetAll() ([]*Post, error)
	Add(*Post) (*Post, error)
	Get(i string) (*Post, error)
	GetInCategory(c string) ([]*Post, error)
	AddComment(id string, text string, author *user.User) (*Post, error)
	DeleteComment(idPost string, idComment int64) (*Post, error)
	UpdateVote(vote int, idPost string, author *user.User) (*Post, error)
	Delete(id string) (bool, error)
	GetFromUser(userName string) ([]*Post, error)
}
