package controllers

import (
    "context"
    "capstone.com/configs"
    "capstone.com/models"
    "capstone.com/responses"
    "net/http"
    "time"
	"fmt"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var mapCollection *mongo.Collection = configs.GetCollection(configs.DB, "map")
var validate = validator.New()

func CreateMap() gin.HandlerFunc {
    return func(c *gin.Context) {

		fmt.println("@@@@@")
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var map models.Map
        defer cancel()

        //validate the request body
        if err := c.BindJSON(&map); err != nil {
            c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadRequest, Message: "error"})
            return
        }

        //use the validator library to validate required fields
        if validationErr := validate.Struct(&map); validationErr != nil {
            c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadFRequest, Message: "error"})
            return
        }

        newMap := models.Map{
            Map_ID:       primitive.NewObjectID(),
            Version:     map.Version,
            Chunk: 		map.Chunk,
        }

        result, err := mapCollection.InsertOne(ctx, newMap)
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.MapResponse{Status: http.StatusInternalServerError, Message: "error"})
            return
        }

        c.JSON(http.StatusCreated, responses.MapResponse{Status: http.StatusCreated, Message: "success"})
    }
}