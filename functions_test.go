package functions

import (
	"context"
	"reflect"
	"testing"
	"time"

	apiv1 "cloud.google.com/go/functions/apiv1"
	"github.com/elvenworks/functions-conector/domain"
	"github.com/elvenworks/functions-conector/internal/driver/functions"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func TestInitFunctions(t *testing.T) {
	type args struct {
		secret Secret
	}
	tests := []struct {
		name    string
		args    args
		wantF   *Functions
		wantErr bool
	}{
		{
			name: "Error loading config",
			args: args{
				secret: Secret{
					JsonCredentials: []byte(`{"type": "service_account","project_id": "","private_key_id": "private_key_id","private_key": "private_key","client_email": "client_email","client_id": "client_id","auth_uri": "auth_uri","token_uri": "token_uri","auth_provider_x509_cert_url": "auth_provider_x509_cert_url","client_x509_cert_url": "client_x509_cert_url"}`),
				},
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				secret: Secret{
					JsonCredentials: []byte(`{"type": "service_account","project_id": "project_id","private_key_id": "private_key_id","private_key": "private_key","client_email": "client_email","client_id": "client_id","auth_uri": "auth_uri","token_uri": "token_uri","auth_provider_x509_cert_url": "auth_provider_x509_cert_url","client_x509_cert_url": "client_x509_cert_url"}`),
				},
			},
			wantF: &Functions{
				&functions.Config{
					Context: context.Background(),
					Credentials: domain.Credentials{
						Type:                    "service_account",
						ProjectID:               "project_id",
						PrivateKeyID:            "private_key_id",
						PrivateKey:              "private_key",
						ClientEmail:             "client_email",
						ClientID:                "client_id",
						AuthURI:                 "auth_uri",
						TokenURI:                "token_uri",
						AuthProviderX509CertURL: "auth_provider_x509_cert_url",
						ClientX509CertURL:       "client_x509_cert_url",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				creds, _ := google.CredentialsFromJSON(tt.wantF.config.Context, tt.args.secret.JsonCredentials, apiv1.DefaultAuthScopes()...)
				tt.wantF.config.Option = option.WithCredentials(creds)
			}

			gotF, err := InitFunctions(tt.args.secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("InitFunctions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("InitFunctions() = %v, want %v", gotF.config, tt.wantF.config)
			}
		})
	}
}

func TestFunctions_GetLastFunctionsRun(t *testing.T) {
	type fields struct {
		config *functions.Config
	}
	type args struct {
		name             string
		validationString string
		seconds          time.Duration
		secret           Secret
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantLastRun *domain.FunctionsLastRun
		wantErr     bool
	}{
		{
			name: "Error NewClient",
			fields: fields{
				config: &functions.Config{},
			},
			args: args{
				name:             "function-1",
				validationString: "a",
				seconds:          time.Duration(60),
			},
			wantErr: true,
		},
		{
			name: "Success",
			fields: fields{
				config: &functions.Config{
					Context: context.Background(),
					Credentials: domain.Credentials{
						Type:                    "service_account",
						ProjectID:               "project_id",
						PrivateKeyID:            "private_key_id",
						PrivateKey:              "private_key",
						ClientEmail:             "client_email",
						ClientID:                "client_id",
						AuthURI:                 "auth_uri",
						TokenURI:                "token_uri",
						AuthProviderX509CertURL: "auth_provider_x509_cert_url",
						ClientX509CertURL:       "client_x509_cert_url",
					},
				},
			},
			args: args{
				name:             "function-1",
				validationString: "",
				seconds:          time.Duration(60),
				secret: Secret{
					JsonCredentials: []byte(`{"type": "service_account","project_id": "project_id","private_key_id": "private_key_id","private_key": "private_key","client_email": "client_email","client_id": "client_id","auth_uri": "auth_uri","token_uri": "token_uri","auth_provider_x509_cert_url": "auth_provider_x509_cert_url","client_x509_cert_url": "client_x509_cert_url"}`),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			creds, _ := google.CredentialsFromJSON(tt.fields.config.Context, tt.args.secret.JsonCredentials, apiv1.DefaultAuthScopes()...)
			tt.fields.config.Option = option.WithCredentials(creds)

			f := &Functions{
				config: tt.fields.config,
			}

			gotLastRun, err := f.GetLastFunctionsRun(tt.args.name, tt.args.validationString, tt.args.seconds)

			if (err != nil) != tt.wantErr {
				t.Errorf("Functions.GetLastFunctionsRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLastRun, tt.wantLastRun) {
				t.Errorf("Functions.GetLastFunctionsRun() = %v, want %v", gotLastRun, tt.wantLastRun)
			}
		})
	}
}
