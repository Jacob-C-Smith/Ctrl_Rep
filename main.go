// Package declaration
package main

// Imports
import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fit/data"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Structure definitions
type User struct {
	Username           string        `json:"user"`
	Password           string        `json:"pass"`
	ID                 int           `json:"id"`
	Day                int           `json:"day"`
	ActiveExercisePlan *ExercisePlan `json:"-"`
}

type ExercisePlanSchema struct {
	Name      string           `json:"name"`
	Exercises []ExerciseSchema `json:"exercises"`
}

type ExercisePlan struct {
	ExercisePlanSchema
	Date      string     `json:"date"`
	Exercises []Exercise `json:"exercises"`
	T         []string   `json:"timestamps"`
	Ti        int        `json:"-"`
	T1        int        `json:"-"`
}

type ExerciseSchema struct {
	Exercise string `json:"exercise"`
	Sets     int    `json:"sets"`
	Reps     int    `json:"reps"`
	Weight   int    `json:"weight"`
	T0       int    `json:"-"`
	T1       int    `json:"-"`
}

type Exercise struct {
	ExerciseSchema `json:"-"`
	Ts             int `json:"t"`
	Weight         int `json:"w"`
}

type SetRecord struct {
	Time   []string `json:"time"`
	Weight int      `json:"kg"`
	Reps   int      `json:"reps"`
}

type ActRecord struct {
	Name   string
	SetQty int
	Sets   []SetRecord
}

type DayRecord struct {
	Name string      `json:"name"`
	Date string      `json:"date"`
	Acts []ActRecord `json:"-"`
}

// Data
var PlanCatalog map[string]*ExercisePlanSchema = nil
var Sessions map[string]*User = nil
var users_qty int = 0
var Templates *template.Template = nil
var kv *data.KeyValueDB = nil

// Error check
func ok(err error) {

	// Error check
	if err != nil {

		// Error
		panic(err.Error())
	}
}

// Plan
func NewExercisePlanSchema(path string) (exercise_plan_schema *ExercisePlanSchema, err error) {

	// Initialized data
	var file []byte

	// Allocate
	exercise_plan_schema = new(ExercisePlanSchema)

	// Load
	file, err = os.ReadFile(path)
	ok(err)
	ok(json.Unmarshal(file, exercise_plan_schema))

	return
}

func NewExercisePlan(exercise_plan_schema *ExercisePlanSchema) (exercise_plan *ExercisePlan, err error) {

	// Initialized data
	var tA int = 0

	// Allocate
	exercise_plan = new(ExercisePlan)

	// Copy
	exercise_plan.ExercisePlanSchema = *exercise_plan_schema
	for _, v := range exercise_plan_schema.Exercises {
		var T0 int = tA
		var T1 int = T0 + v.Sets*2
		exercise_plan.Exercises = append(exercise_plan.Exercises, Exercise{
			Weight: 0,
			ExerciseSchema: ExerciseSchema{
				Exercise: v.Exercise,
				Sets:     v.Sets,
				Reps:     v.Reps,
				T0:       T0,
				T1:       T1,
			},
		})
		tA = T1
		tA = tA + 1

	}

	exercise_plan.T1 = tA
	exercise_plan.T = make([]string, tA)

	// Success
	return
}

// User
func (user *User) GetActiveExercise() *Exercise {
	var i int
	var a *ExercisePlan = user.ActiveExercisePlan
	for i = 0; i < len(a.Exercises); i++ {
		var ith Exercise = a.Exercises[i]
		if ith.T0 <= a.Ti && ith.T1 >= a.Ti {
			break
		}
		continue
	}
	return &user.ActiveExercisePlan.Exercises[i]
}

func (user *User) GetActiveExerciseIndex() int {
	var i int
	var a *ExercisePlan = user.ActiveExercisePlan
	for i = 0; i < len(a.Exercises); i++ {
		var ith Exercise = a.Exercises[i]
		if ith.T0 <= a.Ti && ith.T1 >= a.Ti {
			break
		}
		continue
	}
	return i
}

func (user *User) GetStateString() (result string) {
	// var tR int = user.GetActiveExercise().T1 - user.GetActiveExercise().T0
	var tS int = user.ActiveExercisePlan.Ti - user.GetActiveExercise().T0
	var sets int = user.GetActiveExercise().Sets
	var pad int = 0
	for tS >= 2 {
		result += "🟢"
		pad++
		tS = tS - 2
	}

	if tS == 1 {
		pad++
		result += "🟡"
	}

	for pad < sets {
		pad++
		result += "⚪️"
	}

	return
}

// Session
func Get(r *http.Request) (user *User, err error) {

	// Initialized data
	var cookie *http.Cookie = nil

	// Store the cookie from the request
	cookie, err = r.Cookie("session")

	// Error check
	if err != nil {
		return nil, err
	}

	// Store the session
	user = Sessions[cookie.Value]

	// Success
	return user, err
}

func Set(w http.ResponseWriter, username string) {

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    username,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})
}

// Routes
func landing_page(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var u *User

	// Get the user from the request
	u, _ = Get(r)

	// Error check
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Execute the landing page template
	ok(Templates.ExecuteTemplate(w, "main", &struct {
		User    *User
		Catalog *map[string]*ExercisePlanSchema
	}{
		User:    u,
		Catalog: &PlanCatalog,
	}))
}

func login_page(w http.ResponseWriter, r *http.Request) {

	// Execute the login page template
	ok(Templates.ExecuteTemplate(w, "login", nil))
}

func login_submit(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var username, password string
	var user *User = nil

	// Parse login form
	r.ParseForm()

	// Store the username and password
	username, password = r.Form.Get("username"), r.Form.Get("password")

	// Lookup the username
	user = Sessions[username]

	// Error check
	if user == nil {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	// Hash password
	h := sha1.New()
	h.Write([]byte(password))
	password = hex.EncodeToString(h.Sum(nil))

	// Check the password against the correct password
	if user.Password == password {

		// Set the user cookie
		Set(w, username)

		// All done
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	// All done
	http.Error(w, "Failed to authenticate!", http.StatusNotFound)
}

func signup_page(w http.ResponseWriter, r *http.Request) {

	// Execute the login page template
	ok(Templates.ExecuteTemplate(w, "signup", nil))
}

func signup_submit(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var username, password string

	// Error check
	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse sign up form
	r.ParseForm()

	// Store the username and password
	username, password = r.Form.Get("username"), r.Form.Get("password")

	// Hash password
	h := sha1.New()
	h.Write([]byte(password))
	hp := hex.EncodeToString(h.Sum(nil))

	kv.Set(fmt.Sprintf("fit:user:%d", users_qty), []byte(fmt.Sprintf("{\"user\":\"%s\",\"pass\": \"%s\",\"id\":%d,\"day\":0}", username, hp, users_qty)))
	users_qty = users_qty + 1
	kv.Set("fit:users", []byte(fmt.Sprintf("%d", users_qty)))
	kv.Write()

	// Construct a new user
	var newUser *User = new(User)
	newUser = &User{
		Username: username,
		Password: password,
	}

	// Add the new account to the session list
	Sessions[username] = newUser

	// Set the user for subsequent requests
	Set(w, username)

	// All done
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func card_page(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var err error = nil
	var user *User = nil
	var plan string

	// Get the user from the request
	user, _ = Get(r)

	// Error check
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Store the starting exercise
	plan = r.URL.Query().Get("start")

	// Construct a new exercise plan by replicating the schema
	user.ActiveExercisePlan, err = NewExercisePlan(PlanCatalog[plan])
	ok(err)

	// Initialize the counter
	user.ActiveExercisePlan.Ti = -1

	// Execute the card template
	ok(Templates.ExecuteTemplate(w, "card", user))
}

func card_advance(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var u *User = nil
	var active_exercise *Exercise = nil

	// Get the user from the request
	u, _ = Get(r)

	// Error check
	if u == nil {

		// Error
		return
	}

	// Parse advance
	r.ParseForm()
	var weight, reps, timestamp string = r.Form.Get("weight"), r.Form.Get("reps"), r.Form.Get("timestamp")

	fmt.Printf("w: %s, r: %s, t: %s\n", weight, reps, timestamp)

	// Increment the timestep
	u.ActiveExercisePlan.Ti++

	// Edge case ( user is done )
	if u.ActiveExercisePlan.Ti >= u.ActiveExercisePlan.T1 {
		ok(Templates.ExecuteTemplate(w, "done", active_exercise))
		return
	}

	u.ActiveExercisePlan.T[u.ActiveExercisePlan.Ti] = timestamp

	current := time.Now()

	d, _ := json.Marshal(DayRecord{
		Name: u.ActiveExercisePlan.Name,
		Date: current.Format("2006-01-02"),
	})

	// Set the day's exercise plan
	kv.Set(fmt.Sprintf("fit:user:%d:day:%d", u.ID, u.Day), d)

	//
	var tS int = u.ActiveExercisePlan.Ti - u.GetActiveExercise().T0

	// Set the specific act
	kv.Set(fmt.Sprintf("fit:user:%d:day:%d:act:%d", u.ID, u.Day, u.GetActiveExerciseIndex()), []byte(fmt.Sprintf("\"%s\"", u.GetActiveExercise().Exercise)))

	// Set completed
	if tS%2 == 0 && tS != 0 {
		var timest []string = u.ActiveExercisePlan.T
		var loadTs string = timest[u.GetActiveExercise().T0+tS-1]
		var restTs string = timest[u.GetActiveExercise().T0+tS]
		var lT time.Time
		var rT time.Time

		lT, _ = time.Parse("3:04:05 PM", loadTs)
		rT, _ = time.Parse("3:04:05 PM", restTs)
		var dT time.Duration = rT.Sub(lT)

		iweight, _ := strconv.Atoi(weight)
		ireps, _ := strconv.Atoi(reps)

		var sr SetRecord = SetRecord{
			Time:   []string{loadTs, restTs, dT.String()},
			Weight: iweight,
			Reps:   ireps,
		}
		srb, _ := json.Marshal(sr)
		aer := u.GetActiveExerciseIndex()
		kv.Set(fmt.Sprintf("fit:user:%d:day:%d:act:%d:set:%d", u.ID, u.Day, aer, (tS/2)-1), srb)
	}

	// Send the next card
	ok(Templates.ExecuteTemplate(w, "card_dynamic", &struct {
		State string
		Card  *Exercise
	}{
		State: u.GetStateString(),
		Card:  u.GetActiveExercise(),
	}))
}

func card_done(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var u *User = nil

	// Get the user from the request
	u, _ = Get(r)

	// Error check
	if u == nil {

		// Error
		return
	}
	//tN := time.Now()
	//u.ActiveExercisePlan.Date = tN.Format("2006-01-02")
	//day, _ := json.Marshal(u.ActiveExercisePlan)

	//kv.Set(fmt.Sprintf("fit:user:%d:day:%d", u.ID, u.Day), []byte(day))

	u.Day = u.Day + 1
	new_user, err := json.Marshal(u)
	ok(err)

	// Update
	kv.Set(fmt.Sprintf("fit:user:%d", u.ID), new_user)
	Sessions[u.Username] = u
	kv.Write()

	// All done
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func template_update() {
	// Create a ticker that ticks every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			// Initialized data
			var err error = nil

			// Parse the Templates
			Templates, err = template.ParseFiles("template/main", "template/login", "template/signup", "template/card", "template/card_dynamic", "template/done", "template/history")
			ok(err)

			// Parse the exercise plans
			ok(filepath.WalkDir("resources/plans/", func(path string, d fs.DirEntry, err error) error {

				// Initialized data
				var schema *ExercisePlanSchema = nil

				// Error check
				ok(err)

				// Skip directories
				if d.IsDir() {

					// Continue
					return nil
				}

				// Construct an exercise plan primary
				schema, err = NewExercisePlanSchema(path)
				ok(err)

				// Add the exercise to the list
				PlanCatalog[schema.Name] = schema

				// Success
				return nil
			}))
		}
	}
}

func history(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var user *User = nil

	// Get the user from the request
	user, _ = Get(r)

	// Error check
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	days := kv.GetAll(fmt.Sprintf("fit:user:%d:day", user.ID))
	var d []DayRecord = make([]DayRecord, user.Day)

	for i := user.Day - 1; i >= 0; i-- {
		var iD DayRecord

		ok(json.Unmarshal([]byte(days[i]), &iD))

		acts := kv.GetAll(fmt.Sprintf("fit:user:%d:day:%d:act", user.ID, i))

		var a []ActRecord = make([]ActRecord, len(acts))

		for j := 0; j < len(acts); j++ {
			var iA ActRecord
			iA.Name = strings.Split(acts[j], "\"")[1]

			sets := kv.GetAll(fmt.Sprintf("fit:user:%d:day:%d:act:%d:set", user.ID, i, j))
			var s []SetRecord = make([]SetRecord, len(sets))
			iA.SetQty = len(sets)
			for k := 0; k < len(sets); k++ {
				var iS SetRecord

				ok(json.Unmarshal([]byte(sets[k]), &iS))

				s[k] = iS
			}
			iA.Sets = s
			a[j] = iA
		}
		iD.Acts = a
		d[i] = iD
	}

	fmt.Printf("%+v\n", d)

	// Execute the history template
	ok(Templates.ExecuteTemplate(w, "history", &struct {
		User *User
		Days []DayRecord
	}{
		User: user,
		Days: d,
	}))
}

// Initialize
func init() {

	// Connect to the database
	kv = data.NewDB("data.g10.app:6764")

	// Get each user
	db_users := kv.GetAll("fit:user")

	// Log
	fmt.Printf("%+v\n", db_users)

	// Initialize Sessions
	Sessions = make(map[string]*User)

	// Construct a session foreach user
	for _, v := range db_users {

		// Allocate a user
		var ith_user *User = new(User)

		// Construct the user
		ok(json.Unmarshal([]byte(v), ith_user))

		// Store the user session
		Sessions[ith_user.Username] = ith_user
	}

	// Store the quantity of users
	users_qty = len(Sessions)

	// Initialize exercise plan catalog
	PlanCatalog = make(map[string]*ExercisePlanSchema)

	// Periodically reload templates
	go template_update()

	// Periodically reload exercise plans
	//

	// etc etc

	// Done
	return
}

// Entry point
func main() {

	// Static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Routes
	http.HandleFunc("/", landing_page)
	http.HandleFunc("/signup/", signup_page)
	http.HandleFunc("/signup/submit", signup_submit)
	http.HandleFunc("/login/", login_page)
	http.HandleFunc("/login/submit", login_submit)
	http.HandleFunc("/card/", card_page)
	http.HandleFunc("/card/done", card_done)
	http.HandleFunc("/card/advance", card_advance)
	http.HandleFunc("/history", history)

	// Log
	fmt.Println("Starting HTTPS server on https://fit.g10.app")

	// Host
	err := http.ListenAndServeTLS(":8082", "ca/g10.app.crt", "ca/g10.app.key", nil)

	//err := http.ListenAndServeTLS(":8082", "/root/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/g10.app/g10.app.crt", "/root/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/g10.app/g10.app.key", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
