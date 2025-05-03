package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/invopop/jsonschema"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type ToolDefinition struct {
	Name        string
	Description string
	// InputSchema now map[string]any to match OpenAI's parameter schema format
	InputSchema map[string]any
	Function    func(input json.RawMessage) (string, error) // Input is JSON string from OpenAI args
}

// -------------------------- 工具实现 --------------------------

// -------------------------- read_file --------------------------
type ReadFileInput struct { // Defines the input structure for the tool
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory." jsonschema:"required"`
}

// 工具定义
var ReadFileDefinition = ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: GenerateSchema[ReadFileInput](),
	Function:    ReadFile, // Function implementation remains the same
}

func ReadFile(input json.RawMessage) (string, error) { // Implementation unchanged
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse input for read_file: %w. Input was: %s", err, string(input))
	}
	if readFileInput.Path == "" {
		return "", fmt.Errorf("missing required parameter 'path' for read_file")
	}
	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %w", readFileInput.Path, err)
	}
	return string(content), nil
}

// -------------------------- list_files --------------------------
type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

var ListFilesDefinition = ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory. Returns a JSON array of strings, directories have a trailing slash.",
	InputSchema: GenerateSchema[ListFilesInput](),
	Function:    ListFiles, // Function implementation remains the same
}

// ListFiles function implementation unchanged
func ListFiles(input json.RawMessage) (string, error) {
	listFilesInput := ListFilesInput{}
	if len(input) > 0 && string(input) != "null" {
		err := json.Unmarshal(input, &listFilesInput)
		if err != nil {
			return "", fmt.Errorf("failed to parse input for list_files: %w. Input was: %s", err, string(input))
		}
	}
	dir := "."
	if listFilesInput.Path != "" {
		dir = listFilesInput.Path
	}
	var files []string
	err := filepath.WalkDir(dir, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, currentPath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", currentPath, err)
		}
		if relPath == "." {
			return nil
		}
		if d.IsDir() {
			files = append(files, relPath+"/")
		} else {
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error listing files in '%s': %w", dir, err)
	}
	result, err := json.Marshal(files)
	if err != nil {
		return "", fmt.Errorf("failed to marshal file list to JSON: %w", err)
	}
	return string(result), nil
}

// GenerateSchema adapted to return map[string]any
func GenerateSchema[T any]() map[string]any {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties:  false,
		DoNotReference:             true, // Keep definitions inline for OpenAI
		RequiredFromJSONSchemaTags: true, // Respect `jsonschema:"required"`
	}
	var v T
	schema := reflector.Reflect(v)

	// Convert the jsonschema.Schema to map[string]any expected by OpenAI
	// This is a simplification; a full conversion might be more complex
	schemaBytes, _ := json.Marshal(schema)
	var schemaMap map[string]any
	_ = json.Unmarshal(schemaBytes, &schemaMap)

	// OpenAI expects parameters schema directly, remove unnecessary outer layers if present
	if props, ok := schemaMap["properties"]; ok {
		schemaMap["properties"] = props
	}
	if req, ok := schemaMap["required"]; ok {
		schemaMap["required"] = req
	}
	schemaMap["type"] = "object" // Ensure root type is object

	return schemaMap
}

// -------------------------- 工具实现 --------------------------
type GetMergeDiffInput struct {
	ProjectId int `json:"project_id" jsonschema_description:"gitlab project id." jsonschema:"required"`
	MergeId   int `json:"merge_id" jsonschema_description:"gitlab merge request id." jsonschema:"required"`
}

var GetMergeDiffDefinition = ToolDefinition{
	Name:        "get_merge_diff",
	Description: "Get the diff of a merge request.",
	InputSchema: GenerateSchema[GetMergeDiffInput](),
	Function:    GetMergeDiff,
}

func GetMergeDiff(input json.RawMessage) (string, error) {
	getMergeDiffInput := GetMergeDiffInput{}
	err := json.Unmarshal(input, &getMergeDiffInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse input for get_merge_diff: %w. Input was: %s", err, string(input))
	}

	gitClient, err := gitlab.NewClient(os.Getenv("GITLAB_TOKEN"), gitlab.WithBaseURL(os.Getenv("GITLAB_API_URL")))
	if err != nil {
		return "", fmt.Errorf("failed to create GitLab client: %w", err)
	}

	page := 1
	diffs := []string{}
	for page > 0 {
		mrWithChanges, resp, err := gitClient.MergeRequests.ListMergeRequestDiffs(getMergeDiffInput.ProjectId, getMergeDiffInput.MergeId, &gitlab.ListMergeRequestDiffsOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to get MR changes: %w", err)
		}
		page = resp.NextPage
		for _, change := range mrWithChanges {
			diffs = append(diffs, change.Diff)
		}
	}
	return strings.Join(diffs, "\n"), nil
}
