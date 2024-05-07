package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// 데이터 구조에 맞게 정의한 구조체(JSON파일 내용을 Go언어로 표현)
type MapData struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Map_ID   string             `json:"map_id,omitempty" validate:"required"` // 여기서 int로 하면 안됨
	Version  string             `json:"version,omitempty" validate:"required"`
	Chunk    string             `json:"chunk,omitempty" validate:"required"`
	MapJson  string             `json:"mapjson,omitempty" validate:"required"`
	MapTotal string             `json:"maptotal"`
}

// type MapData struct {
// 	MapID      int
// 	UserID     string
// 	MapName    string
// 	MapCTime   float32
// 	Version    int
// 	MapSize    MapSize
// 	Tags       []string
// 	ChunkSize  int
// 	ChunkNum   int
// 	ObjCount   int
//  }

//  type MapSize struct {
// 	Horizontal float32
// 	Vertical   float32
// 	Height     float32
//  }

//  type MapObjectData struct {
// 	MapID    int
// 	Version  int
// 	ChunkNum int
// 	ObjList  []ObjectData
//  }

//  type ObjectData struct {
// 	ObjID          int
// 	AstID          int
// 	Transform      ObjTransform
// 	Type           string
// 	IsRigidbody    bool
// 	IsMeshCollider bool
//  }

//  type ObjTransform struct {
// 	Position Position
// 	Rotation Rotation
// 	Scale    Scale
//  }

//  type Position struct {
// 	X float32
// 	Y float32
// 	Z float32
//  }

//  type Rotation struct {
// 	X float32
// 	Y float32
// 	Z float32
//  }

//  type Scale struct {
// 	X float32
// 	Y float32
// 	Z float32
//  }
