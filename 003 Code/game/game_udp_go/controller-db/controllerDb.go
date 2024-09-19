package controllerdb

import (
	"GameServer/controller"
	"context"
	"fmt"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() {
	// 접속할 MongoDB 주소 설정
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")

	// MongoDB 연결
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("mongo connect err :", err)
	}

	// 연결 확인
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("mongo ping err :", err)
	}

	fmt.Println("Connected to MongoDB!")

	controller.DBClient = client.Database("GameServer")
}

func CreatorInit() {
	creatorCollection := controller.DBClient.Collection("creators")

	cursor, err := creatorCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	// 결과 출력
	for cursor.Next(context.TODO()) {
		var result controller.CreatorLists
		if err := cursor.Decode(&result); err != nil {

			log.Fatal(err)
		}
		fmt.Println(result)
		controller.MapidCreatorList[strconv.Itoa(result.Map_id)] = result.Creator_list
	}

	fmt.Printf("Creator List Loaded\n")
}
