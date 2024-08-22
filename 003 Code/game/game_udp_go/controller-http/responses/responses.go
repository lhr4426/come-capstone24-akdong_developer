package responses

type UserInfo struct {
	User_key	string `json:"user_key"`
	User_id		string `json:"user_id"`
	Nickname	string `json:"nickname"`
}

type UserAuth struct {
	User	UserInfo	`json:"user"`
	Room_authority	int	`json:"room_authority"`
	Map_id	int	`json:"map_id"`
}

type CreatorListResponseMessage struct {
	Map_id int `json:"map_id"`
	Creator_list []UserAuth `json:"creator_list"`
}

type CreatorListResponse struct {
	Code int 	`json:"code"`
	Message CreatorListResponseMessage `json:"message"`
}


