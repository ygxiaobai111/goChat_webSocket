package service

import (
	"context"
	"go-chat/goChat_webSocket/conf"
	"go-chat/goChat_webSocket/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 插入文章
func InsertArticle(article model.Article) error {

	collection := conf.MongoDBClient.Database("mydb").Collection("articles")

	_, err := collection.InsertOne(context.Background(), article)

	if err != nil {
		return err
	}

	return nil
}

// 修改文章标题
func UpdateArticleTitle(articleID string, newTitle, newContent string) error {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")
	obj_id, err := primitive.ObjectIDFromHex(articleID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": obj_id}

	updateTitle := bson.M{"$set": bson.M{"title": newTitle}}
	updateContent := bson.M{"$set": bson.M{"topic": newContent}}
	//修改标题
	_, err = collection.UpdateOne(context.Background(), filter, updateTitle)

	if err != nil {
		return err
	}
	//修改内容
	_, err = collection.UpdateOne(context.Background(), filter, updateContent)
	return nil
}

// 修改文章内容
func updateArticleContent(articleID string, newContent string) error {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")

	filter := bson.M{"_id": articleID}

	update := bson.M{"$set": bson.M{"content": newContent}}

	_, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

// 添加评论到文章
func AddCommentToArticle(articleID string, comment model.Comment) error {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")
	objID, err := primitive.ObjectIDFromHex(articleID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	comment.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"comments": comment}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// 添加子评论到评论
func AddSubCommentToComment(commentID string, subComment model.SubComment) error {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")
	objID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	filter := bson.M{"comments._id": objID}
	subComment.ID = primitive.NewObjectID()
	update := bson.M{"$push": bson.M{"comments.$.sub_comments": subComment}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// 文章点赞数增加1
func LikeArticle(articleID string) error {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")
	objID, err := primitive.ObjectIDFromHex(articleID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}

	update := bson.M{"$inc": bson.M{"likes": 1}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return nil
}

// 获取作者所有文章
func GetArticlesByAuthor(authorID string) ([]model.Article, error) {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")

	filter := bson.M{"author_id": authorID}

	var articles []model.Article
	cursor, err := collection.Find(context.Background(), filter)

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// 通过文章id返回文章
func GetArticleByID(articleID string) (*model.Article, error) {
	collection := conf.MongoDBClient.Database("mydb").Collection("articles")

	filter := bson.M{"_id": articleID}

	var article model.Article
	err := collection.FindOne(context.Background(), filter).Decode(&article)

	if err != nil {
		return nil, err
	}

	return &article, nil
}
