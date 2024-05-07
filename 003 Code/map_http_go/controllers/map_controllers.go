package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
    "strconv"

	"capstone.com/configs"
	"capstone.com/models"
	"capstone.com/responses"

	"github.com/gin-gonic/gin"
//	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var mapCollection *mongo.Collection = configs.GetCollection(configs.DB, "map")
// var validate = validator.New() // 동작은 됨

func CreateMap() gin.HandlerFunc {
	return func(c *gin.Context) {

		// fmt.Println("@@@@@@@@@@")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

        // 동적으로 Key-Value 쌍 저장 가능
        mapdata := make(map[string]interface{})
		// var map_data models.MapData // map 함수 이름 존재하므로 map쓰면 안됨


		defer cancel()

		//body 유효성 검증
		if err := c.BindJSON(&mapdata); err != nil {
			fmt.Println("1@@@@@@@@@@@@@@@@@@@@@@@@")
			c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadRequest, Message: "error"})
			return
		}

        // 검증

		// //validation library를 통한 필수 필드 검증
		// if validationErr := validate.Struct(&mapdata); validationErr != nil {
		// 	fmt.Println("2@@@@@@@@@@@@@@@@@@@@@@@@")
		// 	c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadRequest, Message: "error"})
		// 	return
		// }

        // newMap := map[string]interface{}{
        //     "id": primitive.NewObjectID(),
        //     "MapData": mapdata["MapData"],
        // }

		// newMap := models.MapData{
		// 	Id:      primitive.NewObjectID(),
		// 	// Map_ID:  map_data.Map_ID,
		// 	// Version: map_data.Version,
		// 	// Chunk:   map_data.Chunk,
        //     // MapJson: map_data.MapJson,
        //     MapTotal: map_data.MapTotal,
            
		// }

		_, err := mapCollection.InsertOne(ctx, mapdata)
		if err != nil {
			fmt.Println("3@@@@@@@@@@@@@@@@@@@@@@@@")
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Status: http.StatusInternalServerError, Message: "error"})
			return
		}

		c.JSON(http.StatusCreated, responses.MapResponse{Status: http.StatusCreated, Message: "success"})
	}
}

func GetAMap() gin.HandlerFunc {
	return func(c *gin.Context) {

        // --> 진행중
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        // url에서 매개변수 추출
		mapId := c.Query("mapId")
        mapVersion, _ := strconv.Atoi(c.Query("version"))
		mapChunk, _ := strconv.Atoi(c.Query("chunk"))

        // MongoDB에서 데이터 조회
		var mapData models.MapData
		filter := bson.M{
			"map_id":  mapId,
			"version": mapVersion,
			"chunk":   mapChunk,
		}
		err := mapCollection.FindOne(ctx, filter).Decode(&mapData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 데이터 반환
		c.JSON(http.StatusOK, mapData)
	}

}
