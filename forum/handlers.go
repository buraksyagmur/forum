package forum

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type mainPageData struct {
	Userinfo    user
	Posts       []post
	ForumUnames []string
}

var (
	urlPost     string
	duplicateIP bool
	changingPos bool
	changingCom bool
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Status Not Found", 404)
		return
	}

	changingPos = false
	fmt.Println("looking for main page method", r.Method)
	curUser := obtainCurUserFormCookie(r)
	if curUser.Username != "" {
		users := AllForumUsers()
		for i := 0; i < len(users); i++ {
			if users[i].Username == curUser.Username {
				curUser.LikedPost = users[i].LikedPost
				curUser.DislikedPost = users[i].DislikedPost
				curUser.DislikedComments2 = users[i].DislikedComments2
				curUser.LikedComments2 = users[i].LikedComments2

			}
		}
		changingPos = true

	}

	// // test
	// var whichUser string
	// var logInOrNot bool
	// rows, err := db.Query("SELECT username, loggedIn FROM users WHERE username = ?;", curUser.Username)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	rows.Scan(&whichUser, &logInOrNot)
	// }

	// fmt.Printf("HomeHandler:: login user: %s, login status: %v\n", whichUser, logInOrNot)
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/header2.gohtml", "./templates/footer.gohtml", "./templates/index.gohtml", "./templates/index2.gohtml")
		// tpl, err := template.ParseFiles("./templates/index.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
			return
		}

		filCat := r.FormValue("category-filter")
		filAuthor := r.FormValue("author-filter")
		filLiked := r.FormValue("liked-post")
		filCatFromButton := r.FormValue("categoryOfPost")
		fmt.Println("************************ ", filCatFromButton, "catname")
		fmt.Println(filCatDisplayPostsAndComments(filCatFromButton))
		var pos []post
		if filCat != "" {
			pos = filCatDisplayPostsAndComments(filCat)
		} else if filAuthor != "" {
			pos = filAuthorDisplayPostsAndComments(filAuthor)
		} else if filLiked != "" {
			pos = filLikedDisplayPostsAndComments(curUser)
		} else if filCatFromButton != "" {
			pos = filCatDisplayPostsAndComments(filCatFromButton)
		} else {
			pos = displayPostsAndComments()
		}
		AllLikes, AllDislikes := SumOfAllLikes(AllForumUsers())
		pos = DistLikesToPosts(pos, AllLikes, AllDislikes)
		for i := 0; i < len(pos); i++ {
		}
		if changingPos {
			userLikes := CountLikesByUser(curUser, "l")
			userDislikes := CountLikesByUser(curUser, "d")

			for i := 0; i < len(pos); i++ {
				for k := 0; k < len(userLikes); k++ {
					if pos[i].PostID == userLikes[k] {
						pos[i].LikedByCur = true
					}
				}
			}
			for i := 0; i < len(pos); i++ {
				for k := 0; k < len(userDislikes); k++ {
					if pos[i].PostID == userDislikes[k] {
						pos[i].DislikedByCur = true
					}
				}
			}
		}
		allForumUnames := allForumUnames()
		data := mainPageData{
			Posts:       pos,
			Userinfo:    curUser,
			ForumUnames: allForumUnames,
		}

		// fmt.Println("---------", forumUser)
		if changingPos {
			err = tpl.ExecuteTemplate(w, "index2.gohtml", data)
			if err != nil {
				http.Error(w, "Executing Error", http.StatusInternalServerError)
				return
			}
		} else {
			err = tpl.ExecuteTemplate(w, "index.gohtml", data)
			if err != nil {
				http.Error(w, "Executing Error", http.StatusInternalServerError)
				return
			}
		}

	}
	if r.Method == http.MethodPost {
		processPost(r, curUser)
		processComment(r, curUser)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("logged in", loggedIn(r))
	if loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/login.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
			return
		}
		err = tpl.ExecuteTemplate(w, "login.gohtml", nil)
		if err != nil {
			http.Error(w, "Executing Error", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodPost {
		processLogin(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if loggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/register.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
			return
		}
		err = tpl.ExecuteTemplate(w, "register.gohtml", nil)
		if err != nil {
			http.Error(w, "Executing Error", http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodPost {
		regNewUser(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//	if loggedIn(r) {
	processLogout(w, r)
	//	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	changingCom = false
	curUser := obtainCurUserFormCookie(r)
	if curUser.Username != "" {
		users := AllForumUsers()
		for i := 0; i < len(users); i++ {
			if users[i].Username == curUser.Username {
				curUser.LikedPost = users[i].LikedPost
				curUser.DislikedPost = users[i].DislikedPost
				curUser.DislikedComments2 = users[i].DislikedComments2
				curUser.LikedComments2 = users[i].LikedComments2

			}
		}
		changingCom = true

	}
	if r.Method == "GET" {
		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/header2.gohtml", "./templates/post.gohtml", "./templates/post2.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
			return
		}
		strID := r.FormValue("postdetails")
		PostIdFromHTML, err := strconv.Atoi(strID)
		if err != nil {
			os.Exit(0)
		}
		// fmt.Println(PostIdFromHTML, "---------")
		var pos []post
		pos = displayPostsAndComments()

		AllLikes, AllDislikes := SumOfAllLikes(AllForumUsers())
		fmt.Println("ALLLIKES", AllLikes, "ALLDISLIKES", AllDislikes)
		pos = DistLikesToPosts(pos, AllLikes, AllDislikes)
		allForumUnames := allForumUnames()
		var Chosen []post
		for i := 0; i < len(pos); i++ {
			if pos[i].PostID == PostIdFromHTML {
				Chosen = append(Chosen, pos[i])
			}
		}
		//********* IP ********
		duplicateIP = false
		if Chosen[0].IPs == "" {
			Chosen[0].IPs = GetOutboundIP().String()
			duplicateIP = true
		}
		if Chosen[0].IPs == GetOutboundIP().String() {
			duplicateIP = true
		}

		if !duplicateIP {
			Chosen[0].IPs += "-" + GetOutboundIP().String()
		}
		allIp := (strings.Split(Chosen[0].IPs, "-"))
		keys := make(map[string]bool)
		list := []string{}
		for _, entry := range allIp {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}
		Chosen[0].View = len(list)
		stmt, err := db.Prepare("UPDATE posts SET ips = ?	WHERE postID = ?;")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(Chosen[0].IPs, Chosen[0].PostID)
		//********* IP ********
		urlPost = "postpage?postdetails=" + strID + "&postdetails=" + Chosen[0].Title
		Alllikes, Alldislikes := CommentSumOfAllLikes(AllForumUsers())

		Chosen[0].Comments = DistLikesToComments(Chosen[0].Comments, Alllikes, Alldislikes)
		if changingCom {
			userLikes := CountLikesByUser(curUser, "l")
			userDislikes := CountLikesByUser(curUser, "d")
			userComLikes := CommentCountLikesByUser(curUser, "l")
			userComDislikes := CommentCountLikesByUser(curUser, "d")
			for k := 0; k < len(userLikes); k++ {
				if Chosen[0].PostID == userLikes[k] {
					Chosen[0].LikedByCur = true
				}
			}
			for k := 0; k < len(userDislikes); k++ {
				if Chosen[0].PostID == userDislikes[k] {
					Chosen[0].DislikedByCur = true
				}
			}

			for i := 0; i < len(Chosen[0].Comments); i++ {
				for k := 0; k < len(userComLikes); k++ {
					if Chosen[0].Comments[i].CommentID == userComLikes[k] {
						Chosen[0].Comments[i].LikedByCur = true
					}
				}
			}
			for i := 0; i < len(Chosen[0].Comments); i++ {
				for k := 0; k < len(userComDislikes); k++ {
					if Chosen[0].Comments[i].CommentID == userComDislikes[k] {
						Chosen[0].Comments[i].DislikedByCur = true
					}
				}
			}
		}
		data := mainPageData{
			Posts:       Chosen,
			Userinfo:    curUser,
			ForumUnames: allForumUnames,
		}
		if changingCom {
			err = tpl.ExecuteTemplate(w, "post2.gohtml", data)
		} else {
			err = tpl.ExecuteTemplate(w, "post.gohtml", data)
		}

		if err != nil {
			http.Error(w, "Executing Error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {

		processPost(r, curUser)
		processComment(r, curUser)
		http.Redirect(w, r, urlPost, http.StatusSeeOther)
	}
}

func CategoryPageHandler(w http.ResponseWriter, r *http.Request) {
	curUser := obtainCurUserFormCookie(r)
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/header2.gohtml", "./templates/footer.gohtml", "./templates/categories.gohtml")
		if err != nil {
			http.Error(w, "Parsing Error", http.StatusInternalServerError)
			return
		}

		var pos []post
		category := r.FormValue("categoryAllPosts")
		pos = filCatDisplayPostsAndComments(category)

		allForumUnames := allForumUnames()
		data := mainPageData{
			Posts:       pos,
			Userinfo:    curUser,
			ForumUnames: allForumUnames,
		}
		// fmt.Println("---------", forumUser)
		err = tpl.ExecuteTemplate(w, "categories.gohtml", data)
		if err != nil {
			http.Error(w, "Executing Error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

// func DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	// for testing purpose
// 	if r.Method == http.MethodGet {
// 		tpl, err := template.ParseFiles("./templates/delete.gohtml", "./templates/footer.gohtml", "./templates/header.gohtml")
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		tpl.ExecuteTemplate(w, "delete.gohtml", nil)
// 	}
// 	if r.Method == http.MethodPost {
// 		deleteUser(r)
// 	}
// }

// func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
// 	tpl, err := template.ParseFiles("./templates/header.gohtml", "./templates/footer.gohtml", "./templates/notFound.gohtml")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tpl.ExecuteTemplate(w, "notFound.gohtml", nil)
// }
