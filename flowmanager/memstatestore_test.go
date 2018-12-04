package flowmanager

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MemStateStoreSuite struct {
	StateStoreSuite
}

func (s *MemStateStoreSuite) SetupTest() {
	s.stateStore = NewMemStateStore()
	s.Require().NotNil(s.stateStore)
}

func TestMemStateStore(t *testing.T) {
	suite.Run(t, new(MemStateStoreSuite))
}
