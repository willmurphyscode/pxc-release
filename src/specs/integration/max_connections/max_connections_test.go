package max_connections_test

import (
	"context"
	"database/sql"
	"log"

	helpers "specs/test_helpers"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MaxConnections", func() {
	var (
		db *sql.DB
	)

	BeforeEach(func() {
		var err error

		firstProxy, err := helpers.FirstProxyHost(helpers.BoshDeployment)
		Expect(err).NotTo(HaveOccurred())

		mysqlUsername := "root"

		mysqlPassword, err := helpers.GetMySQLAdminPassword()
		Expect(err).NotTo(HaveOccurred())

		db = helpers.DbConnWithUser(mysqlUsername, mysqlPassword, firstProxy)
	})

	It("supports client connections near max_connections", func() {
		currentConnectionCount, err := threadsConnected(db)
		Expect(err).NotTo(HaveOccurred())

		// TODO: Redeploy with specific max_connections
		maxConnections := 1500
		ctx := context.Background()
		var connections []*sql.Conn

		log.Printf("Found %d current connections", currentConnectionCount)
		log.Printf("Expecting to establishing %d total connections to proxy", maxConnections - currentConnectionCount)

		for i := 0; i < maxConnections - currentConnectionCount; i ++ {
			conn, err := db.Conn(ctx)
			Expect(err).NotTo(HaveOccurred())

			connections = append(connections, conn)
		}

		log.Printf("Established %d connections", len(connections))

		for _, conn := range connections {
			Expect(conn.Close()).To(Succeed())
		}
	})

})

func threadsConnected(db *sql.DB) (int, error) {
	var (
		unused           string
		threadsConnected int
	)

	err := db.QueryRow(`SHOW GLOBAL STATUS LIKE 'Threads_connected'`).
		Scan(&unused, &threadsConnected)
	if err != nil {
		return -1, err
	}

	return threadsConnected, nil
}
