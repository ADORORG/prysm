package stategen

import (
	"context"
	"encoding/hex"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

// This saves a pre finalized beacon state in the cold section of the DB. The returns an error
// and not store anything if the state does not lie on an archive point boundary.
func (s *State) saveColdState(ctx context.Context, blockRoot [32]byte, state *state.BeaconState) error {
	ctx, span := trace.StartSpan(ctx, "stateGen.saveColdState")
	defer span.End()

	if state.Slot()%s.slotsPerArchivedPoint != 0 {
		return errSlotNonArchivedPoint
	}

	if err := s.beaconDB.SaveState(ctx, state, blockRoot); err != nil {
		return err
	}
	archivedIndex := state.Slot() / s.slotsPerArchivedPoint
	if err := s.beaconDB.SaveArchivedPointRoot(ctx, blockRoot, archivedIndex); err != nil {
		return err
	}

	log.WithFields(logrus.Fields{
		"slot":      state.Slot(),
		"blockRoot": hex.EncodeToString(bytesutil.Trunc(blockRoot[:]))}).Info("Saved full state on archived point")

	return nil
}

// This loads the cold state by block root, it decides whether to load from archived point (faster) or
// somewhere between archived points (slower) because it requires replaying blocks.
// This method is more efficient than load cold state by slot.
func (s *State) loadColdStateByRoot(ctx context.Context, blockRoot [32]byte) (*state.BeaconState, error) {
	ctx, span := trace.StartSpan(ctx, "stateGen.loadColdStateByRoot")
	defer span.End()

	summary, err := s.stateSummary(ctx, blockRoot)
	if err != nil {
		return nil, errors.Wrap(err, "could not get state summary")
	}

	return s.ComputeStateUpToSlot(ctx, summary.Slot)
}

// This loads a cold state by slot where the slot lies between the archived point.
// This is a slower implementation given there's no root and slot is the only argument. It requires fetching
// all the blocks between the archival points.
func (s *State) loadColdIntermediateStateBySlot(ctx context.Context, slot uint64) (*state.BeaconState, error) {
	ctx, span := trace.StartSpan(ctx, "stateGen.loadColdIntermediateStateBySlot")
	defer span.End()

	return s.ComputeStateUpToSlot(ctx, slot)
}
