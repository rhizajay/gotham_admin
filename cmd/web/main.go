package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rhizajay/gotham_admin"
)

// Context keep context of the running application
type Context struct {
	static_content string
	customer_db    map[string]*sql.DB
}

type RhizaWebUser struct {
	Id        int
	Username  string
	Email     string
	Is_active int
	// created time.Time
	// lastLogin time.Time
	Groups map[string]bool
}

func createFromRhizaUser(ru gotham_admin.RhizaUser) RhizaWebUser {
	var r RhizaWebUser
	r.Id = ru.Id
	r.Username = ru.Username
	r.Email = ru.Email
	r.Is_active = ru.Is_active

	groupmap := make(map[string]bool)

	for key, value := range ru.Groups {
		groupmap[strconv.Itoa(key)] = value
	}
	r.Groups = groupmap
	return r
}

var customers = [4]string{"comcast_spotlight", "bbc", "coxauto", "cox"}

func (c *Context) routes() *mux.Router {
	r := mux.NewRouter()

	r.Path("/customer").Methods("GET").HandlerFunc(c.getCustomers)
	r.Path("/customer/{token}").Methods("GET").HandlerFunc(c.customer)

	r.Path("/customer/{token}/users").Methods("GET").HandlerFunc(c.getUsers)
	r.Path("/customer/{token}/users/{userid}").Methods("GET").HandlerFunc(c.getUserById)
	// r.Path("/customer/{token}/users/{userid}").Methods("PUT").HandlerFunc(c.activateUserById)
	// r.Path("/customer/{token}/users/{userid}").Methods("DELETE").HandlerFunc(c.deactivateUserById)
	r.Path("/customer/{token}/users/e/{email}").Methods("GET").HandlerFunc(c.getUserByEmail)
	// r.Path("/customer/{token}/users/e/{email}").Methods("PUT").HandlerFunc(c.activateUserByEmail)
	// r.Path("/customer/{token}/users/e/{email}").Methods("DELETE").HandlerFunc(c.deactivateUserByEmail)

	r.Path("/customer/{token}/groups").Methods("GET").HandlerFunc(c.getGroups)
	r.Path("/customer/{token}/groups/{groupid}").Methods("GET").HandlerFunc(c.getGroupMembers)
	// r.Path("/customer/{token}/groups/{groupid}").Methods("POST").HandlerFunc(c.addMember)
	// r.Path("/customer/{token}/groups/{groupid}").Methods("DELETE").HandlerFunc(c.removeMember)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(c.static_content)))
	return r
}

func main() {
	// var (
	// 	port    = flag.String("port", "7449", "web server port")
	// 	static  = flag.String("static", "/opt/buildbot/static/", "static folder")
	// 	baseURL = flag.String("baseurl", os.Getenv("BASE_URL"), "local base url")
	// )
	flag.Parse()

	dbmap := setupDB()
	var context = Context{"/opt/gothadmin/static", dbmap}
	r := context.routes()

	http.ListenAndServe(":8080", r)
}

func (c *Context) getCustomers(w http.ResponseWriter, r *http.Request) {
	b, _ := json.Marshal(customers)

	w.Write([]byte(b))
}

func UsersToJSON(members []gotham_admin.RhizaUser) []RhizaWebUser {
	webUsers := make([]RhizaWebUser, len(members))
	for key, value := range members {
		webUsers[key] = createFromRhizaUser(value)
	}

	return webUsers
}

func (c *Context) getUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer := vars["token"]

	admin := gotham_admin.GothamDB{c.customer_db[customer]}

	accounts := UsersToJSON(admin.GetUsers())
	b, _ := json.Marshal(accounts)

	w.Write([]byte(b))
}

func (c *Context) getUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer := vars["token"]
	uid, err := strconv.Atoi(vars["userid"])
	if err != nil {
		println("Error : Not an number")
	} else {

		admin := gotham_admin.GothamDB{c.customer_db[customer]}
		user := createFromRhizaUser(admin.GetUserById(uid))
		b, _ := json.Marshal(user)
		w.Write([]byte(b))

	}
}

func (c *Context) getUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer := vars["token"]
	user_email := vars["email"]

	admin := gotham_admin.GothamDB{c.customer_db[customer]}
	user := createFromRhizaUser(admin.GetUserByEmail(user_email))
	b, _ := json.Marshal(user)
	w.Write([]byte(b))

}

func (c *Context) customer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer := vars["token"]
	println(vars)
	println(customer)
	resp := "Welcome to " + customer

	w.Write([]byte(resp))
	w.Write([]byte("Welcome Admin!\n"))
}

func (c *Context) getGroups(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	customer := vars["token"]

	admin := gotham_admin.GothamDB{c.customer_db[customer]}
	b, _ := json.Marshal(admin.GetGroupNames())

	w.Write([]byte(b))
}

func (c *Context) getGroupMembers(w http.ResponseWriter, r *http.Request) {

	println("get group members")

	vars := mux.Vars(r)
	customer := vars["token"]
	group_param := vars["groupid"]

	println("going for" + group_param)
	admin := gotham_admin.GothamDB{c.customer_db[customer]}
	groupid, err := strconv.Atoi(group_param)
	if err != nil {
		println("Error : Not an number")
	}

	members := admin.GetGroupMembersByGroupId(groupid)
	webUsers := make([]RhizaWebUser, len(members))
	for key, value := range members {
		webUsers[key] = createFromRhizaUser(value)
	}
	b, err := json.Marshal(webUsers)
	if err != nil {
		println("Error in Marshal")
		panic(err)
	}
	w.Write([]byte(b))
}

func setupDB() map[string]*sql.DB {
	database := make(map[string]*sql.DB)
	for _, customer := range customers {
		println(customer)
		println("connectDB:" + customer)

		var username string = os.Getenv("RHIZA_ADMIN_USER")
		var password string = os.Getenv("RHIZA_ADMIN_PASS")
		var dbhost string = os.Getenv("RHIZA_USERDB")
		var db *sql.DB

		s := []string{customer, "_user"}
		customerdb := strings.Join(s, "")

		s = []string{username, ":", password, "@tcp(", dbhost, ")/", customerdb}
		connectString := strings.Join(s, "")

		println(connectString)

		var err error
		db, err = sql.Open("mysql", connectString)
		if err != nil {
			panic(err)
		}

		// Open doesn't open a connection. Validate DSN data:
		err = db.Ping()
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		database[customer] = db
	}
	return database
}
