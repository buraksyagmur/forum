package forum

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func findAuthor(posID int) (string, user) {
	var SelectedUser user
	var authorName string
	var msg string

	posIDstr := strconv.Itoa(posID)
	po := displayPostsAndComments()
	usr := AllForumUsers()
	for i := 0; i < len(po); i++ {
		if po[i].PostID == posID {
			authorName = po[i].Author
			msg = "#" + "localhost:8080/postpage?postdetails=" + posIDstr + "&postdetails=" + po[i].Title + "#"
		}
	}
	for i := 0; i < len(usr); i++ {
		if usr[i].Username == authorName {
			SelectedUser = usr[i]
		}
	}

	return msg, SelectedUser
}

func showNotifications(usr user) user {
	msg := usr.notif.message
	view := usr.notif.view
	msg2 := strings.Split(msg, "#")
	view2 := strings.Split(view, "#")
	for i := 0; i < len(msg2); i++ {
		for k := 0; k < len(view2); k++ {
			if msg2[i] == view2[k] {
				msg2[i-2] = ""
				msg2[i-1] = ""
			}
		}
		usr.notif.message = strings.Join(msg2, "#")
	}
	return usr
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
