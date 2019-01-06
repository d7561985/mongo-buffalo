package mongobuf_test

import (
	"fmt"
	"github.com/d7561985/mongo-buffalo"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MongoSuit struct {
	suite.Suite
	db *mongobuf.Mongodb
}

func (m *MongoSuit) SetupTest() {
	v, err := mongobuf.NewMongo(&pop.ConnectionDetails{Host: "127.0.0.1", Database: "testing"})
	assert.NoError(m.T(), err)
	m.db = v
	fmt.Println("Setup!")
}

func TestMongodb_AfterOpen(t *testing.T) {
	suite.Run(t, new(MongoSuit))
}
