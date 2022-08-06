package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"redditclone/pkg/post"
)

type PostMongoRepo struct {
	data *mongo.Collection
	//lastIndex int64
}

func NewMongoRepo(db *mongo.Collection) *PostMongoRepo {
	return &PostMongoRepo{data: db}
}

func (repo *PostMongoRepo) Len() int64 {
	len, err := repo.data.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		log.Println("err in Len:", err)
		return -1
	}
	return len
}

func (repo *PostMongoRepo) Add(c *post.Post) (*post.Post, error) {
	_, err := repo.data.InsertOne(context.TODO(), c)
	if err != nil {
		log.Println(err)
	}
	return c, nil
}

func (repo *PostMongoRepo) Get(id string) (*post.Post, error) {
	post := &post.Post{}
	err := repo.data.FindOne(context.TODO(), bson.M{"id": id}).Decode(&post)
	if err != nil {
		log.Println("err in Get Post:", err)
		return nil, err
	}
	return post, nil
}

func (repo *PostMongoRepo) GetFilter(ctx context.Context, filter interface{}) ([]*post.Post, error) {
	var arr []*post.Post
	cur, err := repo.data.Find(ctx, filter)
	if err != nil {
		log.Println("err in read from DB:", err)
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem post.Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return arr, nil
}

func (repo *PostMongoRepo) AddComm(post *post.Post) (*post.Post, error) {
	_, err := repo.data.UpdateOne(context.TODO(), bson.M{"id": post.ID}, bson.D{
		{"$set", bson.D{{"comments", post.Comments}}},
	})
	if err != nil {
		log.Println("err in update bd in AddComment:", err)
	}
	return post, nil
}
func (repo *PostMongoRepo) DeleteComm(post *post.Post) (*post.Post, error) {
	_, err := repo.data.ReplaceOne(context.TODO(), bson.M{"id": post.ID}, &post)
	if err != nil {
		log.Println("err in update bd in DelateComment:", err)
	}
	return post, nil
}

func (repo *PostMongoRepo) UpVote(ans *post.Post) (*post.Post, error) {
	res := repo.data.FindOneAndUpdate(context.TODO(), bson.M{"id": ans.ID}, bson.D{
		{"$set", bson.M{"votes": ans.Votes, "score": ans.Score, "upvotePercentage": ans.UpvotePercentage}},
	})
	if res.Err() != nil {
		log.Println(res.Err())
		return nil, res.Err()
	}
	return ans, nil
}

func (repo *PostMongoRepo) Delete(id string) (bool, error) {
	_, err := repo.data.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}

/*
func (repo *PostMongoRepo) UpdateScore(post *Post){
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



func (repo *PostMongoRepo) AddComment(num string, text string, author *user.User) (*Post, error) {
	post, err := repo.Get(num)
	if err != nil {
		log.Println("err in AddComment:", err)
	}

	post.Comments, err = comment.Create(post.Comments, text, author)
	if err != nil {
		return nil, err
	}
	_, err = repo.data.UpdateOne(context.TODO(), bson.M{"id": post.ID}, bson.D{
		{"$set", bson.D{{"comments", post.Comments}}},
	})
	if err != nil {
		log.Println("err in update bd in AddComment:", err)
	}
	return post, nil
}

func (repo *PostMongoRepo) DeleteComment(idPost string, idComment int64) (*Post, error) {
	post, err := repo.Get(idPost)
	if err != nil {
		log.Println(err, "err in AddComment")
	}

	post.Comments = comment.Delete(post.Comments, idComment)
	_, err = repo.data.ReplaceOne(context.TODO(), bson.M{"id": post.ID}, &post)
	if err != nil {
		log.Println("err in update bd in DelateComment:", err)
	}
	return post, nil
}

func (repo *PostMongoRepo) UpdateVote(coin int, idPost string, author *user.User) (*Post, error) {
	var ans *Post
	err := repo.data.FindOne(context.TODO(), bson.M{"id": idPost}).Decode(&ans)
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
	repo.UpdateScore(ans)
	res := repo.data.FindOneAndUpdate(context.TODO(), bson.M{"id": ans.ID}, bson.D{
		{"$set", bson.M{"votes": ans.Votes, "score": ans.Score, "upvotePercentage": ans.UpvotePercentage}},
	})
	if res.Err() != nil {
		log.Println(res.Err())
		return nil, res.Err()
	}

	return ans, nil
}
func (repo *PostMongoRepo) UpdateScore(post *Post){
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



func (repo *PostMongoRepo) GetFromUser(userID string) ([]*Post, error) {
	var arr []*Post
	cur, err := repo.data.Find(context.TODO(), bson.M{"authorID": userID})
	if err != nil {
		log.Println("err in read from DB:", err)
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return arr, nil
}

func (repo *PostMongoRepo) GetAll() ([]*Post, error) {
	var arr []*Post
	cur, err := repo.data.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("err in read from DB:", err)
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		arr = append(arr, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return arr, nil
}
func (repo *PostMongoRepo) GetInCategory(category string) ([]*Post, error) {
	var arr []*Post
	cur, err := repo.data.Find(context.TODO(), bson.M{"category": category})
	if err != nil {
		log.Println("err in read from DB:", err)
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var elem Post
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		arr = append(arr, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	cur.Close(context.TODO())
	return arr, nil
}

*/
