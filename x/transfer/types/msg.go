package types

import "fmt"

// MsgSend moves Amount from FromAddress to ToAddress.
type MsgSend struct {
	FromAddress string
	ToAddress   string
	Amount      Coin
}

// ValidateBasic performs stateless validation of a MsgSend. Stateful checks,
// such as whether the sender can afford the transfer, belong to the keeper.
func (m MsgSend) ValidateBasic() error {
	if m.FromAddress == "" {
		return fmt.Errorf("from address cannot be empty")
	}
	if m.ToAddress == "" {
		return fmt.Errorf("to address cannot be empty")
	}
	if m.FromAddress == m.ToAddress {
		return fmt.Errorf("cannot send to self: %s", m.FromAddress)
	}
	if m.Amount.Denom == "" {
		return fmt.Errorf("amount denom cannot be empty")
	}
	return nil
}
