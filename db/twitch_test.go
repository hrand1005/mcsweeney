package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gotest.tools/assert"
	"os"
	"testing"
)

func TestInsert(t *testing.T) {
	db, err := NewTwitchDB("test-insert.db")
	defer os.Remove("test-insert.db")
	assert.Assert(t, err == nil)
	err = db.Insert("www.herb.com")
	assert.Assert(t, err == nil)
	selectClip := `SELECT url FROM twitch WHERE url=?`

	testCases := []struct {
		name    string
		wantUrl string
		wantErr error
	}{
		{
			name:    "Success: inserted clip 'www.herb.com' retrieved",
			wantUrl: "www.herb.com",
			wantErr: nil,
		},
		{
			name:    "Failure: clip 'www.herb2.com' not found",
			wantUrl: "",
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		statement, err := db.dbHandle.Prepare(selectClip)
		assert.Assert(t, err == nil)

		var gotUrl string
		gotErr := statement.QueryRow(tc.wantUrl).Scan(&gotUrl)
		assert.Equal(t, tc.wantUrl, gotUrl)
		assert.Equal(t, tc.wantErr, gotErr)
	}
}
