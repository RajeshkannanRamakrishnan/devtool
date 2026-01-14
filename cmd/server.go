package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var serverPort int
var serverSSL bool

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a static file server",
	Long: `Start a static HTTP file server for the current directory.
You can specify the port using the --port (or -p) flag.
Use --ssl to serve over HTTPS with a self-signed certificate.`,
	Example: `  devtool server
  devtool server --port 9090
  devtool server --ssl`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Error getting current directory: %v", err)
		}

		fs := http.FileServer(http.Dir(cwd))
		mux := http.NewServeMux()
		mux.Handle("/", fs)

		addr := fmt.Sprintf(":%d", serverPort)

		if serverSSL {
			certPEM, keyPEM, err := generateSelfSignedCert()
			if err != nil {
				log.Fatalf("Failed to generate self-signed certificate: %v", err)
			}

			cert, err := tls.X509KeyPair(certPEM, keyPEM)
			if err != nil {
				log.Fatalf("Failed to load generic key pair: %v", err)
			}

			server := &http.Server{
				Addr:    addr,
				Handler: mux,
				TLSConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			}

			fmt.Printf("Serving %s at https://localhost%s (Self-Signed SSL)\n", cwd, addr)
			if err := server.ListenAndServeTLS("", ""); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		} else {
			fmt.Printf("Serving %s at http://localhost%s\n", cwd, addr)
			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Port to listen on")
	serverCmd.Flags().BoolVar(&serverSSL, "ssl", false, "Serve over HTTPS with a self-signed certificate")
}

func generateSelfSignedCert() ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Devtool Safe Self-Signed"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return certPEM, keyPEM, nil
}
