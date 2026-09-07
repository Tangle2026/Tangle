package chain

import (
 "testing"
 "github.com/ICFTProtocol/tangle-testnet/internal/model"
 "github.com/ICFTProtocol/tangle-testnet/internal/store"
)
func TestFee(t *testing.T){if fee(1_000_000)!=100{t.Fatalf("fee mismatch: %d",fee(1_000_000))}}
func TestTargetRange(t *testing.T){for i:=0;i<100;i++{x:=targetFragments();if x<model.MinFragments||x>model.MaxFragments{t.Fatalf("target %d",x)}}}
func TestTransferLifecycle(t *testing.T){s,_:=store.New(t.TempDir()+"/state.json","test");e:=New(s);if err:=e.Faucet("alice",1_000_000_000);err!=nil{t.Fatal(err)};if _,err:=e.CreateVault("alice",1000);err!=nil{t.Fatal(err)};tx,err:=e.SubmitTx("alice","bob","TGL",100_000);if err!=nil{t.Fatal(err)};e.Mine();if s.State.Tx[tx.ID].Status!="SETTLED"{t.Fatal("not settled")};if s.State.Accounts["bob"].Balance!=100_000{t.Fatal("recipient balance mismatch")}}
