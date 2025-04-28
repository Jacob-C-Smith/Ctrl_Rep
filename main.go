// Package declaration
package main

// Imports
import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Structure definitions
type User struct {
	Username           string `json:"username"`
	Password           string `json:"plaintext password"`
	ActiveExercisePlan *ExercisePlan
}

type ExerciseTimestamp struct {
	Real    int `json:"real"`
	Logical int `json:"logical"`
}

type ExercisePlanSchema struct {
	Name      string           `json:"name"`
	Exercises []ExerciseSchema `json:"exercises"`
}

type ExercisePlan struct {
	ExercisePlanSchema
	Exercises []Exercise          `json:"exercises"`
	T         []ExerciseTimestamp `json:"timestamps"`
	Ti        int
}

type ExerciseSchema struct {
	Exercise string `json:"exercise"`
	Sets     int    `json:"sets"`
	Reps     int    `json:"reps"`
}

type Exercise struct {
	ExerciseSchema
	Ts     *ExerciseTimestamp
	weight int
}

// Data
var PlanCatalog map[string]*ExercisePlanSchema = nil
var Sessions map[string]*User = nil
var Templates *template.Template = nil

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

	// Allocate
	exercise_plan = new(ExercisePlan)

	// Copy
	exercise_plan.ExercisePlanSchema = *exercise_plan_schema
	for _, v := range exercise_plan_schema.Exercises {
		exercise_plan.Exercises = append(exercise_plan.Exercises, Exercise{
			weight: 0,
			ExerciseSchema: ExerciseSchema{
				Exercise: v.Exercise,
				Sets:     v.Sets,
				Reps:     v.Reps,
			},
		})
	}

	// Success
	return
}

// Exercise
//

// User
func (user *User) GetActiveExercise() *Exercise {
	return &user.ActiveExercisePlan.Exercises[user.ActiveExercisePlan.Ti]
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
		http.Error(w, "Not Found: User not found", http.StatusNotFound)
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

	// Initialized data
	var err error
	var u *User

	// Get the user from the request
	u, err = Get(r)

	// Error check
	if err != nil {
		Set(w, "anon")
		u, _ = Get(r)
	}

	// Execute the login page template
	ok(Templates.ExecuteTemplate(w, "login", &struct {
		User    *User
		Catalog *map[string]*ExercisePlanSchema
	}{
		User:    u,
		Catalog: &PlanCatalog,
	}))
}

func login_submit(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var username, password string
	var user *User = nil

	// Error check
	if r.Method != "POST" {
		return
	}

	// Parse login form
	r.ParseForm()

	// Store the username and password
	username, password = r.Form.Get("username"), r.Form.Get("password")

	// Lookup the username
	user = Sessions[username]

	// Error check
	if user == nil {
		fmt.Fprintf(w, "who?\n")
	}

	// Check the password against the correct password
	if user.Password == password {

		// Set the user cookie
		Set(w, username)

		// Welcome, [username]
		fmt.Fprintf(w, "welcome, %s\n", username)
	}

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
		fmt.Fprintf(w, "Who?")
	}

	// Store the starting exercise
	plan = r.URL.Query().Get("start")

	// Construct a new exercise plan by replicating the schema
	user.ActiveExercisePlan, err = NewExercisePlan(PlanCatalog[plan])
	ok(err)

	// Execute the card template
	ok(Templates.ExecuteTemplate(w, "card", &struct {
		User *User
		Card *Exercise
	}{
		User: user,
		Card: user.GetActiveExercise(),
	}))
}

func card_advance(w http.ResponseWriter, r *http.Request) {

	// Initialized data
	var u *User = nil
	var active_exercise *Exercise = nil
	// var weight, reps, timestamp string = r.Form.Get("weight"), r.Form.Get("reps"), r.Form.Get("timestamp")

	// Get the user from the request
	u, _ = Get(r)

	// Error check
	if u == nil {

		// Error
		return
	}

	// Store the active exercise
	active_exercise = u.GetActiveExercise()

	// Increment the timestep
	u.ActiveExercisePlan.Ti++

	// Edge case ( user is done )
	if u.ActiveExercisePlan.Ti >= len(u.ActiveExercisePlan.Exercises) {
		ok(Templates.ExecuteTemplate(w, "done", active_exercise))
		return
	}

	// Send the next card
	ok(Templates.ExecuteTemplate(w, "card_dynamic", active_exercise))
}

// Initialize
func init() {

	// Initialized data
	var err error = nil

	// Initialize Sessions
	Sessions = make(map[string]*User)

	Sessions["jake"] = &User{
		Username: "jake",
		Password: "j",
	}

	Sessions["alice"] = &User{
		Username: "alice",
		Password: "a",
	}

	Sessions["bob"] = &User{
		Username: "bob",
		Password: "b",
	}

	Sessions["charlie"] = &User{
		Username: "charlie",
		Password: "c",
	}

	// Initialize exercise plan catalog
	PlanCatalog = make(map[string]*ExercisePlanSchema)

	// Parse the Templates
	Templates, err = template.ParseFiles("template/main", "template/login", "template/card", "template/card_dynamic", "template/done")
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

	// Done
	return
}

// Entry point
func main() {

	// Static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Routes
	http.HandleFunc("/", landing_page)
	http.HandleFunc("/login/", login_page)
	http.HandleFunc("/login/submit", login_submit)
	http.HandleFunc("/card/", card_page)
	http.HandleFunc("/card/advance", card_advance)

	// Host
	http.ListenAndServe(":8080", nil)
}
