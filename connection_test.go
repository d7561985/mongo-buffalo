package mongobuf_test

import (
	"context"
	"fmt"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"github.com/stretchr/testify/assert"
	"time"
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

func (m *MongoSuit) TestPing() {
	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	collection := m.db.Client.Database("testing").Collection("col1")
	_, err := collection.InsertOne(ctx, &A{A: "e", ID: primitive.NewObjectID()})
	assert.NoError(m.T(), err)

	err = m.db.Client.Ping(ctx, readpref.Primary())
	assert.NoError(m.T(), err)

	res := &A{}
	err = collection.FindOne(ctx, bson.M{"a": "e"}).Decode(res)
	assert.NoError(m.T(), err)

	all := []A{}
	cursor, err := collection.Find(ctx, bson.M{})
	assert.NoError(m.T(), err)
	for cursor.Next(ctx) {
		c := A{}
		err = cursor.Decode(&c)
		assert.NoError(m.T(), err)
		all = append(all, c)
	}
	fmt.Printf("%+[1]v %[1]v", all)
}

func (m *MongoSuit) TestCreate() {
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

	in.A = "Update"
	verr, err = m.db.Update(&in)
	assert.False(m.T(), verr.HasAny())
	assert.NoError(m.T(), err)
}
