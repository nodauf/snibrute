package cmd

import (
	"fmt"
	"os"

	"github.com/nodauf/SNIBrute/brute"
	"github.com/spf13/cobra"
)

var domain string
var wordlist string
var ip string
var port int
var option brute.Options

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "SNIBrute",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		option.Brute()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&option.Domain, "domain", "d", "", "Domain to bruteforce. Ex: example.org")
	rootCmd.Flags().StringVarP(&option.Wordlist, "wordlist", "w", "", "Wordlist which contains all the subdomain to try")
	rootCmd.Flags().StringVarP(&option.Ip, "ip", "i", "", "IP address if the host targeted host")
	rootCmd.Flags().IntVarP(&option.Port, "port", "p", 443, "https port")

	// Match
	rootCmd.Flags().StringVar(&option.MatchStatus, "matchStatus", "", "Match only on specified status code. Ex: 200,403")
	rootCmd.Flags().StringVar(&option.MatchSize, "matchSize", "", "Match only on specified size. Ex: 1337,7331")

	// Filter
	rootCmd.Flags().StringVar(&option.FilterStatus, "filterStatus", "", "Filter only on specified status code. Ex: 404,500")
	rootCmd.Flags().StringVar(&option.FilterSize, "filterSize", "", "Filter only on specified size. Ex: 0,10")

	rootCmd.Flags().BoolVarP(&option.Verbose, "verbose", "v", false, "Verbose mode")

	rootCmd.MarkFlagRequired("domain")
	rootCmd.MarkFlagRequired("wordlist")
	rootCmd.MarkFlagRequired("IP")
}
