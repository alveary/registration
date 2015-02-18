package stack

// Stack is a list of Registrations
type Stack []interface{}

// Len of the registration stack
func (s Stack) Len() int {
	return len(s)
}

// Push adds a new registration to the stack
func (s *Stack) Push(entry interface{}) {
	*s = append(*s, entry)
}
