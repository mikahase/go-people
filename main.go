package main

// Original available from https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Person struct {
	ID        int      `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	ID      int    `json:"id,omitempty"`
	Street  string `json:"street,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

var people []Person

// our main function
func main() {
	// manually enter the test data
	people = append(people, Person{ID: 1, Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: 2, Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: 3, Firstname: "Francis", Lastname: "Sunday"})

	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "passw0rd"
	dbName := "rest_api"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err != nil {
		panic(err.Error())
	}

	return db
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	var peeps []Person
	db := dbConn()

	selPerson, err := db.Query("SELECT p.person_id, p.firstname, p.lastname FROM PERSON p")

	if err != nil {
		panic(err.Error())
	}

	for selPerson.Next() {

		var personId int
		var firstname, lastname string

		err = selPerson.Scan(&personId, &firstname, &lastname)
		if err != nil {
			panic(err.Error())
		}

		curPerson := Person{}
		curPerson.ID = personId
		curPerson.Firstname = firstname
		curPerson.Lastname = lastname
		fmt.Printf("Result %v %v\n", firstname, lastname)

		selAddress, err := db.Query("SELECT a.address_id, a.street, a.zipcode, a.city, a.state, a.country FROM ADDRESS a WHERE a.person_id = ?", personId)
		if err != nil {
			panic(err.Error())
		}

		for selAddress.Next() {
			var addressId int
			var street, zipcode, city, state, country string

			err := selAddress.Scan(&addressId, &street, &zipcode, &city, &state, &country)
			if err != nil {
				panic(err.Error())
			}

			curAddress := Address{}
			curAddress.ID = addressId
			curAddress.Street = street
			curAddress.Zipcode = zipcode
			curAddress.City = city
			curAddress.State = state
			curAddress.Country = country
			curPerson.Address = &curAddress
		}

		peeps = append(peeps, curPerson)
	}

	json.NewEncoder(w).Encode(peeps)
	defer db.Close()
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	personId, _ := strconv.Atoi(params["id"])

	for _, item := range people {
		if item.ID == personId {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID, _ = strconv.Atoi(params["id"])
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	personId, _ := strconv.Atoi(params["id"])
	for index, item := range people {
		if item.ID == personId {
			// ... unpacks the slice
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
