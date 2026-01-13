package utils

import "learning-core-api/internal/persistance/store"

type GenerationType string

const GenerationTypeClassification GenerationType = "CLASSIFICATION"

func (g GenerationType) String() string {
	return string(g)
}

func (g GenerationType) DB() store.GenerationType {
	return store.GenerationType(g)
}
