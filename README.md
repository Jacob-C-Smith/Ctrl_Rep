# https://fit.g10.app

Plan and track workouts with minimal interaction

See [PLANNING.md](PLANNING.md)

>  [Data](#data)
>> [Initial](#initial)
>>
>> [Users](#users)
>>
>> [Days](#days)
>>
>> [Acts](#acts)
>>
>> [Sets](#sets)
>>
>  [Routes](#routes)
>> [Static](#static)
>> 
>> [Dynamic](#dynamic)

## Data
Data is organized in my [Key/Value database](https://github.com/Jacob-C-Smith/data). Keys are **strings**, and values are **JSON values**. The database is persistent and networked. I've organized data by carefully structuring keys. I'll expound on this structuring later. The structure of the keys determines the type of data encoded in the value. I encode 4 types of data; *Users*, *Days*, *Acts*, and *Sets*. 

The database maintains two metadata keys. The **"fit"** key indicates that this database is being used for this application. The **"fit:users"** is an atomic cardinal number tracking the quantity of users that have registered. When a user registers, they are assigned the user ID of this key. The key is then incremented for the next (new) user. 

### Commands
The **"get"** command gets you a value from a key. The **"set"** keyword updates a property. 

```go
// Key value database package
package data

type KeyValueDB struct {
    ...
}

// Constructor
func NewDB(hostPort string) (key_value_db *KeyValueDB)

// Get
func (kv *KeyValueDB) Get(key string) (results []string)

// Get all
func (kv *KeyValueDB) GetAll(key string) (result []string)

// Set
func (kv *KeyValueDB) Set(key string, value []byte) (result string)

// Write (the database to the disk)
func (kv *KeyValueDB) Write() (result string)

// Close (the connection cleanly)
func (kv *KeyValueDB) Close()
```

### Initial
When the service is launched for the first time, the database looks like this. 

```
 ✔  $ list
fit:users                        : 0
fit                              : true

Connection from 45.32.223.17:46902
```

### Users

[Image](img/signup.png)

When a user registers an account, the database is updated.
```
 ✔  $ list
fit                              : true
fit:users                        : 4
fit:user:0                       : {"user":"jake","pass":"..........34b9a31214c9651e2f3f..........","id":0,"day":0}
fit:user:1                       : {"user":"alice","pass":"522b276a356bdf39013dfabea2cd43e141ecc9e8","id":1,"day":0}
fit:user:2                       : {"user":"bob","pass":"48181acd22b3edaebc8a447868a7df7ce629920a","id":2,"day":0}
fit:user:3                       : {"user":"charlie","pass":"d8cd10b920dcbdb5163ca0185e402357bc27c265","id":3,"day":0}
```

### Days

[Image](img/days.png)

When a user begins an exercise plan, the database records which plan was selected.
```
fit                              : true
fit:users                        : 4
fit:user:0                       : {"user":"jake","pass":"..........34b9a31214c9651e2f3f..........","id":0,"day":0}
fit:user:0:day:0                 : {"name":"Chest and Triceps","date":"2025-05-05"}
fit:user:0:day:1                 : {"name":"Back and Bicep","date":"2025-05-06"}
fit:user:1                       : {"user":"alice","pass":"522b276a356bdf39013dfabea2cd43e141ecc9e8","id":1,"day":0}
...
```

### Acts and sets
Specific details and timestamps are encoded for each set, of each act.
Some data has been omitted / augmented. See [Full](full.md) for a dump from a test run

```
fit                              : true
fit:users                        : 4
fit:user:0                       : {"user":"jake","pass":"..........34b9a31214c9651e2f3f..........","id":0,"day":0}
fit:user:0:day:0                 : {"name":"Chest and Triceps","date":"2025-05-05"}
fit:user:0:day:0:act:0           : "Barbell Bench Press"
fit:user:0:day:0:act:0:set:0     : {"time":["9:42:03 PM","9:42:44 PM","41s"],"kg":120,"reps":6}
fit:user:0:day:0:act:0:set:1     : {"time":["9:43:09 PM","9:44:08 PM","59s"],"kg":130,"reps":6}
fit:user:0:day:0:act:0:set:2     : {"time":["9:45:17 PM","9:46:18 PM","1m1s"],"kg":150,"reps":6}
fit:user:0:day:0:act:1           : "Incline Dumbbell Press"
fit:user:0:day:0:act:2           : "Chest Dips"
fit:user:0:day:0:act:3           : "Cable Chest Flyes"
fit:user:0:day:0:act:4           : "Tricep Rope Pushdowns"
fit:user:0:day:0:act:5           : "Overhead Dumbbell Tricep Extensions"
fit:user:0:day:0:act:6           : "Close Grip Barbell Bench Press"
fit:user:0:day:0:act:7           : "Skull Crushers"
fit:user:1                       : {"user":"alice","pass":"522b276a356bdf39013dfabea2cd43e141ecc9e8","id":1,"day":0}
...
```

## Routes
### Static route
```
static/
└── style.css
```
### Dynamic routes
```
/                      : GET
├── /signup            : GET
│   └── /signup/submit : POST
├── /login             : GET
│   └── /login/submit  : POST
├── /card              : GET
│   ├── /card/done     : POST
│   └── /card/advance  : POST
└── /history           : GET
```

The signup routes are used for registering accounts. The login route is for existing users. The card routes are the most interesting. GET()ing a card adds context to the users session, and presents the exercise page. The user enters information about their exercise into the form, while front end javascript collects timestamps between form submissions. Form data is POST()ed to the advance route. The advance route responds with a new card for the end user. This process repeats until the exercise is done. The final POST() to done finalizes the data entry and updates the state of the application. 


## Structural 
I've made a cut down version of the code to save you time. Many omissions have been made for brevity, however the full source code can be found in [main.go](main.go)

main.go
```go
// Package declaration
package main

// Imports
import (
    ...
)

// Structure definitions
type User struct { ... }

type ExercisePlanSchema struct { ... }

type ExercisePlan struct { ... }

type ExerciseSchema struct { ... }

type Exercise struct { ... }

type SetRecord struct { ... }

type ActRecord struct { ... }

type DayRecord struct { ... }

// Data
//

// Error check
func ok(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Plan
func NewExercisePlanSchema(path string) (exercise_plan_schema *ExercisePlanSchema, err error) {

	// Allocate
    //
      
	// Load
    //

	return
}

func NewExercisePlan(exercise_plan_schema *ExercisePlanSchema) (exercise_plan *ExercisePlan, err error) {


	// Allocate
    //

	// Copy from schema
    //    

	return
}

// User
func (user *User) GetActiveExercise() *Exercise {
    // This function is very ugly
}

func (user *User) GetActiveExerciseIndex() int {
    // This function is very ugly
}

func (user *User) GetStateString() (result string) {
    // This function is very ugly
}

// Session
func Get(r *http.Request) (user *User, err error) {

	// Initialized data
	var cookie *http.Cookie = nil

	// Store the cookie from the request
	cookie, err = r.Cookie("session")

	// Store the session
	user = Sessions[cookie.Value]

	// Success
	return user, err
}

func Set(w http.ResponseWriter, username string) {

	// Set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    username, // Could do better. Will do better. 
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	})
}

// Routes
func landing_page(w http.ResponseWriter, r *http.Request) {

	//
    // Respond
}

func login_page(w http.ResponseWriter, r *http.Request) {

	//
    // Respond
}

func login_submit(w http.ResponseWriter, r *http.Request) {

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

        // Error
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	// Error
	http.Error(w, "Failed to authenticate!", http.StatusNotFound)
}

func signup_page(w http.ResponseWriter, r *http.Request) {

	// Execute the login page template
	ok(Templates.ExecuteTemplate(w, "signup", nil))
}

func signup_submit(w http.ResponseWriter, r *http.Request) {

	// Hash password
	h := sha1.New()
	h.Write([]byte(password))
	hp := hex.EncodeToString(h.Sum(nil))

    // Store new user
	kv.Set("fit:user:%d")

    // Update user quantity
	kv.Set("fit:users")

    // Commit
	kv.Write()

	// Add the new account to the session list
    //

	// Set the user cookie for subsequent requests
    //

    // Result
}

func card_page(w http.ResponseWriter, r *http.Request) {

	// Store the starting exercise
	plan = r.URL.Query().Get("start")

	// Construct a new exercise plan by replicating the schema
	user.ActiveExercisePlan, err = NewExercisePlan(PlanCatalog[plan])

	// Initialize the counter
	user.ActiveExercisePlan.Ti = -1

	// Result
    // 
}

func card_advance(w http.ResponseWriter, r *http.Request) {

	// Parse advance
    //
    
	// Increment the logical timestamp
	u.ActiveExercisePlan.Ti++

	// Edge case ( user is done )
	if u.ActiveExercisePlan.Ti >= u.ActiveExercisePlan.T1 {

        // Respond
		ok(Templates.ExecuteTemplate(w, "done", active_exercise))
		return
	}

	// Set the day's exercise plan
	kv.Set("fit:user:%d:day:%d")

	// Global time - Start of this set 
	var timestamp_logical_set int = u.ActiveExercisePlan.Ti - u.GetActiveExercise().T0

	// Set the specific act
	kv.Set("fit:user:%d:day:%d:act:%d")
    
	// Set completed 
    // NOTE: by testing the parity of timestamp_logical_set, one can deduce
    //       if the timestamp was taken from the start or end of the rest period
	if timestamp_logical_set%2 == 0 && timestamp_logical_set != 0 {
        
        // Store a set
	    kv.Set("fit:user:%d:day:%d:act:%d:set:%d")

	}

	// Send the next card
    //
}

func card_done(w http.ResponseWriter, r *http.Request) {

	// Update DB
	kv.Set("fit:user:%d", u.ID)

    // Commit changes
	kv.Write()

    // Update session
	Sessions[u.Username] = u

	// All done
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func history(w http.ResponseWriter, r *http.Request) {

    // Get the days
	days := kv.GetAll("fit:user:%d:day")

    // Foreach day
	for i := user.Day - 1; i >= 0; i-- {

        // Get the acts
		acts := kv.GetAll("fit:user:%d:day:%d:act")

        // Foreach act
		for j := 0; j < len(acts); j++ {

            // Get the sets
            sets := kv.GetAll("fit:user:%d:day:%d:act:%d:set")
            
            // Foreach set
			for k := 0; k < len(sets); k++ {
                // Decode each set
			}
		}
	}

	// Execute the history template with the result of the for for for loop 
    // 
}

// Initialize
func init() {

	// Connect to the database
	kv = data.NewDB("data.g10.app:6764")

	// Get each user
	db_users := kv.GetAll("fit:user")

	// Construct a session foreach user
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
	err := http.ListenAndServeTLS(":8082", "/root/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/g10.app/g10.app.crt", "/root/.local/share/caddy/certificates/acme-v02.api.letsencrypt.org-directory/g10.app/g10.app.key", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
```