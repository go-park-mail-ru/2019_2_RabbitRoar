package session

import (
	"database/sql"
	"fmt"
	_ "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/config"
	_grpc "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/delivery/grpc"
	_sessionRepository "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/session/repository"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	consulapi "github.com/hashicorp/consul/api"
	"strconv"
	"time"
)

var log = logging.MustGetLogger("grpc_session")

func Start() {
	sessionAddress := viper.GetString("session.address")

	var err error
	config := consulapi.DefaultConfig()
	config.Address = *sessionAddress
	consul, err := consulapi.NewClient(config)

	serviceID := "SESSION_" + sessionAddress
	grpcPort := viper.GetString("session.port")

	err = consul.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "session-api",

		Port:    *grpcPort,
		Address: sessionAddress,
	})
	if err != nil {
		fmt.Println("cant add session service to consul", err)
		return
	}
	fmt.Println("registered in consul", serviceID)


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

func runOnlineServiceDiscovery(servers []string) {
	currAddrs := make(map[string]struct{}, len(servers))
	for _, addr := range servers {
		currAddrs[addr] = struct{}{}
	}
	ticker := time.Tick(5 * time.Second)
	for _ = range ticker {
		health, _, err := consul.Health().Service("session-api", "", false, nil)
		if err != nil {
			log.Fatalf("cant get alive services")
		}

		newAddrs := make(map[string]struct{}, len(health))
		for _, item := range health {
			addr := item.Service.Address +
				":" + strconv.Itoa(item.Service.Port)
			newAddrs[addr] = struct{}{}
		}

		var updates []*naming.Update
		// проверяем что удалилось
		for addr := range currAddrs {
			if _, exist := newAddrs[addr]; !exist {
				updates = append(updates, &naming.Update{
					Op:   naming.Delete,
					Addr: addr,
				})
				delete(currAddrs, addr)
				fmt.Println("remove", addr)
			}
		}
		// проверяем что добавилось
		for addr := range newAddrs {
			if _, exist := currAddrs[addr]; !exist {
				updates = append(updates, &naming.Update{
					Op:   naming.Add,
					Addr: addr,
				})
				currAddrs[addr] = struct{}{}
				fmt.Println("add", addr)
			}
		}
		if len(updates) > 0 {
			nameResolver.w.inject(updates)
		}
	}
}
