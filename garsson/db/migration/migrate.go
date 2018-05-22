// Database migrations help evolve the database schema automatically over time. Frameworks like flyway handle this really well
// for Java projects. We could still use it, but that would require a JRE during deployment to migrate/patch the database.
//
// This migration package takes a simplified approach by embedding all the statements required to build the database inside
// the application itself (see database_definition.go). Each time the application is started, it first runs the Migrate function
// which brings the database to the most recent version. If the database is already on that version, nothing changes. If
// the database is outdated, only the missing statements are executed to bring it up to date. The current schema version
// is tracked a separate table 'SchemaVersion'. That table contains the ID (=index) and hash of each executed statement.
// The highest ID is the version of the schema.
package migration

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gocraft/dbr"
	"github.com/toefel18/garsson-api/garsson/log"

)

// SchemaVersionDDL contains the DDL for the SchemaVersion table
const schemaVersionDDL = `CREATE TABLE IF NOT EXISTS schemaversion 
						 (statement_id INTEGER PRIMARY KEY, statement_hash VARCHAR(64), insert_time VARCHAR(64), statement TEXT)`

// SchemaVersion is a record containing a database version. StatementID is the version, which
// is just an incremented index, and StatementHash is the SHA256 of the executed query. If the
// migration finds that its statement hash differs from the loaded version, Migrate should stop error.
type SchemaVersion struct {
	StatementID   int
	StatementHash string
	InsertTime    string
	Statement     string
}

// Migrate brings the database to the latest version based on the statements defined in database_definition.go.
func MigrateDatabase(sess *dbr.Session) error {
	return Migrate(sess, completeDatabaseDefinition)
}

// Migrate brings the database to the latest version based on the statements defined in dbDefinition. See package
// documentation for a more complete description of the algorithm.
func Migrate(sess *dbr.Session, dbDefinition []string) error {
	log.WithField("totalMigrations", len(dbDefinition)).Info("migrating database")
	if _, err := sess.Exec(schemaVersionDDL); err != nil {
		return err
	}
	if initialDbVersions, err := FetchDbVersion(sess); err != nil {
		return err
	} else {
		return migrate(sess, initialDbVersions, dbDefinition)
	}
}

func FetchDbVersion(sess *dbr.Session) ([]SchemaVersion, error) {
	var dbVersions = []SchemaVersion{}
	if _, err := sess.Select("*").From("schemaversion").OrderBy("statement_id").Load(&dbVersions); err != nil {
		log.WithError(err).Fatal("Error retrieving current schema version")
		return nil, err
	}
	return dbVersions, nil
}

// migrate analyzes the difference between the current database version and the versionQueries packaged in
// the application
func migrate(sess *dbr.Session, initialDbVersions []SchemaVersion, dbDefinition []string) error {
	for i, query := range dbDefinition {
		queryHash := hash(query)
		version := i + 1 // which is more human readable and matches the database_definition.go
		if len(initialDbVersions) > i && initialDbVersions[i].StatementHash != queryHash {
			return errors.New("Database integrity check failed, query at index index " + strconv.Itoa(i) + " has a different hash. query:" + query)
		} else if len(initialDbVersions) <= i {
			if err := updateToVersionInTx(sess, version, query, queryHash); err != nil {
				return err
			}
			log.WithField("version", version).WithField("hash", queryHash).Info("migrated")
		}
	}
	var err error
	if currentDbVersions, err := FetchDbVersion(sess); err == nil {
		logDatabaseVersion(initialDbVersions, currentDbVersions, dbDefinition)
	}
	return err
}

func logDatabaseVersion(initialDBVersions []SchemaVersion, currentDBVersions []SchemaVersion, dbDefinition []string) {
	if len(currentDBVersions) > len(dbDefinition) {
		log.WithFields(log.Fields{
			"initialDbVersion": len(initialDBVersions),
			"applicationVersion": len(dbDefinition),
			"currentDbVersion": len(currentDBVersions),
			"warn": "database is at newer version",
		}).Warn("database migration complete")
	} else {
		log.WithFields(log.Fields{
			"initialDbVersion": len(initialDBVersions),
			"applicationVersion": len(dbDefinition),
			"currentDbVersion": len(currentDBVersions),
		}).Info("database migration complete")
	}
}

// updateToVersionInTx executes one update statement one version statement and records it into the SchemaVersion
func updateToVersionInTx(session *dbr.Session, version int, query, queryHash string) error {
	tx, err := session.Begin()
	if err != nil {
		log.WithError(err).Fatal("Could not start transaction")
	}
	defer tx.RollbackUnlessCommitted()
	if _, err := tx.Exec(query); err != nil {
		log.WithField("version", version).WithField("statement", query).Fatal("Error migrating to version")
		return err
	}
	if _, err := tx.InsertInto("schemaversion").Columns("statement_id", "statement_hash", "insert_time", "statement").Record(SchemaVersion{version, queryHash, now(), query}).Exec(); err != nil {
		log.WithError(err).WithField("version", version).WithField("statement", query).Fatal("Error writing version to schema version table")
		return err
	}
	return tx.Commit()
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

// hash calculates the sha256 of the query, after stripping it from whitespace to avoid change in formatting.
func hash(query string) string {
	queryWithoutSpaces := strings.Replace(query, " ", "", -1)
	queryWithoutSpacesAndTabs := strings.Replace(queryWithoutSpaces, "\t", "", -1)
	queryWithoutWhitespace := strings.Replace(queryWithoutSpacesAndTabs, "\n", "", -1)
	sha := sha256.Sum256([]byte(queryWithoutWhitespace))
	return base64.URLEncoding.EncodeToString(sha[:])
}
