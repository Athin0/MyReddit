package repo

import (
	"errors"
	"redditclone/pkg/post"
	"sync"
)

var (
	ErrNoPost = errors.New("no post found")
)

// PostMemoryRepo is old realization with storage in memory
type PostMemoryRepo struct {
	lastIndex int
	data      []*post.Post
	mutex     sync.Mutex
}

/*
func NewMemoryRepo() *PostMemoryRepo {
	return new(PostMemoryRepo)
}

func (repo *PostMemoryRepo) Add(c *post.Post) (*post.Post, error) {
	repo.mutex.Lock()
	c.Created = time.Now().Format(time.RFC3339)
	c.UpvotePercentage = 100
	c.Views = 0
	c.Score = 1
	c.Comments = *comment.NewMemoryRepo()
	c.Votes = *vote.NewMemoryRepo()
	repo.lastIndex++
	repo.data = append(repo.data, c)
	repo.mutex.Unlock()
	u, err := repo.Get(c.ID)
	if err != nil {
		log.Println(err)
	}
	return u, nil
}
func (repo *PostMemoryRepo) Get(id string) (*Post, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for _, elem := range repo.data {
		if elem.ID == id {
			return elem, nil
		}
	}
	return nil, ErrNoPost
}

func (repo *PostMemoryRepo) GetAll() ([]*post.Post, error) {
	arr := make([]*post.Post, 0, 10)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	arr = append(arr, repo.data...)
	return arr, nil
}

func (repo *PostMemoryRepo) GetInCategory(category string) ([]*post.Post, error) {
	arr := make([]*post.Post, 0, 10)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for _, elem := range repo.data {
		if elem.Category == category {
			arr = append(arr, elem)
		}
	}
	return arr, nil
}
func (repo *PostMemoryRepo) AddComment(num int64, text string, author *user.User) (*post.Post, error) {
	idPost := strconv.Itoa(int(num))
	post, err := repo.Get(idPost)
	if err != nil {
		log.Println("err in AddComment:", err)
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, err = post.Comments.Create(text, author)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *PostMemoryRepo) DeleteComment(num int64, idComment int64) (*post.Post, error) {
	idPost := strconv.Itoa(int(num))
	post, err := repo.Get(idPost)
	if err != nil {
		log.Println(err, "err in AddComment")
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	post.Comments.Delete(idComment)
	return post, nil
}

func (repo *PostMemoryRepo) UpdateVote(vote int, num int64, author *user.User) (*post.Post, error) {
	idPost := strconv.Itoa(int(num))
	post, err := repo.Get(idPost)
	if err != nil {
		log.Println(err, "err in UpdateVote:", err)
	}
	repo.mutex.Lock()
	_, err = post.Votes.Make(vote, author)
	if err != nil {
		return nil, err
	}
	repo.mutex.Unlock()
	repo.UpdateScore(post)
	return post, nil
}
func PreSendingConvert(repo []*post.Post) []*PostForSending {
	arr := make([]*PostForSending, 0, 10)
	for _, elem := range repo {
		arr = append(arr, elem.Convert())
	}
	return arr
}
func (repo *PostMemoryRepo) UpdateScore(post *post.Post) {
	score := 0
	upvotes := 0
	votes := post.Votes.NumVotes()
	for _, item := range post.Votes.GetAll() {
		score += item.Vote
		if item.Vote == 1 {
			upvotes++
		}
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	post.Score = score
	if votes == 0 {
		post.UpvotePercentage = 0
		return
	}
	post.UpvotePercentage = int(math.Abs(float64(upvotes) / float64(votes) * 100))
}
func (repo *PostMemoryRepo) Delete(id string) (bool, error) {
	var k int
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for i, item := range repo.data {
		if item.ID == id {
			k = i
			break
		}
	}
	repo.data = append(repo.data[:k], repo.data[k+1:]...)
	return true, nil
}
func (repo *PostMemoryRepo) GetFromUser(userID int64) ([]*post.Post, error) {
	arr := make([]*post.Post, 0, 10)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for _, elem := range repo.data {
		if elem.Author.ID == userID {
			arr = append(arr, elem)
		}
	}
	return arr, nil
}

type PostForSending struct {
	Author           user.User          `json:"author" bson:"author"`
	Category         string             `json:"category" bson:"category"`
	Comments         []*comment.Comment `json:"comments" bson:"comments"`
	Created          string             `json:"created" bson:"created"`
	ID               string             `json:"id" bson:"id"`
	Score            int                `json:"score" bson:"score"`
	Text             string             `json:"text,omitempty" bson:"text"`
	URL              string             `json:"url,omitempty" bson:"usl"`
	Title            string             `json:"title" bson:"title"`
	Type             string             `json:"type" bson:"type"`
	UpvotePercentage int                `json:"upvotePercentage" bson:"upvotePercentage"`
	Views            int                `json:"views" bson:"views"`
	Votes            []*vote.Vote       `json:"votes" bson:"votes"`
}

*/
