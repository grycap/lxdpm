package apilxd

import (
	"crypto/x509"
	"encoding/base64"
	//"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"

	//"github.com/gorilla/mux"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	"github.com/lxc/lxd/shared/version"
)

func certificatesGet(lx *LxdpmApi, r *http.Request) Response {

	body := []string{}
	for _, crt := range fsCertsGetSingle() {
		cert, err := x509.ParseCertificate([]byte(crt.Certificate))
		if err !=nil {
			fmt.Println(err)
		} else {
			fingerprint := fmt.Sprintf("/%s/certificates/%s", version.APIVersion, shared.CertFingerprint(cert))
			body = append(body, fingerprint)
		}
	}

	return SyncResponse(true, body)
}

func readSavedClientCAList(lx *LxdpmApi) {
	lx.clientCerts = []x509.Certificate{}

	fsCerts, err := fsCertsGet()
	if err != nil {
		shared.LogInfof("Error reading certificates from dir: %s", err)
		return
	}

	for _, fsCert := range fsCerts {
		certBlock, _ := pem.Decode([]byte(fsCert.Certificate))
		if certBlock == nil {
			shared.LogInfof("Error decoding certificate for %s: %s", fsCert.Name, err)
			continue
		}

		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			shared.LogInfof("Error reading certificate for %s: %s", fsCert.Name, err)
			continue
		}
		lx.clientCerts = append(lx.clientCerts, *cert)
	}
}

func saveCert(host string, cert *x509.Certificate) error {
	baseCert := new(fsCertInfo)
	baseCert.Fingerprint = shared.CertFingerprint(cert)
	baseCert.Type = 1
	baseCert.Name = host
	baseCert.Certificate = string(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}),
	)

	return fsCertSave(baseCert)
}

func certificatesPost(lx *LxdpmApi, r *http.Request) Response {
	// Parse the request
	req := api.CertificatesPost{}
	if err := shared.ReadToJSON(r.Body, &req); err != nil {
		return BadRequest(err)
	}
	/*
	// Access check
	if !d.isTrustedClient(r) && d.PasswordCheck(req.Password) != nil {
		return Forbidden
	}
	*/
	if req.Type != "client" {
		return BadRequest(fmt.Errorf("Unknown request type %s", req.Type))
	}
	

	// Extract the certificate
	var cert *x509.Certificate
	var name string
	if req.Certificate != "" {
		data, err := base64.StdEncoding.DecodeString(req.Certificate)
		if err != nil {
			return BadRequest(err)
		}

		cert, err = x509.ParseCertificate(data)
		if err != nil {
			return BadRequest(err)
		}
		name = req.Name
	} else if r.TLS != nil {
		if len(r.TLS.PeerCertificates) < 1 {
			return BadRequest(fmt.Errorf("No client certificate provided"))
		}
		cert = r.TLS.PeerCertificates[len(r.TLS.PeerCertificates)-1]

		remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return InternalError(err)
		}

		name = remoteHost
	} else {
		return BadRequest(fmt.Errorf("Can't use TLS data on non-TLS link"))
	}

	fingerprint := shared.CertFingerprint(cert)
	for _, existingCert := range lx.clientCerts {
		if fingerprint == shared.CertFingerprint(&existingCert) {
			return BadRequest(fmt.Errorf("Certificate already in trust store"))
		}
	}

	err := saveCert(name, cert)
	if err != nil {
		return SmartError(err)
	}

	lx.clientCerts = append(lx.clientCerts, *cert)

	return SyncResponseLocation(true, nil, fmt.Sprintf("/%s/certificates/%s", version.APIVersion, fingerprint))
}

var certificatesCmd = Command{name: "certificates", untrustedPost: true, get: certificatesGet, post: certificatesPost}