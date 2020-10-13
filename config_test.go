package mongoStore_test

import (
	"errors"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/assert"

	mongoStore "github.com/freemen-app/mongo-store"
)

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		host        string
		port        string
		database    string
		username    string
		password    string
		connTimeout time.Duration
		poolSize    uint64
	}
	tests := []struct {
		name       string
		fields     fields
		wantErrKey string
	}{
		{
			name: "valid",
			fields: fields{
				host:        "localhost",
				port:        "27017",
				database:    "test",
				username:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
		},
		{
			name: "required host",
			fields: fields{
				port:        "27017",
				database:    "test",
				username:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "host",
		},
		{
			name: "required port",
			fields: fields{
				host:     "localhost",
				database: "test",
				username: "test",

				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "port",
		},
		{
			name: "required database",
			fields: fields{
				host:        "localhost",
				port:        "27017",
				username:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "db",
		},
		{
			name: "required username",
			fields: fields{
				host:        "localhost",
				port:        "27017",
				database:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "username",
		},
		{
			name: "required password",
			fields: fields{
				host:        "localhost",
				port:        "27017",
				database:    "test",
				username:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "password",
		},
		{
			name: "required timeout",
			fields: fields{
				host:     "localhost",
				port:     "27017",
				database: "test",
				password: "test",
				poolSize: 1,
			},
			wantErrKey: "conn_timeout",
		},
		{
			name: "required pool size",
			fields: fields{
				host:        "localhost",
				port:        "27017",
				database:    "test",
				username:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
			},
			wantErrKey: "pool_size",
		},
		{
			name: "invalid host",
			fields: fields{
				host:        "test@gmail.com",
				port:        "27017",
				database:    "test",
				username:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "host",
		},
		{
			name: "invalid port",
			fields: fields{
				host:        "localhost",
				port:        "999999999",
				database:    "test",
				username:    "test",
				password:    "test",
				connTimeout: time.Nanosecond,
				poolSize:    1,
			},
			wantErrKey: "port",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := mongoStore.Config{
				Host:        tt.fields.host,
				Port:        tt.fields.port,
				DB:          tt.fields.database,
				Username:    tt.fields.username,
				Password:    tt.fields.password,
				ConnTimeout: tt.fields.connTimeout,
				PoolSize:    tt.fields.poolSize,
			}
			err := c.Validate()
			if tt.wantErrKey == "" {
				assert.Nil(t, err, err)
			} else {
				var validationErr validation.Errors
				assert.True(t, errors.As(err, &validationErr), err)
				assert.Contains(t, validationErr, tt.wantErrKey)
			}
		})
	}
}
