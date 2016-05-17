package cassandra

import (
	"log"
	"time"

	"github.com/gocql/gocql"
)

const host string = "cassandra"

// Init will seed mock data into the database.
func Init() {
	log.Println("Seeding mock data to Cassandra")

	session := Get()
	defer session.Close()

	// Drop existing keyspace.
	err := session.Query(`DROP KEYSPACE IF EXISTS graphql`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Create new keyspace.
	err = session.Query(`CREATE KEYSPACE graphql WITH REPLICATION = {'class' : 'SimpleStrategy','replication_factor' : 1}`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Create new table.
	err = session.Query(`CREATE TABLE graphql.author (id uuid PRIMARY KEY, name varchar, firstName varchar)`).Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Insert some data...
	err = session.Query(`INSERT INTO graphql.author (id, name, firstName) VALUES (?, ?, ?)`, gocql.TimeUUID(), "Pike", "Rob").Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = session.Query(`INSERT INTO graphql.author (id, name, firstName) VALUES (?, ?, ?)`, gocql.TimeUUID(), "Gerrand", "Andrew").Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = session.Query(`INSERT INTO graphql.author (id, name, firstName) VALUES (?, ?, ?)`, gocql.TimeUUID(), "Cox", "Russ").Exec()
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
	cluster := gocql.NewCluster(host)
	// set ProtoVersion to use Cassandra 3.x
	cluster.ProtoVersion = 4
	// Cassandra keyspace and table has to be cerated before, e.g.:
	// 	CREATE KEYSPACE graphql WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
	// 	CREATE TABLE graphql.author (id uuid PRIMARY KEY, name varchar, firstName varchar);
	cluster.Keyspace = "system"
	cluster.Timeout = 20 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	return session
}
