package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	// "fmt"
	"log"
	"net/http"
	"time"

	// "GameServer/controller"
	"game_http_go/responses"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

/*
func GetCreatorList() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 표시되는 필드만 표현
		projection := bson.M{
			"_id":     0,
			"map_id":  1,
			"mapName": 1,
		}

		// filter 모두 존재하는 경우에만 출력
		filter := bson.M{
			"$and": []bson.M{
				{"map_id": bson.M{"$exists": true}},
				{"mapName": bson.M{"$exists": true}},
				{"chunkNum": 0},
			},
		}

		mapCollection := controller.DBClient.Collection("creators")
		// collection.Find(context:취소 시그널 및 타임아웃 전달, 빈맵 : 모든 문서 선택, setProjection(검색 옵션 설정))
		// cursor : 결과 집합의 다음 항목 가져올 수 있음
		//cursor, err := mapCollection.Find(ctx, bson.M{}, options.Find().SetProjection(projection)) // rejection이랑 filter 정확하게 알기
		cursor, err := mapCollection.Find(ctx, filter, options.Find().SetProjection(projection))
		if err != nil {

			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: "error"})
			log.Println("(list)FindErr :", err)
			return
		}

		fmt.Println(cursor)

		// cursor에서 반환된 모든 값을 가져와 map[string]interface{} 슬라이스로 변환
		var results []map[string]interface{} // 여러개라서 []map[string]interface{}
		if err = cursor.All(ctx, &results); err != nil {

			c.JSON(http.StatusInternalServerError, responses.MapResponse{Code: 0, Message: "error"})
			log.Println("(list)FindReturnErr :", err)
			return
		}

		// JSON으로 결과 반환

		c.JSON(http.StatusOK, responses.MapResponse_list{Code: 1, Message: results})
		log.Println("(list)Success")
	}
}*/

var LoginServerEndpoint = "http://127.0.0.1:8000"

func HttpGet(url string) []byte {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// GET 요청 생성
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Failed to create HTTP request: %s\n", err)
		return nil
	}

	// 필요 시 헤더 추가
	// req.Header.Set("Key", "Value")

	// 요청 보내기
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP request failed: %s\n", err)
		return nil
	}
	defer response.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return nil
	}

	// 응답 상태 코드 및 본문 출력
	// fmt.Printf("Response status code: %d\n", response.StatusCode)
	// fmt.Printf("Response body: %s\n", body)

	return body
}

func GetUserInfo(userid string) responses.CreatorListResponse {
	body := HttpGet(LoginServerEndpoint + "/info?userid=" + userid)
	var result responses.CreatorListResponse
	err := json.Unmarshal(body, &result)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response : %s\n", err)
		return responses.CreatorListResponse{}
	}
	return result
}

type TestType struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func GetTest() gin.HandlerFunc {
	return func(c *gin.Context) {

		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c.JSON(http.StatusOK, TestType{Code: 1, Message: "good"})
		log.Println("(HTTP) Test OK")
	}
}
