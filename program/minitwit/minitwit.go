package minitwit

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	dto "github.com/lauralunddahl/DevOps-GroupF/program/dto"
	helper "github.com/lauralunddahl/DevOps-GroupF/program/helper"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"tawesoft.co.uk/go/dialog"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)

	layout      = "./templates/layout.html"
	tmp         = "./templates/tmp.html"
	login       = "./templates/login.html"
	register    = "./templates/register.html"
	noUserFound = "User not found"
	notAuth     = "Not authorized"
)

func BeforeRequest(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session1")
		session.Values["authenticated"] = false
		session.Values["userId"] = 0
		err := session.Save(r, w)
		if err != nil {
			println(err.Error())
		}

		handler(w, r)
	}
}

func PrivateTimeline(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session1")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/public", http.StatusFound)
	} else {
		userId := session.Values["userId"].(int)
		var timelines = dto.GetPrivateTimeline(userId)

		templ := template.Must(template.ParseFiles(layout, tmp))
		err := templ.Execute(w, map[string]interface{}{
			"timeline":  timelines,
			"public":    false,
			"type":      "default",
			"loggedin":  true,
			"sess_u_id": userId,
			"username":  dto.GetUsername(userId),
		})
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
}

func Loginpage(w http.ResponseWriter, r *http.Request) {
	loginp, err := template.ParseFiles(layout, login)
	if err != nil {
		println(err.Error())
	}
	err = loginp.Execute(w, nil)
	if err != nil {
		println(err.Error())
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("Loggin in with: " + username)
	var user = dto.GetUser(username)
	if user.Username == "" {
		fmt.Fprintln(w, "invalid username")
	}
	//check password hash from database against input password from user
	byteHash := []byte(user.PwHash)
	bytePw := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePw)
	if err != nil {
		println(err.Error())
		fmt.Fprintln(w, "invalid password")
	}

	session, _ := store.Get(r, "session1")
	session.Values["authenticated"] = true
	session.Values["userId"] = user.UserId
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusOK)
}

func userLoggedin(r *http.Request) bool {
	var userLoggedin = false
	session, _ := store.Get(r, "session1")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		userLoggedin = false
	} else {
		userId := session.Values["userId"].(int)
		if userId != 0 {
			userLoggedin = true
		}
	}
	return userLoggedin
}

func PublicTimeline(w http.ResponseWriter, r *http.Request) {
	var timelines = dto.GetPublicTimeline()
	templ := template.Must(template.ParseFiles(layout, tmp))

	err := templ.Execute(w, map[string]interface{}{
		"timeline": timelines,
		"public":   true,
		"type":     "public",
		"loggedin": userLoggedin(r),
	})
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func UserTimeline(w http.ResponseWriter, r *http.Request) {
	userId := 0
	vars := mux.Vars(r)

	username := vars["username"]
	if username == "metrics" {
		return
	}

	session, _ := store.Get(r, "session1")

	if auth, _ := session.Values["authenticated"].(bool); auth {
		userId = session.Values["userId"].(int)
	}
	profileuser := dto.GetUser(username)
	user := dto.GetUsername(userId)

	if profileuser.Username != "" {
		followed := dto.IsFollowing(userId, profileuser.UserId)

		var timelines = dto.GetUserTimeline(profileuser.UserId)

		templ := template.Must(template.ParseFiles(layout, tmp))

		err := templ.Execute(w, map[string]interface{}{
			"timeline":     timelines,
			"public":       false,
			"loggedin":     userLoggedin(r),
			"profileuser":  profileuser,
			"followed":     followed,
			"visiting":     true,
			"type":         "user",
			"sess_u_id":    userId,
			"loggedinuser": user == profileuser.Username,
		})
		if err != nil {
			fmt.Fprintln(w, err)
		}

	} else {
		http.Error(w, noUserFound, 404)
	}
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	userId := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		userId = session.Values["userId"].(int)
	}
	if userId == 0 {
		http.Error(w, notAuth, http.StatusUnauthorized)
	} else {
		whomId := dto.GetUserID(username)
		if whomId == 0 {
			http.Error(w, noUserFound, 404)
		}
		dto.FollowUser(userId, whomId)
		dialog.Alert("You are now following %s", username)
		http.Redirect(w, r, "/"+username, http.StatusOK)
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userId := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		userId = session.Values["userId"].(int)
	}
	if userId == 0 {
		http.Error(w, notAuth, http.StatusUnauthorized)
	} else {
		whomId := dto.GetUserID(username)
		if whomId == 0 {
			http.Error(w, noUserFound, 404)
		}
		dto.UnfollowUser(userId, whomId)
		dialog.Alert("You are no longer following %s", username)
		http.Redirect(w, r, "/"+username, http.StatusOK)
	}
}

func AddMessage(w http.ResponseWriter, r *http.Request) {
	userId := 0
	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		userId = session.Values["userId"].(int)
	}
	if userId == 0 {
		http.Error(w, notAuth, http.StatusUnauthorized)
	} else {
		text := r.FormValue("text")
		dto.AddMessage(strconv.Itoa(userId), text, time.Now(), 0)
		dialog.Alert("Your message was recorded")
		http.Redirect(w, r, "/", http.StatusOK)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	register, err := template.ParseFiles(layout, register)
	if err != nil {
		println(err.Error())
	}
	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		http.Redirect(w, r, "/", http.StatusOK)
	} else {
		err = register.Execute(w, nil)
		if err != nil {
			println(err.Error())
		}
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	err := ""
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password2 := r.FormValue("password2")

	if len(username) == 0 {
		err = "You have to enter a username\n"
	} else if len(email) == 0 || !strings.Contains(email, "@") {
		err += "You have to enter a valid email address\n"
	} else if len(password) == 0 {
		err += "You have to enter a password\n"
	} else if password != password2 {
		err += "The two passwords do not match\n"
	} else if dto.GetUserID(username) > 0 {
		err += "The username is already taken"
	} else {
		pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			println(err.Error())
		} else {
			image := helper.GravatarUrl(email)
			dto.RegisterUser(username, email, string(pwHash), image)
			fmt.Println(w, "You were successfully registered and can login now")
			http.Redirect(w, r, "/login", http.StatusOK)
		}
	}

	register, err2 := template.ParseFiles(layout, register)
	if err2 != nil {
		println(err2.Error())
	}

	e := register.Execute(w, map[string]interface{}{
		"error": err,
	})
	if e != nil {
		fmt.Fprintln(w, err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session1")
	session.Values["authenticated"] = false
	session.Values["userId"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusOK)
}
