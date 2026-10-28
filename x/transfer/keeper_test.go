package transfer_test

import (
	"testing"

	"github.com/devcon-dutralab/bor-lab/store"
	"github.com/devcon-dutralab/bor-lab/x/transfer"
	"github.com/devcon-dutralab/bor-lab/x/transfer/types"
)

func newKeeper() *transfer.Keeper {
	return transfer.NewKeeper(store.NewMemStore())
}

func TestSendMovesFunds(t *testing.T) {
	cases := []struct {
		name               string
		alice, bob, amount int64
		wantAlice, wantBob int64
	}{
		{name: "basic transfer", alice: 100, bob: 0, amount: 30, wantAlice: 70, wantBob: 30},
		{name: "full balance", alice: 50, bob: 50, amount: 50, wantAlice: 0, wantBob: 100},
		{name: "zero amount", alice: 5, bob: 5, amount: 0, wantAlice: 5, wantBob: 5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			k := newKeeper()
			k.SetBalance("alice", tc.alice)
			k.SetBalance("bob", tc.bob)

			err := k.Send(types.MsgSend{
				FromAddress: "alice",
				ToAddress:   "bob",
				Amount:      types.Coin{Denom: types.DefaultDenom, Amount: tc.amount},
			})
			if err != nil {
				t.Fatalf("Send: %v", err)
			}
			if got := k.Balance("alice"); got != tc.wantAlice {
				t.Errorf("alice balance = %d, want %d", got, tc.wantAlice)
			}
			if got := k.Balance("bob"); got != tc.wantBob {
				t.Errorf("bob balance = %d, want %d", got, tc.wantBob)
			}
		})
	}
}

func TestSendRejectsSelfTransfer(t *testing.T) {
	k := newKeeper()
	k.SetBalance("alice", 100)

	err := k.Send(types.MsgSend{
		FromAddress: "alice",
		ToAddress:   "alice",
		Amount:      types.Coin{Denom: types.DefaultDenom, Amount: 1},
	})
	if err == nil {
		t.Fatal("Send to self: want error, got nil")
	}
}

// TestSendAllowsOverdraft documents the current behaviour: Send does not check
// that the sender holds the amount, so an account can be driven negative.
func TestSendAllowsOverdraft(t *testing.T) {
	k := newKeeper()
	k.SetBalance("alice", 10)

	err := k.Send(types.MsgSend{
		FromAddress: "alice",
		ToAddress:   "bob",
		Amount:      types.Coin{Denom: types.DefaultDenom, Amount: 25},
	})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}
	if got := k.Balance("alice"); got != -15 {
		t.Fatalf("alice balance = %d, want -15", got)
	}
}
