package main

import (
	"context"
	"fmt"
)

const apiKey = "API_KEY"

func main() {
	ctx := context.Background()
	client, err := NewPlaceClient(ctx, apiKey)
	if err != nil {
		panic(err)
	}

	input := AutoCompleteInput{
		CountryCode:  []string{"vn"},
		Language:     "vi",
		KeyWord:      "132 ham nghi",
		Limit:        5,
		SessionToken: "",
	}

	output, err := client.PlaceAutocomplete(ctx, input)
	if err != nil {
		panic(err)
	}

	if len(output.Places) == 0 {
		fmt.Println("No result")
		return
	}
	for _, place := range output.Places {
		fmt.Printf("%s - %s\n", place.ID, place.Address)
	}

	// Get place detail
	inputDetail := GetPlaceDetailInput{
		PlaceID:      output.Places[0].ID,
		Language:     "vi",
		SessionToken: output.SessionToken,
	}
	outputDetail, err := client.GetPlaceDetail(ctx, inputDetail)
	if err != nil {
		panic(err)
	}
	fmt.Printf("- %s\n- %s\n- %s\n", outputDetail.ID, outputDetail.Name, outputDetail.Address)
}
