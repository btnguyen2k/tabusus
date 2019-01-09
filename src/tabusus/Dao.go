package tabusus

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"strings"
	"tabusus/utils"
	"time"
)

// Application defines app's attribute
type Application struct {
	Data map[string]interface{} // app's data
}

func NewApp(id string) *Application {
	app := &Application{
		Data: map[string]interface{}{},
	}
	now := time.Now()
	app.SetId(id).SetTimeCreated(now).SetTimeUpdated(now).SetStatus(0)
	return app
}

func NewAppFromJson(json bson.M) *Application {
	return &Application{Data: json}
}

const (
	attrId          = "id"
	attrStatus      = "status"
	attrDesc        = "description"
	attrRsaPubKey   = "rsa_pubkey"
	attrTimeCreated = "tc"
	attrTimeUpdated = "tu"
	tableApps       = "apps"
)

func (app *Application) ToJson() ([]byte, error) {
	return bson.MarshalExtJSON(app.Data, false, false)
}

func (app *Application) GetId() string {
	return app.Data[attrId].(string)
}

func (app *Application) SetId(value string) *Application {
	app.Data[attrId] = strings.ToLower(strings.TrimSpace(value))
	return app
}

func (app *Application) GetDescription() string {
	return app.Data[attrDesc].(string)
}

func (app *Application) SetDescription(value string) *Application {
	app.Data[attrDesc] = strings.TrimSpace(value)
	return app
}

func (app *Application) GetRsaPubKey() string {
	return app.Data[attrRsaPubKey].(string)
}

func (app *Application) SetRsaPubKey(value string) *Application {
	app.Data[attrRsaPubKey] = strings.TrimSpace(value)
	return app
}

func (app *Application) GetStatus() int32 {
	v, ok := utils.ToInt32(app.Data[attrStatus])
	if ok {
		return v
	} else {
		return 0
	}
}

func (app *Application) GetStatusStr() string {
	switch app.GetStatus() {
	case 0:
		return "Disabled"
	case 1:
		return "Enabled"
	default:
		return "Unknown"
	}
}

func (app *Application) SetStatus(value int32) *Application {
	app.Data[attrStatus] = value
	return app
}

func (app *Application) GetTimeCreated() *time.Time {
	v := app.Data[attrTimeCreated]
	switch v.(type) {
	case time.Time:
		t := v.(time.Time)
		return &t
	case *time.Time:
		return v.(*time.Time)
	case int64:
		t := time.Unix(0, v.(int64)*int64(time.Millisecond))
		return &t
	case uint64:
		t := time.Unix(0, int64(v.(uint64)*uint64(time.Millisecond)))
		return &t
	case int:
		t := time.Unix(0, int64(v.(int)*int(time.Millisecond)))
		return &t
	}
	return nil
}

func (app *Application) SetTimeCreated(value time.Time) *Application {
	app.Data[attrTimeCreated] = value.UnixNano() / 1000000
	return app
}

func (app *Application) GetTimeUpdated() *time.Time {
	v := app.Data[attrTimeUpdated]
	switch v.(type) {
	case time.Time:
		t := v.(time.Time)
		return &t
	case *time.Time:
		return v.(*time.Time)
	case int64:
		t := time.Unix(0, v.(int64)*int64(time.Millisecond))
		return &t
	case uint64:
		t := time.Unix(0, int64(v.(uint64)*uint64(time.Millisecond)))
		return &t
	}
	return nil
}

func (app *Application) SetTimeUpdated(value time.Time) *Application {
	app.Data[attrTimeUpdated] = value.UnixNano() / 1000000
	return app
}

func (app *Application) UrlEdit() string {
	return "/editApp/" + app.GetId()
}

func (app *Application) UrlDelete() string {
	return "/deleteApp/" + app.GetId()
}

/*----------------------------------------------------------------------*/

type ApplicationDao interface {
	List() []Application
	Delete(app *Application) error
	Get(string) (*Application, error)
	Save(app *Application) error
}

type MongoApplicationDao struct {
	url    string        // connection url
	db     string        // database name
	client *mongo.Client // client instance
}

func NewMongoApplicationDao(url, db string) ApplicationDao {
	m := &MongoApplicationDao{
		url: url,
		db:  db,
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	c, err := mongo.Connect(ctx, url)
	if err != nil {
		panic(err)
	}
	m.client = c
	return m
}

func (dao *MongoApplicationDao) List() []Application {
	collection := dao.client.Database(dao.db).Collection(tableApps)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.M{})
	defer cur.Close(ctx)
	if err != nil {
		log.Warn(err)
		return nil
	}
	var result []Application
	for cur.Next(ctx) {
		var row bson.M
		err := cur.Decode(&row)
		if err != nil {
			log.Error(err)
		} else {
			result = append(result, *NewAppFromJson(row))
		}
	}
	return result
}

func (dao *MongoApplicationDao) Delete(app *Application) error {
	collection := dao.client.Database(dao.db).Collection(tableApps)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := collection.DeleteOne(ctx, bson.M{attrId: app.GetId()})
	return err
}

func (dao *MongoApplicationDao) Get(id string) (*Application, error) {
	collection := dao.client.Database(dao.db).Collection(tableApps)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	dbResult := collection.FindOne(ctx, bson.M{attrId: strings.ToLower(strings.TrimSpace(id))})
	if dbResult.Err() != nil {
		log.Error(dbResult.Err())
		return nil, dbResult.Err()
	}
	var row bson.M
	err := dbResult.Decode(&row)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Error(err)
		return nil, err
	}
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return NewAppFromJson(row), nil
}

func (dao *MongoApplicationDao) Save(app *Application) error {
	json, err := app.ToJson()
	if err != nil {
		return err
	}
	var m bson.M
	err = bson.UnmarshalExtJSON(json, false, &m)
	if err != nil {
		return err
	}
	collection := dao.client.Database(dao.db).Collection(tableApps)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_ = collection.FindOneAndReplace(ctx, bson.M{attrId: app.GetId()}, m)
	return nil
}
