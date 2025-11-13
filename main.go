package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

const apiBaseURL = "https://www.dictionaryapi.com/api/v3/references/collegiate/json"

// DictionaryEntry represents a single entry from the Merriam-Webster API
type DictionaryEntry struct {
	Meta struct {
		ID string `json:"id"`
		FL string `json:"fl"` // functional label (part of speech)
	} `json:"meta"`
	Hwi struct {
		Hw string `json:"hw"` // headword
		Prs []struct {
			Mw string `json:"mw"` // pronunciation
		} `json:"prs"`
	} `json:"hwi"`
	Shortdef []string `json:"shortdef"` // short definitions
}

func main() {
	if len(os.Args) < 2 {
		color.Red("Usage: %s <word>", os.Args[0])
		os.Exit(1)
	}

	word := os.Args[1]
	apiKey := os.Getenv("DICTIONARY_KEY")

	if apiKey == "" {
		color.Red("Error: DICTIONARY_KEY environment variable not set")
		color.Yellow("Please set your Merriam-Webster API key:")
		color.Yellow("  export DICTIONARY_KEY=your-api-key-here")
		os.Exit(1)
	}

	definitions, err := fetchDefinition(word, apiKey)
	if err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}

	printDefinitions(definitions, word)
}

func fetchDefinition(word, apiKey string) ([]DictionaryEntry, error) {
	url := fmt.Sprintf("%s/%s?key=%s", apiBaseURL, word, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for API errors (API returns 200 OK with plain text error messages)
	bodyStr := string(body)
	if strings.Contains(bodyStr, "Invalid API key") {
		return nil, fmt.Errorf("invalid API key - please check your DICTIONARY_KEY environment variable")
	}

	// Try to parse as dictionary entries
	var entries []DictionaryEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		// If parsing fails, it might be a list of suggestions
		var suggestions []string
		if err := json.Unmarshal(body, &suggestions); err == nil {
			color.Red("Word not found. Did you mean:")
			for _, suggestion := range suggestions {
				color.Red("  - %s", suggestion)
			}
			os.Exit(1)
		}
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no definitions found")
	}

	return entries, nil
}

func printDefinitions(entries []DictionaryEntry, word string) {
	cyan := color.New(color.FgCyan, color.Bold)
	yellow := color.New(color.FgYellow)
	green := color.New(color.FgGreen)

	for i, entry := range entries {
		// Print word and part of speech
		_, _ = cyan.Printf("\n%s", strings.ReplaceAll(entry.Hwi.Hw, "*", "·"))
		if entry.Meta.FL != "" {
			fmt.Printf(" [%s]", entry.Meta.FL)
		}
		fmt.Println()

		// Print pronunciation if available
		if len(entry.Hwi.Prs) > 0 && entry.Hwi.Prs[0].Mw != "" {
			_, _ = yellow.Printf("  /%s/\n", entry.Hwi.Prs[0].Mw)
		}

		// Print definitions
		for j, def := range entry.Shortdef {
			_, _ = green.Printf("  %d. %s\n", j+1, def)
		}

		// Add spacing between entries
		if i < len(entries)-1 {
			fmt.Println()
		}
	}
}
