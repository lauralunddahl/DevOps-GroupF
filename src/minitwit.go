package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	api "github.com/lauralunddahl/DevOps-GroupF/src/api"
	dto "github.com/lauralunddahl/DevOps-GroupF/src/dto"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"tawesoft.co.uk/go/dialog"
)

const database string = "../minitwit.db"
const per_page int = 30
const debug bool = true
const secret_key string = "development key"

var router = mux.NewRouter()

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type User struct {
	Username string
	UserId   string
	Email    string
	PwHash   string
}

type Message struct {
	MessageId string
	AuthorId  string
	Text      string
	PubDate   string
	Flagged   int
}

type Timeline struct {
	Username string
	UserId   int
	Email    string
	PwHash   string

	MessageId int
	AuthorId  int
	Text      string
	PubDate   time.Time
	Flagged   int
}

func format_datetime(timestamp string) string {
	const layout = "2016-03-28 @ 08:30"
	t, err := time.Parse(layout, timestamp)
	checkErr(err)
	return t.String()
}

func gravatar_url(email string) string {
	size := 80
	h := sha1.New()
	h.Write([]byte(strings.ToLower(strings.TrimSpace(email))))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon&s=%d", sha1_hash, size)
}

func before_request(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session1")
		session.Values["authenticated"] = false
		session.Values["userid"] = 0
		err := session.Save(r, w)
		if err != nil {
			println(err.Error())
		}

		handler(w, r)
	}
}

//Not sure if we need the after_request function

func timeline(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session1")

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, r, "/public", 302)
	} else {
		user_id := session.Values["userid"].(int)
		var timelines = dto.GetPrivateTimeline(user_id)
		println(user_id)

		templ := template.Must(template.ParseFiles("./templates/layout.html", "./templates/tmp.html"))

		err := templ.Execute(w, map[string]interface{}{
			"timeline":  timelines,
			"public":    false,
			"type":      "default",
			"sess_u_id": user_id,
			"username":  dto.GetUsername(user_id),
			"loggedin":  true,
		})
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}
}

func loginpage(w http.ResponseWriter, r *http.Request) {
	loginp, err := template.ParseFiles("./templates/layout.html", "./templates/login.html")
	if err != nil {
		println(err.Error())
	}
	err = loginp.Execute(w, nil)
	if err != nil {
		println(err.Error())
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	var user = dto.GetUser(username)
	if user.Username == "" {
		fmt.Fprintln(w, "invalid username")
	}
	println(user.Username)
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
	http.Redirect(w, r, "/", 302)
}

func public_timeline(w http.ResponseWriter, r *http.Request) {
	var timelines = dto.GetPublicTimeline()
	templ := template.Must(template.ParseFiles("./templates/layout.html", "./templates/tmp.html"))

	err := templ.Execute(w, map[string]interface{}{
		"timeline": timelines,
		"public":   true,
		"type":     "public",
	})
	if err != nil {
		fmt.Fprintln(w, err)
	}
}

func user_timeline(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")

	//session.Values["authenticated"] = false
	//println(session.Values["authenticated"].(bool))
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	profileuser := dto.GetUser(username)

	if profileuser.Username != "" {
		followed := dto.IsFollowing(user_id, profileuser.UserId)

		var timelines = dto.GetUserTimeline(profileuser.UserId)

		templ := template.Must(template.ParseFiles("./templates/tmp.html", "./templates/layout.html"))

		err := templ.Execute(w, map[string]interface{}{
			"timeline":    timelines,
			"public":      false,
			"profileuser": profileuser,
			"followed":    followed,
			"type":        "user",
			"sess_u_id":   user_id,
		})
		if err != nil {
			fmt.Fprintln(w, err)
		}

	} else {
		http.NotFound(w, r)
	}
}

func follow_user(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	if user_id == 0 {
		http.Error(w, "not authorized", 401)
	} else {
		whom_id := dto.GetUserID(username)
		if whom_id == 0 {
			http.NotFound(w, r)
		}
		dto.FollowUser(user_id, whom_id)
		dialog.Alert("You are now following %s", username)
		http.Redirect(w, r, "/{username}", 302)
	}
}

func unfollow_user(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	vars := mux.Vars(r)

	username := vars["username"]

	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	if user_id == 0 {
		http.Error(w, "not authorized", 401)
	} else {
		whom_id := dto.GetUserID(username)
		if whom_id == 0 {
			http.NotFound(w, r)
		}
		dto.UnfollowUser(user_id, whom_id)
		dialog.Alert("You are no longer following %s", username)
		http.Redirect(w, r, "/{username}", 302)
	}
}

func add_message(w http.ResponseWriter, r *http.Request) {
	user_id := 0
	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		user_id = session.Values["userid"].(int)
	}
	if user_id == 0 {
		http.Error(w, "not authorized", 401)
	} else {
		text := r.FormValue("text")
		dto.AddMessage(strconv.Itoa(user_id), text, time.Now(), 0)
		dialog.Alert("Your message was recorded")
		http.Redirect(w, r, "/", 302)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	register, err := template.ParseFiles("./templates/layout.html", "./templates/register.html")
	if err != nil {
		println(err.Error())
	}
	session, _ := store.Get(r, "session1")
	if auth, _ := session.Values["authenticated"].(bool); auth {
		http.Redirect(w, r, "/", 302)
	} else {
		err = register.Execute(w, nil)
		if err != nil {
			println(err.Error())
		}
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	println("handle register")
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
	} else if dto.GetUserID(username) > 0 { //this might have to be another check at some point
		err += "The username is already taken"
	} else {
		pw_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			println(err.Error())
		}
		dto.RegisterUser(username, email, string(pw_hash))
		fmt.Println(w, "You were successfully registered and can login now")
		http.Redirect(w, r, "/login", 302)
	}

	register, err2 := template.ParseFiles("./templates/layout.html", "./templates/register.html")
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

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session1")
	session.Values["authenticated"] = false
	session.Values["userId"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}

func printSlice(s []Timeline) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	router.HandleFunc("/", timeline)
	router.HandleFunc("/register", before_request(register)).Methods("GET")
	router.HandleFunc("/registerfunc", handleRegister).Methods("POST")
	router.HandleFunc("/", before_request(timeline))
	router.HandleFunc("/login", loginpage)
	router.HandleFunc("/loginfunc", handleLogin).Methods("POST")
	router.HandleFunc("/public", public_timeline)
	api.HandleApiRequest(router)
	router.HandleFunc("/{username}", user_timeline)

	router.HandleFunc("/{username}/follow", follow_user)
	router.HandleFunc("/{username}/unfollow", unfollow_user)
	router.HandleFunc("/add_message", add_message).Methods("POST")
	router.HandleFunc("/logout", logout)

	log.Fatal(http.ListenAndServe(":8080", router))
}
