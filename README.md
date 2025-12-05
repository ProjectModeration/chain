# CHAIN - A Roblox Moderation Chaining Tool

CHAIN is a Go-based tool designed to automate content moderation on Roblox. It works by "chaining" through user profiles, starting from a given user ID, and analyzing their description and avatar for inappropriate content. If NSFW content is detected, it sends a notification to a specified Discord webhook.

## How It Works

The tool operates in a continuous loop, performing the following steps for each user in the chain:

1.  **Fetch User Data**: Retrieves the user's profile information, including their description, avatar, and friends list, using the Roblox API.
2.  **Text Moderation**: The user's description is sent to a local text classification service to check for inappropriate language.
3.  **Image Moderation**: The user's avatar is downloaded and sent to a local image classification service to detect NSFW content.
4.  **Report**: If either the text or image moderation service returns a positive result, a detailed alert is sent to a configured Discord webhook. This alert includes the moderation results, the user's description, and their avatar.
5.  **Chain**: The tool randomly selects a user from the current user's friend list and begins the process again with the new user ID, creating a "chain" of moderation.

## Features

*   **Roblox Integration**: Directly interfaces with Roblox APIs to fetch user data.
*   **Text and Image Analysis**: Supports both text and image-based content moderation.
*   **Discord Webhook Notifications**: Provides real-time alerts for moderation events.
*   **Chaining Mechanism**: Autonomously traverses the Roblox social graph to discover new content to moderate.
*   **Local and Modular**: Relies on local classification services, which can be customized or replaced.

## Setup and Configuration

1.  **Dependencies**: Project dependencies are managed with Go Modules. Run `go mod tidy` to ensure you have everything you need.

2.  **Classification Server**: This tool requires a separate server running locally on `http://localhost:5000` that provides the following endpoints for content moderation. See https://github.com/ProjectModeration/mlsidecar to setup.

3.  **Environment Variables**: You must create a `.env` file in the `interaction` directory with the following variable:
    ```
    DISCORD_WEBHOOK_URL=your_discord_webhook_url_here
    ```

## Usage

To run the application, execute the following command from the root of the project:

```bash
go run .
```

By default, the starting user ID is hardcoded in the `main.go` file. You can change this value to start the chain from a different user.