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
	"time"
)

var log = logging.MustGetLogger("grpc_session")

func Start() {
	time.Sleep(5 * time.Second)

	var err error

	config := consulapi.DefaultConfig()
	config.Address = viper.GetString("consul.address")
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("Err initializing consulapi client:", err)
	}

	serviceAddress := os.Getenv("HOSTNAME")
	servicePort := viper.GetInt("session.port")
	serviceBind := fmt.Sprintf("%s:%d", serviceAddress, servicePort)
	serviceID := "SESSION_" + serviceAddress

	err = consul.Agent().ServiceRegister(
		&consulapi.AgentServiceRegistration{
			Kind:    "",
			ID:      serviceID,
			Name:    "session-api",
			Tags:    nil,
			Port:    servicePort,
			Address: serviceAddress,
			Check: &consulapi.AgentServiceCheck{
				CheckID:                        "session",
				Name:                           "Session service heath status",
				Interval:                       "10s",
				TCP:                            serviceBind,
				DeregisterCriticalServiceAfter: "20s",
			},
		},
	)

	if err != nil {
		log.Fatal("cant add session service to consul:", err)
	}

	log.Infof(
		"registered session microservice %s:%d in consul with id: %s",
		serviceAddress,
		servicePort,
		serviceID,
	)

	defer func() {
		err := consul.Agent().ServiceDeregister(serviceID)
		if err != nil {
			log.Error("Error Deristering service:", err)
		}
	}()

	log.Info("Staring grpc session service.")

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", servicePort))
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
