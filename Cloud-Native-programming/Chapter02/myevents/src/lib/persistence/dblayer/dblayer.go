package dblayer

import (
	"github.com/ImagineDevOps DevOps/Cloud-Native-programming-with-Golang/chapter03/myevents/src/lib/persistence"
	"github.com/ImagineDevOps DevOps/Cloud-Native-programming-with-Golang/chapter03/myevents/src/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB    DBTYPE = "mongodb"
	DOCUMENTDB DBTYPE = "documentdb"
	DYNAMODB   DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {

	switch options {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}
