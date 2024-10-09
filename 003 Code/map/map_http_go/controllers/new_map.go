package controllers

import (
	"context"
	"fmt"
	"log"
	"map_http_go/responses"
	"net"
	"net/http"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mapCTime 전송하기 (timestamp)
// 새 맵 만들기
func NewMap() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		existMapidList, err := mapCollection.Distinct(ctx, "map_id", bson.M{})
		if (err != nil) && (err != mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: err.Error()})
			log.Fatal()
		}

		for _, value := range existMapidList {
			fmt.Println(value)
		}

		creatorId := c.Query("creator_id")
		mapName := c.Query("map_name")

		var newMapId int

		distinctFlag := 0

		if len(existMapidList) != 0 {
			for distinctFlag == 0 {
				newMapId = rand.Intn(90000) + 10000
				for _, existMapid := range existMapidList {
					if existMapid == newMapId {
						break
					}
				}
				distinctFlag = 1
			}
		}

		newMapMetaData := bson.M{
			"map_id":   newMapId,
			"user_id":  creatorId,
			"mapName":  mapName,
			"mapCTime": time.Now().UTC().Format(time.RFC3339Nano),
			"version":  0,
			"mapSize": bson.M{
				"horizontal": 100.0,
				"vertical":   100.0,
				"height":     40.0,
			},
			"Tags":      bson.A{},
			"chunkSize": 20,
			"chunkNum":  0,
			"objCount":  1,
		}

		newMapChunkData := bson.M{
			"map_id":   newMapId,
			"version":  0,
			"chunkNum": 1,
			"objList": []bson.M{
				{
					"obj_id": 0,
					"ast_id": "",
					"transform": bson.M{
						"position": bson.M{"x": 0.0, "y": 0.0, "z": 0.0},
						"rotation": bson.M{"x": 0.0, "y": 0.0, "z": 0.0},
						"scale":    bson.M{"x": 0.0, "y": 0.0, "z": 0.0},
					},

					"type":           "Null",
					"isRigidbody":    false,
					"isMeshCollider": false},
			},
		}

		newMapDocs := []interface{}{
			newMapMetaData,
			newMapChunkData,
		}

		_, err = mapCollection.InsertMany(context.TODO(), newMapDocs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: err.Error()})
			log.Fatal()
		}

		responseData := make(map[string]interface{})
		responseData["map_id"] = newMapId

		c.JSON(http.StatusOK, responses.MapResponse_map{Code: 1, Message: responseData})

		udpAddr, err := net.ResolveUDPAddr("udp", "localhost:8050") // 게임 UDP서버 연결
		if err != nil {
			fmt.Println("UDP resolve error:", err)
			return
		}

		conn, err := net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			fmt.Println("UDP connection error:", err)
		}
		defer conn.Close()

		message := fmt.Sprintf("CreateNewMap$%v$%v$%v", creatorId, time.Now().UTC().Format(time.RFC3339Nano), newMapId)
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Send UDP Error :", err)
			return
		}

		fmt.Println("Send UDP Successed")

	}
}
