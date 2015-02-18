package stack

import (
	"testing"

	. "gopkg.in/check.v1"
)

type StackSuite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&StackSuite{})

func (s *StackSuite) TestStackLength(c *C) {
	stack := Stack{}
	c.Assert(stack.Len(), Equals, 0)

	entry := struct {
		firstname string
		lastname  string
	}{"firstname", "lastname"}

	stack = append(stack, entry)

	c.Assert(stack.Len(), Equals, 1)
}

func (s *StackSuite) TestRegistrationsPush(c *C) {
	stack := Stack{}
	c.Assert(stack.Len(), Equals, 0)

	entry := struct {
		firstname string
		lastname  string
	}{"firstname", "lastname"}

	stack.Push(entry)

	c.Assert(stack.Len(), Equals, 1)
}
