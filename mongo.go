package mongoStore

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type (
	MongoClient interface {
		Ping(ctx context.Context, rp *readpref.ReadPref) error
		Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
	}

	Sessioner interface {
		StartSession(opts ...*options.SessionOptions) (mongo.Session, error)
	}

	Store interface {
		MongoClient
		Sessioner
		Start() error
		IsRunning() bool
		Shutdown()
	}

	store struct {
		*mongo.Client
		*mongo.Database
		isRunning bool
		conf      *Config
		options   *options.ClientOptions
	}
)

func New(conf *Config) Store {
	clientOptions := options.Client().
		ApplyURI(conf.DSN()).
		SetConnectTimeout(conf.ConnTimeout).
		SetMinPoolSize(conf.PoolSize)
	return &store{
		conf:    conf,
		options: clientOptions,
	}
}

func (s *store) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.conf.ConnTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, s.options)
	if err != nil {
		return err
	} else if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	s.isRunning = true
	s.Client = client
	s.Database = client.Database(s.conf.DB)

	return nil
}

func (s *store) IsRunning() bool {
	return s.isRunning
}

func (s *store) Shutdown() {
	if s.Client != nil {
		if err := s.Client.Disconnect(context.Background()); err != nil {
			log.Error().Err(err)
		}
	}
	s.isRunning = false
}
