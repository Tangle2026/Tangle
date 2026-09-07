package model

import "time"

const (
	MaxSupply       uint64 = 10_000_000_000_000 // 10m TGL with 6 decimals
	FeeBPS          uint64 = 1                  // 0.01%
	VaultFeePct     uint64 = 30
	ValidatorFeePct uint64 = 30
	TreasuryFeePct  uint64 = 30
	InsuranceFeePct uint64 = 10
	MinFragments           = 100
	MaxFragments           = 1000
	MaxChildren            = 10
)

type Account struct {
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}
type Vault struct {
	ID     string `json:"vault_id"`
	Owner  string `json:"owner"`
	Stake  uint64 `json:"stake"`
	Work   uint64 `json:"work"`
	Reward uint64 `json:"reward"`
	Active bool   `json:"active"`
}
type Fragment struct {
	ID           string `json:"fragment_id"`
	TxID         string `json:"tx_id"`
	Parent       string `json:"parent_commitment"`
	Commitment   string `json:"value_commitment"`
	NextHop      string `json:"next_hop"`
	State        string `json:"state"`
	CreatedBlock uint64 `json:"created_block"`
	Expiry       uint64 `json:"expiry"`
}
type RoutingTask struct {
	ID         string `json:"task_id"`
	FragmentID string `json:"fragment_id"`
	VaultID    string `json:"vault_id"`
	Operation  string `json:"operation"`
	NextHop    string `json:"next_hop"`
	Nonce      uint64 `json:"nonce"`
	Expiry     uint64 `json:"expiry"`
}
type PoR struct {
	ID        string `json:"por_id"`
	TaskID    string `json:"task_id"`
	VaultID   string `json:"vault_id"`
	In        string `json:"c_in"`
	Out       string `json:"c_out"`
	Operation string `json:"operation"`
	Proof     string `json:"proof"`
	Valid     bool   `json:"valid"`
	Block     uint64 `json:"block"`
}
type Tx struct {
	ID              string    `json:"id"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Asset           string    `json:"asset"`
	Amount          uint64    `json:"amount"`
	Fee             uint64    `json:"fee"`
	Created         time.Time `json:"created"`
	TargetFragments int       `json:"target_fragments"`
	Status          string    `json:"status"`
	FragmentIDs     []string  `json:"fragment_ids"`
}
type Block struct {
	Height       uint64    `json:"height"`
	Time         time.Time `json:"time"`
	PreviousHash string    `json:"previous_hash"`
	Hash         string    `json:"hash"`
	TxIDs        []string  `json:"tx_ids"`
	PoRIDs       []string  `json:"por_ids"`
}
type State struct {
	ChainID       string                  `json:"chain_id"`
	Height        uint64                  `json:"height"`
	Supply        uint64                  `json:"supply"`
	Treasury      uint64                  `json:"treasury"`
	Insurance     uint64                  `json:"insurance"`
	ValidatorPool uint64                  `json:"validator_pool"`
	VaultPool     uint64                  `json:"vault_pool"`
	Accounts      map[string]*Account     `json:"accounts"`
	Vaults        map[string]*Vault       `json:"vaults"`
	Fragments     map[string]*Fragment    `json:"fragments"`
	Tasks         map[string]*RoutingTask `json:"tasks"`
	PoRs          map[string]*PoR         `json:"pors"`
	Tx            map[string]*Tx          `json:"tx"`
	Blocks        []Block                 `json:"blocks"`
	Nullifiers    map[string]bool         `json:"nullifiers"`
	Pending       []string                `json:"pending"`
}
