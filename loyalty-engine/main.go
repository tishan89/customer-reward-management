package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	UserID    string `json:"userID"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserReward struct {
	UserID               string `json:"userID"`
	SelectedRewardDealID string `json:"selectedRewardDealID"`
	Timestamp            string `json:"timestamp"` // Consider using time.Time if you need date-time operations
	AcceptedTnC          bool   `json:"acceptedTnC"`
	UserAgent            string `json:"userAgent"`
	IPAddress            string `json:"ipAddress"`
}

var userReward []UserReward
var users []User

func getRewards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userReward)
}

func getUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range users {
		if item.UserID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&User{})
}

func main() {
	r := mux.NewRouter()
	userReward = append(userReward, UserReward{"U451298", "RWD34589", "2023-09-04T14:32:21Z", true, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36", "192.168.1.101"})
	userReward = append(userReward, UserReward{"U451299", "RWD34590", "2023-09-04T14:32:21Z", true, "Mozilla/5.0 (Linux; Android 10; SM-A205F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36", "192.168.1.102"})

	users = append(users, User{"U451298", "John", "Doe", "john@example.com"})
	users = append(users, User{"U451299", "Katie", "Smith", "katie@example.com"})
	users = append(users, User{"U451300", "Peter", "Parker", "peter@example.com"})

	r.HandleFunc("/user-rewards", getRewards).Methods("GET")
	r.HandleFunc("/user/{id}", getUserDetails).Methods("GET")
	http.ListenAndServe(":8080", r)
}
