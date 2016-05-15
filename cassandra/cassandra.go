package cassandra

import (
	"log"

	"github.com/gocql/gocql"
)

const host string = "cassandra"

// Init will seed mock data into the database.
func Init() {
	log.Println("Seeding mock data to Cassandra")

	session := Get()
	defer session.Close()

	err := session.Query(`INSERT INTO author (id, name, firstName) VALUES (?, ?, ?)`, gocql.TimeUUID(), "Norris", "Chuck").Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = session.Query(`INSERT INTO author (id, name, firstName) VALUES (?, ?, ?)`, gocql.TimeUUID(), "MacGyver", "Angus").Exec()
	if err != nil {
		log.Fatal(err)
	}
}

// Cleanup will remove all mock data from the database.
func Cleanup() {
	log.Println("Cleaning up Cassandra...")
	session := Get()
	defer session.Close()

	err := session.Query(`TRUNCATE graphql.author;`).Exec()
	if err != nil {
		log.Fatal(err)
	}
}

// Get creates a new session to the CassandraHost and returns it.
func Get() *gocql.Session {
	// Cassandra keyspace and table has to be cerated before, e.g.:
	// 	CREATE KEYSPACE graphql WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
	// 	CREATE TABLE author (id uuid PRIMARY KEY, name varchar, firstName varchar);
	cluster := gocql.NewCluster(host)
	// set ProtoVersion to use Cassandra 3.x
	cluster.ProtoVersion = 4
	cluster.Keyspace = "graphql"

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	return session
}
