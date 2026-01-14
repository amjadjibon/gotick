package yfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GetEarningsCalendar fetches upcoming earnings events
func GetEarningsCalendar(ctx context.Context, params CalendarParams) ([]EarningsEvent, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return GetEarningsCalendarWithClient(ctx, client, params)
}

// GetEarningsCalendarWithClient fetches earnings calendar using a specific client
func GetEarningsCalendarWithClient(ctx context.Context, client *Client, params CalendarParams) ([]EarningsEvent, error) {
	queryParams := buildCalendarParams(params, "earnings")
	data, err := client.Get(ctx, CalendarURL, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Rows []struct {
					Symbol           string  `json:"ticker"`
					CompanyShortName string  `json:"companyshortname"`
					StartDateTime    string  `json:"startDateTime"`
					EpsEstimate      float64 `json:"epsestimate,omitempty"`
				} `json:"rows"`
			} `json:"result"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var events []EarningsEvent
	if len(response.Finance.Result) > 0 {
		for _, row := range response.Finance.Result[0].Rows {
			event := EarningsEvent{Symbol: row.Symbol, CompanyShortName: row.CompanyShortName}
			if t, err := time.Parse(time.RFC3339, row.StartDateTime); err == nil {
				event.EarningsDate = t.Unix()
			}
			events = append(events, event)
		}
	}
	return events, nil
}

// GetIPOCalendar fetches upcoming IPO events
func GetIPOCalendar(ctx context.Context, params CalendarParams) ([]IPOEvent, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	return GetIPOCalendarWithClient(ctx, client, params)
}

// GetIPOCalendarWithClient fetches IPO calendar using a specific client
func GetIPOCalendarWithClient(ctx context.Context, client *Client, params CalendarParams) ([]IPOEvent, error) {
	queryParams := buildCalendarParams(params, "ipo")
	data, err := client.Get(ctx, CalendarURL, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Rows []struct {
					Symbol      string `json:"ticker"`
					CompanyName string `json:"companyName"`
					Exchange    string `json:"exchange"`
					PricingDate string `json:"pricingDate"`
				} `json:"rows"`
			} `json:"result"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var events []IPOEvent
	if len(response.Finance.Result) > 0 {
		for _, row := range response.Finance.Result[0].Rows {
			event := IPOEvent{Symbol: row.Symbol, CompanyName: row.CompanyName, Exchange: row.Exchange}
			if t, err := time.Parse("2006-01-02", row.PricingDate); err == nil {
				event.PricingDate = t.Unix()
			}
			events = append(events, event)
		}
	}
	return events, nil
}

// GetSplitsCalendar fetches upcoming stock split events
func GetSplitsCalendar(ctx context.Context, params CalendarParams) ([]SplitEvent, error) {
	client, err := getDefaultClient()
	if err != nil {
		return nil, err
	}
	queryParams := buildCalendarParams(params, "splits")
	data, err := client.Get(ctx, CalendarURL, queryParams)
	if err != nil {
		return nil, err
	}

	var response struct {
		Finance struct {
			Result []struct {
				Rows []struct {
					Symbol     string `json:"ticker"`
					SplitDate  string `json:"date"`
					SplitRatio string `json:"splitRatio"`
				} `json:"rows"`
			} `json:"result"`
		} `json:"finance"`
	}

	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	var events []SplitEvent
	if len(response.Finance.Result) > 0 {
		for _, row := range response.Finance.Result[0].Rows {
			event := SplitEvent{Symbol: row.Symbol, SplitRatio: row.SplitRatio}
			if t, err := time.Parse("2006-01-02", row.SplitDate); err == nil {
				event.SplitDate = t.Unix()
			}
			events = append(events, event)
		}
	}
	return events, nil
}

func buildCalendarParams(params CalendarParams, calendarType string) url.Values {
	queryParams := url.Values{}
	if params.Start.IsZero() {
		params.Start = time.Now()
	}
	if params.End.IsZero() {
		params.End = params.Start.AddDate(0, 0, 7)
	}
	queryParams.Set("startDate", params.Start.Format("2006-01-02"))
	queryParams.Set("endDate", params.End.Format("2006-01-02"))
	if params.Size > 0 {
		queryParams.Set("size", strconv.Itoa(params.Size))
	}
	queryParams.Set("type", calendarType)
	return queryParams
}
