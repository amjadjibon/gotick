package yfinance

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

// Stream represents a real-time WebSocket connection for streaming quotes
type Stream struct {
	symbols  []string
	conn     *websocket.Conn
	messages chan StreamMessage
	errors   chan error
	done     chan struct{}
	mu       sync.Mutex
	running  bool
}

// NewStream creates a new WebSocket stream for the given symbols
func NewStream(symbols []string) *Stream {
	return &Stream{
		symbols:  symbols,
		messages: make(chan StreamMessage, 100),
		errors:   make(chan error, 10),
		done:     make(chan struct{}),
	}
}

// Connect establishes a WebSocket connection
func (s *Stream) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, WebSocketURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	s.conn = conn
	s.running = true

	// Subscribe to symbols
	if len(s.symbols) > 0 {
		if err := s.subscribe(s.symbols); err != nil {
			s.conn.Close()
			s.running = false
			return err
		}
	}

	// Start reading messages
	go s.readLoop()

	return nil
}

// subscribe sends a subscription message
func (s *Stream) subscribe(symbols []string) error {
	msg := map[string]interface{}{
		"subscribe": symbols,
	}
	return s.conn.WriteJSON(msg)
}

// unsubscribe sends an unsubscription message
func (s *Stream) unsubscribe(symbols []string) error {
	msg := map[string]interface{}{
		"unsubscribe": symbols,
	}
	return s.conn.WriteJSON(msg)
}

// Subscribe adds symbols to the subscription
func (s *Stream) Subscribe(symbols ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		s.symbols = append(s.symbols, symbols...)
		return nil
	}

	if err := s.subscribe(symbols); err != nil {
		return err
	}
	s.symbols = append(s.symbols, symbols...)
	return nil
}

// Unsubscribe removes symbols from the subscription
func (s *Stream) Unsubscribe(symbols ...string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	if err := s.unsubscribe(symbols); err != nil {
		return err
	}

	// Remove from symbols list
	newSymbols := make([]string, 0, len(s.symbols))
	removeSet := make(map[string]bool)
	for _, sym := range symbols {
		removeSet[sym] = true
	}
	for _, sym := range s.symbols {
		if !removeSet[sym] {
			newSymbols = append(newSymbols, sym)
		}
	}
	s.symbols = newSymbols
	return nil
}

// Messages returns a channel for receiving stream messages
func (s *Stream) Messages() <-chan StreamMessage {
	return s.messages
}

// Errors returns a channel for receiving errors
func (s *Stream) Errors() <-chan error {
	return s.errors
}

// readLoop continuously reads messages from the WebSocket
func (s *Stream) readLoop() {
	defer func() {
		close(s.messages)
		close(s.errors)
	}()

	for {
		select {
		case <-s.done:
			return
		default:
			_, data, err := s.conn.ReadMessage()
			if err != nil {
				s.errors <- err
				return
			}

			msg, err := parseStreamMessage(data)
			if err != nil {
				s.errors <- err
				continue
			}

			select {
			case s.messages <- *msg:
			default:
				// Channel full, skip message
			}
		}
	}
}

// parseStreamMessage parses a WebSocket message using protobuf
func parseStreamMessage(data []byte) (*StreamMessage, error) {
	// Yahoo Finance sends base64-encoded protobuf messages
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		// Try parsing as JSON fallback
		var msg StreamMessage
		if jsonErr := json.Unmarshal(data, &msg); jsonErr != nil {
			return nil, fmt.Errorf("failed to parse message: %w", err)
		}
		return &msg, nil
	}

	// Parse using protobuf
	pricingData := &PricingData{}
	if err := proto.Unmarshal(decoded, pricingData); err != nil {
		// Fallback to empty message if proto fails
		return &StreamMessage{}, nil
	}

	// Convert PricingData to StreamMessage
	msg := &StreamMessage{
		ID:            pricingData.GetId(),
		Price:         float64(pricingData.GetPrice()),
		Time:          pricingData.GetTime(),
		Currency:      pricingData.GetCurrency(),
		Exchange:      pricingData.GetExchange(),
		MarketHours:   int(pricingData.GetMarketHours()),
		ChangePercent: float64(pricingData.GetChangePercent()),
		DayVolume:     pricingData.GetDayVolume(),
		DayHigh:       float64(pricingData.GetDayHigh()),
		DayLow:        float64(pricingData.GetDayLow()),
		Change:        float64(pricingData.GetChange()),
		PreviousClose: float64(pricingData.GetPreviousClose()),
		Bid:           float64(pricingData.GetBid()),
		BidSize:       pricingData.GetBidSize(),
		Ask:           float64(pricingData.GetAsk()),
		AskSize:       pricingData.GetAskSize(),
	}

	return msg, nil
}

// Close closes the WebSocket connection
func (s *Stream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	close(s.done)
	s.running = false

	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// IsConnected returns whether the stream is connected
func (s *Stream) IsConnected() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// Symbols returns the currently subscribed symbols
func (s *Stream) Symbols() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]string, len(s.symbols))
	copy(result, s.symbols)
	return result
}
