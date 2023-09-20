package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type User struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type UserReward struct {
	UserId               string `json:"userId"`
	SelectedRewardDealId string `json:"selectedRewardDealId"`
	Timestamp            string `json:"timestamp"` // Consider using time.Time if you need date-time operations
	AcceptedTnC          bool   `json:"acceptedTnC"`
}

var logger *zap.Logger
var userRewards []UserReward
var users []User

func getRewards(w http.ResponseWriter, r *http.Request) {
	logger.Info("get all rewards")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRewards)
}

func getUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range users {
		if item.UserId == params["id"] {
			json.NewEncoder(w).Encode(item)
			logger.Info("get user details", zap.Any("user", item))
			return
		}
	}

	logger.Info("user not found", zap.String("user id", params["id"]))
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&User{})
}

func main() {

	defer logger.Sync() // Ensure all buffered logs are written

	logger.Info("Starting the loyalty engine...")

	r := mux.NewRouter()
	userRewards = append(userRewards, UserReward{"U451298", "RWD34589", "2023-09-04T14:32:21Z", true})
	userRewards = append(userRewards, UserReward{"U451299", "RWD34590", "2023-09-04T14:32:21Z", true})

	users = append(users, User{"U451298", "John", "Doe", "john@example.com"})
	users = append(users, User{"U451299", "Katie", "Smith", "katie@example.com"})
	users = append(users, User{"U451300", "Peter", "Parker", "peter@example.com"})

	r.HandleFunc("/user-rewards", getRewards).Methods("GET")
	r.HandleFunc("/user/{id}", getUserDetails).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
