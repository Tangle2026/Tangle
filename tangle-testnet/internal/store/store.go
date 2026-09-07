package store

import (
	"encoding/json"
	"github.com/ICFTProtocol/tangle-testnet/internal/model"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	Path  string
	State model.State
}

func New(path, chainID string) (*Store, error) {
	s := &Store{Path: path}
	b, err := os.ReadFile(path)
	if err == nil {
		if err = json.Unmarshal(b, &s.State); err == nil {
			return s, nil
		}
	}
	s.State = model.State{ChainID: chainID, Supply: model.MaxSupply, Accounts: map[string]*model.Account{}, Vaults: map[string]*model.Vault{}, Fragments: map[string]*model.Fragment{}, Tasks: map[string]*model.RoutingTask{}, PoRs: map[string]*model.PoR{}, Tx: map[string]*model.Tx{}, Nullifiers: map[string]bool{}}
	s.State.Accounts["faucet"] = &model.Account{Address: "faucet", Balance: model.MaxSupply / 2}
	s.State.Supply = model.MaxSupply
	return s, s.Save()
}
func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0755); err != nil {
		return err
	}
	b, e := json.MarshalIndent(s.State, "", "  ")
	if e != nil {
		return e
	}
	tmp := s.Path + ".tmp"
	if e = os.WriteFile(tmp, b, 0644); e != nil {
		return e
	}
	return os.Rename(tmp, s.Path)
}
func (s *Store) RLock()   { s.mu.RLock() }
func (s *Store) RUnlock() { s.mu.RUnlock() }
func (s *Store) Lock()    { s.mu.Lock() }
func (s *Store) Unlock()  { s.mu.Unlock() }
