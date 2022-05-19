package forum

import (
	"time"
)

type comment struct {
	CommentID      int
	Author         string
	PostID         int
	Content        string
	CommentTime    time.Time
	CommentTimeStr string
	Likes          int
	Dislikes       int
	LikedByCur     bool
	DislikedByCur  bool
}

type post struct {
	PostID        int
	Author        string // author
	Image         string
	Title         string
	Content       string
	Category      string
	PostTime      time.Time
	PostTimeStr   string
	Likes         int
	Dislikes      int
	Comments      []comment
	IPs           string
	View          int
	LikedByCur    bool
	DislikedByCur bool
}

type user struct {
	Username          string
	Email             string
	Access            int // 0 means no access, not logged in
	LoggedIn          bool
	Image             string
	Posts             []post
	Comments          []comment
	LikedPost         string
	DislikedPost      string
	LikedComments2    string
	DislikedComments2 string
	LikedComments     []comment
	Password          string
}
