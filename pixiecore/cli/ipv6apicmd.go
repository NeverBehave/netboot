package cli

import (
	"github.com/spf13/cobra"
	"fmt"
	"go.universe.tf/netboot/pixiecorev6"
	"go.universe.tf/netboot/dhcp6"
	"time"
	"net"
)

var ipv6ApiCmd = &cobra.Command{
	Use:   "ipv6api",
	Short: "Boot a kernel and optional init ramdisks over IPv6 using api",
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := cmd.Flags().GetString("listen-addr")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}
		apiUrl, err := cmd.Flags().GetString("api-request-url")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}
		apiTimeout, err := cmd.Flags().GetDuration("api-request-timeout")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}

		s := pixiecorev6.NewServerV6()
		s.Log = logWithStdFmt
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			s.Debug = logWithStdFmt
		}
		if debug { s.Debug = logWithStdFmt }

		if addr == "" {
			fatalf("Please specify address to bind to")
		} else {
		}
		if apiUrl == "" {
			fatalf("Please specify ipxe config file url")
		}
		s.Address = addr
		preference, err := cmd.Flags().GetUint8("preference")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}
		s.BootConfig = dhcp6.MakeApiBootConfiguration(apiUrl, apiTimeout, preference, cmd.Flags().Changed("preference"))

		addressPoolStart, err := cmd.Flags().GetString("address-pool-start")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}
		addressPoolSize, err := cmd.Flags().GetUint64("address-pool-size")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}
		addressPoolValidLifetime, err := cmd.Flags().GetUint32("address-pool-lifetime")
		if err != nil {
			fatalf("Error reading flag: %s", err)
		}
		s.AddressPool = dhcp6.NewRandomAddressPool(net.ParseIP(addressPoolStart), addressPoolSize, addressPoolValidLifetime)
		s.PacketBuilder = dhcp6.MakePacketBuilder(s.Duid, addressPoolValidLifetime - addressPoolValidLifetime*3/100, addressPoolValidLifetime)

		fmt.Println(s.Serve())
	},
}

func serverv6ApiConfigFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("listen-addr", "", "", "IPv6 address to listen on")
	cmd.Flags().StringP("api-request-url", "", "", "Ipv6-specific API server url")
	cmd.Flags().Duration("api-request-timeout", 5*time.Second, "Timeout for request to the API server")
	cmd.Flags().Bool("debug", false, "Enable debug-level logging")
	cmd.Flags().Uint8("preference", 255, "Set dhcp server preference value")
	cmd.Flags().StringP("address-pool-start", "", "2001:db8:f00f:cafe:ffff::100", "Starting ip of the address pool, e.g. 2001:db8:f00f:cafe:ffff::100")
	cmd.Flags().Uint64("address-pool-size", 50, "Address pool size")
	cmd.Flags().Uint32("address-pool-lifetime", 1850, "Address pool ip address valid lifetime in seconds")
}

func init() {
	rootCmd.AddCommand(ipv6ApiCmd)
	serverv6ApiConfigFlags(ipv6ApiCmd)
}

