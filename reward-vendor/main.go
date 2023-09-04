package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2/clientcredentials"
)

type Reward struct {
	RewardId  string `json:"rewardId"`
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type RewardConfirmation struct {
	RewardConfirmationNumber string `json:"rewardConfirmationNumber"`
	UserId                   string `json:"userId"`
	RewardId                 string `json:"rewardId"`
}

var logger *zap.Logger
var rewards []Reward

var rewardMgtClientId = os.Getenv("REWARD_MGT_CLIENT_ID")
var rewardMgtClientSecret = os.Getenv("REWARD_MGT_CLIENT_SECRET")
var rewardMgtTokenUrl = os.Getenv("REWARD_MGT_TOKEN_URL")
var rewardMgtApiUrl = os.Getenv("REWARD_MGT_API_URL")

var config = clientcredentials.Config{
	ClientID:     rewardMgtClientId,
	ClientSecret: rewardMgtClientSecret,
	TokenURL:     rewardMgtTokenUrl,
}

func RespondWithRewardConfirmation(rewardId string, userId string) {
	logger.Info("responding with reward confirmation")
	client := config.Client(context.Background())

	// Generate the 16-digit number and encapsulate in an anonymous struct
	rewardConfirmation := RewardConfirmation{
		RewardConfirmationNumber: Generate16DigitNumber(),
		UserId:                   userId,
		RewardId:                 rewardId,
	}

	// Convert the anonymous struct to JSON
	data, err := json.Marshal(rewardConfirmation)
	if err != nil {
		logger.Error("Failed to marshal data: %v", err)
	}

	resp, err := client.Post(rewardMgtApiUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logger.Error("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response: %v", err)
	}

	logger.Printf("Response: %s\n", body)

}

func Generate16DigitNumber() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// To ensure it's always 16 digits, generate a number between 1000_0000_0000_0000 and 9999_9999_9999_9999
	number := r.Int63n(9000_0000_0000_0000) + 1000_0000_0000_0000

	return fmt.Sprintf("%016d", number)
}

func CreateReward(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reward Reward
	_ = json.NewDecoder(r.Body).Decode(&reward)

	logger.Info("creating a reward", zap.Any("reward", reward))

	rewards = append(rewards, reward)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reward)

	logger.Info("responding with reward confirmation", zap.Any("reward", reward))
	//RespondWithRewardConfirmation(reward.RewardId, reward.UserId)
}

func main() {
	defer logger.Sync() // Ensure all buffered logs are written

	logger.Info("Starting the reward vender...")

	r := mux.NewRouter()

	// Sample Data
	rewards = append(rewards, Reward{RewardId: "RWD0000", UserId: "U0000", FirstName: "John", LastName: "Doe", Email: "john@example.com"})
	r.HandleFunc("/rewards", CreateReward).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
