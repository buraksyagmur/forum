package forum

func FillActivity(usr user) Activity {
	po := displayPostsAndComments()
	var Likes string
	var Dlikes string
	var LikeCom string
	var DlikeCom string
	var act Activity
	for i := 0; i < len(po); i++ {
		Likes = Likes + "$" + CheckLikesAndDislikes(usr, po[i].PostID, "l")
		Dlikes = Dlikes + "$" + CheckLikesAndDislikes(usr, po[i].PostID, "d")
		LikeCom = LikeCom + "$" + CheckLikesAndDislikes(usr, po[i].PostID, "ComL")
		DlikeCom = DlikeCom + "$" + CheckLikesAndDislikes(usr, po[i].PostID, "ComD")
	}
	act.Likes = Likes
	act.Dislikes = Dlikes
	act.LikesCom = LikeCom
	act.DlikesCom = DlikeCom
	return act
}
