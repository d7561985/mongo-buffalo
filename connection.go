package mongobuf

import (
	"context"
	"fmt"
	"github.com/gobuffalo/fizz"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/kataras/iris/core/errors"
	"github.com/markbates/going/defaults"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
	"io"
	"sync"
	"time"
)

const nameMongoDB = "mongo"
const portMongo = "27017"

type Mongodb struct {
	mu                sync.Mutex
	ConnectionDetails *pop.ConnectionDetails
	Client            *mongo.Client
}

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

	res.Client = client

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return res, nil
}

func finalizerMongoDB(cd *pop.ConnectionDetails) {
	cd.Port = defaults.String(cd.Port, portMongo)
}

func (m *Mongodb) Name() string {
	return m.ConnectionDetails.Database
}

func (m *Mongodb) URL() string {
	return "mongodb://" + m.ConnectionDetails.Host + ":" + m.ConnectionDetails.Port
}

func (m *Mongodb) Details() *pop.ConnectionDetails {
	return m.ConnectionDetails
}

// Save model inside db.
func (m *Mongodb) Create(model interface{}) (*validate.Errors, error) {
	mdl := M(model)
	verr := mdl.Validate()
	if verr.HasAny() {
		return verr, nil
	}

	mdl.UpdateObjectID()

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	collection := m.Client.Database(m.Details().Database).Collection(mdl.TableName())
	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		return verr, err
	}

	return verr, nil
}

func (m *Mongodb) Update(model interface{}) (*validate.Errors, error) {
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

	collection := m.Client.Database(m.Details().Database).Collection(mdl.TableName())
	result, err := collection.ReplaceOne(ctx, filter, model)
	if err != nil {
		return verr, err
	}

	// by the good this is never happen
	if result.ModifiedCount != 1 {
		return verr, fmt.Errorf("expect 1 update but got %d", result.ModifiedCount)
	}

	return verr, nil
}

func (Mongodb) TranslateSQL(string) string {
	panic("implement me")
}

func (Mongodb) CreateDB() error {
	panic("implement me")
}

func (Mongodb) DropDB() error {
	panic("implement me")
}

func (Mongodb) DumpSchema(io.Writer) error {
	panic("implement me")
}

func (Mongodb) LoadSchema(io.Reader) error {
	panic("implement me")
}

func (Mongodb) FizzTranslator() fizz.Translator {
	panic("implement me")
}

func (Mongodb) Lock(func() error) error {
	panic("implement me")
}

func (Mongodb) TruncateAll(*pop.Connection) error {
	panic("implement me")
}

func (Mongodb) AfterOpen(*pop.Connection) error {
	panic("implement me")
}
