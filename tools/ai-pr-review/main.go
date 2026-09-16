// Command ai-pr-review sends a git diff to the Claude API along with this
// repo's CLAUDE.md conventions and asks for a first-pass review: does the
// diff violate any stated architecture rule or convention?
//
// This is a first-pass filter, not a replacement for human review. It
// will miss things a human catches (business-logic correctness, whether
// the change actually solves the right problem) and it can be wrong about
// what it flags — treat its output as "worth a second look," not verdict.
//
// Usage:
//
//	git diff main... | go run ./tools/ai-pr-review
//
// Requires ANTHROPIC_API_KEY in the environment.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const claudeMdPath = "CLAUDE.md"

const systemPromptTemplate = `You are reviewing a git diff against this repository's stated coding
conventions. Only flag things the conventions below explicitly call out —
do not invent new style opinions. For each violation, cite the specific
rule from the conventions and quote the offending line(s) from the diff.
If nothing violates the stated conventions, say so plainly and briefly.

Repository conventions (CLAUDE.md):
---
%s
---
`

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "ai-pr-review: error:", err)
		os.Exit(1)
	}
}

func run() error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY is not set")
	}

	diff, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading diff from stdin: %w", err)
	}
	if len(bytes.TrimSpace(diff)) == 0 {
		return fmt.Errorf("empty diff on stdin; try: git diff main... | go run ./tools/ai-pr-review")
	}

	conventions, err := os.ReadFile(claudeMdPath)
	if err != nil {
		return fmt.Errorf("reading %s (run from repo root): %w", claudeMdPath, err)
	}

	reqBody := anthropicRequest{
		Model:     "claude-sonnet-4-6",
		MaxTokens: 1024,
		System:    fmt.Sprintf(systemPromptTemplate, conventions),
		Messages: []anthropicMessage{
			{Role: "user", Content: string(diff)},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling Anthropic API: %w", err)
	}
	defer resp.Body.Close()

	var parsed anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	if parsed.Error != nil {
		return fmt.Errorf("API error: %s", parsed.Error.Message)
	}

	for _, block := range parsed.Content {
		if block.Type == "text" {
			fmt.Println(block.Text)
		}
	}
	return nil
}
