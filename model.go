package mongobuf

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"reflect"
	"strings"
)

type Model struct {
	pop.Model
}

func M(in interface{}) *Model {
	return &Model{Model: pop.Model{Value: in}}
}

type ValidateAble interface {
	Validate() *validate.Errors
}

func (m *Model) Validate() *validate.Errors {
	if x, ok := m.Value.(ValidateAble); ok {
		return x.Validate()
	}
	return validate.NewErrors()
}

// UpdateObjectID generate new BSON ObjectID for first found in struct
func (m *Model) UpdateObjectID() {
	v := reflect.Indirect(reflect.ValueOf(m.Value).Elem())
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}

		if _, ok := v.Field(i).Interface().(primitive.ObjectID); ok {
			v.Field(i).Set(reflect.ValueOf(primitive.NewObjectID()))
			break
		}
	}
}

// GetObjectID find bson ObjectID in structure and return bson.M
// for searching model linked to this ID
func (m *Model) GetObjectID() bson.M {
	v := reflect.Indirect(reflect.ValueOf(m.Value).Elem())
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}

		if f, ok := v.Field(i).Interface().(primitive.ObjectID); ok {
			if t, ok := v.Type().Field(i).Tag.Lookup("bson"); ok {
				return bson.M{t: f}
			}
			return bson.M{strings.ToLower(v.Type().Field(i).Name): f}
		}
	}
	return nil
}
