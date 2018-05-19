package migration

import (
	"strconv"
	"strings"
	"testing"

	"github.com/toefel18/garsson-api/garsson/db/dbtest"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	v1CreateTestTable    = "CREATE TABLE testtable (huisnummer INTEGER, naam VARCHAR(16)"
	v2CreateAnotherTable = "CREATE TABLE anothertable (x INTEGER, y INTEGER)"
	v3CreateIndex        = "CREATE INDEX idx_naam ON testtable (naam)"
	ValidHashes          = true
	InvalidHashes        = false
)

var test_definition = []string{v1CreateTestTable, v2CreateAnotherTable, v3CreateIndex}

// TestMigrateFull tests a full migration
// no records exist yet in schema version
// v1, v2 and v3 are expected to be executed, each in a different transaction
func TestMigrateFull(t *testing.T) {
	dao, mock := dbtest.NewDbMock(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schemaversion .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(0, ValidHashes))
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE testtable .*").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO `schemaversion` .*" + hash(v1CreateTestTable) + ".*").WillReturnResult(sqlmock.NewErrorResult(nil))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("CREATE TABLE anothertable .*").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO `schemaversion` .*" + hash(v2CreateAnotherTable) + ".*").WillReturnResult(sqlmock.NewErrorResult(nil))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("CREATE INDEX idx_naam .*").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO `schemaversion` .*" + hash(v3CreateIndex) + ".*").WillReturnResult(sqlmock.NewErrorResult(nil))
	mock.ExpectCommit()
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, ValidHashes))

	if err := Migrate(dao.NewSession(), test_definition); err != nil {
		t.Error("Migrate failed with error: ", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

// TestMigratePartial tests that a partial migration works as expected
// v1CreateTestTable and v2CreateAnotherTable already exist the schemaversion table
// v3CreateIndex should be executed and a new record with hash should be written to schemaversion
func TestMigratePartial(t *testing.T) {
	dao, mock := dbtest.NewDbMock(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schemaversion .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(2, ValidHashes))
	mock.ExpectBegin()
	mock.ExpectExec("CREATE INDEX idx_naam .*").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO `schemaversion` .*" + hash(v3CreateIndex) + ".*").WillReturnResult(sqlmock.NewErrorResult(nil))
	mock.ExpectCommit()
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, ValidHashes))

	if err := Migrate(dao.NewSession(), test_definition); err != nil {
		t.Error("Migrate failed with error: ", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

// TestMigrateAlreadyOnVersion tests that a partial migration works as expected
// v1CreateTestTable, v2CreateAnotherTable and v3CreateIndex exist the schemaversion table
// nothing should be executed besides the creation and retrieval of database version
func TestMigrateAlreadyOnVersion(t *testing.T) {
	dao, mock := dbtest.NewDbMock(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schemaversion .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, ValidHashes))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, ValidHashes))

	if err := Migrate(dao.NewSession(), test_definition); err != nil {
		t.Error("Migrate failed with error: ", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

// TestMigrateDatabaseNewer tests that nothing happens if the database is on a newer version
func TestMigrateDatabaseNewer(t *testing.T) {
	dao, mock := dbtest.NewDbMock(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schemaversion .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, ValidHashes))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, ValidHashes))

	//note the slice taken from test_definition, were we remove the 3th item from the definition!
	if err := Migrate(dao.NewSession(), test_definition[:2]); err != nil {
		t.Error("Migrate failed with error: ", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

// TestMigrateHashesMismatch tests that an error is thrown if the database that exists contains version hashes that
// differ from the database definition, which indicates a mismatch in versions.
func TestMigrateHashesMismatch(t *testing.T) {
	dao, mock := dbtest.NewDbMock(t)

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS schemaversion .*").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`SELECT \* FROM schemaversion ORDER BY statement_id`).WillReturnRows(schemaVersionRowsToVersion(3, InvalidHashes))

	//note the slice taken from test_definition, were we remove the 3th item from the definition!
	if err := Migrate(dao.NewSession(), test_definition); err == nil {
		t.Error("Migrate didn't error, but it should because hashes were invalid!")
	} else if !strings.Contains(err.Error(), "integrity check failed") {
		t.Error("The error did not contain text 'integrity check failed' but gave: ", err.Error())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err.Error())
	}
}

func schemaVersionRowsToVersion(toVersion int, validHashes bool) sqlmock.Rows {
	resultSet := sqlmock.NewRows([]string{"statement_id", "statement_hash"})
	for i := 0; i < toVersion; i++ {
		if validHashes {
			resultSet.AddRow(strconv.Itoa(i), hash(test_definition[i]))
		} else {
			resultSet.AddRow(strconv.Itoa(i), "someinvalidhash")
		}
	}
	return resultSet
}

const (
	Same1 = `hello
my name		is earl`
	Same2     = "hello my name is 	earl"
	Same3     = "hellomynameisearl"
	Different = "hellomynameisearly"
)

func TestHashOf(t *testing.T) {
	if hash(Same1) != hash(Same2) {
		t.Error(Same1 + " hash is not equal to " + Same2)
	}
	if hash(Same2) != hash(Same3) {
		t.Error(Same2 + " hash is not equal to " + Same3)
	}
	if hash(Same3) == hash(Different) {
		t.Error(Same3 + " hash is equal to " + Different)
	}
}
