package mongoStore_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"golang.org/x/xerrors"

	mongoStore "github.com/freemen-app/mongo-store"
)

var (
	conf *mongoStore.Config
)

func getWrappedErr(err error) error {
	for err != nil {
		if wrapper, ok := err.(xerrors.Wrapper); !ok {
			break
		} else {
			if wrapper.Unwrap() != nil {
				err = wrapper.Unwrap()
			} else {
				break
			}
		}
	}
	return err
}

func TestMain(m *testing.M) {
	conf = &mongoStore.Config{
		Host:        os.Getenv("MONGO_HOST"),
		Port:        os.Getenv("MONGO_PORT"),
		Username:    os.Getenv("MONGO_USERNAME"),
		Password:    os.Getenv("MONGO_PASSWORD"),
		DB:          os.Getenv("MONGO_DB"),
		ConnTimeout: time.Second / 2,
		PoolSize:    1,
	}
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	conf := &mongoStore.Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "test",
	}
	store := mongoStore.New(conf)
	assert.False(t, store.IsRunning())
}

func TestStore_Start(t *testing.T) {
	type args struct {
		conf *mongoStore.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "succeed",
			args: args{conf: conf},
		},
		{
			name: "timeout",
			args: args{
				conf: &mongoStore.Config{
					Host:        conf.Host,
					Port:        conf.Port,
					DB:          conf.DB,
					Username:    conf.Username,
					Password:    conf.Password,
					ConnTimeout: time.Nanosecond,
					PoolSize:    conf.PoolSize,
				},
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "invalid host/port",
			args: args{
				conf: &mongoStore.Config{
					Host:        conf.Host,
					Port:        "12345",
					DB:          conf.DB,
					Username:    conf.Username,
					Password:    conf.Password,
					ConnTimeout: conf.ConnTimeout,
					PoolSize:    conf.PoolSize,
				},
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "wrong auth",
			args: args{
				conf: &mongoStore.Config{
					Host:        conf.Host,
					Port:        conf.Port,
					DB:          conf.DB,
					Username:    conf.Username,
					Password:    "error",
					ConnTimeout: conf.ConnTimeout,
					PoolSize:    conf.PoolSize,
				},
			},
			wantErr: driver.Error{Code: 18, Name: "AuthenticationFailed", Message: "Authentication failed."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mongoStore.New(tt.args.conf)
			gotErr := store.Start()
			assert.EqualValues(t, tt.wantErr, getWrappedErr(gotErr))
			assert.EqualValues(t, store.IsRunning(), gotErr == nil)
			t.Cleanup(store.Shutdown)
		})
	}
}

func TestStore_Shutdown(t *testing.T) {
	store := mongoStore.New(conf)

	err := store.Start()
	assert.NoError(t, err)
	assert.True(t, store.IsRunning())

	store.Shutdown()
	assert.False(t, store.IsRunning())
}

func TestStore_Ping(t *testing.T) {
	store := mongoStore.New(conf)

	assert.NoError(t, store.Start())
	assert.True(t, store.IsRunning())

	assert.NoError(t, store.Ping(context.Background(), nil))

	t.Cleanup(store.Shutdown)
}
