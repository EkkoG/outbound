package shadowtls

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/daeuniverse/outbound/dialer"
	"github.com/daeuniverse/outbound/netproxy"
	shadowtls "github.com/sagernet/sing-shadowtls"
	// "github.com/sagernet/sing/common/metadata"
)

var (
	DefaultALPN = []string{"h2", "http/1.1"}
)

func init() {
	dialer.FromLinkRegister("shadowtls", NewShadowsocksFromLink)
}

type ShadowTLS struct {
	dialer        netproxy.Dialer
	Name          string `json:"name"`
	Server        string `json:"server"`
	Port          int    `json:"port"`
	Password      string `json:"password"`
	SNI           string `json:"sni"`
	Version       int    `json:"version"`
	AllowInsecure bool   `json:"allow_insecure"`
}

func NewShadowsocksFromLink(option *dialer.ExtraOption, nextDialer netproxy.Dialer, link string) (npd netproxy.Dialer, property *dialer.Property, err error) {
	s, err := ParseShadowTLSURL(link)
	fmt.Println("ahhhh", s)
	if err != nil {
		return nil, nil, err
	}
	s.dialer = nextDialer
	return s, &dialer.Property{
		Name:     s.Name,
		Address:  fmt.Sprintf("%s:%d", s.Server, s.Port),
		Protocol: "shadowtls",
		Link:     link,
	}, nil

}

func ParseShadowTLSURL(link string) (*ShadowTLS, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	password := u.User.Username()
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}
	if host == "" {
		return nil, fmt.Errorf("empty host")
	}
	if port == "" {
		return nil, fmt.Errorf("empty port")
	}
	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	sni := q.Get("sni")
	version := 3
	if v := q.Get("version"); v != "" {
		version, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return &ShadowTLS{
		Name:     u.Fragment,
		Server:   host,
		Port:     portInt,
		Password: password,
		SNI:      sni,
		Version:  version,
	}, nil
}

func (s *ShadowTLS) DialContext(ctx context.Context, network, addr string) (c netproxy.Conn, err error) {
	fmt.Println("dialing ahhhhh", network, addr, s)
	magicNetwork, err := netproxy.ParseMagicNetwork(network)
	if err != nil {
		return nil, err
	}

	switch magicNetwork.Network {
	case "tcp":
		tlsConfig := &tls.Config{
			NextProtos:         DefaultALPN,
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false,
			ServerName:         s.SNI,
		}

		tlsHandshake := shadowtls.DefaultTLSHandshakeFunc(s.Password, tlsConfig)
		client, err := shadowtls.NewClient(shadowtls.ClientConfig{
			Version:      s.Version,
			Password:     s.Password,
			TLSHandshake: tlsHandshake,
			Logger:       SingLogger,
			// Server: metadata.Socksaddr{
			// 	Addr: metadata.AddrFromIP(net.ParseIP(s.Server)),
			// 	Port: uint16(s.Port),
			// },
		})
		if err != nil {
			fmt.Println("error", err)
			return nil, err
		}
		shadowtlsAddr := fmt.Sprintf("%s:%d", s.Server, s.Port)
		rc, err := s.dialer.DialContext(ctx, network, shadowtlsAddr)
		if err != nil {
			fmt.Println("error", err)
			return nil, err
		}
		conn, err := client.DialContextConn(ctx, &netproxy.FakeNetConn{Conn: rc})
		if err != nil {
			fmt.Println("error", err)
			return nil, err
		}
		return conn, nil
	default:
		return nil, fmt.Errorf("%w: %v", netproxy.UnsupportedTunnelTypeError, network)
	}
}
