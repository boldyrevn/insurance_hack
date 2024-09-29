package db

import (
	"context"
	"fmt"
)

func (suite *TestSuite) TestGetUserByLogin() {
	user, err := suite.db.GetUserByLogin(context.Background(), "sussy")
	suite.Require().NoError(err)
	fmt.Println(user)
}
