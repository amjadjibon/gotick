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
// s = current stock price, k = strike, r = risk-free rate, t = time to expiry (years)
// sigma = implied volatility, isCall = true for call, false for put
func CalculateGreeks(s, k, r, t, sigma float64, isCall bool) *Greeks {
	if t <= 0 || sigma <= 0 {
		return nil
	}

	sqrtT := math.Sqrt(t)
	d1 := (math.Log(s/k) + (r+sigma*sigma/2)*t) / (sigma * sqrtT)
	d2 := d1 - sigma*sqrtT

	// Standard normal CDF
	nd1Val := normalCDF(d1)
	nd2Val := normalCDF(d2)
	nd2Neg := normalCDF(-d2)

	// Standard normal PDF
	nd1PDF := normalPDF(d1)

	g := &Greeks{}

	if isCall {
		// Call option Greeks
		g.Delta = nd1Val
		g.Theta = -(s * nd1PDF * sigma / (2 * sqrtT)) - r*k*math.Exp(-r*t)*nd2Val
		g.Rho = k * t * math.Exp(-r*t) * nd2Val / 100 // Per 1% change
	} else {
		// Put option Greeks
		g.Delta = nd1Val - 1
		g.Theta = -(s * nd1PDF * sigma / (2 * sqrtT)) + r*k*math.Exp(-r*t)*nd2Neg
		g.Rho = -k * t * math.Exp(-r*t) * nd2Neg / 100 // Per 1% change
	}

	// Common Greeks
	g.Gamma = nd1PDF / (s * sigma * sqrtT)
	g.Vega = s * sqrtT * nd1PDF / 100 // Per 1% change in IV

	// Convert theta to daily
	g.Theta /= 365

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

// ImpliedVolatility calculates implied volatility using Newton-Raphson method
func ImpliedVolatility(marketPrice, s, k, r, t float64, isCall bool) float64 {
	const maxIterations = 100
	const tolerance = 0.0001

	sigma := 0.3 // Initial guess

	for i := 0; i < maxIterations; i++ {
		price := blackScholesPrice(s, k, r, t, sigma, isCall)
		vega := blackScholesVega(s, k, r, t, sigma)

		if vega == 0 {
			break
		}

		diff := marketPrice - price
		if math.Abs(diff) < tolerance {
			return sigma
		}

		sigma += diff / vega
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
func blackScholesPrice(s, k, r, t, sigma float64, isCall bool) float64 {
	if t <= 0 {
		if isCall {
			return math.Max(s-k, 0)
		}
		return math.Max(k-s, 0)
	}

	sqrtT := math.Sqrt(t)
	d1 := (math.Log(s/k) + (r+sigma*sigma/2)*t) / (sigma * sqrtT)
	d2 := d1 - sigma*sqrtT

	if isCall {
		return s*normalCDF(d1) - k*math.Exp(-r*t)*normalCDF(d2)
	}
	return k*math.Exp(-r*t)*normalCDF(-d2) - s*normalCDF(-d1)
}

// blackScholesVega calculates vega for IV calculation
func blackScholesVega(s, k, r, t, sigma float64) float64 {
	if t <= 0 {
		return 0
	}
	sqrtT := math.Sqrt(t)
	d1 := (math.Log(s/k) + (r+sigma*sigma/2)*t) / (sigma * sqrtT)
	return s * sqrtT * normalPDF(d1)
}
