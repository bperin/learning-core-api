package subjects

import (
	"time"

	"github.com/google/uuid"
)

type Subject struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	SubSubjects []SubSubject `json:"sub_subjects,omitempty"`
}

type SubSubject struct {
	ID        uuid.UUID `json:"id"`
	SubjectID uuid.UUID `json:"subject_id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubjectForSelection struct {
	ID          uuid.UUID  `json:"id"`
	DisplayName string     `json:"display_name"`
	FullName    string     `json:"full_name"`
	URL         string     `json:"url"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
}
