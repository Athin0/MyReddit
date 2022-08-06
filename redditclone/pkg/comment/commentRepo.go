package comment

import (
	"errors"
	"redditclone/pkg/user"
	"sync"
	"time"
)

var (
	ErrNoComment = errors.New("no comment found")
)

type CommentMemoryRepo struct {
	lastIndex int        `json:"last_index" bson:"lastIndex"`
	data      []*Comment `bson:"data" json:"data"`
	mutex     sync.Mutex
}

func NewMemoryRepo() *CommentMemoryRepo {
	return &CommentMemoryRepo{
		lastIndex: 0,
		data:      make([]*Comment, 0, 10),
		mutex:     sync.Mutex{},
	}
}

func (c *CommentMemoryRepo) GetAll() []*Comment {
	return c.data
}
func (c *CommentMemoryRepo) Get(id int64) (*Comment, error) {
	for _, item := range c.data {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, ErrNoComment
}

func (c *CommentMemoryRepo) Create(text string, author *user.User) (*Comment, error) {
	item := new(Comment)
	item.ID = int64(c.lastIndex)
	item.Author = *author
	item.Created = time.Now().Format(time.RFC3339)
	item.Body = text
	c.mutex.Lock()
	c.lastIndex++
	c.data = append(c.data, item)
	c.mutex.Unlock()
	return item, nil
}

func Create(c []Comment, text string, author *user.User) ([]Comment, error) {
	item := new(Comment)
	item.ID = int64(len(c))
	item.Author = *author
	item.Created = time.Now().Format(time.RFC3339)
	item.Body = text
	c = append(c, *item)
	return c, nil
}
func Delete(c []Comment, id int64) []Comment {
	var k int
	for i, item := range c {
		if item.ID == id {
			k = i
			break
		}
	}
	c = append(c[:k], c[k+1:]...)
	return c
}

/*
func (repo *CommentMemoryRepo) AddCommentInRepo(c *comment) (*comment, error) {
	repo.mutex.Lock()
	repo.data[c.ID] = c

	repo.mutex.Unlock()
	u, ok := repo.data[c.ID]
	if !ok {
		return nil, ErrNoComment
	}
	return u, nil
}
*/
