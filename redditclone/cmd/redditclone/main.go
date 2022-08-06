package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"redditclone/pkg/handler"
	"redditclone/pkg/middleware"
	"redditclone/pkg/repo"
	"redditclone/pkg/session"
	"redditclone/pkg/user"
)

func main() {
	// основные настройки к базе
	dsn := "root:love@tcp(localhost:3306)/golang?"
	dsn += "charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	db.SetMaxOpenConns(10)
	err = db.Ping() // проверяем подключение
	if err != nil {
		panic(err)
	}
	log.Println("Connected to MySQL!")

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	collection := client.Database("coursera").Collection("posts")

	zapLogger, errZap := zap.NewProduction() //create logger
	if errZap != nil {
		log.Println("Error in creation zapLogger")
	}
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			log.Println(err)
		}
	}(zapLogger)
	logger := zapLogger.Sugar()

	userRepo := user.NewMysqlRepo(db)
	postRepo := repo.NewMongoRepo(collection)
	sessRepo := session.NewSessionsRepo(db)
	userHandler := &handler.UserHandler{
		UserRepo: userRepo,
		Logger:   logger,
		Sessions: sessRepo,
	}
	postHandler := &handler.PostHandler{
		PostRepo: repo.NewPostDB(postRepo),
		Logger:   logger,
		Sessions: sessRepo,
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/posts", postHandler.Add).Methods("POST")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postHandler.Category).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postHandler.Get).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postHandler.AddComment).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postHandler.DeleteComment).Methods("DELETE")
	r.HandleFunc("/api/posts/", postHandler.List).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/upvote", postHandler.Upvote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/downvote", postHandler.Downvote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/unvote", postHandler.Unvote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postHandler.DeletePost).Methods("DELETE")
	r.HandleFunc("/api/user/{USER_LOGIN}", postHandler.GetPostsOfUser).Methods("GET")

	r.HandleFunc("/api/login", userHandler.Re).Methods("POST")
	r.HandleFunc("/api/register", userHandler.RegisterPage).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Handle("/", http.FileServer(http.Dir("./static/html/")))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, errReadFile := ioutil.ReadFile("./static/html/index.html")
		if errReadFile != nil {
			log.Println(errReadFile)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(file)
		if err != nil {
			log.Println("err in Write", err)
			return
		}
	})

	mux0 := middleware.Auth(r)
	mux0 = middleware.AccessLog(logger, mux0)
	mux0 = middleware.Panic(mux0)

	fmt.Println("starting server at :8080")
	errListen := http.ListenAndServe(":8080", mux0)
	if errListen != nil {
		log.Println("err in listen and serve", errListen)
		return
	}
}
