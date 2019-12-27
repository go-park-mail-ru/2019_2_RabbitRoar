package balancer

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/op/go-logging"
	"strconv"
	"time"

	"google.golang.org/grpc/naming"
)

var log = logging.MustGetLogger("server")

type Watcher struct {
	update chan *naming.Update
	side chan int
	readDone chan int
}

func (w *Watcher) Next() (updates []*naming.Update, err error) {
	n := <-w.side

	if n == 0 {
		log.Error("w.side is closed")
		return nil, fmt.Errorf("w.side is closed")
	}

	for i := 0; i < n; i++ {
		u := <-w.update
		if u != nil {
			updates = append(updates, u)
		}
	}

	w.readDone <- 0

	return
}

func (w *Watcher) Close() {
	close(w.side)
}

func (w *Watcher) inject(updates []*naming.Update) {
	w.side <- len(updates)
	for _, u := range updates {
		w.update <- u
	}
	<-w.readDone
}

type NameResolver struct {
	w    *Watcher
	Addr string
}

func (r *NameResolver) Resolve(target string) (naming.Watcher, error) {
	r.w = &Watcher{
		update:   make(chan *naming.Update, 1),
		side:     make(chan int, 1),
		readDone: make(chan int),
	}
	r.w.side <- 1
	r.w.update <- &naming.Update{
		Op:   naming.Add,
		Addr: r.Addr,
	}
	go func() {
		<-r.w.readDone
	}()
	return r.w, nil
}

func RunOnlineSD(servers []string, nameResolver *NameResolver, consul *consulapi.Client) {
	currAddrs := make(map[string]struct{}, len(servers))
	for _, addr := range servers {
		currAddrs[addr] = struct{}{}
	}

	ticker := time.Tick(5 * time.Second)

	for range ticker {
		log.Info("Checking session service online.")
		health, _, err := consul.Health().Service("session-api", "", false, nil)
		if err != nil {
			log.Fatalf("Cant get alive services")
		}

		newAddrs := make(map[string]struct{}, len(health))
		for _, item := range health {
			addr := item.Service.Address +
				":" + strconv.Itoa(item.Service.Port)
			newAddrs[addr] = struct{}{}
		}
		log.Info("Online:", newAddrs)

		var updates []*naming.Update
		for addr := range currAddrs {
			if _, exist := newAddrs[addr]; !exist {
				updates = append(updates, &naming.Update{
					Op:   naming.Delete,
					Addr: addr,
				})
				delete(currAddrs, addr)
				log.Warning("Removing:", addr)
			}
		}

		for addr := range newAddrs {
			if _, exist := currAddrs[addr]; !exist {
				updates = append(updates, &naming.Update{
					Op:   naming.Add,
					Addr: addr,
				})
				currAddrs[addr] = struct{}{}
				log.Info("Adding: ", addr)
			}
		}

		if len(updates) > 0 {
			log.Info("Injecting update:", newAddrs)
			nameResolver.w.inject(updates)
		}
	}
}
