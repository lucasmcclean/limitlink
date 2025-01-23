package server

import (
	"crypto/tls"
	"fmt"

	"github.com/lucasmcclean/url-shortener/config"
)

func getTLSConfig(srvCfg *config.Server) (*tls.Config, error) {
	cert, err := getCertificate(srvCfg.CertPath)
	if err != nil {
		return nil, err
	}

	// All configurations follow reccomendations by:
	// https://ssl-config.mozilla.org/#server=go&version=1.22.0&config=intermediate&guideline=5.7
	tlsCfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		Certificates:             []tls.Certificate{*cert},
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.X25519, // Go 1.8+
			tls.CurveP256,
			tls.CurveP384,
			// tls.x25519Kyber768Draft00, // Go 1.23+
		},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	return tlsCfg, nil
}

func getCertificate(certDirectory string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(certDirectory+"/server.crt", certDirectory+"/server.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load X509 keypair: %s", err)
	}
	return &cert, nil
}
