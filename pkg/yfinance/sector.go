package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetSectors fetches available sectors
func GetSectors(ctx context.Context) ([]Sector, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return GetSectorsWithClient(ctx, client)
}

// GetSectorsWithClient fetches sectors using a specific client
func GetSectorsWithClient(ctx context.Context, client *Client) ([]Sector, error) {
	data, err := client.Get(ctx, SectorURL, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Sectors []Sector `json:"sectors"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse sectors response: %w", err)
	}

	if response.Finance.Error != nil {
		return nil, &APIError{
			Code:        response.Finance.Error.Code,
			Description: response.Finance.Error.Description,
		}
	}

	if len(response.Finance.Result) > 0 {
		return response.Finance.Result[0].Sectors, nil
	}

	return []Sector{}, nil
}

// GetIndustries fetches available industries
func GetIndustries(ctx context.Context) ([]Industry, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}

	return GetIndustriesWithClient(ctx, client)
}

// GetIndustriesWithClient fetches industries using a specific client
func GetIndustriesWithClient(ctx context.Context, client *Client) ([]Industry, error) {
	data, err := client.Get(ctx, IndustryURL, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Industries []Industry `json:"industries"`
			} `json:"result"`
			Error *struct {
				Code        string `json:"code"`
				Description string `json:"description"`
			} `json:"error"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse industries response: %w", err)
	}

	if response.Finance.Error != nil {
		return nil, &APIError{
			Code:        response.Finance.Error.Code,
			Description: response.Finance.Error.Description,
		}
	}

	if len(response.Finance.Result) > 0 {
		return response.Finance.Result[0].Industries, nil
	}

	return []Industry{}, nil
}

// GetIndustriesBySector fetches industries for a specific sector
func GetIndustriesBySector(ctx context.Context, sector string) ([]Industry, error) {
	industries, err := GetIndustries(ctx)
	if err != nil {
		return nil, err
	}

	var result []Industry
	for _, ind := range industries {
		if ind.Sector == sector {
			result = append(result, ind)
		}
	}

	return result, nil
}

// Predefined sector keys
const (
	SectorTechnology        = "technology"
	SectorHealthcare        = "healthcare"
	SectorFinancials        = "financial-services"
	SectorConsumerCyclical  = "consumer-cyclical"
	SectorConsumerDefensive = "consumer-defensive"
	SectorIndustrials       = "industrials"
	SectorEnergy            = "energy"
	SectorUtilities         = "utilities"
	SectorBasicMaterials    = "basic-materials"
	SectorRealEstate        = "real-estate"
	SectorCommunication     = "communication-services"
)
