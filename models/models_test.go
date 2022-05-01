package models_test

import (
	"testing"

	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	*suite.Model
}

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.NewBox("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}
