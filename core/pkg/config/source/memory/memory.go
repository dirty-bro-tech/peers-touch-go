// Package memory is a memory source
package memory

import (
	"sync"
	"time"

	"github.com/dirty-bro-tech/peers-touch-go/core/option"
	"github.com/dirty-bro-tech/peers-touch-go/core/pkg/config/source"
	"github.com/google/uuid"
)

type memory struct {
	opts *source.Options

	sync.RWMutex
	ChangeSet *source.ChangeSet
	Watchers  map[string]*watcher
}

func (s *memory) Read() (*source.ChangeSet, error) {
	s.RLock()
	cs := &source.ChangeSet{
		Format:    s.ChangeSet.Format,
		Timestamp: s.ChangeSet.Timestamp,
		Data:      s.ChangeSet.Data,
		Checksum:  s.ChangeSet.Checksum,
		Source:    s.ChangeSet.Source,
	}
	s.RUnlock()
	return cs, nil
}

func (s *memory) Watch() (source.Watcher, error) {
	w := &watcher{
		Id:      uuid.New().String(),
		Updates: make(chan *source.ChangeSet, 100),
		Source:  s,
	}

	s.Lock()
	s.Watchers[w.Id] = w
	s.Unlock()
	return w, nil
}

// Update allows manual updates of the config data.
func (s *memory) Update(c *source.ChangeSet) {
	// don't process nil
	if c == nil {
		return
	}

	// hash the file
	s.Lock()
	// update changeset
	s.ChangeSet = &source.ChangeSet{
		Data:      c.Data,
		Format:    c.Format,
		Source:    "memory",
		Timestamp: time.Now(),
	}
	s.ChangeSet.Checksum = s.ChangeSet.Sum()

	// update watchers
	for _, w := range s.Watchers {
		select {
		case w.Updates <- s.ChangeSet:
		default:
		}
	}
	s.Unlock()
}

func (s *memory) String() string {
	return "memory"
}

func NewSource(opts ...option.Option) source.Source {
	s := &memory{
		opts:     source.GetOptions(option.GetOptions(opts...)),
		Watchers: make(map[string]*watcher),
	}

	if s.opts.Ctx() != nil {
		c, ok := s.opts.Ctx().Value(changeSetKey{}).(*source.ChangeSet)
		if ok {
			s.Update(c)
		}
	}

	return s
}
