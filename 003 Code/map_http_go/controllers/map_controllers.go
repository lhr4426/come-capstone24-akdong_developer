package controllers

import (
	"context"
	"fmt"
	"capstone.com/configs"
	"capstone.com/models"
	"capstone.com/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mapCollection *mongo.Collection = configs.GetCollection(configs.DB, "map")
var validate = validator.New()

func CreateMap() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("@@@@@") // 여긴 넘어가지도 않음
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var map_data models.MapData // map 함수 이름 존재하므로 map쓰면 안됨
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&map_data); err != nil {
			c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadRequest, Message: "error"})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&map_data); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadRequest, Message: "error"})
			return
		}

		newMap := models.MapData{
			Map_ID:  primitive.NewObjectID(),
			Version: map_data.Version,
			Chunk:   map_data.Chunk,
		}

		_, err := mapCollection.InsertOne(ctx, newMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Status: http.StatusInternalServerError, Message: "error"})
			return
		}

		c.JSON(http.StatusCreated, responses.MapResponse{Status: http.StatusCreated, Message: "success"})
	}
}

func GetAMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		mapId := c.Param("mapId")
		var user models.MapData
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(mapId)

		err := mapCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error"})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success"})
	}
}
