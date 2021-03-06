package blockchain

import (
	"context"
	"testing"

	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/cache"
	testDB "github.com/prysmaticlabs/prysm/beacon-chain/db/testing"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"
)

func TestHeadSlot_DataRace(t *testing.T) {
	db := testDB.SetupDB(t)
	defer testDB.TeardownDB(t, db)
	s := &Service{
		beaconDB: db,
	}
	go func() {
		if err := s.saveHead(context.Background(), [32]byte{}, ); err != nil {
			t.Fatal(err)
		}
	}()
	s.HeadSlot()
}

func TestHeadRoot_DataRace(t *testing.T) {
	db := testDB.SetupDB(t)
	defer testDB.TeardownDB(t, db)
	s := &Service{
		beaconDB: db,
		head:     &head{root: [32]byte{'A'}},
		stateGen: stategen.New(db, cache.NewStateSummaryCache()),
	}
	go func() {
		if err := s.saveHead(context.Background(), [32]byte{}, ); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := s.HeadRoot(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestHeadBlock_DataRace(t *testing.T) {
	db := testDB.SetupDB(t)
	defer testDB.TeardownDB(t, db)
	s := &Service{
		beaconDB: db,
		head:     &head{block: &ethpb.SignedBeaconBlock{}},
		stateGen: stategen.New(db, cache.NewStateSummaryCache()),
	}
	go func() {
		if err := s.saveHead(context.Background(), [32]byte{}, ); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := s.HeadBlock(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestHeadState_DataRace(t *testing.T) {
	db := testDB.SetupDB(t)
	defer testDB.TeardownDB(t, db)
	s := &Service{
		beaconDB: db,
		stateGen: stategen.New(db, cache.NewStateSummaryCache()),
	}
	go func() {
		if err := s.saveHead(context.Background(), [32]byte{}, ); err != nil {
			t.Fatal(err)
		}
	}()
	if _, err := s.HeadState(context.Background()); err != nil {
		t.Fatal(err)
	}
}
