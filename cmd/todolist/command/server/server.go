package server

import (
	"fmt"
	"os"
	"time"

	"github.com/odacremolbap/rest-demo/pkg/db"
	"github.com/odacremolbap/rest-demo/pkg/log"
	"github.com/odacremolbap/rest-demo/pkg/server"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	serverPort      int
	shutdownTimeout time.Duration

	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	dbSSL      bool
)

func init() {
	ServerCmd.PersistentFlags().IntVar(&serverPort, "port", 8080, "insecure listen port")
	ServerCmd.PersistentFlags().DurationVar(&shutdownTimeout, "shutdown-timeout", 10, "graceful shutdown timeout for API server")

	ServerCmd.PersistentFlags().StringVar(&dbHost, "db-host", "", "database host")
	ServerCmd.PersistentFlags().IntVar(&dbPort, "db-port", 5432, "database port")
	ServerCmd.PersistentFlags().StringVar(&dbUser, "db-user", "", "database user")
	ServerCmd.PersistentFlags().StringVar(&dbPassword, "db-password", "", "database password")
	ServerCmd.PersistentFlags().StringVar(&dbName, "db-name", "", "database name")
	ServerCmd.PersistentFlags().BoolVar(&dbSSL, "db-ssl", false, "set database SSL connection support")
}

// ServerCmd TODO list server command
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Simple TODO Server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("starting server")

		err := validate()
		if err != nil {
			log.Error(err, "")
			_ = cmd.Usage()
			os.Exit(-1)
		}

		connDB, err := db.ConnectPostgressDB(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSL)
		if err != nil {
			log.Error(err, "")
			os.Exit(-1)
		}
		db.Manager = db.NewTODOPersistenceManager(connDB)

		s := server.NewServer(serverPort, shutdownTimeout)

		log.Info(fmt.Sprintf("listening on port %d", serverPort))
		s.Run()
	},
}

// validate server flags
func validate() error {

	if len(dbHost) == 0 {
		return errors.New("a database host is needed")
	}

	if len(dbUser) == 0 {
		return errors.New("a database user is needed")
	}

	if len(dbName) == 0 {
		return errors.New("a database instance name is needed")
	}

	return nil
}
