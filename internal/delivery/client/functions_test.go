package client

import (
	"testing"
)

func TestClient_Close(t *testing.T) {
	type fields struct {
		client MockLogging
	}
	tests := []struct {
		name     string
		fields   fields
		errClose error
	}{
		{
			name: "Test Close",
			fields: fields{
				client: MockLogging{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.client.On("Close").Return(tt.errClose)

			c := &Client{
				client: tt.fields.client,
			}

			c.Close()
		})
	}
}
