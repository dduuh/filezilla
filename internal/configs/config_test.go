package configs

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	type env struct {
		host string
		port string
		appEnv string
	}

	type args struct {
		path string
		env env
	}

	setEnv := func(env env){
		os.Setenv("HTTP_HOST", env.host)
		os.Setenv("APP_ENV", env.appEnv)
		os.Setenv("HTTP_PORT", env.port)
	}

	tests := []struct{
		name string
		args args
		want *Config
		wantErr bool
	} {
		{
			name: "OK",
			args: args{
				path: "fixtures",
				env: env{
					host: "localhost",
					appEnv: "local",
				},
			},
			want: &Config{
				HttpCfg: HTTPConfig{
					Host: "localhost",
					Port: "8082",
					MaxHeaderBytes: 1,
					ReadTimeout: 10 * time.Second,
					WriteTimeout: 10 * time.Second,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.args.env)

			got, err := Init(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Want: %+v\n: %+v\n", tt.wantErr, err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Got: %+v\nWant: %+v\n", got, tt.want)
			}
		})
	}
}
