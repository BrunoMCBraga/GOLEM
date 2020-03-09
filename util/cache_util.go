package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const ipRegexString string = "^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"

func GetDNSIPs(useLocalCache bool, cacheFilePath string, dnsServersSourceURLs []string) ([]string, error) {

	var dnsServerIPs []string

	if useLocalCache {
		if _, osStatError := os.Stat(cacheFilePath); osStatError == nil {
			readFileResult, readFileError := ioutil.ReadFile(cacheFilePath)
			if readFileError != nil {
				return dnsServerIPs, errors.New(fmt.Sprintf("%s: %s", "GOLEM->cache_util->GetDNSIPs->ioutil.ReadFile:", readFileError.Error()))
			}
			dnsServerIPs = strings.Split(strings.Trim(string(readFileResult), "\n"), "\n")
		}
	} else {

		for _, dnsServersSourceURL := range dnsServersSourceURLs {

			getResult, getError := http.Get(dnsServersSourceURL)
			if getError != nil {
				return dnsServerIPs, errors.New(fmt.Sprintf("%s: %s", "GOLEM->cache_util->GetDNSIPs->ioutil.ReadFile:", getError.Error()))
			}
			defer getResult.Body.Close()

			readAllResult, readAllError := ioutil.ReadAll(getResult.Body)
			if getError != nil {
				return dnsServerIPs, errors.New(fmt.Sprintf("%s: %s", "GOLEM->cache_util->GetDNSIPs->ioutil.ReadFile:", readAllError.Error()))
			}

			dnsServerIPsTemp := strings.Split(strings.Trim(string(readAllResult), "\n"), "\n")
			dnsServerIPs = append(dnsServerIPs, dnsServerIPsTemp...)

		}

	}

	deduplicatedIPsMap := make(map[string]bool, 0)
	regexIPv4CompileResult, regexIPv4CompileError := regexp.Compile(ipRegexString)

	if regexIPv4CompileError != nil {
		return dnsServerIPs, errors.New(fmt.Sprintf("%s: %s", "GOLEM->cache_util->GetDNSIPs->ioutil.ReadFile->regexp.Compile:", regexIPv4CompileError.Error()))
	}

	for _, dnsServerIP := range dnsServerIPs {
		if regexIPv4CompileResult.MatchString(dnsServerIP) {
			deduplicatedIPsMap[dnsServerIP] = true
		}
	}

	deduplicatedIPsSlice := make([]string, 0)
	for ip := range deduplicatedIPsMap {
		deduplicatedIPsSlice = append(deduplicatedIPsSlice, ip)
	}

	return deduplicatedIPsSlice, nil

}
