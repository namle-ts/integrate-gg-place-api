package main

import (
	places "cloud.google.com/go/maps/places/apiv1"
	"cloud.google.com/go/maps/places/apiv1/placespb"
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/grpc/metadata"
	"googlemaps.github.io/maps"
)

const ComponentCountryMaxSize = 5

type PlaceClient struct {
	oldClient *maps.Client
	newClient *places.Client
}

func NewPlaceClient(ctx context.Context, apiKey string) (*PlaceClient, error) {
	oldClient, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	newClient, err := places.NewRESTClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &PlaceClient{
		oldClient: oldClient,
		newClient: newClient,
	}, nil
}

func (c *PlaceClient) PlaceAutocomplete(ctx context.Context, input AutoCompleteInput) (*AutoCompleteOutput, error) {
	//build components
	components := make(map[maps.Component][]string)
	if len(input.CountryCode) > 0 {
		if len(input.CountryCode) > ComponentCountryMaxSize {
			components[maps.ComponentCountry] = input.CountryCode[0:ComponentCountryMaxSize]
		} else {
			components[maps.ComponentCountry] = input.CountryCode
		}
	} else {
		components[maps.ComponentCountry] = []string{"vn"}
	}

	//parse token
	var token maps.PlaceAutocompleteSessionToken
	if len(input.SessionToken) == 0 {
		token = maps.NewPlaceAutocompleteSessionToken()
	} else {
		parseUUID, err := uuid.Parse(input.SessionToken)
		if err != nil {
			return nil, err
		}
		token = maps.PlaceAutocompleteSessionToken(parseUUID)
	}

	req := &maps.PlaceAutocompleteRequest{
		Components:   components,
		Language:     input.Language,
		Input:        input.KeyWord,
		SessionToken: token,
	}

	resp, err := c.oldClient.PlaceAutocomplete(ctx, req)
	if err != nil {
		return nil, err
	}

	output := NewAutoCompleteOutput(resp, uuid.UUID(req.SessionToken).String())
	return &output, nil
}

func (c *PlaceClient) GetPlaceDetail(ctx context.Context, input GetPlaceDetailInput) (*GetPlaceDetailOutput, error) {
	parseUUID, err := uuid.Parse(input.SessionToken)
	if err != nil {
		return nil, err
	}

	reqDetail := &maps.PlaceDetailsRequest{
		Language: input.Language,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskName,
			maps.PlaceDetailsFieldMaskFormattedAddress,
			maps.PlaceDetailsFieldMaskAddressComponent,
		},
		PlaceID:      input.PlaceID,
		SessionToken: maps.PlaceAutocompleteSessionToken(parseUUID),
	}
	respDetail, err := c.oldClient.PlaceDetails(ctx, reqDetail)
	if err != nil {
		return nil, err
	}

	output := NewGetPlaceDetailOutput(respDetail)
	return &output, nil
}

func (c *PlaceClient) GetPlaceDetailNew(ctx context.Context, input GetPlaceDetailInput) (*GetPlaceDetailOutput, error) {
	req := &placespb.GetPlaceRequest{
		Name:         fmt.Sprintf("places/%s", input.PlaceID),
		LanguageCode: input.Language,
	}

	md := metadata.Pairs("X-Goog-FieldMask", "displayName,formattedAddress,addressComponents")
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := c.newClient.GetPlace(ctx, req)
	if err != nil {
		return nil, err
	}

	output := NewGetPlaceDetailOutputFromNew(resp)
	return &output, nil
}
