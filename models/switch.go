package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

type Switch struct {
	ID          uuid.UUID `json:"-" db:"id"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
	Name        string    `json:"name" db:"name"`
	Location    string    `json:"location" db:"location"`
	IpAddr      string    `json:"ip_addr" db:"ip_addr"`
	ModelID     uuid.UUID `json:"-" db:"model_id"`
	CommunityID uuid.UUID `json:"-" db:"community_id"`
}

// String is not required by pop and may be deleted
func (s Switch) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Switches is not required by pop and may be deleted
type Switches []Switch

// String is not required by pop and may be deleted
func (s Switches) String() string {
	js, _ := json.Marshal(s)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (s *Switch) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: s.Name, Name: "Name"},
		&validators.StringIsPresent{Field: s.Location, Name: "Location"},
		&validators.StringIsPresent{Field: s.IpAddr, Name: "IpAddr"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (s *Switch) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (s *Switch) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
