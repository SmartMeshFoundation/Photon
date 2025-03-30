package restful

/*import (
	"context"
	"time"

	logx "log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectTo(uri, dbname string, timeout time.Duration) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	opt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		logx.Fatal(err)
		return nil, err
	}
	//check the health of server
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		logx.Fatal(err)
		return nil, err
	}
	return client.Database(dbname), nil
}

func Insert(dat interface{}) (lastid int64, err error) {
	objid, err := mongo.Collection.InsertOne(context.TODO(), &dat)
	if err != nil {
		logx.Println(err)
		return
	}
	return int64(objid), nil
}
*/
