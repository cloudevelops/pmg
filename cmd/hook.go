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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"encoding/json"

	//	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var wg sync.WaitGroup

// hookCmd represents the hook command
var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Starting hook api")
		httpApi()
	},
}

func init() {
	RootCmd.AddCommand(hookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func httpApi() (err error) {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homePage)
	r.HandleFunc("/hook", doHook)
	err = http.ListenAndServe(":8666", r)
	return err
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the pmg hook homepage!")
	fmt.Println("Refer to documentation on how to use it...")
}

func doHook(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Start processing hook !")
	// Read request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Failed to read request: " + err.Error())
		return
	}
	// Unmarshal request data
	var t map[string]interface{}
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Errorf("Failed to parse request json: " + err.Error())
	}
	//spew.Dump(t)
	// Determine if repository is Puppet module
	repo := t["repository"].(map[string]interface{})
	repoFullName := repo["full_name"].(string)
	gitHome := viper.GetString("hook.githome")
	//spew.Dump(gitHome)
	cmd := exec.Command("git", "--git-dir="+gitHome+"/"+repoFullName+".git", "cat-file", "blob", "HEAD:metadata.json")
	//spew.Dump(cmd)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	err = cmd.Run()
	var metadata map[string]interface{}
	var module string
	if err == nil {
		// Puppet module with metadata.json found
		log.Debugf("Git found metadata.json, parsing !")
		// Unmarshall JSON into plain interface
		err = json.Unmarshal(outb.Bytes(), &metadata)
		if err == nil {
			log.Debugf("Module medatada sucessfuly parsed !")
			moduleFullName := metadata["name"].(string)
			moduleFullNameSplit := strings.Split(moduleFullName, "-")
			if len(moduleFullNameSplit) < 2 {
				log.Debugf("Invalid module name '" + moduleFullName + "', should be [organization]-[module_name] !!!")
			} else {
				module = moduleFullNameSplit[1]
				updateModule(module)
			}
		}
	} else {
		cmd := exec.Command("git", "--git-dir="+gitHome+"/"+repoFullName+".git", "cat-file", "blob", "HEAD:Modulefile")
		var outb bytes.Buffer
		cmd.Stdout = &outb
		err = cmd.Run()
		if err == nil {
			// Puppet module with Modulefile found
			log.Debugf("Git found Modulefile, parsing. PLEASE UPDATE THE MODULE, Modulefile is long time obsolete !!!")
			// Unmarshall JSON into plain interface
			moduleFullName := outb.String()
			moduleFullNameSplit := strings.Split(moduleFullName, "'")
			module = strings.Split(moduleFullNameSplit[1], "-")[1]
			updateModule(module)
		}
	}
	// Determine if repository is hiera data
	repoName := strings.Split(repoFullName, "/")
	if strings.Contains(repoName[1], "_hiera") {
		log.Debugf("Found hiera repo")
		updateHiera(repoName[1])
	}
	// Determine if repository is puppet r10k repository
	if repoName[1] == "puppet_r10k" {
		log.Debugf("Found puppet r10k repo")
		updatePuppetR10k(repoName[1])
	}
	// Determine if repository is hiera r10k repository
	if repoName[1] == "hiera_r10k" {
		log.Debugf("Found hiera r10k repo")
		updateHieraR10k(repoName[1])
	}
}

func updateModule(module string) {
	log.Debugf("Updating puppet module: " + module)
	puppetServers := viper.GetStringSlice("hook.puppetservers")
	for _, puppetServer := range puppetServers {
		wg.Add(1)
		go executeSshCommand(puppetServer, "r10k deploy module "+module+" --config /etc/r10k/puppet_r10k.yaml")
	}
	wg.Wait()
}

func updateHiera(module string) {
	log.Debugf("Updating hiera data: " + module)
	puppetServers := viper.GetStringSlice("hook.puppetservers")
	for _, puppetServer := range puppetServers {
		wg.Add(1)
		go executeSshCommand(puppetServer, "r10k deploy module "+module+" --config /etc/r10k/hiera_r10k.yaml")
	}
	wg.Wait()
}

func updatePuppetR10k(module string) {
	log.Debugf("Updating puppet r10k data: " + module)
	puppetServers := viper.GetStringSlice("hook.puppetservers")
	for _, puppetServer := range puppetServers {
		wg.Add(1)
		go executeSshCommand(puppetServer, "r10k deploy environment -p --config /etc/r10k/puppet_r10k.yaml")
	}
	wg.Wait()
}

func updateHieraR10k(module string) {
	log.Debugf("Updating hiera r10k data: " + module)
	puppetServers := viper.GetStringSlice("hook.puppetservers")
	for _, puppetServer := range puppetServers {
		wg.Add(1)
		go executeSshCommand(puppetServer, "r10k deploy environment -p --config /etc/r10k/hiera_r10k.yaml")
	}
	wg.Wait()
}

func executeSshCommand(host string, command string) {
	defer wg.Done()
	cmd := exec.Command("ssh", "root@"+host, command)
	var outb bytes.Buffer
	cmd.Stdout = &outb
	err := cmd.Run()
	if err == nil {
		// Puppet module with Modulefile found
		log.Debugf("Command: " + command + " on host:" + host + " did run OK, result:" + outb.String())
	} else {
		log.Debugf("Command: " + command + " on host:" + host + " did run with Error:" + err.Error() + ", result: " + outb.String())
	}
}
