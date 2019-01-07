package mongobuf

import (
	"context"
	"fmt"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/kataras/iris/core/errors"
	"github.com/markbates/going/defaults"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"sync"
	"time"
)

const nameMongoDB = "mongo"
const portMongo = "27017"

// Mongodb is helper for access to MongoDB.
// Current approach used here not support concurrent access => used Mutex for protect this.
// Goal of this adapter help to develop simple low/mid loaded systems.
type Mongodb struct {
	mu                sync.Mutex
	ConnectionDetails *pop.ConnectionDetails
	client            *mongo.Client
}

// NewMongo
func NewMongo(deets *pop.ConnectionDetails) (*Mongodb, error) {
	finalizerMongoDB(deets)

	deets.Dialect = nameMongoDB

	res := &Mongodb{ConnectionDetails: deets}
	client, err := mongo.NewClient(res.URL())
	if err != nil {
		return nil, err
	}

	// make long server connection
	if err = client.Connect(context.Background()); err != nil {
		return nil, err
	}

	res.client = client
	return res, res.Ping()
}

func finalizerMongoDB(cd *pop.ConnectionDetails) {
	cd.Port = defaults.String(cd.Port, portMongo)
}

func (m *Mongodb) Name() string {
	return m.ConnectionDetails.Database
}

func (m *Mongodb) URL() string {
	if m.ConnectionDetails.URL != "" {
		return m.ConnectionDetails.URL
	}
	return "mongodb://" + m.ConnectionDetails.Host + ":" + m.ConnectionDetails.Port
}

func (m *Mongodb) Details() *pop.ConnectionDetails {
	return m.ConnectionDetails
}

// GetCollection create collection object suitable for current instance of model.
func (m *Mongodb) GetCollection(mdl *Model) *mongo.Collection {
	return m.client.Database(m.Details().Database).Collection(mdl.TableName())
}

// Ping simple
func (m *Mongodb) Ping() error {
	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	return m.client.Ping(ctx, readpref.Primary())
}

// Save model inside db.
func (m *Mongodb) Create(model interface{}) (*validate.Errors, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mdl := M(model)
	verr := mdl.Validate()
	if verr.HasAny() {
		return verr, nil
	}

	mdl.UpdateObjectID()

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	_, err := m.GetCollection(mdl).InsertOne(ctx, model)
	if err != nil {
		return verr, err
	}

	return verr, nil
}

// Update update model with presented
func (m *Mongodb) Update(model interface{}) (*validate.Errors, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	mdl := M(model)

	verr := mdl.Validate()
	if verr.HasAny() {
		return verr, nil
	}

	filter := mdl.GetObjectID()
	if filter == nil {
		return verr, errors.New("model not contain primitive.ObjectID")
	}

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	result, err := m.GetCollection(mdl).ReplaceOne(ctx, filter, model)
	if err != nil {
		return verr, err
	}

	// by the good this is never happen
	if result.ModifiedCount != 1 {
		return verr, fmt.Errorf("expect 1 update but got %d", result.ModifiedCount)
	}

	return verr, nil
}

// Get make search one document corresponding to filter rules.
// @filter - simple map[string]interface{}, where string is in model var. name (key) and interface value of this key
func (m *Mongodb) Get(in interface{}, filter bson.M) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	mdl := M(in)

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	return m.GetCollection(mdl).FindOne(ctx, filter).Decode(in)
}

// All pollute @in parameter of All interface with data witch was founded.
// @filter - simple map[string]interface{}, where string is in model var. name (key) and interface value of this key
func (m *Mongodb) All(in All, filter bson.M) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := checkAll(in); err != nil {
		return err
	}

	mdl := M(in.T())

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	cursor, err := m.GetCollection(mdl).Find(ctx, filter)
	if err != nil {
		return err
	}

	for cursor.Next(ctx) {
		v := in.T()
		if err := cursor.Decode(v); err != nil {
			return err
		}

		if err := in.Add(v); err != nil {
			return err
		}
	}
	return err
}
