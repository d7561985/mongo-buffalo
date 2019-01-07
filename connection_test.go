package mongobuf_test

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/icrowley/fake"
	"github.com/kataras/iris/core/errors"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/stretchr/testify/assert"
	"sync"
)

type A struct {
	ID         primitive.ObjectID `bson:"_id"`
	A          string
	AnnnaLirra int
}

func (a *A) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.A, Name: "A"},
	)
}

type As []A

func (As) T() interface{} {
	return &A{}
}

func (a *As) Add(in interface{}) error {
	link, ok := in.(*A)
	if !ok {
		return errors.New("bad cast")
	}
	*a = append(*a, *link)
	return nil
}

func (m *MongoSuit) TestPing() {
	err := m.db.Ping()
	assert.NoError(m.T(), err)
}

// create
func (m *MongoSuit) TestCRUD1() {
	in := A{}

	// first check validation
	verr, err := m.db.Create(&in)
	assert.True(m.T(), verr.HasAny())
	assert.NoError(m.T(), err)

	// for not validation should pass
	in.A = "not blank"
	verr, err = m.db.Create(&in)
	assert.False(m.T(), verr.HasAny())
	assert.NoError(m.T(), err)
	assert.NotEqual(m.T(), in.ID, primitive.NilObjectID)

	m.Store["CRUD"] = in.ID
}

// update
func (m *MongoSuit) TestCRUD2() {
	if !assert.NotEqual(m.T(), primitive.NilObjectID, m.Store["CRUD"]) {
		return
	}
	in := A{A: "Update", ID: m.Store["CRUD"]}

	verr, err := m.db.Update(&in)
	assert.False(m.T(), verr.HasAny())
	assert.NoError(m.T(), err)
}

func (m *MongoSuit) TestAll() {
	g := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		g.Add(1)
		go func() {
			verr, err := m.db.Create(&A{A: fake.Brand()})
			assert.NoError(m.T(), err)
			assert.False(m.T(), verr.HasAny())
			g.Done()
		}()
	}
	g.Wait()
	list := make(As, 0)
	err := m.db.All(&list, bson.M{})
	assert.NoError(m.T(), err)
	assert.NotEmpty(m.T(), list)

}
