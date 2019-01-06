package mongobuf

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/stretchr/testify/assert"

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
