package vote

import "redditclone/pkg/user"

type Vote struct {
	User int64 `json:"user" bson:"user"`
	Vote int   `json:"vote" bson:"vote"`
}
type VoteRepo interface {
	GetAll() []*Vote
	Make(coin int, author *user.User) (*Vote, error)
	Delete(user *user.User) []*Vote
	NumVotes() int
}
