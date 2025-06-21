package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"main/chat"
	"main/logger"
	"main/tools"
)

var (
	convo []chat.Message
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Log.Println("MCP started.")

	client := anthropic.NewClient(
		option.WithAPIKey("claude-api-key"),
	)

	for {

		reader := bufio.NewScanner(os.Stdin)
		fmt.Print("[User]: ")
		reader.Scan()
		userInput := strings.TrimSpace(reader.Text())

		if userInput == "exit" {
			break
		}

		convo = append(convo, chat.Message{Role: "user", Content: userInput})

		messages := []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userInput)),
		}

		toolParams := []anthropic.ToolParam{
			{
				Name:        "get_coordinates",
				Description: anthropic.String("Accepts a place as an address, then returns the latitude and longitude coordinates."),
				InputSchema: tools.GetCoordinatesInputSchema,
			},
			{
				Name:        "get_weather",
				Description: anthropic.String("Fetches weather for a given location."),
				InputSchema: tools.GetWeatherInputSchema,
			},
			{
				Name:        "get_emp_leaves",
				Description: anthropic.String("Fetches leave information for an employee based on their employee ID."),
				InputSchema: tools.GetEmpLeavesInputSchema,
			},
		}

		toolsList := make([]anthropic.ToolUnionParam, len(toolParams))
		for i, tool := range toolParams {
			toolsList[i] = anthropic.ToolUnionParam{OfTool: &tool}
		}

		for {
			message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
				Model:     anthropic.ModelClaude3_7SonnetLatest,
				MaxTokens: 1024,
				Messages:  messages,
				Tools:     toolsList,
			})
			if err != nil {
				logger.Log.Printf("API error: %v\n", err)
				panic(err)
			}

			var aiResponse string
			fmt.Print("[assistant]: ")
			for _, block := range message.Content {
				switch b := block.AsAny().(type) {
				case anthropic.TextBlock:
					fmt.Println(b.Text)
					aiResponse += b.Text + "\n"
				case anthropic.ToolUseBlock:
					inputJSON, _ := json.Marshal(b.Input)
					fmt.Println(b.Name + ": " + string(inputJSON))
				}
			}
			convo = append(convo, chat.Message{Role: "assistant", Content: aiResponse})

			messages = append(messages, message.ToParam())
			toolResults := []anthropic.ContentBlockParamUnion{}

			for _, block := range message.Content {
				switch b := block.AsAny().(type) {
				case anthropic.ToolUseBlock:
					logger.Log.Printf("Tool requested: %s\n", block.Name)
					var response interface{}

					switch block.Name {
					case "get_coordinates":
						var input tools.GetCoordinatesInput
						json.Unmarshal([]byte(b.JSON.Input.Raw()), &input)
						response = tools.GetCoordinates(input.Location)
					case "get_weather":
						var input tools.GetWeatherInput
						json.Unmarshal([]byte(b.JSON.Input.Raw()), &input)
						response = tools.GetWeather(input.Location)
					case "get_emp_leaves":
						var input tools.GetEmpLeavesInput
						json.Unmarshal([]byte(b.JSON.Input.Raw()), &input)
						response = tools.GetEmpLeaves(input.EmployeeID)
					}

					resultJSON, _ := json.Marshal(response)
					fmt.Printf("[user (%s)]: %s\n", block.Name, string(resultJSON))

					toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, string(resultJSON), false))
					convo = append(convo, chat.Message{Role: "user", Content: string(resultJSON)})
				}
			}

			if len(toolResults) == 0 {
				break
			}

			messages = append(messages, anthropic.NewUserMessage(toolResults...))
		}
	}

	if err := chat.SaveConversation(convo); err != nil {
		logger.Log.Printf("Failed to save conversation: %v\n", err)
	} else {
		logger.Log.Println("Conversation saved successfully.")
	}
}
