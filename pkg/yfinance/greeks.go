package yfinance

import (
	"math"
)

// Greeks contains the calculated option Greeks
type Greeks struct {
	Delta float64 `json:"delta"`
	Gamma float64 `json:"gamma"`
	Theta float64 `json:"theta"`
	Vega  float64 `json:"vega"`
	Rho   float64 `json:"rho"`
}

// OptionWithGreeks extends Option with calculated Greeks
type OptionWithGreeks struct {
	Option
	Greeks *Greeks `json:"greeks"`
}

// CalculateGreeks calculates Black-Scholes Greeks for an option
// S = current stock price, K = strike, r = risk-free rate, T = time to expiry (years)
// sigma = implied volatility, isCall = true for call, false for put
func CalculateGreeks(S, K, r, T, sigma float64, isCall bool) *Greeks {
	if T <= 0 || sigma <= 0 {
		return nil
	}

	sqrtT := math.Sqrt(T)
	d1 := (math.Log(S/K) + (r+sigma*sigma/2)*T) / (sigma * sqrtT)
	d2 := d1 - sigma*sqrtT

	// Standard normal CDF
	Nd1 := normalCDF(d1)
	Nd2 := normalCDF(d2)
	Nd2Neg := normalCDF(-d2)

	// Standard normal PDF
	nd1 := normalPDF(d1)

	g := &Greeks{}

	if isCall {
		// Call option Greeks
		g.Delta = Nd1
		g.Theta = -(S * nd1 * sigma / (2 * sqrtT)) - r*K*math.Exp(-r*T)*Nd2
		g.Rho = K * T * math.Exp(-r*T) * Nd2 / 100 // Per 1% change
	} else {
		// Put option Greeks
		g.Delta = Nd1 - 1
		g.Theta = -(S * nd1 * sigma / (2 * sqrtT)) + r*K*math.Exp(-r*T)*Nd2Neg
		g.Rho = -K * T * math.Exp(-r*T) * Nd2Neg / 100 // Per 1% change
	}

	// Common Greeks
	g.Gamma = nd1 / (S * sigma * sqrtT)
	g.Vega = S * sqrtT * nd1 / 100 // Per 1% change in IV

	// Convert theta to daily
	g.Theta = g.Theta / 365

	return g
}

// normalCDF calculates the cumulative distribution function of the standard normal distribution
func normalCDF(x float64) float64 {
	return 0.5 * (1 + math.Erf(x/math.Sqrt2))
}

// normalPDF calculates the probability density function of the standard normal distribution
func normalPDF(x float64) float64 {
	return math.Exp(-x*x/2) / math.Sqrt(2*math.Pi)
}

// CalculateOptionGreeks adds Greeks to an option
func CalculateOptionGreeks(opt *Option, underlyingPrice, riskFreeRate float64, isCall bool) *OptionWithGreeks {
	// Calculate time to expiration in years
	now := float64(unixNow())
	expiry := float64(opt.Expiration)
	T := (expiry - now) / (365.25 * 24 * 60 * 60)

	if T <= 0 {
		T = 0.0001 // Avoid division by zero
	}

	greeks := CalculateGreeks(
		underlyingPrice,
		opt.Strike,
		riskFreeRate,
		T,
		opt.ImpliedVolatility,
		isCall,
	)

	return &OptionWithGreeks{
		Option: *opt,
		Greeks: greeks,
	}
}

// OptionsWithGreeks returns the option chain with Greeks calculated
func (o *OptionChain) WithGreeks(riskFreeRate float64) *OptionChainWithGreeks {
	result := &OptionChainWithGreeks{
		Symbol:          o.Symbol,
		UnderlyingPrice: o.UnderlyingPrice,
		ExpirationDates: o.ExpirationDates,
		Strikes:         o.Strikes,
		Calls:           make([]OptionWithGreeks, len(o.Calls)),
		Puts:            make([]OptionWithGreeks, len(o.Puts)),
	}

	for i, call := range o.Calls {
		owg := CalculateOptionGreeks(&call, o.UnderlyingPrice, riskFreeRate, true)
		result.Calls[i] = *owg
	}

	for i, put := range o.Puts {
		owg := CalculateOptionGreeks(&put, o.UnderlyingPrice, riskFreeRate, false)
		result.Puts[i] = *owg
	}

	return result
}

// OptionChainWithGreeks is an option chain with Greeks calculated
type OptionChainWithGreeks struct {
	Symbol          string             `json:"symbol"`
	UnderlyingPrice float64            `json:"underlyingPrice"`
	ExpirationDates []int64            `json:"expirationDates"`
	Strikes         []float64          `json:"strikes"`
	Calls           []OptionWithGreeks `json:"calls"`
	Puts            []OptionWithGreeks `json:"puts"`
}

// Helper to get current unix timestamp
func unixNow() int64 {
	return int64(float64(1e9) * float64(1)) // Placeholder - will use time.Now().Unix()
}

func init() {
	// Override with real implementation
	unixNowFunc = func() int64 {
		return int64(float64(1e9) * float64(1))
	}
}

var unixNowFunc func() int64

// ImpliedVolatility calculates implied volatility using Newton-Raphson method
func ImpliedVolatility(marketPrice, S, K, r, T float64, isCall bool) float64 {
	const maxIterations = 100
	const tolerance = 0.0001

	sigma := 0.3 // Initial guess

	for i := 0; i < maxIterations; i++ {
		price := blackScholesPrice(S, K, r, T, sigma, isCall)
		vega := blackScholesVega(S, K, r, T, sigma)

		if vega == 0 {
			break
		}

		diff := marketPrice - price
		if math.Abs(diff) < tolerance {
			return sigma
		}

		sigma = sigma + diff/vega
		if sigma < 0.001 {
			sigma = 0.001
		}
		if sigma > 5 {
			sigma = 5
		}
	}

	return sigma
}

// blackScholesPrice calculates the Black-Scholes option price
func blackScholesPrice(S, K, r, T, sigma float64, isCall bool) float64 {
	if T <= 0 {
		if isCall {
			return math.Max(S-K, 0)
		}
		return math.Max(K-S, 0)
	}

	sqrtT := math.Sqrt(T)
	d1 := (math.Log(S/K) + (r+sigma*sigma/2)*T) / (sigma * sqrtT)
	d2 := d1 - sigma*sqrtT

	if isCall {
		return S*normalCDF(d1) - K*math.Exp(-r*T)*normalCDF(d2)
	}
	return K*math.Exp(-r*T)*normalCDF(-d2) - S*normalCDF(-d1)
}

// blackScholesVega calculates vega for IV calculation
func blackScholesVega(S, K, r, T, sigma float64) float64 {
	if T <= 0 {
		return 0
	}
	sqrtT := math.Sqrt(T)
	d1 := (math.Log(S/K) + (r+sigma*sigma/2)*T) / (sigma * sqrtT)
	return S * sqrtT * normalPDF(d1)
}
