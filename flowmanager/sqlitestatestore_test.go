package flowmanager

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type SQLiteStateStoreSuite struct {
	StateStoreSuite

	tmpDir string
}

func (s *SQLiteStateStoreSuite) SetupTest() {
	s.StateStoreSuite.SetupTest()

	var err error
	s.tmpDir, err = ioutil.TempDir("", "test")
	s.Require().NoError(err)

	s.stateStore, err = NewSQLiteStateStore(path.Join(s.tmpDir, "test.db"))
	s.Require().NoError(err)
}

func (s *SQLiteStateStoreSuite) TearDownTest() {
	s.NoError(s.stateStore.(*SQLiteStateStore).Close())

	if s.T().Failed() {
		s.log.Info("Test failed. Skipping temp dir cleanup.", zap.String("tmpDir", s.tmpDir))
	} else {
		s.NoError(os.RemoveAll(s.tmpDir))
	}
}

func TestSQLiteStateStore(t *testing.T) {
	suite.Run(t, new(SQLiteStateStoreSuite))
}
