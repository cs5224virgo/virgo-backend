package googleaiclient

import (
	"context"
	"fmt"

	"github.com/cs5224virgo/virgo/logger"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GoogleAiClient struct {
	client *genai.Client
}

func New(APIKey string) (*GoogleAiClient, error) {
	if APIKey == "" {
		return nil, fmt.Errorf("google ai api key is blank")
	}
	newClient := GoogleAiClient{}
	var err error
	newClient.client, err = genai.NewClient(context.Background(), option.WithAPIKey(APIKey))
	if err != nil {
		return nil, fmt.Errorf("unable to init google ai client: %w", err)
	}
	return &newClient, nil
}

func (c *GoogleAiClient) Close() {
	err := c.client.Close()
	if err != nil {
		logger.Error("cannot close google ai client:", err)
	}
}

func (c *GoogleAiClient) GetSummary(conversation string) (string, error) {
	model := c.client.GenerativeModel("gemini-1.0-pro")
	promptTemplate := fmt.Sprintf(`Summarise the conversation below in a way that the user is able to understand the general idea but not missing out on important details too. 
The summary should be on two levels.  First level is on overall level of the conversation. Second level is individual level, summarise what everyone said. 

Conversation:
%s`, conversation)
	resp, err := model.GenerateContent(context.Background(), genai.Text(promptTemplate))
	if err != nil {
		return "", fmt.Errorf("error querying gemini: %w", err)
	}

	fullResp := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fullResp = fullResp + fmt.Sprintln(part)
			}
		}
	}

	return fullResp, nil
}
