package entity

type Profile struct {
	ID        string  `json:"id" validate:"required,omitempty"`
	FirstName string  `json:"firstName" validate:"omitempty,alpha"`
	LastName  string  `json:"lastName" validate:"omitempty,alpha"`
	Address   Address `json:"address" validate:"omitempty,dive"`
	Email     string  `json:"email" validate:"omitempty,email"`
}

type Address struct {
	Street  string `json:"street" validate:"omitempty"`
	City    string `json:"city" validate:"omitempty,alpha"`
	State   string `json:"state" validate:"omitempty,alpha"`
	ZipCode string `json:"zipCode" validate:"omitempty,numeric"`
}

type CreateProfileResult struct {
	ProfileID string
}
