package controllers

import (
	"context"
	"fmt"
	"log"
	"map_http_go/responses"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// map_id값으로 timestamp 변경하기
func Get_time() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mapId, _ := strconv.ParseFloat(c.Query("mapId"), 64)
		fmt.Println(mapId)

		filter := bson.M{
			"map_id":   mapId,
			"chunkNum": 0,
		}

		projection := bson.M{
			"_id":      0,
			"mapCTime": 1,
		}

		cursor, err := mapCollection.Find(ctx, filter, options.Find().SetProjection(projection))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: "error"})
			return
		}

		fmt.Println(cursor)

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: "error"})
			return
		}

		// 결과가 없으면 에러 반환
		// if len(results) == 0 {
		// 	c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		// 	return
		// }

		// floattimestamp, ok := results[0]["mapCTime"].(float64)
		// if !ok {
		// 	fmt.Println("float64아님")
		// }

		// seconds := int64(floattimestamp)
		// fmt.Println(seconds)
		// nanoseconds := int64((floattimestamp - float64(seconds)) * 1e9)

		// timestamp := time.Unix(seconds, nanoseconds)
		// fmt.Println(timestamp)

		// JSON으로 변환된 시간 반환
		c.JSON(http.StatusOK, gin.H{"time": results})

	}

}
