package types

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultDenom is the simulation chain's staking denomination.
const DefaultDenom = "ubor"

// Coin is a denomination-qualified amount of tokens.
type Coin struct {
	Denom  string
	Amount int64
}

// String renders a Coin in the canonical "<amount><denom>" form.
func (c Coin) String() string {
	return strconv.FormatInt(c.Amount, 10) + c.Denom
}

// ParseCoin parses a coin literal such as "1000ubor" or "12ubor".
func ParseCoin(s string) (Coin, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Coin{}, fmt.Errorf("invalid coin: empty string")
	}

	// Split the leading sign and digits from the trailing denomination.
	i := 0
	if i < len(s) && (s[i] == '-' || s[i] == '+') {
		i++
	}
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	amtStr, denom := s[:i], s[i:]
	if amtStr == "" || denom == "" {
		return Coin{}, fmt.Errorf("invalid coin %q", s)
	}

	amount, err := strconv.Atoi(amtStr)
	if err != nil {
		return Coin{}, fmt.Errorf("invalid amount %q: %w", amtStr, err)
	}

	return Coin{Denom: denom, Amount: int64(amount)}, nil
}
