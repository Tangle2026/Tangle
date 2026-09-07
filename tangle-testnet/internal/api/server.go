package api

import (
	"encoding/json"
	"github.com/ICFTProtocol/tangle-testnet/internal/chain"
	"github.com/ICFTProtocol/tangle-testnet/internal/store"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct{ E *chain.Engine }

func New(e *chain.Engine) *Server { return &Server{E: e} }
func jsonw(w http.ResponseWriter, v any) {
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(v)
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "frontend/index.html")
		return
	}
	if strings.HasPrefix(r.URL.Path, "/assets/") {
		http.StripPrefix("/assets/", http.FileServer(http.Dir("frontend/assets"))).ServeHTTP(w, r)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/api/status" {
		s.E.RLock()
		defer s.E.RUnlock()
		jsonw(w, map[string]any{"chain_id": s.E.S.State.ChainID, "height": s.E.S.State.Height, "supply": s.E.S.State.Supply, "accounts": len(s.E.S.State.Accounts), "vaults": len(s.E.S.State.Vaults), "fragments": len(s.E.S.State.Fragments), "transactions": len(s.E.S.State.Tx)})
		return
	}
	if r.Method == "GET" && r.URL.Path == "/api/blocks" {
		s.E.RLock()
		defer s.E.RUnlock()
		jsonw(w, s.E.S.State.Blocks)
		return
	}
	if r.Method == "GET" && r.URL.Path == "/api/vaults" {
		s.E.RLock()
		defer s.E.RUnlock()
		jsonw(w, s.E.S.State.Vaults)
		return
	}
	if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/account/") {
		addr := strings.TrimPrefix(r.URL.Path, "/api/account/")
		s.E.RLock()
		defer s.E.RUnlock()
		jsonw(w, s.E.S.State.Accounts[addr])
		return
	}
	if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/api/tx/") {
		id := strings.TrimPrefix(r.URL.Path, "/api/tx/")
		s.E.RLock()
		defer s.E.RUnlock()
		jsonw(w, s.E.S.State.Tx[id])
		return
	}
	if r.Method == "POST" && r.URL.Path == "/api/faucet" {
		var x struct {
			Address string `json:"address"`
			Amount  uint64 `json:"amount"`
		}
		json.NewDecoder(r.Body).Decode(&x)
		if x.Amount == 0 {
			x.Amount = 1_000_000_000
		}
		if x.Address == "" {
			http.Error(w, "address required", 400)
			return
		}
		if err := s.E.Faucet(x.Address, x.Amount); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonw(w, map[string]any{"ok": true, "address": x.Address, "amount": x.Amount})
		return
	}
	if r.Method == "POST" && r.URL.Path == "/api/vaults" {
		var x struct {
			Owner string `json:"owner"`
			Stake uint64 `json:"stake"`
		}
		json.NewDecoder(r.Body).Decode(&x)
		v, err := s.E.CreateVault(x.Owner, x.Stake)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonw(w, v)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/api/transactions" {
		var x struct {
			From, To, Asset string
			Amount          uint64
		}
		json.NewDecoder(r.Body).Decode(&x)
		tx, err := s.E.SubmitTx(x.From, x.To, x.Asset, x.Amount)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsonw(w, tx)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/api/mine" {
		s.E.Mine()
		jsonw(w, map[string]any{"ok": true, "time": time.Now().UTC()})
		return
	}
	if r.Method == "GET" && r.URL.Path == "/api/health" {
		jsonw(w, map[string]string{"status": "ok", "version": "0.1.0-alpha"})
		return
	}
	http.NotFound(w, r)
}
func Start(st *store.Store, addr string) {
	e := chain.New(st)
	srv := New(e)
	http.Handle("/", srv)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			e.Mine()
		}
	}()
	http.ListenAndServe(addr, nil)
}

var _ = strconv.Itoa
