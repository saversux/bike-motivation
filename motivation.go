package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/vertexai/genai"
	"github.com/bwmarrin/discordgo"
)

var model = "gemini-1.0-pro-001"

var projectID string
var region string
var botToken string
var channelID string

// generate generates content using the specified project ID, region, and model name.
// It returns the generated content as a string.
func generate(projectID string, region string, modelName string) string {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, projectID, region)
	if err != nil {
		log.Fatalf("error creating genai client: %v", err)
	}
	gemini := client.GenerativeModel(modelName)

	prompt := genai.Text("Write a new poem to encourage someone to go outside and cycle, instead of staying inside and code on the computer.")
	resp, err := gemini.GenerateContent(ctx, prompt)
	if err != nil {
		log.Fatalf("error generating content: %v", err)
	}

	return getFirstPart(resp)
}

// sendMessage sends a message to the specified Discord channel using the provided bot token.
func sendMessage(botToken string, channelID string, message string) {
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("error sending message: %v", err)
	}

	discord.ChannelMessageSend(channelID, message)
}

// getFirstPart returns the first part of the generated content response as a string.
func getFirstPart(resp *genai.GenerateContentResponse) string {
	s := ""
	buf := bytes.NewBufferString(s)
	fmt.Fprint(buf, resp.Candidates[0].Content.Parts[0])

	return buf.String()
}

func main() {
	botToken = os.Getenv("DISCORD_BOT_TOKEN")
	channelID = os.Getenv("DISCORD_CHANNEL_ID")
	projectID = os.Getenv("RUN_PROJECT_ID")
	region = os.Getenv("RUN_REGION")

	if botToken == "" || channelID == "" || projectID == "" || region == "" {
		log.Fatalf("missing environment variables DISCORD_BOT_TOKEN, DISCORD_CHANNEL_ID, RUN_PROJECT_ID, RUN_REGION")
	}

	sendMessage(botToken, channelID, generate(projectID, region, model))
}
