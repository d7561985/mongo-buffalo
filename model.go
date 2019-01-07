package mongobuf

import (
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"reflect"
	"strings"
)

type (
	Model struct {
		pop.Model
	}

	ValidateAble interface {
		Validate() *validate.Errors
	}

	// All represent slice
	All interface {
		// return empty instance of slice type
		// PTR only!!!!!
		T() interface{}

		// add to back
		Add(interface{}) error
	}
)

func M(in interface{}) *Model {
	return &Model{Model: pop.Model{Value: in}}
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

// small helper for check All interface is it correct
func checkAll(in All) error {
	if reflect.Indirect(reflect.ValueOf(in)).Kind() != reflect.Slice {
		return fmt.Errorf("@in is not slice type (%s)", reflect.ValueOf(in).Kind())
	}

	rfl := reflect.ValueOf(in.T())
	if rfl.Kind() != reflect.Ptr {
		return fmt.Errorf("@in.T() is not ptr type (%s)", rfl.Kind().String())
	}
	return nil
}
