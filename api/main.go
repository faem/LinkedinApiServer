package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

//----------------------------------------Structures---------------------------------------------------
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

//-------------------------------------Global Variables---------------------------------------------------
var profilesDB map[string]Profile //map for our demo database, used map key as ID
var apiUser map[string]string //stores the username (key) and pass (value) of the user of the API
var srvr http.Server
var bypassAuthentication bool
//var Stop chan bool
var stopTime int8
var router = mux.NewRouter()

//------------------------------------Handler Functions---------------------------------------------------

//used to get all the profile info using GET request
func GetProfiles(w http.ResponseWriter, r *http.Request) {
	if l, err := BasicAuth(r); !err {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: " + l)))
		return
	}

	var profiles []Profile

	for _, profile := range profilesDB {
		profiles = append(profiles, profile)
	}
	//shows error if no content in the DB
	if profiles == nil {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(fmt.Sprintf("No content found!")))
		return
	}

	json.NewEncoder(w).Encode(profiles)
	//w.Write([]byte(fmt.Sprintf("%v", profiles)))
}

//used to get a specific profile info using GET request
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if l, err := BasicAuth(r); !err {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: " + l)))
		return
	}

	params := mux.Vars(r)

	if _, flag := profilesDB[params["id"]]; !flag {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Profile of " + string(params["id"]) + " not found!")))
		return
	}

	json.NewEncoder(w).Encode(profilesDB[params["id"]])
}

//used to delete a profile using DELETE request
func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if l, err := BasicAuth(r); !err {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: " + l)))
		return
	}

	params := mux.Vars(r)
	if _, flag := profilesDB[params["id"]]; !flag {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Profile of " + string(params["id"]) + " not found!")))
		return
	}

	delete(profilesDB, params["id"])
}

//used to add new profile using POST request
func AddProfile(w http.ResponseWriter, r *http.Request) {
	if l, err := BasicAuth(r); !err {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: " + l)))
		return
	}
	var profile Profile
	json.NewDecoder(r.Body).Decode(&profile)

	if _, flag := profilesDB[profile.ID]; flag {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf("Username " + string(profile.ID) + " already exists!")))
		return
	}

	profilesDB[profile.ID] = profile
}

//used to update a profile using PUT request
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if l, err := BasicAuth(r); !err {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Error: " + l)))
		return
	}

	params := mux.Vars(r)

	if _, flag := profilesDB[params["id"]]; !flag {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Profile of " + string(params["id"]) + " not found!\n")))
		return
	}

	var profile Profile
	//log.Println(r.Body)
	decodedValue := json.NewDecoder(r.Body).Decode(&profile)
	//log.Println(profile)
	if decodedValue != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	//profilesDB[params["id"]] = profile
}

//Shutdown handler function
func ShutDown(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Server is shutting down")))
	err := srvr.Shutdown(context.Background())
	if err!= nil{
		log.Fatal("Error shutting down server")
	}
}

//Get JWT token for authorization
func GetToken(w http.ResponseWriter, r *http.Request){
	if msg, flag := BasicAuth(r); !flag{
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(msg))
		return
	}

	mySigningKey := []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	user, _, _ := GetUserPass(r)

	claims["admin"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, _ := token.SignedString(mySigningKey)

	w.Write([]byte(tokenString))
}

//-------------------------------------Other Functions---------------------------------------------------

func GetUserPass(r *http.Request) (user string, pass string, err string){
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		err = "Need authorization!\n"
	}

	decodedStr, e := base64.StdEncoding.DecodeString(strings.Split(authHeader, " ")[1])
	if e != nil {
		err = "Base64 decoding error!\n"
	}

	userPass := strings.Split(string(decodedStr), ":")

	if len(userPass) != 2 {
		err = "Authorization header format error!\n"
	}

	return userPass[0], userPass[1], err
}

//basic authentication function
func BasicAuth(r *http.Request) (string, bool) {
	if bypassAuthentication{
		return "", true
	}

	user, pass, err := GetUserPass(r)

	if err != ""{
		return err, false
	}


	if p, err := apiUser[user]; err {
		if p == pass {
			return "", true
		} else {
			return "Username and Password doesn't match!\n", false
		}
	} else {
		return "User doesn't exist!\n", false
	}
}



//this function creates a demo DB for our server
func CreateDemoDB() {
	//Creating profiles database
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

	profilesDB["mohan"] = Profile{
		"mohan",
		"Tahsin Rahman",
		"AppsCode Inc.",
		"Software Engineer",
		[]Skill{
			{
				"C",
				100,
			},
			{
				"C++",
				110,
			},
			{
				"Linux",
				100,
			},
		}}

	//creating API user info
	apiUser = make(map[string]string)

	apiUser["fahim"] = "1234"
	apiUser["admin"] = "admin"
}

func SetValues(port string, bpa bool, stop int8){
	srvr.Addr = ":"+port
	bypassAuthentication = bpa
	stopTime = stop
}

func init(){
	//Creating demo database
	CreateDemoDB()
	router.Use(jwtMiddleware)
	//setting router handler functions
	router.HandleFunc("/in", GetProfiles).Methods("GET")
	router.HandleFunc("/in/{id}", GetProfile).Methods("GET")
	router.HandleFunc("/in/{id}", UpdateProfile).Methods("PUT")
	router.HandleFunc("/in", AddProfile).Methods("POST")
	router.HandleFunc("/in/{id}", DeleteProfile).Methods("DELETE")

	router.HandleFunc("/shutdown", ShutDown).Methods("GET")
	router.HandleFunc("/token", GetToken).Methods("GET")
}
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTQzMjU5OTQ2LCJ1c2VyIjoiZmFoaW0ifQ.qRTYLq4en4MMRZdNs3XjhOAOHSrkt_UqZM-xmpnoXIo

func StartServer() {
	log.Println("---------------------Starting server---------------------")

	srvr.Handler = router

	if stopTime!= 0{
		fmt.Println(stopTime)
		go StopServer(stopTime)
	}
	stop := make(chan os.Signal,1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Fatal(srvr.ListenAndServe())

	}()

	<-stop

	StopServer(0)
}

func StopServer(x int8)  {
	if x == 0{
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := srvr.Shutdown(ctx)
		if err != nil{
			log.Println("Error in shutting down server!")
		}
		fmt.Println("")
		log.Println("------------------Shutting down server-----------------------\n")
		return
	}
	timer := time.NewTimer(time.Duration(x)*time.Minute)
	fmt.Println("----------------------Shutting Down server in",x,"min---------------------")
	<-timer.C
	err := srvr.Shutdown(context.Background())
	if err!=nil {
		log.Fatal("Error shutting down server!")
	}
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Middleware")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Need authorization!\n"))
		}

		decodedStr, e := jwt.DecodeSegment(strings.Split(authHeader, " ")[1])
		if e != nil {
			//w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Base64 decoding error!\n"))
		}

		log.Println(string(decodedStr))

		next.ServeHTTP(w, r)
	})
}