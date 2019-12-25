package session

import (
	"database/sql"
	"fmt"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/config"
	_grpc "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/delivery/grpc"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	consulapi "github.com/hashicorp/consul/api"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"os"
)

var log = logging.MustGetLogger("grpc_session")

func Start() {
	var err error

	config := consulapi.DefaultConfig()
	config.Address = viper.GetString("consul.address")
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("Err initializing consulapi client:", err)
	}

	serviceAddress := os.Getenv("HOSTNAME")
	serviceID := "SESSION_" + serviceAddress

	grpcPort := viper.GetInt("session.port")

	err = consul.Agent().ServiceRegister(
		&consulapi.AgentServiceRegistration{
			ID:      serviceID,
			Name:    "session-api",
			Port:    grpcPort,
			Address: serviceAddress,
		},
	)

	if err != nil {
		log.Fatal("cant add session service to consul:", err)
	}
	log.Info("registered in consul with id:", serviceID)

	defer func() {
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			log.Error("Error Deristering service:", err)
		}
	} ()

	log.Info("Staring grpc session service.")

	lis, err := net.Listen("tcp", viper.GetString("session.address"))
	if err != nil {
		log.Fatal("cant listen port", err)
	}

	dbDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.pass"),
		viper.GetString("database.db"),
	)
	db, err := sql.Open(
		"postgres",
		dbDSN,
	)
	if err != nil {
		log.Fatal("error init db: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("error connecting db: ", err)
	}
	defer db.Close()

	sessionRepo := _sessionRepository.NewSqlSessionRepository(db)

	server := grpc.NewServer()

	_grpc.RegisterSessionServiceServer(server, _grpc.NewManager(sessionRepo))

	fmt.Println("starting grpc session server")
	log.Fatal(server.Serve(lis))
}
