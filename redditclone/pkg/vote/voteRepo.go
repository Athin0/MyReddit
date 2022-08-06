package vote

import (
	"errors"
	"redditclone/pkg/user"
	"sync"
)

var ErrBadVote = errors.New("Bab vote number for Vote")

type VoteMemoryRepo struct {
	num   int     `bson:"num" json:"num"`
	data  []*Vote `json:"data" bson:"data"`
	mutex sync.Mutex
}

func NewMemoryRepo() *VoteMemoryRepo {
	return &VoteMemoryRepo{
		num:   0,
		data:  make([]*Vote, 0, 10),
		mutex: sync.Mutex{},
	}
}

func (c *VoteMemoryRepo) GetAll() []*Vote {
	return c.data
}

func (c *VoteMemoryRepo) Make(coin int, user *user.User) (*Vote, error) {
	if coin == 0 {
		c.Delete(user)
		return nil, nil
	}
	for _, item := range c.data {
		if item.User == user.ID {
			item.Vote = coin
			return item, nil
		}
	}
	item := new(Vote)
	item.User = user.ID
	item.Vote = coin
	c.mutex.Lock()
	c.num++
	c.data = append(c.data, item)
	c.mutex.Unlock()
	return item, nil
}
func (c *VoteMemoryRepo) Delete(user *user.User) []*Vote {
	var k int
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i, item := range c.data {
		if item.User == user.ID {
			k = i
			break
		}
	}
	c.num--
	c.data = append(c.data[:k], c.data[k+1:]...)
	return c.data
}
func (c *VoteMemoryRepo) NumVotes() int {
	return c.num
}
func MakeVoteArr(c []Vote, coin int, user *user.User) ([]Vote, error) {
	if (coin != 1) && coin != 0 && coin != -1 {
		return nil, ErrBadVote
	}
	if coin == 0 {
		c = DeleteArr(c, user)
		return c, nil
	}
	for _, item := range c {
		if item.User == user.ID {
			c = DeleteArr(c, user)
			if coin == item.Vote {
				return c, nil
			}
		}
	}
	item := new(Vote)
	item.User = user.ID
	item.Vote = coin
	c = append(c, *item)
	return c, nil
}
func DeleteArr(c []Vote, user *user.User) []Vote {
	var k int
	for i, item := range c {
		if item.User == user.ID {
			k = i
			break
		}
	}
	c = append(c[:k], c[k+1:]...)
	return c
}
