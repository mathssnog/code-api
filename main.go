package main 

/*Import the dependencies */
import (
	"database/sql" 
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func main (){
	/* Connection to a Database*/
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		Log.fatal(err)
	}
	defer db.Close()

	/*Let's create routers*/
	router := mux.NewRouter()
	
	/*Route to get all the users*/
	router.HandleFunc("/users", getUsers(db).Methods("GET"))
	
	/*Route to get a specific user*/
	router.HandleFunc("/users/{id}", getUsers(db).Methods("GET"))
	
	/*Route to create users*/
	router.HandleFunc("/users", createUsers(db).Methods("POST"))

	/*Route to update users*/
	router.HandleFunc("/users/{id}", updateUsers(db).Methods("PUT")

	/*Route to delete users*/
	router.HandleFunc("/users/{id}", deleteUsers(db).Methods("DELETE")

	/*let's start the server*/
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))

}

/*here we have the middleware function and http.handler as a parameter */
func jsonContentTypeMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w,r)
	})
}

/*Let's get all users */
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
				log.Fatal(err)
			}
			users = append(users,u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}


/*Let's get user by id*/
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars ["id"]

		var u User
		err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			//remember that the application crush with this error. we need to fix this
			log.Fatal(err)
		}
	}
}

/*Let's create a user*/
func createUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u User
        json.NewDecoder(r.Body).Decode(&u)

        err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
        if err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(u)
    }
}

/*Let's update a user*/
func updateUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var u User
        json.NewDecoder(r.Body).Decode(&u)

        vars := mux.Vars(r)
        id := vars["id"]

        _, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", u.Name, u.Email, id)
        if err != nil {
            log.Fatal(err)
        }

        json.NewEncoder(w).Encode(u)
    }
}

/*Let's delete a user*/
func deleteUser(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        var u User
        err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        } else {
            _, err := db.Exec("DELETE FROM users WHERE id = $1", id)
            if err != nil {
                //remember that the application crush with this error. we need to fix this
                w.WriteHeader(http.StatusNotFound)
                return
            }

            json.NewEncoder(w).Encode("User deleted")
        }
    }
}