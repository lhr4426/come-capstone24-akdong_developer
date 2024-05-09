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
		defer cancel()

		// 동적으로 Key-Value 쌍 저장 가능
		mapdata := make(map[string]interface{}) // 빈 맵 생성(동적으로 key-value쌍 저장 가능)
		// var map_data models.MapData // map 함수 이름 존재하므로 map쓰면 안됨 + 동적 저장이므로 struct 필요 없음

		//body 유효성 검증
		if err := c.BindJSON(&mapdata); err != nil {
			c.JSON(http.StatusBadRequest, responses.MapResponse{Code: 0, Message: "body error"})
			return
		}

		fmt.Println("2222222222222", mapdata) // 어떤거 나오는지 비교, json 받아온 값 확인

		// mapdata에서 map_id, version, chunkNum 확인
		mapId := mapdata["map_id"].(float64)
		mapVersion := mapdata["version"].(float64)
		mapChunk := mapdata["chunkNum"].(float64)

		fmt.Println("3333333333333", mapId, mapVersion, mapChunk) 

		// 중복 확인을 위한 filter
		filter := bson.M{
			"map_id":   mapId,
			"version":  mapVersion,
			"chunkNum": mapChunk,
		}

		err := mapCollection.FindOne(ctx, filter).Decode(&mapdata)

		// 중복 없을 때 Insert
		if err != nil {
			if err == mongo.ErrNoDocuments {
				_, err := mapCollection.InsertOne(ctx, mapdata) // DB에 바로 저장
				if err != nil {
					c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: "insert error"})
					return
				}
				c.JSON(http.StatusCreated, responses.MapResponse{Code: 1, Message: "insert success"})
				return
			}
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: err.Error()})
			return
		}

		// Update(모두 삭제 하고 Insert)
		// 모두 삭제(if chunkNum이 더 있으면)
		_, err = mapCollection.ReplaceOne(ctx, filter, mapdata)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: "error"})
			return
		}

		// 만약 map_id, version, chunk(개수), mapdata 제외하고 map_id, version있는거 모두 삭제
			// 만약 짧은 정보가 Update 된다면 나머지 Chunk_Num을 삭제
			// 첫번째가 중복되면 이후 chunk_num 모두 삭제하기, 그러면 다음부터는 create가능

		chk_filter := bson.M{
			"chunkNum": mapChunk
		}

		chk_err := mapCollection.DeleteMany(ctx, ).Decode(&mapdata)

		c.JSON(http.StatusOK, responses.MapResponse{Code: 1, Message: "update success"})

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

		// ObjectId 제거
		err := mapCollection.FindOne(ctx, filter).Decode(&mapinfo)
		fmt.Println("mapinfo", mapinfo)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, responses.MapResponse{Code: 0, Message: "No documents found"})
				return
			}
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: err.Error()})
			return
		}

		// 특정 key 제외
		filtermapinfo := make(map[string]interface{})

		for key, value := range mapinfo {
			if key != "_id" {
				filtermapinfo[key] = value
			}
		}

		// 데이터 반환
		c.JSON(http.StatusOK, responses.MapResponse_map{Code: 1, Message: filtermapinfo}) // 형변환
	}
}
