package registrations

// Registration is a user registration, that is not submitted
type Registration struct {
	Firstname string
	Lastname  string
}

// Registrations is a list of Registrations
type Registrations []Registration

// Len of the registration stack
func (r Registrations) Len() int {
	return len(r)
}

// Push adds a new registration to the stack
func (r *Registrations) Push(registration Registration) {
	*r = append(*r, registration)
}
