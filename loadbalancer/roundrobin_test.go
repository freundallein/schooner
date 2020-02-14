package loadbalancer

import (
	"github.com/freundallein/schooner/proxy"
	"net/url"
	"strings"
	"testing"
	"time"
)

type MockTarget struct {
	address     *url.URL
	isAvailable bool
	ping        bool
}

func (mt *MockTarget) IsAvailable() bool {
	return mt.isAvailable
}

func (mt *MockTarget) SetAvailable(status bool) {
	mt.isAvailable = status
}

func (mt *MockTarget) Address() *url.URL {
	return mt.address
}

func (mt *MockTarget) ReverseProxy() proxy.Proxy {
	addr, _ := url.Parse("http://localhost:8000/")
	proxy := proxy.New(proxy.DefaultStrategy, addr)
	return proxy
}

func (mt *MockTarget) LastSeen() int64 {
	return 1
}

func (mt *MockTarget) Ping() bool {
	return mt.ping
}

func TestAddTarget(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}

	trg, _ := NewTarget("http://testhost:8000")
	bckt.AddTarget(trg)
	if len(bckt.targets) != 1 {
		t.Error("Expected", 1, "got", len(bckt.targets))
	}
	if trg.IsAvailable() {
		t.Error("Expected", false, "got", trg.IsAvailable())
	}
}

func TestgetNextTarget(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}
	addrs := []string{"http://testhost1:8000", "http://testhost2:8000", "http://testhost3:8000"}
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		trg := &MockTarget{
			address:     addr,
			isAvailable: true,
			ping:        true,
		}
		bckt.AddTarget(trg)
	}
	for i := 0; i < 6; i++ {
		trg, err := bckt.getNextTarget()
		if err != nil {
			t.Error(err)
		}
		host := strings.Split(addrs[i%3], "/")[2]
		if trg.Address().Host != host {
			t.Error("Expected", host, "got", trg.Address().Host)
		}
	}

}
func TestgetNextTargetEmpty(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}
	trg, err := bckt.getNextTarget()
	if err == nil {
		t.Error("Expected", ErrNoTargetsAvailable, "got", nil)
	}
	if trg != nil {
		t.Error("Expected", nil, "got", trg)
	}
}
func TestgetNextTargetUnreachable(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}
	addr, _ := url.Parse("http://testhost1:8000")

	bckt.AddTarget(&MockTarget{address: addr, isAvailable: false})
	trg, err := bckt.getNextTarget()
	if err == nil {
		t.Error("Expected", ErrAllTargetsUnreachable, "got", nil)
	}
	if trg != nil {
		t.Error("Expected", nil, "got", trg)
	}
}

func TestHealthcheck(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}
	addrs := []string{"http://testhost7:8000", "http://testhost8:8000", "http://testhost9:8000"}
	flag := true
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		trg := &MockTarget{
			address:     addr,
			isAvailable: true,
			ping:        flag,
		}
		bckt.AddTarget(trg)
		flag = !flag
	}
	bckt.Healthcheck()
	flag = true
	for _, trg := range bckt.targets {
		if trg.IsAvailable() != flag {
			t.Error("Expected", trg, "got", trg.IsAvailable())
		}
		flag = !flag
	}

}

func TestRemoveStale(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}
	addrs := []string{"http://testhost1:8000", "http://testhost2:8000", "http://testhost3:8000"}
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		trg := &MockTarget{
			address:     addr,
			isAvailable: false,
		}
		bckt.AddTarget(trg)
	}
	bckt.RemoveStale(time.Second * 0)
	if len(bckt.targets) != 0 {
		t.Error("Expected", 0, "got", len(bckt.targets))
	}
}
func TestRemoveStaleDifferent(t *testing.T) {
	bckt := &RoundRobinBucket{
		targets: []Target{},
	}
	addrs := []string{"http://testhost3:8000", "http://testhost4:8000", "http://testhost5:8000"}
	flag := true
	for i := 0; i < 3; i++ {
		addr, _ := url.Parse(addrs[i])
		trg := &MockTarget{
			address:     addr,
			isAvailable: flag,
			ping:        flag,
		}
		bckt.AddTarget(trg)
		flag = !flag
	}
	bckt.RemoveStale(time.Second * 0)
	if len(bckt.targets) != 2 {
		t.Error("Expected", 2, "got", len(bckt.targets))
	}
}
