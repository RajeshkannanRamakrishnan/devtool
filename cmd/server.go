package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var serverPort int
var serverSSL bool

const maxLoggedRequestBodyBytes = 64 * 1024

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an HTTP server that responds 200 OK to any request",
	Long: `Start an HTTP server that responds with 200 OK to any URL path.
You can specify the port using the --port (or -p) flag.
Use --ssl to serve over HTTPS with a self-signed certificate.

Every request to any path will receive an HTTP 200 OK response
with a JSON body: {"status": "ok"}.`,
	Example: `  devtool server
  devtool server --port 9090
  devtool server --ssl`,
	Run: func(cmd *cobra.Command, args []string) {
		// Catch-all handler: respond 200 OK with JSON for any path
		catchAll := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, `{"status": "ok"}`)
		})
		handler := requestLogger(catchAll)

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
				Handler: handler,
				TLSConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			}

			fmt.Printf("Listening at https://localhost%s (Self-Signed SSL) — responds 200 OK to all paths\n", addr)
			if err := server.ListenAndServeTLS("", ""); err != nil {
				log.Fatalf("Server failed: %v", err)
			}
		} else {
			fmt.Printf("Listening at http://localhost%s — responds 200 OK to all paths\n", addr)
			if err := http.ListenAndServe(addr, handler); err != nil {
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

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

type bodyCaptureReadCloser struct {
	body      io.ReadCloser
	buf       strings.Builder
	limit     int
	truncated bool
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.status == 0 {
		lrw.status = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}

func (b *bodyCaptureReadCloser) Read(p []byte) (int, error) {
	n, err := b.body.Read(p)
	if n > 0 && b.limit > 0 {
		remaining := b.limit - b.buf.Len()
		if remaining > 0 {
			if n <= remaining {
				_, _ = b.buf.Write(p[:n])
			} else {
				_, _ = b.buf.Write(p[:remaining])
				b.truncated = true
			}
		} else {
			b.truncated = true
		}
	}

	return n, err
}

func (b *bodyCaptureReadCloser) Close() error {
	return b.body.Close()
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyCapture := &bodyCaptureReadCloser{
			body:  r.Body,
			limit: maxLoggedRequestBodyBytes,
		}
		r.Body = bodyCapture

		headers := make(map[string]string, len(r.Header))
		keys := make([]string, 0, len(r.Header))
		for key := range r.Header {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			headers[key] = strings.Join(r.Header.Values(key), ", ")
		}

		logData := map[string]any{
			"timestamp":   time.Now().Format(time.RFC3339),
			"remote_addr": r.RemoteAddr,
			"method":      r.Method,
			"url":         r.URL.String(),
			"proto":       r.Proto,
			"host":        r.Host,
			"headers":     headers,
		}

		lrw := &loggingResponseWriter{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(lrw, r)
		if lrw.status == 0 {
			lrw.status = http.StatusOK
		}

		logData["status"] = lrw.status
		logData["response_bytes"] = lrw.bytes
		logData["duration_ms"] = time.Since(start).Milliseconds()
		logData["body"] = bodyCapture.buf.String()
		if bodyCapture.truncated {
			logData["body_truncated"] = true
		}

		payload, marshalErr := json.Marshal(logData)
		if marshalErr != nil {
			log.Printf("Failed to marshal request log: %v", marshalErr)
			return
		}
		log.Println(string(payload))
	})
}
