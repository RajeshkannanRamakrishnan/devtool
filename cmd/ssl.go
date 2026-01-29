package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
)

// sslCmd represents the ssl command
var sslCmd = &cobra.Command{
	Use:   "ssl [domain]",
	Short: "Check SSL certificate expiration",
	Long: `Checks the SSL certificate for a given domain and prints the issuer,
expiry date, and days remaining without needing to remember openssl commands.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		domain := args[0]
		
		// Add default port if missing
		target := domain
		if _, _, err := net.SplitHostPort(domain); err != nil {
			target = domain + ":443"
		}

		conf := &tls.Config{
			InsecureSkipVerify: true,
		}

		fmt.Printf("Connecting to %s...\n", target)
		conn, err := tls.Dial("tcp", target, conf)
		if err != nil {
			fmt.Printf("Error connecting to %s: %v\n", target, err)
			return
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			fmt.Printf("No certificates found for %s\n", domain)
			return
		}

		// The first certificate is the leaf
		cert := certs[0]
		
		fmt.Printf("\nDomain: \033[1;36m%s\033[0m\n", domain)
		fmt.Printf("Issuer: %s\n", cert.Issuer.CommonName)
		if cert.Issuer.Organization != nil && len(cert.Issuer.Organization) > 0 {
			fmt.Printf("        (%s)\n", cert.Issuer.Organization[0])
		}
		
		fmt.Printf("Expiry: %s\n", cert.NotAfter.Format("2006-01-02 15:04:05 MST"))
		
		daysRemaining := time.Until(cert.NotAfter).Hours() / 24
		
		// Color code days remaining
		daysStr := fmt.Sprintf("%.0f", daysRemaining)
		if daysRemaining < 0 {
			// Expired
			fmt.Printf("Status: \033[1;31mEXPIRED (%s days ago)\033[0m\n", stringsTrimPrefix(daysStr, "-"))
		} else if daysRemaining < 14 {
			// Critical (less than 2 weeks)
			fmt.Printf("Status: \033[1;31mExpiring soon! (%s days remaining)\033[0m\n", daysStr)
		} else if daysRemaining < 30 {
			// Warning (less than a month)
			fmt.Printf("Status: \033[1;33mExpiring in %s days\033[0m\n", daysStr)
		} else {
			// Good
			fmt.Printf("Status: \033[1;32mValid (%s days remaining)\033[0m\n", daysStr)
		}
	},
}

func init() {
	rootCmd.AddCommand(sslCmd)
}

func stringsTrimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[0:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
