package robloxapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AvatarResponse struct {
	Data []struct {
		State    string `json:"state"`
		ImageURL string `json:"imageUrl"`
	} `json:"data"`
}

func GetUserAvatar(userID int) string {
	result, err := GetAvatarDefault(userID)
	if err != nil {
		fmt.Println("Error fetching avatar:", err)
		return ""
	}

	if len(result.Data) == 0 {
		fmt.Println("No avatar data returned")
		return ""
	}

	if result.Data[0].State != "Completed" {
		fmt.Println("Avatar generation not completed")
		return ""
	}

	return result.Data[0].ImageURL
}

func GetAvatar(userID int) (AvatarResponse, error) {
	resp, err := http.Get("https://thumbnails.roblox.com/v1/users/avatar-bust?userIds=" + fmt.Sprint(userID) + "&size=420x420&format=Png&isCircular=false")

	if err != nil {
		fmt.Println("Error:", err)
		return AvatarResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return AvatarResponse{}, fmt.Errorf("non-200 status: %d", resp.StatusCode)
	}

	var result AvatarResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return AvatarResponse{}, err
	}

	return result, nil
}
func GetAvatarDefault(userID int) (AvatarResponse, error) {
	// https://thumbnails.roblox.com/v1/users/avatar?userIds=8671884543&size=720x720&format=Png&isCircular=false
	resp, err := http.Get("https://thumbnails.roblox.com/v1/users/avatar?userIds=" + fmt.Sprint(userID) + "&size=420x420&format=Png&isCircular=false")

	if err != nil {
		fmt.Println("Error:", err)
		return AvatarResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Non-200 response (%d): %s\n", resp.StatusCode, string(body))
		return AvatarResponse{}, fmt.Errorf("non-200 status: %d", resp.StatusCode)
	}

	var result AvatarResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return AvatarResponse{}, err
	}

	return result, nil
}
