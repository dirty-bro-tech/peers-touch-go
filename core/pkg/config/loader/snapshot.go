package loader

import "github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/source"

// Snapshot is a merged ChangeSet
type Snapshot struct {
	// The merged ChangeSet
	ChangeSet *source.ChangeSet
	// Deterministic and comparable version of the snapshot
	Version string
}

// Copy snapshot
func (s *Snapshot) Clone() *Snapshot {
	cs := *(s.ChangeSet)

	return &Snapshot{
		ChangeSet: &cs,
		Version:   s.Version,
	}
}
