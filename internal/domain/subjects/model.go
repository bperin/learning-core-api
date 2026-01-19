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
