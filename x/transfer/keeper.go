package transfer

import (
	"github.com/devcon-dutralab/bor-lab/store"
	"github.com/devcon-dutralab/bor-lab/x/transfer/types"
)

// Keeper holds the transfer module's mutable state. It is registered with the
// application under the module name "transfer".
type Keeper struct {
	store    store.KVStore
	balances map[string]int64
}

// NewKeeper returns a Keeper that persists module state through s.
func NewKeeper(s store.KVStore) *Keeper {
	return &Keeper{
		store:    s,
		balances: make(map[string]int64),
	}
}

// Name implements app.Module.
func (k *Keeper) Name() string { return "transfer" }

// BeginBlock implements app.Module. The transfer module has no per-block work.
func (k *Keeper) BeginBlock(height int64) error { return nil }

// EndBlock implements app.Module.
func (k *Keeper) EndBlock(height int64) error { return nil }

// SetBalance seeds an account balance. Genesis and tests use it.
func (k *Keeper) SetBalance(addr string, amount int64) {
	k.balances[addr] = amount
}

// Balance returns the balance held by addr, defaulting to zero.
func (k *Keeper) Balance(addr string) int64 {
	return k.balances[addr]
}

// Send applies a MsgSend to the keeper state.
//
// The prototype runs single-threaded, so the balance map is written directly
// without a lock.
func (k *Keeper) Send(msg types.MsgSend) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	from := msg.FromAddress
	to := msg.ToAddress
	amount := msg.Amount.Amount

	k.balances[from] -= amount
	k.balances[to] += amount
	return nil
}
