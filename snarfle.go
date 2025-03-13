/*
SNARFLE

Read text files and convert IP addresses to domain names
*/

package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"snarfle/filelib"
	"strings"
)

const (
	ipRegexStr = `([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`
)

var (
	fileIn = ""
)

func main() {
	/** OPEN INPUT FILE **/

	cfg, cfgErr := filelib.ReadConfig("snarfle.cfg")
	if cfgErr != nil {
		log.Fatal("Could not read config file.")
	}

	srcDir := cfg["srcDir"]
	outputFmt := cfg["outputFmt"]

	/*  GET COMMAND LINE FLAGS  */
	flag.StringVar(&fileIn, "f", fileIn, "Name of input file")
	flag.StringVar(&srcDir, "d", cfg["srcDir"], "Input file directory")
	flag.StringVar(&outputFmt, "o", cfg["outputFmt"], "Output format")
	flag.Parse()

	fileInH, err := os.Open(filepath.Join(srcDir, fileIn))
	if err != nil {
		log.Fatal(err)
	}
	defer fileInH.Close()

	/** OPEN OUTPUT FILE **/
	fileBase := strings.TrimSuffix(filepath.Base(fileIn), filepath.Ext(fileIn))

	fileOutH, err := os.Create(fileBase + "-" + cfg["outputSuffix"] + "." + outputFmt)
	if err != nil {
		log.Fatal(err)
	}
	defer fileOutH.Close() // ensure file gets closed

	knownHosts := make(map[string]string)
	localHosts, err := filelib.ReadConfig("localhosts.cfg")
	if err != nil {
		log.Fatal("Couldn't read localhosts.cfg file.")
	}

	ipRegex, _ := regexp.Compile(ipRegexStr)

	// Create a scanner to read input file line by line
	scanner := bufio.NewScanner(fileInH)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		elements := strings.Split(line, " ")
		newElements := make([]string, 0) // to save changed items

		note := ""

		for _, item := range elements {
			item = strings.TrimSpace(item)
			if item != "" {
				// Check to see if this item is an IP address
				if ipRegex.MatchString(item) { // It is!
					//quads := strings.Split(item, ".")
					//if quads[0] != "10" {
					if strings.HasPrefix(item, "10.") || strings.HasPrefix(item, "192.168") {
						// This is a local address. Let's look whether it's
						// in the localHosts map. Doesn't deal with 172. private addresses
						for dev, ip := range localHosts {
							if ip == item {
								item = dev
							}
						}
					} else {
						// This is a remote address
						name, found := knownHosts[item] // Have we seen this already?
						if found {
							note = name
						} else {
							// If it's new to us, look it up
							cmd := exec.Command("dig", "+short", "-x", item)
							out, err := cmd.CombinedOutput()
							if err != nil {
								log.Fatal(err)
							}
							domain := strings.TrimSpace(string(out))
							domain = strings.TrimSuffix(domain, ".")
							knownHosts[item] = domain // Add to list of found domains
							note = domain
						}
					}
				}
				if item != ":" {
					newElements = append(newElements, item)
				}
			}
		}
		newElements = append(newElements, note)

		logline := ""
		switch outputFmt {
		case "csv":
			logline = strings.Join(newElements, ",")
		case "log", "txt":
			logline = strings.Join(newElements, " ")
		}
		_, ferr := fileOutH.WriteString(logline + "\n")
		if ferr != nil {
			log.Fatal(ferr)
		}

	}
	// for ip, val := range knownHosts {
	// 	fmt.Println(ip, " ", val)
	// }
}
