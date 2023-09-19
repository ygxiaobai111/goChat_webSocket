package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Article struct {
	Title     string    `bson:"title"`
	AuthorID  string    `bson:"author_id"`
	BuildTime time.Time `bson:"build_time"`
	Topic     string    `bson:"topic"`
	Comments  []Comment `bson:"comments"`
	Likes     int       `bson:"likes"`
}

type Comment struct {
	ID          primitive.ObjectID `bson:"_id"`
	Text        string             `bson:"text"`
	AuthorID    string             `bson:"author_id"`
	BuildTime   time.Time          `bson:"build_time"`
	Likes       int                `bson:"likes"`
	SubComments []SubComment       `bson:"sub_comments"`
}

type SubComment struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `bson:"text"`
	AuthorID  string             `bson:"author_id"`
	BuildTime time.Time          `bson:"build_time"`
	Likes     int                `bson:"likes"`
}
