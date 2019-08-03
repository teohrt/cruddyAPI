package entity

type Profile struct {
	ID        int
	FirstName string
	LastName  string
	Address   Address
	Email     string
}

type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
}
