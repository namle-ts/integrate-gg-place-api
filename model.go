package main

import (
	"cloud.google.com/go/maps/places/apiv1/placespb"
	"googlemaps.github.io/maps"
)

type AutoCompleteInput struct {
	CountryCode  []string `json:"country_code"`
	Language     string   `json:"language"`
	KeyWord      string   `json:"key_word"`
	Limit        int      `json:"limit"`
	SessionToken string   `json:"session_token"`
}

type Place struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type AutoCompleteOutput struct {
	Places       []Place `json:"places"`
	SessionToken string  `json:"session_token"`
}

func NewAutoCompleteOutput(resp maps.AutocompleteResponse, sessionToken string) AutoCompleteOutput {
	var places []Place
	for _, pred := range resp.Predictions {
		places = append(places,
			Place{
				ID:      pred.PlaceID,
				Name:    pred.Description,
				Address: pred.Description,
			},
		)
	}
	return AutoCompleteOutput{
		Places:       places,
		SessionToken: sessionToken,
	}
}

type GetPlaceDetailInput struct {
	PlaceID      string `json:"place_id"`
	Language     string `json:"language"`
	SessionToken string `json:"session_token"`
}

type AddressComponent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type GetPlaceDetailOutput struct {
	ID                string             `json:"id"`
	Address           string             `json:"address"`
	Name              string             `json:"name"`
	AddressComponents []AddressComponent `json:"address_components"`
}

func NewGetPlaceDetailOutputFromNew(resp *placespb.Place) GetPlaceDetailOutput {
	var addressComponents []AddressComponent

	for _, comp := range resp.AddressComponents {
		for _, typeComp := range comp.Types {
			addressComponents = append(addressComponents,
				AddressComponent{
					Type:  typeComp,
					Value: comp.LongText,
				},
			)
		}
	}

	return GetPlaceDetailOutput{
		ID:                resp.Id,
		Address:           resp.FormattedAddress,
		Name:              resp.DisplayName.Text,
		AddressComponents: addressComponents,
	}
}

func NewGetPlaceDetailOutput(resp maps.PlaceDetailsResult) GetPlaceDetailOutput {
	var addressComponents []AddressComponent

	for _, comp := range resp.AddressComponents {
		for _, typeComp := range comp.Types {
			addressComponents = append(addressComponents,
				AddressComponent{
					Type:  typeComp,
					Value: comp.LongName,
				},
			)
		}
	}

	return GetPlaceDetailOutput{
		ID:                resp.PlaceID,
		Address:           resp.FormattedAddress,
		Name:              resp.Name,
		AddressComponents: addressComponents,
	}
}
