/*
Copyright © 2022 Martin Marsh martin@marshtrio.com

*/

package cmd

import (
	"fmt"
	"navmux/version"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "navmux",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.navmux.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(runCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of navmux",
	Long:  `Version Number of NavData Boat General and Navigation data processing`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NavMux - Multiplexing data sources and data recording\nVersion: " + version.BuildVersion + "\nBuild: " + version.BuildTime)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Navmux starts processing",
	Long:  `Start Navmux Multiplexing data sources and recording of Boat data - runs until aborted`,
	Run: func(cmd *cobra.Command, args []string) {
		// to move up to start of oneline up use \033[F

		if len(args) > 0 {
			fmt.Printf("\nStarting Navmux using %s\n\nruns until aborted\n", args[0])
		} else {
			fmt.Println("\nStarting Navmux\nruns until aborted")
		}

		load_and_run()

	},
}
