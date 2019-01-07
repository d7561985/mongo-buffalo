package mongobuf_test

import (
	"github.com/d7561985/mongo-buffalo"
	"github.com/d7561985/mongo-buffalo/config"
	"github.com/gobuffalo/envy"
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

	if err := config.LoadConfigFile(); err != nil {
		return err
	}

	env := envy.Get("GO_ENV", "test")
	v, err := mongobuf.NewMongo(config.ConnectionDetails[env])
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
