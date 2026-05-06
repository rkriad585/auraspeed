package speedtest

import (
	"fmt"
	"time"

	st "github.com/showwin/speedtest-go/speedtest"
)

// MockServer mimics st.Server for testing
type MockServer struct {
	DLSpeed float64
	ULSpeed float64
	Latency time.Duration
	Name    string
	ID      int
}

func (m *MockServer) DownloadTest() error {
	return nil
}

func (m *MockServer) UploadTest() error {
	return nil
}

// MockClient mimics st.Speedtest for testing
type MockClient struct {
	User             *st.User
	Servers          st.Servers
	FetchUserInfoErr error
	FetchServersErr  error
}

func (m *MockClient) FetchUserInfo() (*st.User, error) {
	return m.User, m.FetchUserInfoErr
}

func (m *MockClient) FetchServers() (st.Servers, error) {
	return m.Servers, m.FetchServersErr
}

// MockUser creates a mock st.User for testing
func MockUser(isp, ip string) *st.User {
	return &st.User{
		Isp: isp,
		IP:  ip,
	}
}

// MockServers creates mock st.Servers for testing
func MockServers(servers ...*MockServer) st.Servers {
	ss := make(st.Servers, len(servers))
	for i, s := range servers {
		id := fmt.Sprintf("%d", s.ID)
		ss[i] = &st.Server{
			ID:      id,
			Name:    s.Name,
			Latency: s.Latency,
			DLSpeed: st.ByteRate(s.DLSpeed),
			ULSpeed: st.ByteRate(s.ULSpeed),
		}
	}
	return ss
}
