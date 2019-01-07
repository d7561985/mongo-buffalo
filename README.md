# MongoBuffalo
Package for simple usage MongoDB database with buffalo framework

#### Current approach
Not support concurrent access => used Mutex for protect this.

#### Goal of this adapter 
Help to develop simple low/mid loaded systems!

## Usage
```bash
$ go get -u github.com/d7561985/mongo-buffalo
```

### Init client:
```
    "github.com/d7561985/mongo-buffalo"
    "github.com/gobuffalo/pop"
```
```
if err := config.LoadConfigFile(); err != nil{
    return err
}
env := envy.Get("GO_ENV", "development")
client, err := mongobuf.NewMongo(config.ConnectionDetails[env])
if err != nil {
    return err
}
```

### Validation
If u need validation(it's just optional, without this everything will work) for create and update  instance of `T` (structure)
Yout need to implement following interface with requirement of validate pkg of gobuffalo:
```
"github.com/gobuffalo/validate"
```
```
ValidateAble interface {
    Validate() *validate.Errors
}
``` 

This open a lot of flexible ways of usage. For example where is packege with prebuilded validators:

```
import(
    "github.com/gobuffalo/validate/validators"	
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
```

### Requirements to structure type

#### Update
For using update we should know ID of current object in MongoDB. One simple way follow it in out structure while create and update.

Because of it, u should have not closed var ( First letter is High ) of type primitive.ObjectID
Example:
```
type A struct {
	AnyName  primitive.ObjectID `bson:"_id", json:"-"`
	...
}
```

#### All
For getting list of all struct instances in db where is requirement of slice interface:
```
	// All represent slice
	All interface {
		// return empty instance of slice type
		// PTR only!!!!!
		T() interface{}

		// add to back
		Add(interface{}) error
	}
```

receiver of Add method should be link (ptr). Example:
```
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
```
or as pointer in slice:
```
type As []*A

func (As) T() interface{} {
	return &A{}
}

func (a *As) Add(in interface{}) error {
	link, ok := in.(*A)
	if !ok {
		return errors.New("bad cast")
	}
	*a = append(*a, link)
	return nil
}
```