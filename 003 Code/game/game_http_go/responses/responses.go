package responses

type UserInfo struct {
	User_key string `bson:"user_key"`
	User_id  string `bson:"user_id"`
	Nickname string `bson:"nickname"`
}

type UserAuth struct {
	User           UserInfo `bson:"user"`
	Room_authority int      `bson:"room_authority"`
}

type CreatorListResponseMessage struct {
	Map_id       int        `bson:"map_id" json:"map_id"`
	Creator_list []UserAuth `bson:"creator_list" json:"creator_list"`
}

type CreatorListResponse struct {
	Code    int                        `bson:"code" json:"code"`
	Message CreatorListResponseMessage `bson:"message" json:"message"`
}

type DefaultReponse struct {
	Code    int    `bson:"code" json:"code"`
	Message string `bson:"message" json:"message"`
}
