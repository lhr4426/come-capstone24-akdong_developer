package controllers

import (
	"asset_http_go/configs"
	"asset_http_go/models"
	"asset_http_go/responses"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var assetCollection *mongo.Collection = configs.GetCollection(configs.DB, "assets")
var validate = validator.New()

func CreateAsset() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var asset models.Asset
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&asset); err != nil {
			c.JSON(http.StatusBadRequest, responses.AssetResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&asset); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.AssetResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newAsset := models.Asset{
			ID:            primitive.NewObjectID(),
			Name:          asset.Name,
			CategoryID:    asset.CategoryID,
			Thumbnail:     asset.Thumbnail,
			File:          asset.File,
			UploadDate:    asset.UploadDate,
			DownloadCount: asset.DownloadCount,
			Price:         asset.Price,
			IsDisable:     asset.IsDisable,
		}

		result, err := assetCollection.InsertOne(ctx, newAsset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.AssetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.AssetResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetAsset() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		assetID := c.Param("assetID")
		var asset models.Asset
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(assetID)

		err := assetCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&asset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.AssetResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.AssetResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": asset}})
	}
}
