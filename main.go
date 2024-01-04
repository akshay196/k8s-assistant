package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var initPrompt = `You are k8s-assistant, an AI agent that assist and manage Kubernetes Cluster via bash commands.
Generate the bash command scripts in response to each of my Instructions.
I will then run those bash commands for you and let you know if there were errors.
You should modify the current list of commands based on my instructions.
You should not start from scratch unless asked.

Do not use the local filesystem.  Do not use kube-config.
`

var (
	messages []azopenai.ChatRequestMessageClassification
)

func main() {

	keyCredential := azcore.NewKeyCredential(os.Getenv("OPENAI_API_KEY"))
	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	model := "gpt-3.5-turbo"

	// very first message
	messages = append(messages, &azopenai.ChatRequestUserMessage{
		Content: azopenai.NewChatRequestUserMessageContent(initPrompt),
	})

	resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &model,
	}, nil)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
	// fmt.Println(*resp.Choices[0].Message.Content)
	messages = append(messages, &azopenai.ChatRequestAssistantMessage{
		Content: resp.Choices[0].Message.Content,
	})

	fmt.Println("I am k8s-assistant, your AI-powered Kubernetes assistant. I will assist you in managing your Kubernetes Cluster. Please provide me with your instruction.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		query := scanner.Text()
		switch query {
		case "!quit":
			return
		default:
			messages = append(messages, &azopenai.ChatRequestUserMessage{
				Content: azopenai.NewChatRequestUserMessageContent(query),
			})
			resp, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
				Messages:       messages,
				DeploymentName: &model,
			}, nil)
			if err != nil {
				log.Fatalf("ERROR: %s", err)
			}
			// act on the message
			fmt.Println(*resp.Choices[0].Message.Content)

			messages = append(messages, &azopenai.ChatRequestAssistantMessage{
				Content: resp.Choices[0].Message.Content,
			})

		}
	}

}
