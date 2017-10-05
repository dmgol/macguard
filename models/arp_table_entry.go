package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

type ArpTableEntry struct {
	ID        uuid.UUID `json:"-" db:"id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
	MacAddr   string    `json:"mac_addr" db:"mac_addr"`
	IpAddr    string    `json:"ip_addr" db:"ip_addr"`
	TableID   uuid.UUID `json:"-" db:"table_id"`
}

// String is not required by pop and may be deleted
func (a ArpTableEntry) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// ArpTableEntries is not required by pop and may be deleted
type ArpTableEntries []ArpTableEntry

// String is not required by pop and may be deleted
func (a ArpTableEntries) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *ArpTableEntry) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *ArpTableEntry) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *ArpTableEntry) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
