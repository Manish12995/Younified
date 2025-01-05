package model

import (
	"time"
	"younified-backend/contracts/user/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type News struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Content      string               `json:"content,omitempty" bson:"content,omitempty"`
	CreatedOn    time.Time            `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	UserID       primitive.ObjectID   `json:"userID,omitempty" bson:"userID"`
	Unit         []string             `json:"unit,omitempty" bson:"unit,omitempty"`
	Likes        []primitive.ObjectID `json:"likes" bson:"likes"`
	Dislikes     []primitive.ObjectID `json:"dislikes" bson:"dislikes"`
	Images       []string             `json:"images,omitempty" bson:"images,omitempty"`
	Deleted      bool                 `json:"deleted,omitempty" bson:"deleted"`
	CommentIDs   []primitive.ObjectID `json:"commentIDs" bson:"commentIDs"`
	Creator      model.User
	LikedBy      []model.User
	Comments     []*Comment
	Documents    []*Document `json:"documents,omitempty" bson:"documents,omitempty"`
	Pinned       bool        `json:"pinned,omitempty" bson:"pinned,omitempty"`
	Show         bool        `json:"show,omitempty" bson:"show,omitempty"`
	Private      bool        `json:"private,omitempty" bson:"private,omitempty"`
	ShowLikes    bool        `json:"showLikes,omitempty" bson:"showLikes,omitempty"`
	ShowComments bool        `json:"showComments,omitempty" bson:"showComments,omitempty"`
	AsUnion      bool        `json:"asUnion,omitempty" bson:"asUnion,omitempty"`
	Category     string      `json:"category,omitempty" bson:"category,omitempty"`
	CommentCount int64
}

type Comment struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	NewsID    primitive.ObjectID   `json:"newsID,omitempty" bson:"newsID,omitempty"`
	Content   string               `json:"content,omitempty" bson:"content,omitempty"`
	UserID    primitive.ObjectID   `json:"userID,omitempty" bson:"userID,omitempty"`
	Creator   *model.User          `json:"creator,omitempty" bson:"creator,omitempty"`
	Deleted   bool                 `json:"deleted" bson:"deleted"`
	Likes     []primitive.ObjectID `json:"likes" bson:"likes"`
	Dislikes  []primitive.ObjectID `json:"dislikes" bson:"dislikes"`
	CreatedOn time.Time            `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
}

type Report struct {
	Data   []*News
	Pinned []*News
	Total  int64
}

type Document struct {
	Url  string `json:"url,omitempty" bson:"url,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

type Blog struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedBy string             `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
	Header    string             `json:"header,omitempty" bson:"header,omitempty"`
	SubHeader string             `json:"subHeader,omitempty" bson:"subHeader,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	Images    []string           `json:"images,omitempty" bson:"images,omitempty"`
	CreatedOn time.Time          `json:"createdOn,omitempty" bson:"createdOn,omitempty"`
	Deleted   bool               `json:"deleted" bson:"deleted"`
	Featured  bool               `json:"featured" bson:"featured"`
}