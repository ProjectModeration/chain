package interaction

import (
	"fmt"
	"os"

	api "ProjectModeration/chain/chain/robloxapi"

	"github.com/gtuk/discordwebhook"
)

func SendModerationResults(imageRes string, textRes int, UserID int) {

	url := os.Getenv("DISCORD_WEBHOOK_URL")

	info, err := api.GetUserInfo(UserID)
	if err != nil {
		fmt.Println(err)
	}

	textResult := ""
	if imageRes == "nsfw" || textRes == 1 {
		fmt.Println("not safe")
		switch imageRes {
		case "nsfw":
			imageRes = "`NSFW DETECTED` Source : `/M.T.D.I.R.A/`"
		case "normal":
			imageRes = "Safe"
		}
		switch textRes {
		case 1:
			textResult = "`NSFW BIO DETECTED` Source : `/sybauML/`"
		default:
			textResult = "Safe"
		}

	} else {
		fmt.Println("safe")
		return
	}

	desc := fmt.Sprintf("Description moderation result : %s\nAvatar moderation result: %s.\nUser Description : %s", textResult, imageRes, info.Description)
	title := "LIVE | Moderation CHAIN result"
	color := "16711680"
	footertext := "Powered by Project Moderation."
	iconfooter := "https://avatars.githubusercontent.com/u/240123181?s=200&v=4" // great
	imageUrl := api.GetUserAvatar(UserID)

	image := discordwebhook.Image{Url: &imageUrl}
	footer := discordwebhook.Footer{
		Text:    &footertext,
		IconUrl: &iconfooter,
	}

	embed := discordwebhook.Embed{
		Title:       &title,
		Description: &desc,
		Color:       &color,
		Footer:      &footer,
		Image:       &image,
	}

	content := "WARN : This data might not be accurate."

	msg := discordwebhook.Message{
		Username: nil,
		Content:  &content,
		Embeds:   &[]discordwebhook.Embed{embed},
	}

	err = discordwebhook.SendMessage(url, msg)
	if err != nil {
		fmt.Println("error sending message:", err)
	}
}
