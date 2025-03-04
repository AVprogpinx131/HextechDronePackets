package repository

import (
    "database/sql"
    "os"
    "testing"
    _ "github.com/lib/pq"
	"hextech_interview_project/internal/testutils"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
    var cleanup func()
    testDB, cleanup = testutils.SetupTestDB()
    code := m.Run()
    cleanup()
    os.Exit(code)
}