package functions

import (
	"context"
	"reflect"
	"testing"

	apiv1 "cloud.google.com/go/functions/apiv1"
	"github.com/elvenworks/functions-conector/domain"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		jsonCredentials []byte
		authScopes      []string
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Config
		wantErr bool
	}{
		{
			name: "Error project not found",
			args: args{
				jsonCredentials: []byte(`{"type": "service_account","project_id": "","private_key_id": "private_key_id","private_key": "private_key","client_email": "client_email","client_id": "client_id","auth_uri": "auth_uri","token_uri": "token_uri","auth_provider_x509_cert_url": "auth_provider_x509_cert_url","client_x509_cert_url": "client_x509_cert_url"}`),
				authScopes:      apiv1.DefaultAuthScopes(),
			},
			wantErr: true,
		},
		{
			name: "Error Unmarshal",
			args: args{
				jsonCredentials: []byte(`{"teste": }`),
				authScopes:      apiv1.DefaultAuthScopes(),
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				jsonCredentials: []byte(`{"type": "service_account","project_id": "project_id","private_key_id": "private_key_id","private_key": "private_key","client_email": "client_email","client_id": "client_id","auth_uri": "auth_uri","token_uri": "token_uri","auth_provider_x509_cert_url": "auth_provider_x509_cert_url","client_x509_cert_url": "client_x509_cert_url"}`),
				authScopes:      apiv1.DefaultAuthScopes(),
			},
			wantC: &Config{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if !tt.wantErr {
				creds, _ := google.CredentialsFromJSON(tt.wantC.Context, tt.args.jsonCredentials, tt.args.authScopes...)
				tt.wantC.Option = option.WithCredentials(creds)
			}

			gotC, err := NewConfig(tt.args.jsonCredentials)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
