package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"capstone.com/configs"
	"capstone.com/responses"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mapCollection *mongo.Collection = configs.GetCollection(configs.DB, "map")


// json 파일 DB 저장
func CreateMap() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// 동적으로 Key-Value 쌍 저장 가능
		mapdata := make(map[string]interface{}) // 빈 맵 생성(동적으로 key-value쌍 저장 가능)
		// var map_data models.MapData // map 함수 이름 존재하므로 map쓰면 안됨 + 동적 저장이므로 struct 필요 없음

		defer cancel()

        // 예외처리(기존에 있는 정보인지 확인 필요, 기존에 있는 정보라면 변경)
        


		//body 유효성 검증
		if err := c.BindJSON(&mapdata); err != nil {
			c.JSON(http.StatusBadRequest, responses.MapResponse{Status: http.StatusBadRequest, Message: "error"})
			return
		}

		_, err := mapCollection.InsertOne(ctx, mapdata) // DB에 바로 저장
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Status: http.StatusInternalServerError, Message: "error"})
			return
		}

		c.JSON(http.StatusCreated, responses.MapResponse{Status: http.StatusCreated, Message: "success"})
	}
}

// GET은 반대로 map안에 데이터를 저장하고 json으로 보내주면 되지 않을까?
// query는 어떻게 전송해야하나(router에서는 /만처리해주고 Query를 통해서 매개변수 확인)
// 자꾸 collection 찾을 수 없다고 뜸(형 변환 문제)
func GetMap() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// // url에서 매개변수 추출
		mapId, _ := strconv.ParseFloat(c.Query("mapID"), 64) // double로 저장되어 있는데 .Atoi(int변환)으로 해서 반환 못하는 문제 발생 (해결완)
		mapVersion, _ := strconv.ParseFloat(c.Query("version"), 64)
		mapChunk, _ := strconv.ParseFloat(c.Query("chunk"), 64)

		// MongoDB에서 데이터 조회
		mapinfo := make(map[string]interface{})

		filter := bson.M{
			"map_id":   mapId,
			"version":  mapVersion,
			"chunkNum": mapChunk,
		}

		fmt.Println("filter", filter)

		err := mapCollection.FindOne(ctx, filter).Decode(&mapinfo)
		fmt.Println("mapinfo", mapinfo)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "No documents found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 데이터 반환
		c.JSON(http.StatusOK, mapinfo)
	}
}
