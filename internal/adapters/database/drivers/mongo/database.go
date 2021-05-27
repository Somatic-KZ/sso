package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/JetBrainer/sso/internal/adapters/database/drivers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 20 * time.Second
	retries           = 1
	CollectionRoles   = "roles"
)

type Mongo struct {
	MongoURL string
	client   *mongo.Client
	dbname   string

	DB      *mongo.Database
	Context context.Context

	rolesRepository *RolesRepository
	retries         int

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration
}

func (m *Mongo) Users(ctx context.Context, paging *interface{}, filters *interface{}) ([]interface{}, error) {
	panic("implement me")
}

func (m *Mongo) Name() string { return "Mongo" }

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	return &Mongo{
		MongoURL:          conf.URL,
		dbname:            conf.DataBaseName,
		retries:           retries,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}

func (m *Mongo) Connect() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	fmt.Printf("Connecting to %s\n", m.dbname)
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.MongoURL))
	if err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbname)

	return nil
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close() error {
	return m.client.Disconnect(m.Context)
}

func (m *Mongo) Roles() drivers.RolesRepository {
	if m.rolesRepository == nil {
		m.rolesRepository = &RolesRepository{
			collection: m.DB.Collection(CollectionRoles),
		}
	}

	return m.rolesRepository
}

// убеждается что все индексы построены
func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	if err := m.ensureUsersIndexes(ctx); err != nil {
		return err
	}

	return nil
}

// ensureUsersIndexes строит индексы для коллекции users
func (m *Mongo) ensureUsersIndexes(ctx context.Context) error {
	col := m.DB.Collection(CollectionUsers)

	models := []mongo.IndexModel{
		{Keys: bson.M{"email": 1}, Options: options.Index().SetUnique(true)},
		{Keys: bson.M{"phone": 1}, Options: options.Index().SetUnique(true)},
	}

	var err error
	var exists bool
	// проверяем существование индекса с именем restore_token_send_status
	if exists, err = m.indexExistsByName(ctx, col, "restore_token_send_status"); err != nil {
		return err
	}
	if !exists {
		idx := mongo.IndexModel{Keys: bson.M{"restore.token": 1, "restore.send": 1, "restore.status": 1}, Options: options.Index().SetName("restore_token_send_status")}
		models = append(models, idx)
	}

	if exists, err = m.indexExistsByName(ctx, col, "login"); err != nil {
		return err
	}
	if !exists {
		idx := mongo.IndexModel{Keys: bson.M{"login": 1}, Options: options.Index().SetName("login")}
		models = append(models, idx)
	}

	// проверяем существование индекса с именем verify_token_send_status
	if exists, err = m.indexExistsByName(ctx, col, "verify_token_send_status"); err != nil {
		return err
	}
	if !exists {
		idx := mongo.IndexModel{Keys: bson.M{"verify.token": 1, "verify.send": 1, "verify.status": 1}, Options: options.Index().SetName("verify_token_send_status")}
		models = append(models, idx)
	}

	opts := options.CreateIndexes().SetMaxTime(m.ensureIdxTimeout)
	_, err = col.Indexes().CreateMany(ctx, models, opts)

	return err
}

func (m *Mongo) ensureRolesIndexes(ctx context.Context) error {
	col := m.DB.Collection(CollectionRoles)

	models := []mongo.IndexModel{
		{Keys: bson.M{"name": 1}, Options: options.Index().SetUnique(true)},
	}

	opts := options.CreateIndexes().SetMaxTime(m.ensureIdxTimeout)
	_, err := col.Indexes().CreateMany(ctx, models, opts)

	return err
}

// indexExistsByName проверяет существование индекса с именем name.
func (m *Mongo) indexExistsByName(ctx context.Context, collection *mongo.Collection, name string) (bool, error) {
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	for cur.Next(ctx) {
		if name == cur.Current.Lookup("name").StringValue() {
			return true, nil
		}
	}

	return false, nil
}
