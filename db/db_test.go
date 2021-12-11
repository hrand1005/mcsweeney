package db_test

import (
	"errors"
	"fmt"
	"gotest.tools/assert"
	"mcsweeney/db"
	"os"
	"reflect"
	"testing"
)

func TestNewContentDB(t *testing.T) {
	testCases := []struct {
		name    string
		source  string
		wantDB  db.ContentDB
		wantErr error
	}{
		{
			name:    "Success: New TwitchDB",
			source:  "twitch",
			wantDB:  new(db.TwitchDB),
			wantErr: nil,
		},
		{
			name:    "Failure: DB not found",
			source:  "not-implemented",
			wantDB:  nil,
			wantErr: fmt.Errorf("DB not-implemented not found"),
		},
	}

	for _, tc := range testCases {
		gotDB, gotErr := db.NewContentDB(tc.source, "test-new-content-db.db")
		defer os.Remove("test-new-content-db.db")
		assert.Equal(t, reflect.TypeOf(gotDB), reflect.TypeOf(tc.wantDB))
		assert.Equal(t, errors.Unwrap(gotErr), errors.Unwrap(tc.wantErr))
	}
}
