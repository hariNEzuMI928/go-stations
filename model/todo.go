package model

import "time"

type (
	// A TODO expresses ...
	TODO struct {
		ID          int       `json:"id"`
		Subject     string    `json:"subject"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	// A CreateTODORequest expresses ...
	CreateTODORequest struct{
		TODO `json:"todo"`
	}
	// A CreateTODOResponse expresses ...
	CreateTODOResponse struct{
		TODO `json:"todo"`
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct{TODO}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{TODO}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct{TODO}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct{TODO}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{TODO}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{TODO}
)
