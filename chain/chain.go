package chain

import (
	api "ProjectModeration/chain/chain/robloxapi"
	textprocess "ProjectModeration/chain/chain/textprocess"
	inter "ProjectModeration/chain/interaction"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var visited = map[int]bool{}

func StartChaining(startID int) {
	fmt.Println("chaining on", startID)

	// 1. scan main user
	friends, imageRes, textRes := Chain(startID)

	// send moderation for main user
	inter.SendModerationResults(imageRes, textRes, startID)

	// 2. scan ALL friends (one by one)
	for _, friend := range friends {
		if visited[friend] {
			continue
		}

		fmt.Printf("scanning friend: %d\n", friend)
		visited[friend] = true

		fF, fImg, fTxt := Chain(friend)
		inter.SendModerationResults(fImg, fTxt, friend)

		_ = fF // we don’t chain their friends yet, only scan them
	}

	// 3. after scanning all friends, pick ONE random friend to chain next
	if len(friends) == 0 {
		fmt.Println("no friends bruh, chain ends here")
		return
	}

	rand.Seed(time.Now().UnixNano())
	nextID := friends[rand.Intn(len(friends))]

	fmt.Printf("next random dude to chain: %d\n", nextID)

	fmt.Println()

	// 4. chain again but from that random friend
	StartChaining(nextID)
}

func Chain(userID int) ([]int, string, int) {
	if visited[userID] {
		fmt.Println("already scanned bro:", userID)
	}
	visited[userID] = true

	userInfo, err := api.GetUserInfo(userID)
	if err != nil {
		fmt.Println("error fetching user:", err)
		return nil, "", 0
	}

	if userInfo.IsBanned {
		fmt.Println("user banned. chain stops here.")
		return nil, "", 0
	}

	fmt.Printf("user: %s\ndesc: %s\n", userInfo.DisplayName, userInfo.Description)

	friendIDs := api.GetUserFriendIDs(userID)
	if friendIDs == nil {
		fmt.Println("error fetching friends.")
		return nil, "", 0
	}
	fmt.Printf("user has %d friends.\n", len(friendIDs))

	fmt.Println("trying out textprocess")

	raw := userInfo.Description
	clean := textprocess.NormalizeText(raw)

	rawScore := textprocess.ChiScore(clean)

	decoded := false

	if rawScore >= 120 { // only if text looks non-english
		if ok, _, conf := textprocess.DetectROT13(clean); ok && conf >= 60 {
			decodedText := textprocess.ApplyROT13(raw)

			fmt.Println("rot13 detected")
			fmt.Println("decoded:", decodedText)
			fmt.Println("confidence:", conf, "%")

			userInfo.Description = decodedText
			decoded = true

		} else if ok, shift, _, conf := textprocess.DetectCaesar(clean); ok && conf >= 65 {
			decodedText := textprocess.ApplyCaesar(raw, shift)

			fmt.Println("caesar detected")
			fmt.Println("shift:", shift)
			fmt.Println("decoded:", decodedText)
			fmt.Println("confidence:", conf, "%")

			userInfo.Description = decodedText
			decoded = true
		}
	}

	if !decoded {
		fmt.Println("no reliable cipher detected")
	}

	textResult := inter.TextClassification(userInfo.Description)
	fmt.Println("text classifier:", textResult)

	avatarURL := api.GetUserAvatar(userID)
	if avatarURL == "" {
		fmt.Println("no avatar url returned.")
		return friendIDs, "", textResult
	}

	resp, err := http.Get(avatarURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("failed to fetch avatar.")
		if resp != nil {
			resp.Body.Close()
		}
		return friendIDs, "", textResult
	}
	defer resp.Body.Close()

	outputFile, err := os.Create("./Avatar.png")
	if err != nil {
		fmt.Println("error creating file:", err)
		return friendIDs, "", textResult
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fmt.Println("error saving avatar image:", err)
		return friendIDs, "", textResult
	}

	imageResult := inter.ImageClassification("Avatar.png")
	fmt.Println("image classifier:", imageResult)

	// save feedback
	addNewDataToFeedBack(userInfo.Description, textResult)

	return friendIDs, imageResult, textResult
}
