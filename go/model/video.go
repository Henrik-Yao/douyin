package model

type Video struct { //视频信息
	Id            int64    `json:"id,omitempty"`
	Author        UserInfo `json:"author"`
	PlayUrl       string   `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount int64    `json:"favorite_count,omitempty"`
	CommentCount  int64    `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
}
