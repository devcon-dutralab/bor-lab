// Package app wires the simulation's modules into a single application object
// and drives the ABCI-style begin/end block lifecycle, mirroring the shape of a
// cosmos-sdk simapp without importing the SDK.
package app

import (
	"fmt"

	"github.com/devcon-dutralab/bor-lab/store"
	"github.com/devcon-dutralab/bor-lab/x/transfer"
)

// Module is the subset of the Cosmos SDK AppModule lifecycle this simulation
// models: a name and begin/end blockers.
type Module interface {
	Name() string
	BeginBlock(height int64) error
	EndBlock(height int64) error
}

// App is the root application container. It owns the store, the module
// registry, and the begin/end block loops.
type App struct {
	store   store.KVStore
	modules map[string]Module
	order   []string

	// Transfer is the registered transfer keeper, exposed for genesis and
	// command wiring.
	Transfer *transfer.Keeper
}

// NewApp builds a simulation app with the transfer module registered.
func NewApp() *App {
	s := store.NewMemStore()
	a := &App{
		store:   s,
		modules: make(map[string]Module),
	}
	a.Transfer = transfer.NewKeeper(s)
	a.Register(a.Transfer)
	return a
}

// Register adds a module to the app. Modules run in registration order.
func (a *App) Register(m Module) {
	if _, ok := a.modules[m.Name()]; !ok {
		a.order = append(a.order, m.Name())
	}
	a.modules[m.Name()] = m
}

// Module returns a registered module by name.
func (a *App) Module(name string) (Module, bool) {
	m, ok := a.modules[name]
	return m, ok
}

// BeginBlock runs every registered module's begin blocker.
func (a *App) BeginBlock(height int64) error {
	for _, name := range a.order {
		if err := a.modules[name].BeginBlock(height); err != nil {
			return fmt.Errorf("begin block %q: %w", name, err)
		}
	}
	return nil
}

// EndBlock runs every registered module's end blocker.
func (a *App) EndBlock(height int64) error {
	for _, name := range a.order {
		if err := a.modules[name].EndBlock(height); err != nil {
			return fmt.Errorf("end block %q: %w", name, err)
		}
	}
	return nil
}
