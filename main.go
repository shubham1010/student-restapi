package main

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"io/ioutil"

	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

type MigrationLogger struct {
	verbose bool
}

type Student struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Age string `json:"age"`
	Department string `json:"dept"`
	Subject json.RawMessage `json:"subject"`
}

func (ml *MigrationLogger) Printf(format string, v ...interface{}) {
	log.Printf(format,v)
}

func (ml *MigrationLogger) Verbose() bool {
	return ml.verbose
}

var db *sql.DB
var err error
func main() {
	db ,err = sql.Open("mysql", "shubham1010:#Dontknow1010@/demo")


	if err != nil {
	  //panic(err.Error())
	  //log.Fatal("[DBCONNECTION]: FAILED",err)
	  fmt.Println("[DBCONNECTION]: FAILED",err)
	}

	if err := migrateDatabase(db); err!=nil {
		log.Fatal("[MigrateDatabase]:\n", err)
	}

	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/students", getStudents).Methods("GET")
	router.HandleFunc("/student", createStudent).Methods("POST")
	router.HandleFunc("/student/{id}", getStudent).Methods("GET")
	router.HandleFunc("/student/{id}", updateStudent).Methods("PUT")
	router.HandleFunc("/student/{id}", deleteStudent).Methods("DELETE")
	http.ListenAndServe(":8080", router)
	log.Println("Server is Listening...")
}


func migrateDatabase(db *sql.DB) error {
	log.Println("[INSIDE MIGRATION]")
	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		return err
	}

	dir, err := os.Getwd()
	log.Println("[DIR]: ",dir)
	if err != nil {
		log.Fatal(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", dir),
		"mysql",
		driver,
	)

	if err !=nil {
		return err
	}

	migration.Log = &MigrationLogger{}

	migration.Log.Printf("Applying database migrations")
	err = migration.Up()
	if err!=nil && err != migrate.ErrNoChange {
		return err
	}

	version, _, err := migration.Version()
	if err != nil {
		return err
	}
	migration.Log.Printf("Active database version: %d", version)
	return nil
}


func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var students []Student
	result, err := db.Query("SELECT * from students")
	if err != nil {
		//panic(err.Error())
		log.Fatal("No entry present in table: ", err);
	}
	defer result.Close()
	for result.Next() {
		var student Student
		err := result.Scan(&student.ID, &student.Name, &student.Age, &student.Department, &student.Subject)
		if err != nil {
			//panic(err.Error())
			log.Fatal("[GETALLRECORDS] Fatching Failed: ", err)
		}
		students = append(students, student)
	}

	json.NewEncoder(w).Encode(students)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	if(r.Method!="POST") {
		w.Write([]byte("Invalid Request"))
	}
	w.Header().Set("Content-Type", "application/json")


	body, err := ioutil.ReadAll(r.Body)
	var student Student
    err = json.Unmarshal(body, &student)

	if err != nil {
		//panic(err.Error())
		log.Fatal("Unmarshal Error: ", err)
	}
	stmt, err := db.Prepare("INSERT INTO students(id,name,age,dept,subject) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
		log.Fatal("[CREATESTUDENT]: Insertion Failed: ", err)
	}

	_, err = stmt.Exec(student.ID,student.Name,student.Age,student.Department,student.Subject)
	if err != nil {
		//panic(err.Error())
		log.Fatal("[CreateUser]: Insertion Failed", err)
	}

	fmt.Fprintf(w, "New student is created")
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM students WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var student Student
	for result.Next() {
		err := result.Scan(&student.ID, &student.Name, &student.Age, &student.Department, &student.Subject)
		if err != nil {
			//panic(err.Error())
			log.Fatal("[GETSTUDENT]: Fetching Failed: ", err)
		}
	}
	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE students SET name=?, age=?, dept=?, subject=? WHERE id = ?")
	if err != nil {
		//panic(err.Error())
		log.Fatal("[UPDATESTUDENT]: prepare failed: ", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	var student Student
    err = json.Unmarshal(body, &student)

	_, err = stmt.Exec(student.Name, student.Age, student.Department, student.Subject, params["id"])
	if err != nil {
		//panic(err.Error())
		log.Fatal("[UPDATESTUDENT]: Exec Failed: ", err)
	}
	fmt.Fprintf(w, "Student with ID = %s is updated", params["id"])
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		//panic(err.Error())
		log.Fatal("[DELETESTUDENT]: Prepare Failed: ", err)
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
	//	panic(err.Error())
	log.Fatal("[DELETESTUDENT]: Exec Failed: ", err)
	}
	fmt.Fprintf(w, "Student with ID = %s is deleted", params["id"])
}

