package storage

import (
	"fmt"
	"go.etcd.io/bbolt"
)

func (s *Store) Health() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store unavailable")
	}
	return s.db.View(func(tx *bbolt.Tx) error { return nil })
}
func (s *Store) DeleteTicket(id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("tickets")).Delete([]byte(id)) })
}
