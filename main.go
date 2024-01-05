package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var initPrompt = `You are k8s-assistant, an AI agent that assist and manage Kubernetes Cluster via bash commands.
Generate the bash command scripts in response to each of my Instructions.
I will then run those bash commands for you and let you know if there were errors.
You should modify the current list of commands based on my instructions.
You should not start from scratch unless asked.

Do not use the local filesystem.  Do not use kube-config.

Please output in markdown syntax.
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
			messages = append(messages, &azopenai.ChatRequestAssistantMessage{
				Content: resp.Choices[0].Message.Content,
			})

			msg := *resp.Choices[0].Message.Content
			// fmt.Println(msg)

			// act on response received
			blocks := extractBlocks(msg)
			if blocks == nil {
				fmt.Println("No command to run")
				continue
			}

			for _, block := range blocks {
				fmt.Println(block)
				output, err := executeCommand(block)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(output)
			}

		}
	}

}

func extractBlocks(input string) []string {
	regex := regexp.MustCompile(`(?s)\x60\x60\x60bash(.*?)\x60\x60\x60`)
	matches := regex.FindAllStringSubmatch(input, -1)
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}

func extractBlocksA(input string) [][]string {
	regex := regexp.MustCompile(`(?s)\x60\x60\x60bash(.*?)\x60\x60\x60`)
	matches := regex.FindAllStringSubmatch(input, -1)
	return matches
}

func executeCommand(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error executing command: %v", err)
	}
	return string(output), nil
}
