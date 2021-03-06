package kv

import (
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	stateTrie "github.com/prysmaticlabs/prysm/beacon-chain/state"
)

// SaveUnaggregatedAttestation saves an unaggregated attestation in cache.
func (p *AttCaches) SaveUnaggregatedAttestation(att *ethpb.Attestation) error {
	if att == nil {
		return nil
	}
	if helpers.IsAggregated(att) {
		return errors.New("attestation is aggregated")
	}

	r, err := hashFn(att)
	if err != nil {
		return errors.Wrap(err, "could not tree hash attestation")
	}

	p.unAggregateAttLock.Lock()
	defer p.unAggregateAttLock.Unlock()
	p.unAggregatedAtt[r] = stateTrie.CopyAttestation(att) // Copied.

	return nil
}

// SaveUnaggregatedAttestations saves a list of unaggregated attestations in cache.
func (p *AttCaches) SaveUnaggregatedAttestations(atts []*ethpb.Attestation) error {
	for _, att := range atts {
		if err := p.SaveUnaggregatedAttestation(att); err != nil {
			return err
		}
	}

	return nil
}

// UnaggregatedAttestations returns all the unaggregated attestations in cache.
func (p *AttCaches) UnaggregatedAttestations() []*ethpb.Attestation {
	atts := make([]*ethpb.Attestation, 0)

	p.unAggregateAttLock.RLock()
	defer p.unAggregateAttLock.RUnlock()
	for _, att := range p.unAggregatedAtt {
		atts = append(atts, stateTrie.CopyAttestation(att) /* Copied */)
	}

	return atts
}

// DeleteUnaggregatedAttestation deletes the unaggregated attestations in cache.
func (p *AttCaches) DeleteUnaggregatedAttestation(att *ethpb.Attestation) error {
	if att == nil {
		return nil
	}
	if helpers.IsAggregated(att) {
		return errors.New("attestation is aggregated")
	}

	r, err := hashFn(att)
	if err != nil {
		return errors.Wrap(err, "could not tree hash attestation")
	}

	p.unAggregateAttLock.Lock()
	defer p.unAggregateAttLock.Unlock()
	delete(p.unAggregatedAtt, r)

	return nil
}

// UnaggregatedAttestationCount returns the number of unaggregated attestations key in the pool.
func (p *AttCaches) UnaggregatedAttestationCount() int {
	p.unAggregateAttLock.RLock()
	defer p.unAggregateAttLock.RUnlock()
	return len(p.unAggregatedAtt)
}
