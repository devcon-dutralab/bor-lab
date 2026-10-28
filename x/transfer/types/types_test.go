package types

import "testing"

func TestParseCoinValid(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		amount  int64
		denom   string
		wantErr bool
	}{
		{name: "whole amount", in: "1000ubor", amount: 1000, denom: "ubor"},
		{name: "small amount", in: "12ubor", amount: 12, denom: "ubor"},
		{name: "zero amount", in: "0ubor", amount: 0, denom: "ubor"},
		{name: "empty string", in: "", wantErr: true},
		{name: "denom only", in: "ubor", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseCoin(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ParseCoin(%q): want error, got %+v", tc.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseCoin(%q): unexpected error: %v", tc.in, err)
			}
			if got.Amount != tc.amount || got.Denom != tc.denom {
				t.Fatalf("ParseCoin(%q) = %+v, want %d%s", tc.in, got, tc.amount, tc.denom)
			}
		})
	}
}

// TestParseCoinAcceptsNegative pins the current behaviour: ParseCoin does not
// reject a negative amount, so "-50ubor" parses cleanly.
func TestParseCoinAcceptsNegative(t *testing.T) {
	c, err := ParseCoin("-50ubor")
	if err != nil {
		t.Fatalf("ParseCoin(-50ubor): unexpected error: %v", err)
	}
	if c.Amount != -50 {
		t.Fatalf("ParseCoin(-50ubor).Amount = %d, want -50", c.Amount)
	}
}

func TestCoinString(t *testing.T) {
	if got := (Coin{Denom: "ubor", Amount: 42}).String(); got != "42ubor" {
		t.Fatalf("Coin.String() = %q, want %q", got, "42ubor")
	}
}
