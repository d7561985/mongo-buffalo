package mongobuf

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"testing"
)

func TestModel_FindObjectID(t *testing.T) {
	type A struct {
		id int
		V  primitive.ObjectID
	}

	in := &A{}
	mdl := M(in)
	mdl.UpdateObjectID()
	assert.NotEqualf(t, in.V, primitive.ObjectID{}, "")
}

func TestModel_GetObjectID(t *testing.T) {
	type A struct {
		id int
		V  primitive.ObjectID
	}

	in := &A{}
	mdl := M(in)
	mdl.UpdateObjectID()
	K := mdl.GetObjectID()
	assert.NotNil(t, K)
}
