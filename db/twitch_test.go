package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gotest.tools/assert"
	"os"
	"testing"
)

const (
	selectClip = `SELECT url FROM twitch WHERE url=?`
	testDB     = "test-twitch.db"
	testURL    = "www.herb.com"
)

func TestInsert(t *testing.T) {
	// init test db
	db, err := NewTwitchDB(testDB)
	defer os.Remove(testDB)
	assert.Assert(t, err == nil)

	// insert test url
	err = db.Insert(testURL)
	assert.Assert(t, err == nil)

	testCases := []struct {
		name    string
		wantURL string
		wantErr error
	}{
		{
			name:    "Success: inserted clip retrieved",
			wantURL: testURL,
			wantErr: nil,
		},
		{
			name:    "Failure: clip not found",
			wantURL: "",
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		statement, err := db.dbHandle.Prepare(selectClip)
		assert.Assert(t, err == nil)

		var gotURL string
		gotErr := statement.QueryRow(tc.wantURL).Scan(&gotURL)
		assert.Equal(t, tc.wantURL, gotURL)
		assert.Equal(t, tc.wantErr, gotErr)
	}
}

func TestExists(t *testing.T) {
	// init test db
	db, err := NewTwitchDB(testDB)
	defer os.Remove(testDB)
	assert.Assert(t, err == nil)

	//insert test url
	err = db.Insert(testURL)
	assert.Assert(t, err == nil)

	testCases := []struct {
		name       string
		url        string
		wantExists bool
		wantErr    error
	}{
		{
			name:       "Success: clip exists in db",
			url:        testURL,
			wantExists: true,
			wantErr:    nil,
		},
		{
			name:       "Success: clip doesn't exist in db",
			url:        "",
			wantExists: false,
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		gotExists, gotErr := db.Exists(tc.url)
		assert.Equal(t, tc.wantExists, gotExists)
		assert.Equal(t, tc.wantErr, gotErr)
	}
}
