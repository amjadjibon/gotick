package yfinance

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

// TestGreeksCalculation tests Black-Scholes Greeks calculation
func TestGreeksCalculation(t *testing.T) {
	// Test case: AAPL call option
	// S=150, K=150, r=0.05, T=0.25 (3 months), sigma=0.25
	greeks := CalculateGreeks(150, 150, 0.05, 0.25, 0.25, true)

	if greeks == nil {
		t.Fatal("Expected non-nil Greeks")
	}

	// Delta for ATM call should be around 0.5
	if greeks.Delta < 0.45 || greeks.Delta > 0.65 {
		t.Errorf("Expected Delta around 0.5, got %f", greeks.Delta)
	}

	// Gamma should be positive
	if greeks.Gamma <= 0 {
		t.Errorf("Expected positive Gamma, got %f", greeks.Gamma)
	}

	// Theta should be negative for long options
	if greeks.Theta >= 0 {
		t.Errorf("Expected negative Theta, got %f", greeks.Theta)
	}

	// Vega should be positive
	if greeks.Vega <= 0 {
		t.Errorf("Expected positive Vega, got %f", greeks.Vega)
	}
}

// TestGreeksPutOption tests put option Greeks
func TestGreeksPutOption(t *testing.T) {
	greeks := CalculateGreeks(150, 150, 0.05, 0.25, 0.25, false)

	if greeks == nil {
		t.Fatal("Expected non-nil Greeks")
	}

	// Delta for ATM put should be around -0.5
	if greeks.Delta > -0.35 || greeks.Delta < -0.65 {
		t.Errorf("Expected Delta around -0.5, got %f", greeks.Delta)
	}
}

// TestImpliedVolatility tests IV calculation
func TestImpliedVolatility(t *testing.T) {
	S, K, r, T := 150.0, 150.0, 0.05, 0.25
	expectedSigma := 0.25

	// Calculate option price with known sigma
	price := blackScholesPrice(S, K, r, T, expectedSigma, true)

	// Calculate IV from price
	iv := ImpliedVolatility(price, S, K, r, T, true)

	// IV should be close to original sigma
	if math.Abs(iv-expectedSigma) > 0.01 {
		t.Errorf("Expected IV around %f, got %f", expectedSigma, iv)
	}
}

// TestCacheMemory tests memory cache operations
func TestCacheMemory(t *testing.T) {
	cache := NewCache(CacheConfig{
		Type:       CacheTypeMemory,
		DefaultTTL: 1 * time.Minute,
		MaxSize:    100,
	})

	key := "test_key"
	data := []byte("test_data")

	// Test Set and Get
	cache.Set(key, data, 0)
	retrieved, ok := cache.Get(key)

	if !ok {
		t.Error("Expected cache hit")
	}

	if string(retrieved) != string(data) {
		t.Errorf("Expected %s, got %s", string(data), string(retrieved))
	}

	// Test Delete
	cache.Delete(key)
	_, ok = cache.Get(key)

	if ok {
		t.Error("Expected cache miss after delete")
	}
}

// TestCacheExpiration tests cache TTL
func TestCacheExpiration(t *testing.T) {
	cache := NewCache(CacheConfig{
		Type:       CacheTypeMemory,
		DefaultTTL: 50 * time.Millisecond,
		MaxSize:    100,
	})

	key := "expiring_key"
	data := []byte("expiring_data")

	cache.Set(key, data, 50*time.Millisecond)

	// Should exist immediately
	_, ok := cache.Get(key)
	if !ok {
		t.Error("Expected cache hit before expiration")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get(key)
	if ok {
		t.Error("Expected cache miss after expiration")
	}
}

// TestRetryBackoff tests backoff calculation
func TestRetryBackoff(t *testing.T) {
	backoff := calculateBackoff(1*time.Second, 30*time.Second, 0)
	if backoff != 1*time.Second {
		t.Errorf("Expected 1s backoff, got %v", backoff)
	}

	// Test max limit
	backoff = calculateBackoff(60*time.Second, 30*time.Second, 0)
	if backoff != 30*time.Second {
		t.Errorf("Expected 30s max backoff, got %v", backoff)
	}
}

// TestRateLimiter tests rate limiter
func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, 5) // 10 req/s, burst of 5

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Should be able to make burst requests immediately
	for i := 0; i < 5; i++ {
		err := rl.Wait(ctx)
		if err != nil {
			t.Errorf("Expected no error on burst request %d, got %v", i, err)
		}
	}
}

// TestNewTicker tests ticker creation
func TestNewTicker(t *testing.T) {
	ticker, err := NewTicker("AAPL")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ticker.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", ticker.Symbol)
	}
}

// TestNewTickerEmpty tests empty symbol
func TestNewTickerEmpty(t *testing.T) {
	_, err := NewTicker("")
	if !errors.Is(err, ErrInvalidSymbol) {
		t.Errorf("Expected ErrInvalidSymbol, got %v", err)
	}
}

// TestNewTickers tests batch ticker creation
func TestNewTickers(t *testing.T) {
	tickers, err := NewTickers("AAPL", "GOOGL", "MSFT")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(tickers.Symbols()) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(tickers.Symbols()))
	}

	_, ok := tickers.Ticker("AAPL")
	if !ok {
		t.Error("Expected to find AAPL ticker")
	}
}

// TestPeriodConstants tests period constant values
func TestPeriodConstants(t *testing.T) {
	tests := []struct {
		period Period
		want   string
	}{
		{Period1d, "1d"},
		{Period5d, "5d"},
		{Period1mo, "1mo"},
		{Period1y, "1y"},
		{PeriodMax, "max"},
	}

	for _, tt := range tests {
		if string(tt.period) != tt.want {
			t.Errorf("Expected %s, got %s", tt.want, string(tt.period))
		}
	}
}

// TestIntervalConstants tests interval constant values
func TestIntervalConstants(t *testing.T) {
	tests := []struct {
		interval Interval
		want     string
	}{
		{Interval1m, "1m"},
		{Interval5m, "5m"},
		{Interval1h, "1h"},
		{Interval1d, "1d"},
		{Interval1wk, "1wk"},
	}

	for _, tt := range tests {
		if string(tt.interval) != tt.want {
			t.Errorf("Expected %s, got %s", tt.want, string(tt.interval))
		}
	}
}

// TestDefaultModules tests default modules helper
func TestDefaultModules(t *testing.T) {
	modules := DefaultModules()
	if len(modules) == 0 {
		t.Error("Expected non-empty default modules")
	}
}

// TestFinancialModules tests financial modules helper
func TestFinancialModules(t *testing.T) {
	modules := FinancialModules()
	if len(modules) != 6 {
		t.Errorf("Expected 6 financial modules, got %d", len(modules))
	}
}

// TestErrorTypes tests custom error types
func TestErrorTypes(t *testing.T) {
	apiErr := &APIError{
		Code:        "test",
		Description: "test error",
		StatusCode:  400,
	}

	errStr := apiErr.Error()
	if errStr == "" {
		t.Error("Expected non-empty error string")
	}

	symbolErr := NewSymbolError("AAPL", ErrNotFound)
	if symbolErr.Error() == "" {
		t.Error("Expected non-empty symbol error string")
	}
}

// TestTickerSymbolNormalization tests that ticker symbol is uppercased
func TestTickerSymbolNormalization(t *testing.T) {
	ticker, err := NewTicker("aapl")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ticker.Symbol != "AAPL" {
		t.Errorf("Expected symbol AAPL, got %s", ticker.Symbol)
	}
}

// TestTickerWithCustomClient tests ticker with custom client option
func TestTickerWithCustomClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("Expected no error creating client, got %v", err)
	}

	ticker, err := NewTicker("MSFT", WithClient(client))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if ticker.Symbol != "MSFT" {
		t.Errorf("Expected symbol MSFT, got %s", ticker.Symbol)
	}
}

// TestStreamCreation tests WebSocket stream creation
func TestStreamCreation(t *testing.T) {
	stream := NewStream([]string{"AAPL", "GOOGL"})

	symbols := stream.Symbols()
	if len(symbols) != 2 {
		t.Errorf("Expected 2 symbols, got %d", len(symbols))
	}

	// Should not be connected yet since Connect() wasn't called
	if stream.IsConnected() {
		t.Error("Expected stream to not be connected before Connect()")
	}
}

// TestStreamSymbolManagement tests stream subscribe/unsubscribe before connection
func TestStreamSymbolManagement(t *testing.T) {
	stream := NewStream([]string{})

	err := stream.Subscribe("AAPL", "GOOGL")
	if err != nil {
		t.Errorf("Expected no error on subscribe before connect, got %v", err)
	}

	if len(stream.Symbols()) != 2 {
		t.Errorf("Expected 2 symbols, got %d", len(stream.Symbols()))
	}
}

// TestOptionChainWithGreeks tests option chain Greeks calculation
func TestOptionChainWithGreeks(t *testing.T) {
	chain := &OptionChain{
		Symbol:          "AAPL",
		UnderlyingPrice: 150.0,
		Calls: []Option{
			{Strike: 150, ImpliedVolatility: 0.25, Expiration: time.Now().Add(30 * 24 * time.Hour).Unix()},
		},
		Puts: []Option{
			{Strike: 150, ImpliedVolatility: 0.25, Expiration: time.Now().Add(30 * 24 * time.Hour).Unix()},
		},
	}

	withGreeks := chain.WithGreeks(0.05)

	if len(withGreeks.Calls) != 1 {
		t.Errorf("Expected 1 call, got %d", len(withGreeks.Calls))
	}

	if len(withGreeks.Puts) != 1 {
		t.Errorf("Expected 1 put, got %d", len(withGreeks.Puts))
	}
}

// TestCacheKeyGeneration tests cache key generation
func TestCacheKeyGeneration(t *testing.T) {
	params := map[string]string{"symbol": "AAPL", "modules": "price"}
	key := CacheKey("quote", params)

	if key == "" {
		t.Error("Expected non-empty cache key")
	}

	// Same inputs should produce same key
	key2 := CacheKey("quote", params)
	if key != key2 {
		t.Errorf("Expected same cache key, got %s vs %s", key, key2)
	}

	// Different inputs should produce different keys
	params2 := map[string]string{"symbol": "GOOGL", "modules": "price"}
	key3 := CacheKey("quote", params2)
	if key == key3 {
		t.Error("Expected different cache keys for different symbols")
	}
}

// TestCacheClear tests cache clear operation
func TestCacheClear(t *testing.T) {
	cache := NewCache(CacheConfig{
		Type:       CacheTypeMemory,
		DefaultTTL: 1 * time.Minute,
		MaxSize:    100,
	})

	cache.Set("key1", []byte("data1"), 0)
	cache.Set("key2", []byte("data2"), 0)

	cache.Clear()

	_, ok := cache.Get("key1")
	if ok {
		t.Error("Expected cache miss after clear")
	}
}

// TestRequestError tests request error type
func TestRequestError(t *testing.T) {
	reqErr := &RequestError{
		Endpoint: "/test",
		Method:   "GET",
		Err:      errors.New("connection failed"),
	}

	errStr := reqErr.Error()
	if errStr == "" {
		t.Error("Expected non-empty error string")
	}

	// Test Unwrap
	unwrapped := reqErr.Unwrap()
	if unwrapped == nil {
		t.Error("Expected non-nil unwrapped error")
	}
}

// TestGreeksNilForInvalidInput tests Greeks returns nil for invalid inputs
func TestGreeksNilForInvalidInput(t *testing.T) {
	// Zero time to expiry
	greeks := CalculateGreeks(150, 150, 0.05, 0, 0.25, true)
	if greeks != nil {
		t.Error("Expected nil Greeks for zero time to expiry")
	}

	// Negative time to expiry
	greeks = CalculateGreeks(150, 150, 0.05, -0.1, 0.25, true)
	if greeks != nil {
		t.Error("Expected nil Greeks for negative time to expiry")
	}

	// Zero volatility
	greeks = CalculateGreeks(150, 150, 0.05, 0.25, 0, true)
	if greeks != nil {
		t.Error("Expected nil Greeks for zero volatility")
	}
}

// TestBlackScholesPriceEdgeCases tests edge cases in BS pricing
func TestBlackScholesPriceEdgeCases(t *testing.T) {
	// ITM call at expiry should be intrinsic value
	price := blackScholesPrice(160, 150, 0.05, 0, 0.25, true)
	if price != 10.0 {
		t.Errorf("Expected ITM call price 10, got %f", price)
	}

	// OTM call at expiry should be 0
	price = blackScholesPrice(140, 150, 0.05, 0, 0.25, true)
	if price != 0 {
		t.Errorf("Expected OTM call price 0, got %f", price)
	}

	// ITM put at expiry should be intrinsic value
	price = blackScholesPrice(140, 150, 0.05, 0, 0.25, false)
	if price != 10.0 {
		t.Errorf("Expected ITM put price 10, got %f", price)
	}
}
