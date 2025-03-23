package company

import "github.com/google/uuid"

// Company represents company
type Company struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	AmountOfEmployees int       `json:"amount_of_employees"`
	Registered        bool      `json:"registered"`
	Type              string    `json:"type"`
}

// PartialCompany contains optional fields for a company update operation
type PartialCompany struct {
	Name              *string `json:"name,omitempty"`
	Description       *string `json:"description,omitempty"`
	AmountOfEmployees *int    `json:"amount_of_employees,omitempty"`
	Registered        *bool   `json:"registered,omitempty"`
	Type              *string `json:"type,omitempty"`
}
