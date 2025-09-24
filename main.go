package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

// one set entry in the resonee json format
type DataEntry struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	ScryfallURI string `json:"scryfall_uri"`
	ReleasedAt  string `json:"released_at"`
	IconSVGURI  string `json:"icon_svg_uri"`
}

// a full response from the api conatining the sets as an array
type Response struct {
	Object  string      `json:"object"`
	HasMore bool        `json:"has_more"`
	Data    []DataEntry `json:"data"`
}

// csv filename to output the data to
const outFile string = "sets.csv"
const apiURL string = "https://api.scryfall.com/sets"

func main() {
	sets, err := fetchSets()
	if err != nil {
		log.Fatalf("Error fetching sets: %v", err)
	}

	fmt.Printf("Fetched %d sets\n", len(sets))

	sortSetsByReleaseDate(sets)
	fmt.Println("Sorted sets by release date")

	err = writeToCSV(sets, outFile)
	if err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}

	fmt.Printf("Exported %d sets to %s\n", len(sets), outFile)
}

func fetchSets() ([]DataEntry, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return response.Data, nil
}

func sortSetsByReleaseDate(sets []DataEntry) {
	sort.Slice(sets, func(i, j int) bool {
		return sets[i].ReleasedAt < sets[j].ReleasedAt
	})
}

func writeToCSV(sets []DataEntry, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Code", "Name", "API_url", "Released", "Icon_url"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, set := range sets {
		record := []string{
			set.Code,
			set.Name,
			set.ScryfallURI,
			set.ReleasedAt,
			set.IconSVGURI,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}
