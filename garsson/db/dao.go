package db

import (
	"github.com/gocraft/dbr"
	"time"
	log "github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
)

// Dao manages an underlying db connection
type Dao struct {
	db *dbr.Connection
}

// NewDaoWithConfiguredDb creates a new dao with an already configured connection
func NewDaoWithConfiguredDb(db *dbr.Connection) *Dao {
	return &Dao{db}
}

// NewDao creates a new dao configured with the connection string, allows
func NewDao(connectionString string) (*Dao, error) {
	if conn, err := dbr.Open("postgres", connectionString, nil); err != nil {
		return nil, err
	} else {
		conn.SetMaxOpenConns(0) //0 means unlimited
		return &Dao{conn}, nil
	}
}

// NewSession creates a new DBR session
func (dao *Dao) NewSession() *dbr.Session {
	return dao.db.NewSession(nil)
}

// WaitTillAvailable tries to connect to the database and retries until it's available with a back-off interval.
func (dao *Dao) WaitTillAvailable() {
	sess := dao.NewSession()
	intervals := []int{1, 2, 4, 8, 16, 32}
	try := 0
	for {
		if tx, err := sess.Begin(); err != nil {
			backoffTime := intervals[min(try, len(intervals)-1)]
			try += 1
			log.WithFields(log.Fields{"retry": try, "nextRetryInSeconds": backoffTime}).WithError(err).Error("Wait for database")
			time.Sleep(time.Duration(backoffTime) * time.Second)
		} else {
			tx.Rollback()
			log.Info("Connected to database")
			return
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

