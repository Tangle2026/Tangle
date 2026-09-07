package chain

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	tcrypto "github.com/ICFTProtocol/tangle-testnet/internal/crypto"
	"github.com/ICFTProtocol/tangle-testnet/internal/model"
	"github.com/ICFTProtocol/tangle-testnet/internal/store"
	"math"
	"time"
)

type Engine struct{ S *store.Store }

func New(s *store.Store) *Engine { return &Engine{S: s} }
func randInt(max int) int {
	if max <= 1 {
		return 0
	}
	var b [8]byte
	_, _ = rand.Read(b[:])
	return int(binary.LittleEndian.Uint64(b[:]) % uint64(max))
}
func targetFragments() int {
	return model.MinFragments + randInt(model.MaxFragments-model.MinFragments+1)
}
func fee(amount uint64) uint64 { return (amount * model.FeeBPS) / 10000 }
func (e *Engine) Faucet(addr string, amount uint64) error {
	e.S.Lock()
	defer e.S.Unlock()
	a := e.S.State.Accounts[addr]
	if a == nil {
		a = &model.Account{Address: addr}
		e.S.State.Accounts[addr] = a
	}
	if amount == 0 {
		amount = 1_000_000_000
	}
	if e.S.State.Accounts["faucet"].Balance < amount {
		return fmt.Errorf("faucet empty")
	}
	e.S.State.Accounts["faucet"].Balance -= amount
	a.Balance += amount
	return e.S.Save()
}
func (e *Engine) CreateVault(owner string, stake uint64) (*model.Vault, error) {
	e.S.Lock()
	defer e.S.Unlock()
	a := e.S.State.Accounts[owner]
	if a == nil || a.Balance < stake {
		return nil, fmt.Errorf("insufficient balance")
	}
	a.Balance -= stake
	id := tcrypto.Hash("vault", owner, e.S.State.Height, a.Nonce)
	v := &model.Vault{ID: id[:16], Owner: owner, Stake: stake, Active: true}
	e.S.State.Vaults[v.ID] = v
	return v, e.S.Save()
}
func (e *Engine) SubmitTx(from, to, asset string, amount uint64) (*model.Tx, error) {
	e.S.Lock()
	defer e.S.Unlock()
	a := e.S.State.Accounts[from]
	if a == nil || a.Balance < amount {
		return nil, fmt.Errorf("insufficient balance")
	}
	if asset != "TGL" {
		return nil, fmt.Errorf("alpha supports TGL only")
	}
	f := fee(amount)
	if a.Balance < amount+f {
		return nil, fmt.Errorf("insufficient balance for fee")
	}
	a.Balance -= amount + f
	e.S.State.Treasury += (f * 30) / 100
	e.S.State.ValidatorPool += (f * 30) / 100
	e.S.State.VaultPool += (f * 30) / 100
	e.S.State.Insurance += (f * 10) / 100
	e.S.State.Height++
	h := e.S.State.Height
	id := tcrypto.Hash("tx", from, to, amount, a.Nonce, h)
	tx := &model.Tx{ID: id[:24], From: from, To: to, Asset: asset, Amount: amount, Fee: f, Created: time.Now().UTC(), TargetFragments: targetFragments(), Status: "ROUTING"}
	e.S.State.Tx[tx.ID] = tx
	e.S.State.Pending = append(e.S.State.Pending, tx.ID)
	a.Nonce++
	return tx, e.S.Save()
}
func (e *Engine) Mine() {
	e.S.Lock()
	defer e.S.Unlock()
	if len(e.S.State.Pending) == 0 {
		return
	}
	height := e.S.State.Height + 1
	prev := "GENESIS"
	if n := len(e.S.State.Blocks); n > 0 {
		prev = e.S.State.Blocks[n-1].Hash
	}
	var txids, porids []string
	for _, id := range e.S.State.Pending {
		tx := e.S.State.Tx[id]
		e.routeTx(tx)
		txids = append(txids, id)
		for _, p := range e.S.State.PoRs {
			if p.Block == height {
				porids = append(porids, p.ID)
			}
		}
		if a := e.S.State.Accounts[tx.To]; a != nil {
			a.Balance += tx.Amount
		} else {
			e.S.State.Accounts[tx.To] = &model.Account{Address: tx.To, Balance: tx.Amount}
		}
		tx.Status = "SETTLED"
	}
	e.S.State.Pending = nil
	e.S.State.Height = height
	b := model.Block{Height: height, Time: time.Now().UTC(), PreviousHash: prev, TxIDs: txids, PoRIDs: porids}
	b.Hash = tcrypto.Hash(b.Height, b.Time, b.PreviousHash, b.TxIDs, b.PoRIDs)
	e.S.State.Blocks = append(e.S.State.Blocks, b)
	_ = e.S.Save()
}
func (e *Engine) routeTx(tx *model.Tx) {
	v := e.pickVault()
	if v == nil {
		return
	}
	current := &model.Fragment{ID: tcrypto.Hash("fragment", tx.ID)[:24], TxID: tx.ID, Parent: "", Commitment: tcrypto.Hash("commit", tx.ID, tx.Amount), State: "ACTIVE", CreatedBlock: e.S.State.Height, Expiry: e.S.State.Height + 100}
	e.S.State.Fragments[current.ID] = current
	queue := []*model.Fragment{current}
	final := 0
	for len(queue) > 0 && final < tx.TargetFragments {
		f := queue[0]
		queue = queue[1:]
		op := randInt(3)
		switch op {
		case 0:
			f.State = "ROUTED"
			f.NextHop = e.pickVaultID(v.ID)
			e.makePoR(tx, v, f, "PASS")
			queue = append(queue, f)
			final++
		case 1:
			children := 2 + randInt(model.MaxChildren-1)
			for i := 0; i < children; i++ {
				cid := tcrypto.Hash("child", f.ID, i)[:24]
				c := &model.Fragment{ID: cid, TxID: tx.ID, Parent: f.Commitment, Commitment: tcrypto.Hash("childcommit", f.Commitment, i), State: "ACTIVE", NextHop: e.pickVaultID(v.ID), CreatedBlock: e.S.State.Height, Expiry: e.S.State.Height + 100}
				e.S.State.Fragments[cid] = c
				e.makePoR(tx, v, c, "SPLIT")
				queue = append(queue, c)
			}
			f.State = "SPLIT"
		case 2:
			f.State = "POOL"
			e.makePoR(tx, v, f, "POOL")
			f.State = "ACTIVE"
			queue = append(queue, f)
		}
	}
	tx.FragmentIDs = tx.FragmentIDs[:0]
	for id, f := range e.S.State.Fragments {
		if f.TxID == tx.ID && f.State != "SPLIT" {
			tx.FragmentIDs = append(tx.FragmentIDs, id)
		}
	}
	e.distributeVaultRewards()
}
func (e *Engine) pickVault() *model.Vault {
	for _, v := range e.S.State.Vaults {
		if v.Active {
			return v
		}
	}
	return nil
}
func (e *Engine) pickVaultID(exclude string) string {
	for id, v := range e.S.State.Vaults {
		if v.Active && id != exclude {
			return id
		}
	}
	return exclude
}
func (e *Engine) makePoR(tx *model.Tx, v *model.Vault, f *model.Fragment, op string) {
	taskID := tcrypto.Hash("task", tx.ID, f.ID, v.ID, e.S.State.Height, op)[:24]
	task := &model.RoutingTask{ID: taskID, FragmentID: f.ID, VaultID: v.ID, Operation: op, NextHop: f.NextHop, Nonce: uint64(len(e.S.State.Tasks)), Expiry: f.Expiry}
	e.S.State.Tasks[task.ID] = task
	proof := tcrypto.Hash("PoR", task.ID, f.Commitment, f.NextHop, op)
	pid := tcrypto.Hash("por", proof)[:24]
	e.S.State.PoRs[pid] = &model.PoR{ID: pid, TaskID: task.ID, VaultID: v.ID, In: f.Commitment, Out: tcrypto.Hash(f.Commitment, op), Operation: op, Proof: proof, Valid: true, Block: e.S.State.Height}
	v.Work++
}
func (e *Engine) distributeVaultRewards() {
	pool := e.S.State.VaultPool
	if pool == 0 {
		return
	}
	var denom float64
	for _, v := range e.S.State.Vaults {
		denom += float64(v.Work) * math.Sqrt(float64(v.Stake+1))
	}
	if denom == 0 {
		return
	}
	for _, v := range e.S.State.Vaults {
		share := float64(v.Work) * math.Sqrt(float64(v.Stake+1)) / denom
		r := uint64(float64(pool) * share)
		v.Reward += r
		if a := e.S.State.Accounts[v.Owner]; a != nil {
			a.Balance += r
		}
	}
	e.S.State.VaultPool = 0
}
func (e *Engine) RLock()   { e.S.RLock() }
func (e *Engine) RUnlock() { e.S.RUnlock() }
func (e *Engine) Lock()    { e.S.Lock() }
func (e *Engine) Unlock()  { e.S.Unlock() }
