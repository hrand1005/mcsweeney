package content_test

import (
	"mcsweeney/content"
	"reflect"
	"testing"
)

// TestNewGetter calls the NewGetter function, providing a platform type,
// credentials, and query, checking for valid return values.
func TestNewGetter(t *testing.T) {
	type args struct {
		credentials string
		platform    content.Platform
		query       content.Query
	}
	tests := []struct {
		name     string
		args     args
		wantType interface{}
		wantErr  error
	}{
		{
			name: "NewGetter for twitch platform, valid credentials and query.",
			args: args{
				platform:    content.TWITCH,
				credentials: "fake_credentials.yaml",
				query:       content.Query{},
			},
			wantType: reflect.TypeOf(&content.TwitchGetter{}),
			wantErr:  nil,
		},
	}
	for _, tc := range tests {
		got, gotErr := content.NewGetter(tc.args.platform, tc.args.credentials, tc.args.query)
		if reflect.TypeOf(got) != tc.wantType || gotErr != tc.wantErr {
			t.Fatalf("Got:\nType %v, Err %v\nWanted: Type %v, Err %v", reflect.TypeOf(got), gotErr, tc.wantType, tc.wantErr)
		}
	}

}
