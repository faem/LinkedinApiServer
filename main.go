package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Profile struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Company  string  `json:"company"`
	Position string  `json:"position"`
	Skill    []Skill `json:"skill"`
}

type Skill struct {
	Name            string `json:"name"`
	NoOfEndorsement int    `json:"noOfEndorsement"`
}

var profilesDB map[string]Profile

func GetProfiles(w http.ResponseWriter, r *http.Request) {
	if !BasicAuth(){
		return
	}
	var profiles []Profile
	//spew.Dump(profilesDB==nil)
	for _, profile := range profilesDB {
		profiles = append(profiles, profile)
	}

	spew.Dump(json.NewEncoder(w).Encode(profiles))
	//w.Write([]byte(fmt.Sprintf("%v", profiles)))
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	if !BasicAuth(){
		return
	}
	params := mux.Vars(r)
	//spew.Dump(params)
	//spew.Dump("GP")
	json.NewEncoder(w).Encode(profilesDB[params["id"]])
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if !BasicAuth(){
		return
	}
	//spew.Dump("delete Profile")
	params := mux.Vars(r)
	//spew.Dump(params["id"])
	delete(profilesDB, params["id"])
}

func AddProfile(w http.ResponseWriter, r *http.Request) {
	if !BasicAuth(){
		return
	}
	var profile Profile
	json.NewDecoder(r.Body).Decode(&profile)
	profilesDB[profile.ID] = profile
	spew.Dump(profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if !BasicAuth(){
		return
	}
	params := mux.Vars(r)
	spew.Dump(params)
	//spew.Dump(r.Body)
	var profile Profile
	json.NewDecoder(r.Body).Decode(&profile)
	profilesDB[params["id"]] = profile
	spew.Dump(profile)
}

func BasicAuth() bool{
	return true
}

func CreateDemoDB() {
	profilesDB = make(map[string]Profile)
	profilesDB["fahim-abrar"] = Profile{
		"fahim-abrar",
		"Mohammad Fahim Abrar",
		"AppsCode Inc.",
		"Software Engineer",
		[]Skill{
			{
				"C++",
				3,
			},
			{
				"Android",
				4,
			},
		}}

	profilesDB["masud-rahman"] = Profile{
		"masud-rahman",
		"Masudur Rahman",
		"AppsCode Inc.",
		"Software Engineer",
		[]Skill{
			{
				"C",
				3,
			},
			{
				"C++",
				4,
			},
		}}
}

func main() {
	router := mux.NewRouter()
	CreateDemoDB()

	router.HandleFunc("/in", GetProfiles).Methods("GET")
	router.HandleFunc("/in/{id}", GetProfile).Methods("GET")
	router.HandleFunc("/in/{id}", UpdateProfile).Methods("PUT")
	router.HandleFunc("/in", AddProfile).Methods("POST")
	router.HandleFunc("/in/{id}", DeleteProfile).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
