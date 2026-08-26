package storage

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
	"pharmacy-counter/internal/domain"
	"sync"
)

var buckets = []string{"patients", "prescriptions", "tickets", "dispenses", "events", "meta"}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	err = s.init()
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, n := range buckets {
			if _, e := tx.CreateBucketIfNotExists([]byte(n)); e != nil {
				return e
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func encode(v any) ([]byte, error)    { return json.Marshal(v) }
func decode(data []byte, v any) error { return json.Unmarshal(data, v) }
func (s *Store) put(bucket, key string, v any) error {
	data, e := encode(v)
	if e != nil {
		return e
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), data) })
}
func (s *Store) get(bucket, key string, v any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		d := b.Get([]byte(key))
		if d == nil {
			return fmt.Errorf("not found")
		}
		return decode(d, v)
	})
}
func (s *Store) list(bucket string, out func([]byte) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, v []byte) error {
			if v == nil {
				return nil
			}
			return out(v)
		})
	})
}
func (s *Store) SavePatient(v domain.Patient) error { return s.put("patients", v.ID, v) }
func (s *Store) GetPatient(id string) (domain.Patient, error) {
	var v domain.Patient
	e := s.get("patients", id, &v)
	return v, e
}
func (s *Store) ListPatients() ([]domain.Patient, error) {
	out := []domain.Patient{}
	e := s.list("patients", func(d []byte) error {
		var v domain.Patient
		if x := decode(d, &v); x != nil {
			return x
		}
		out = append(out, v)
		return nil
	})
	return out, e
}
func (s *Store) SavePrescription(v domain.Prescription) error { return s.put("prescriptions", v.ID, v) }
func (s *Store) GetPrescription(id string) (domain.Prescription, error) {
	var v domain.Prescription
	e := s.get("prescriptions", id, &v)
	return v, e
}
func (s *Store) ListPrescriptions() ([]domain.Prescription, error) {
	out := []domain.Prescription{}
	e := s.list("prescriptions", func(d []byte) error {
		var v domain.Prescription
		if x := decode(d, &v); x != nil {
			return x
		}
		out = append(out, v)
		return nil
	})
	return out, e
}
func (s *Store) SaveTicket(v domain.PickupTicket) error { return s.put("tickets", v.ID, v) }
func (s *Store) GetTicket(id string) (domain.PickupTicket, error) {
	var v domain.PickupTicket
	e := s.get("tickets", id, &v)
	return v, e
}
func (s *Store) ListTickets() ([]domain.PickupTicket, error) {
	out := []domain.PickupTicket{}
	e := s.list("tickets", func(d []byte) error {
		var v domain.PickupTicket
		if x := decode(d, &v); x != nil {
			return x
		}
		out = append(out, v)
		return nil
	})
	return out, e
}
func (s *Store) SaveDispense(v domain.DispenseRecord) error { return s.put("dispenses", v.ID, v) }
func (s *Store) GetDispense(id string) (domain.DispenseRecord, error) {
	var v domain.DispenseRecord
	e := s.get("dispenses", id, &v)
	return v, e
}
func (s *Store) ListDispenses() ([]domain.DispenseRecord, error) {
	out := []domain.DispenseRecord{}
	e := s.list("dispenses", func(d []byte) error {
		var v domain.DispenseRecord
		if x := decode(d, &v); x != nil {
			return x
		}
		out = append(out, v)
		return nil
	})
	return out, e
}
func (s *Store) SaveEvent(v domain.CounterEvent) error { return s.put("events", v.ID, v) }
func (s *Store) ListEvents() ([]domain.CounterEvent, error) {
	out := []domain.CounterEvent{}
	e := s.list("events", func(d []byte) error {
		var v domain.CounterEvent
		if x := decode(d, &v); x != nil {
			return x
		}
		out = append(out, v)
		return nil
	})
	return out, e
}
func (s *Store) SetMeta(k, v string) error { return s.put("meta", k, v) }
func (s *Store) GetMeta(k string) (string, error) {
	var v string
	e := s.get("meta", k, &v)
	return v, e
}
