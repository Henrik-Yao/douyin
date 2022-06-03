package model

type FavoriteAction struct{
	UserId int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
}

type FavoriteRequest struct{
	UserId int64 `json:"user_id"`
	Token string `json:"token"`
	VideoId int64 `json:"video_id"`
	ActionType int32 `json:"action_type"`
}

type FavoriteListRequest struct{
	UserId int64 `json:"user_id"`
	Token string `json:"token"`
}