package services

import (
	"encoding/json"
	"fmt"

	// "google-trends-api/src/services"
	"strings"
)

type TrendingItem struct {
	Title        string   `json:"title"`
	Views        string   `json:"views,omitempty"`
	Direction    string   `json:"direction,omitempty"`
	Percentage   string   `json:"percentage,omitempty"`
	TimeAgo      string   `json:"time_ago,omitempty"`
	Status       string   `json:"status,omitempty"`
	State        string   `json:"state,omitempty"`
	MoreCount    string   `json:"more_count,omitempty"`
	RelatedTerms []string `json:"related_terms,omitempty"`
}

func SanitizeHTML() {
	htmlContent := RawData
	lines := strings.Split(htmlContent, "\n")
	var items []TrendingItem
	var currentItem TrendingItem
	var relatedTerms []string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" || line == "\t" {
			continue
		}

		// If we find a line that contains "+" and "more", it's the more count
		if strings.Contains(line, "+") && strings.Contains(line, "more") {
			currentItem.MoreCount = line
			if len(relatedTerms) > 0 {
				currentItem.RelatedTerms = make([]string, len(relatedTerms))
				copy(currentItem.RelatedTerms, relatedTerms)
				relatedTerms = []string{}
			}
			items = append(items, currentItem)
			currentItem = TrendingItem{}
			continue
		}

		// If we have a number followed by "+" (like "10M+", "100K+"), it's the views
		if strings.HasSuffix(line, "+") && currentItem.Views == "" {
			currentItem.Views = line
			continue
		}

		// Check for direction indicators
		if line == "arrow_upward" {
			currentItem.Direction = line
			continue
		}

		// Check for percentage
		if strings.HasSuffix(line, "%") {
			currentItem.Percentage = line
			continue
		}

		// Check for time ago
		if strings.Contains(line, "hours ago") || strings.Contains(line, "hour ago") {
			currentItem.TimeAgo = line
			continue
		}

		// Check for status
		if line == "trending_up" || line == "timelapse" {
			currentItem.Status = line
			continue
		}

		// Check for state
		if line == "Active" || strings.HasPrefix(line, "Lasted") {
			currentItem.State = line
			continue
		}

		// If we have a complete item but no more_count, this might be a related term
		if currentItem.Title != "" && currentItem.Views != "" && currentItem.MoreCount == "" {
			if !strings.HasSuffix(line, "+") && !strings.HasSuffix(line, "%") &&
				line != "arrow_upward" && line != "trending_up" &&
				line != "timelapse" && line != "Active" &&
				!strings.Contains(line, "hour") {
				relatedTerms = append(relatedTerms, line)
				continue
			}
		}

		// If we have an empty title, this must be the title
		if currentItem.Title == "" {
			currentItem.Title = line
			continue
		}

		// If we get here and have a title, we're starting a new item
		if currentItem.Title != "" {
			// If we have related terms but no more_count, add them before resetting
			if len(relatedTerms) > 0 && currentItem.MoreCount == "" {
				currentItem.RelatedTerms = make([]string, len(relatedTerms))
				copy(currentItem.RelatedTerms, relatedTerms)
				relatedTerms = []string{}
			}
			items = append(items, currentItem)
			currentItem = TrendingItem{Title: line}
		}
	}

	// Add the last item if it exists
	if currentItem.Title != "" {
		if len(relatedTerms) > 0 {
			currentItem.RelatedTerms = make([]string, len(relatedTerms))
			copy(currentItem.RelatedTerms, relatedTerms)
		}
		items = append(items, currentItem)
	}

	jsonBytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		// return "", err
		fmt.Println("Error: ", err)
	}

	SanitizedData = string(jsonBytes)
}

// func main() {
// 	// Read input from stdin if no file is specified
// 	var input []byte
// 	var err error

// 	if len(os.Args) > 1 {
// 		// Read from file if provided
// 		input, err = ioutil.ReadFile(os.Args[1])
// 	} else {
// 		// Read from stdin
// 		input, err = ioutil.ReadAll(os.Stdin)
// 	}

// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
// 		os.Exit(1)
// 	}

// 	jsonOutput, err := sanitizeHTML(string(input))
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error processing HTML: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println(jsonOutput)
// }
