package db

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestDB(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
