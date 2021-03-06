// Copyright © 2017 Zdenek Janda <zdenek.janda@cloudevelops.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/juju/loggo"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var log = loggo.GetLogger("cmd")

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pmg",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	loggo.ConfigureLoggers("<root>=TRACE")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pmg.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pmg" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pmg")
	}

	viper.SetDefault("hook.githome", "/srv")
	//puppetServers := []string{"puppetca1.infra.prod.ci", "puppetserver1.infra.prod.ci", "puppetserver2.infra.prod.ci", "puppetserver3.infra.prod.ci", "puppetserver4.infra.prod.ci"}
	puppetServers := []string{"puppetserver1.cz2.cloudevelops.lan", "puppetserver2.cz2.cloudevelops.lan", "puppetserver3.cz2.cloudevelops.lan", "puppetserver4.cz2.cloudevelops.lan", "puppetserver5.cz2.cloudevelops.lan", "puppetserver6.cz2.cloudevelops.lan", "puppetserver7.cz2.cloudevelops.lan", "puppetserver8.cz2.cloudevelops.lan", "puppetserver9.cz2.cloudevelops.lan", "puppetserver10.cz2.cloudevelops.lan", "puppetserver11.cz2.cloudevelops.lan", "puppetserver12.cz2.cloudevelops.lan", "puppetserver3-1.cz2.cloudevelops.lan", "puppetserver3-2.cz2.cloudevelops.lan", "puppetserver3-3.cz2.cloudevelops.lan", "puppetserver3-4.cz2.cloudevelops.lan", "puppetserver3-5.cz2.cloudevelops.lan", "puppetserver3-6.cz2.cloudevelops.lan", "puppetserver3-7.cz2.cloudevelops.lan", "puppetserver3-8.cz2.cloudevelops.lan", "puppetserver3-9.cz2.cloudevelops.lan", "puppetserver3-10.cz2.cloudevelops.lan", "puppetserver3-11.cz2.cloudevelops.lan", "puppetserver3-12.cz2.cloudevelops.lan", "puppetserver3-13.cz2.cloudevelops.lan", "puppetserver3-14.cz2.cloudevelops.lan", "puppetserver3-15.cz2.cloudevelops.lan", "puppetserver3-16.cz2.cloudevelops.lan", "puppetserver3-17.cz2.cloudevelops.lan", "puppetserver3-18.cz2.cloudevelops.lan", "puppetserver3-19.cz2.cloudevelops.lan", "puppetserver3-20.cz2.cloudevelops.lan", "puppetserver3-21.cz2.cloudevelops.lan", "puppetserver3-22.cz2.cloudevelops.lan", "puppetserver3-23.cz2.cloudevelops.lan", "puppetserver3-24.cz2.cloudevelops.lan", "puppetserver3-25.cz2.cloudevelops.lan", "puppetserver3-26.cz2.cloudevelops.lan", "puppetserver3-27.cz2.cloudevelops.lan", "puppetserver3-28.cz2.cloudevelops.lan", "puppetserver3-29.cz2.cloudevelops.lan", "puppetserver1.cz2.gtflixtv.lan", "puppetserver2.cz2.gtflixtv.lan", "puppetserver3.cz2.gtflixtv.lan", "puppetserver4.cz2.gtflixtv.lan", "puppetserver5.cz2.gtflixtv.lan", "puppetserver6.cz2.gtflixtv.lan"}
	viper.SetDefault("hook.puppetservers", puppetServers)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
