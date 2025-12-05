package robloxapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type FriendsResponse struct {
	Data []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"data"`
}

func GetUserFriends(userID int) []string {
	result, err := GetFriends(userID)
	if err != nil {
		fmt.Println("Error fetching friends:", err)
		return nil
	}

	var friendNames []string
	for _, friend := range result.Data {
		friendNames = append(friendNames, friend.Name)
	}
	return friendNames
}

func GetUserFriendIDs(userID int) []int {
	result, err := GetFriends(userID)
	if err != nil {
		fmt.Println("Error fetching friends:", err)
		return nil
	}

	var friendIDs []int
	for _, friend := range result.Data {
		friendIDs = append(friendIDs, friend.Id)
	}
	return friendIDs
}

func GetFriends(userID int) (FriendsResponse, error) {
	resp, err := http.Get("https://friends.roblox.com/v1/users/" + fmt.Sprint(userID) + "/friends")

	if err != nil {
		fmt.Println("Error:", err)
		return FriendsResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return FriendsResponse{}, err
	}

	var friendsResponse FriendsResponse
	err = json.Unmarshal(body, &friendsResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return FriendsResponse{}, err
	}

	return friendsResponse, nil
}
