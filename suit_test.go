package mongobuf_test

import (
	"github.com/d7561985/mongo-buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MongoSuit struct {
	suite.Suite
	db    *mongobuf.Mongodb
	Store map[string]primitive.ObjectID
}

func (m *MongoSuit) Initialize() error {
	m.Store = make(map[string]primitive.ObjectID)

	v, err := mongobuf.NewMongo(&pop.ConnectionDetails{Host: "127.0.0.1", Database: "testing"})
	if err != nil {
		return err
	}
	m.db = v
	return nil
}

// SetupTest call before every test
func (m *MongoSuit) SetupTest() {

}

func TestMongodb(t *testing.T) {
	s := new(MongoSuit)
	err := s.Initialize()
	assert.NoError(t, err)
	suite.Run(t, s)
}
