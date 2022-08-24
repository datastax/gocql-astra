package gocqlastra

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/datastax/cql-proxy/astra"
	"github.com/gocql/gocql"
)

const AstraAPIURL = "https://api.astra.datastax.com"

type dialer struct {
	sniProxyAddr      string   // Don't use directly
	contactPoints     []string // Don't use directly
	contactPointIndex int32
	bundle            *astra.Bundle
	dialer            net.Dialer
	mu                sync.Mutex
	timeout           time.Duration
}

func NewDialerFromBundle(path string, timeout time.Duration) (gocql.HostDialer, error) {
	bundle, err := astra.LoadBundleZipFromPath(path)
	if err != nil {
		return nil, err
	}
	return &dialer{
		bundle:  bundle,
		timeout: timeout,
	}, nil
}

func NewDialerFromURL(url, databaseID, token string, timeout time.Duration) (gocql.HostDialer, error) {
	bundle, err := astra.LoadBundleZipFromURL(url, databaseID, token, timeout)
	if err != nil {
		return nil, err
	}
	return &dialer{
		bundle:  bundle,
		timeout: timeout,
	}, nil
}

func (d *dialer) DialHost(ctx context.Context, host *gocql.HostInfo) (*gocql.DialedHost, error) {
	sniAddr, contactPoints, err := d.resolveMetadata(ctx)
	if err != nil {
		return nil, err
	}

	addr, err := lookupHost(sniAddr)
	if err != nil {
		return nil, err
	}

	conn, err := d.dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	hostId := host.HostID()
	if hostId == "" {
		hostId = contactPoints[int(atomic.AddInt32(&d.contactPointIndex, 1))%len(d.contactPoints)]
	}

	tlsConn := tls.Client(conn, copyTLSConfig(d.bundle, hostId))
	if err = tlsConn.HandshakeContext(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return &gocql.DialedHost{
		Conn:            tlsConn,
		DisableCoalesce: true, // See https://github.com/mpenick/gocqlastra/issues/1
	}, nil
}

func (d *dialer) resolveMetadata(ctx context.Context) (string, []string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// TODO: Make this value have a TTL
	if d.sniProxyAddr != "" {
		return d.sniProxyAddr, d.contactPoints, nil
	}

	var metadata *astraMetadata

	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	httpsClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: d.bundle.TLSConfig().Clone(),
		},
	}

	url := fmt.Sprintf("https://%s:%d/metadata", d.bundle.Host(), d.bundle.Port())
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return "", nil, err
	}

	response, err := httpsClient.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("unable to get Astra metadata from %s: %w", url, err)
	}

	body, err := readAllWithTimeout(response.Body, ctx)
	if err != nil {
		return "", nil, err
	}

	err = json.Unmarshal(body, &metadata)
	if err != nil {
		return "", nil, err
	}

	d.sniProxyAddr = metadata.ContactInfo.SniProxyAddress
	d.contactPoints = metadata.ContactInfo.ContactPoints

	return d.sniProxyAddr, d.contactPoints, err
}

func copyTLSConfig(bundle *astra.Bundle, serverName string) *tls.Config {
	tlsConfig := bundle.TLSConfig().Clone()
	tlsConfig.ServerName = serverName
	tlsConfig.InsecureSkipVerify = true
	tlsConfig.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		certs := make([]*x509.Certificate, len(rawCerts))
		for i, asn1Data := range rawCerts {
			cert, err := x509.ParseCertificate(asn1Data)
			if err != nil {
				return errors.New("tls: failed to parse certificate from server: " + err.Error())
			}
			certs[i] = cert
		}

		opts := x509.VerifyOptions{
			Roots:         tlsConfig.RootCAs,
			CurrentTime:   time.Now(),
			DNSName:       bundle.Host(),
			Intermediates: x509.NewCertPool(),
		}
		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}
		var err error
		verifiedChains, err = certs[0].Verify(opts)
		return err
	}
	return tlsConfig
}

func readAllWithTimeout(r io.Reader, ctx context.Context) (bytes []byte, err error) {
	ch := make(chan struct{})

	go func() {
		bytes, err = ioutil.ReadAll(r)
		close(ch)
	}()

	select {
	case <-ch:
	case <-ctx.Done():
		return nil, errors.New("timeout reading data")
	}

	return bytes, err
}

func lookupHost(hostWithPort string) (string, error) {
	host, port, err := net.SplitHostPort(hostWithPort)
	if err != nil {
		return "", err
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	addr := addrs[rand.Intn(len(addrs))]
	if len(port) > 0 {
		addr = net.JoinHostPort(addr, port)
	}
	return addr, nil
}

type contactInfo struct {
	TypeName        string   `json:"type"`
	LocalDc         string   `json:"local_dc"`
	SniProxyAddress string   `json:"sni_proxy_address"`
	ContactPoints   []string `json:"contact_points"`
}

type astraMetadata struct {
	Version     int         `json:"version"`
	Region      string      `json:"region"`
	ContactInfo contactInfo `json:"contact_info"`
}
