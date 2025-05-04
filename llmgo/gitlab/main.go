package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func newClient() (*gitlab.Client, error) {
	fmt.Println(os.Getenv("GITLAB_TOKEN"))
	fmt.Println(os.Getenv("GITLAB_API_URL"))
	client, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(os.Getenv("GITLAB_API_URL")))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 修改为格式化输出
func formatOutput() *openai.ResponseFormat {
	return &openai.ResponseFormat{
		Type: "json_schema",
		JSONSchema: &openai.ResponseFormatJSONSchema{
			Name: "object",
			Schema: &openai.ResponseFormatJSONSchemaProperty{
				Type: "object",
				Properties: map[string]*openai.ResponseFormatJSONSchemaProperty{
					"function_name": {
						Type:        "string",
						Description: "The name of the function",
					},
					"change_content": {
						Type:        "string",
						Description: "content and reason of the change",
					},
					"is_function": {
						Type:        "boolean",
						Description: "Whether the change is a function",
					},
					"suggestion": {
						Type:        "string",
						Description: "suggestion of the change",
					},
				},
				AdditionalProperties: false,
				Required:             []string{"function_name", "change_content", "suggestion", "is_function"},
			},
			Strict: true,
		},
	}

}

func AnalyseMergeRequest(projectId int, mergeRequestId int) (string, error) {
	client, err := newClient()
	if err != nil {
		return "", err
	}
	llmClient, err := llmClient()
	if err != nil {
		return "", err
	}
	page := 1
	diffs := []string{}
	for page > 0 {
		mrWithChanges, resp, err := client.MergeRequests.ListMergeRequestDiffs(projectId, mergeRequestId, &gitlab.ListMergeRequestDiffsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to get MR changes: %w", err)
		}
		page = resp.NextPage
		for _, change := range mrWithChanges {
			if change.Diff == "" {
				continue
			}
			// 如果是新增文件，则跳过
			if change.NewFile {
				continue
			}
			if !strings.HasSuffix(change.NewPath, ".go") {
				continue
			}

			content, err := LLM(llmClient, change.Diff)
			if err != nil {
				return "", err
			}
			diffs = append(diffs, content)
		}
	}
	return strings.Join(diffs, "\n"), nil
}

func llmClient() (*openai.LLM, error) {
	llm, err := openai.New(openai.WithModel("gpt-3.5-turbo"),
		openai.WithBaseURL(os.Getenv("OPENAI_API_BASE")+"/v1"),
		openai.WithToken(os.Getenv("OPENAI_API_KEY")),
		openai.WithResponseFormat(formatOutput()),
	)
	if err != nil {
		return nil, err
	}
	return llm, nil
}

func LLM(llm *openai.LLM, diff string) (string, error) {
	ctx := context.Background()
	prompt := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "作为一个golang专业开发人员，请根据git diff中的内容中获取所有被改动的函数，注意改动的函数必须是在函数体内部进行有效修改;无效修改包含添加注释, 添加日志, 语句中添加空格等相关操作; "),
		llms.TextParts(llms.ChatMessageTypeHuman, diff),
	}
	completion, err := llm.GenerateContent(ctx, prompt)
	if err != nil {
		log.Fatal(err)
	}
	return completion.Choices[0].Content, nil
}

func main() {
	content, err := AnalyseMergeRequest(638, 29)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content)
}
