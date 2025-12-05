package robloxapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserInfoResponse struct {
	Description string `json:"description"`
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
	IsBanned    bool   `json:"isBanned"`
}

func GetUserInfo(userID int) (UserInfoResponse, error) {
	user, err := FetchUserInfo(userID)
	if err != nil {
		return UserInfoResponse{}, err
	}

	return user, nil
}

func FetchUserInfo(userID int) (UserInfoResponse, error) {
	resp, err := http.Get("https://users.roblox.com/v1/users/" + fmt.Sprint(userID))
	if err != nil {
		return UserInfoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return UserInfoResponse{}, fmt.Errorf("Non-200 status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo UserInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return UserInfoResponse{}, err
	}

	return userInfo, nil
}
