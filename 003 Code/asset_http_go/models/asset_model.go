package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Asset struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string             `bson:"name" json:"name" validate:"required"`
	CategoryID    int                `bson:"categoryid" json:"categoryid" validate:"required"`
	Thumbnail     []byte             `bson:"thumbnail" json:"thumbnail" validate:"required"` // 썸네일의
	File          []byte             `bson:"file" json:"file" validate:"required"`           // 파일 DB의 PK
	UploadDate    time.Time          `bson:"uploaddate" json:"uploaddate"`
	DownloadCount int                `bson:"downloadcount" json:"downloadcount"`
	Price         float64            `bson:"price" json:"price"`
	IsDisable     bool               `bson:"isdisable" json:"isdisable"`
}
