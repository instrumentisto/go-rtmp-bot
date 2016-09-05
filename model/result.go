package model

// Scheme of status request result.
type Result struct {
	ID          uint8  // Status id
	Description string // String description.
}

// Returns new instance of result.
//
// param: id          uint8     Result status ID.
//        description string   Description of status.
func NewResult(id uint8, description string) *Result {
	return &Result{
		ID:          id,
		Description: description,
	}
}
