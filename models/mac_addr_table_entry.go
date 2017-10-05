package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

type MacAddrTableEntry struct {
	ID         uuid.UUID `json:"-" db:"id"`
	CreatedAt  time.Time `json:"-" db:"created_at"`
	UpdatedAt  time.Time `json:"-" db:"updated_at"`
	MacAddr    string    `json:"mac_addr" db:"mac_addr"`
	PortNumber int       `json:"port_number" db:"-"`
	PortID     uuid.UUID `json:"-" db:"port_id"`
	VlanID     uuid.UUID `json:"-" db:"vlan_id"`
	TableID    uuid.UUID `json:"-" db:"table_id"`
}

// String is not required by pop and may be deleted
func (m MacAddrTableEntry) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// MacAddrTableEntries is not required by pop and may be deleted
type MacAddrTableEntries []MacAddrTableEntry

// String is not required by pop and may be deleted
func (m MacAddrTableEntries) String() string {
	jm, _ := json.Marshal(m)
	return string(jm)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (m *MacAddrTableEntry) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (m *MacAddrTableEntry) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (m *MacAddrTableEntry) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
