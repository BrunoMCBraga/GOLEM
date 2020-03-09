package engines

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/golem/util"
	"github.com/miekg/dns"
)

var dnsTimeout time.Duration = time.Second * 2

func resolveFQDNARecords(fqdn string, dnsServer string) (r *dns.Msg, rtt time.Duration, err error) {

	fqdnWithDot := fqdn + "."

	dnsClient := dns.Client{}
	dnsClient.ReadTimeout = dnsTimeout
	dnsClient.WriteTimeout = dnsTimeout

	dnsMessage := dns.Msg{}
	dnsMessage.SetQuestion(fqdnWithDot, dns.TypeA)
	dnsMessage.RecursionDesired = true

	return dnsClient.Exchange(&dnsMessage, fmt.Sprintf("%s:%d", dnsServer, 53))

}

func GenerateSliceOfFQDNs(topDomain string, fqdnsFilePath string) ([]string, error) {

	fqdns := make([]string, 0)
	readFileResult, readFileError := ioutil.ReadFile(fqdnsFilePath)
	if readFileError != nil {
		return fqdns, errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->engines->dns_engine->GenerateSliceOfFQDNs:", readFileError.Error()))
	}

	fqdnsSlice := strings.Split(string(readFileResult), "\n")

	for _, fqdn := range fqdnsSlice {
		fqdns = append(fqdns, fmt.Sprintf("%s.%s", fqdn, topDomain))
	}

	return fqdns, nil
}

func CheckIfFQDNsResolve(fqdns []string, dnsServer string) []string {

	resolvedFQDNs := make([]string, 0)

	for _, fqdn := range fqdns {
		doesFQDNResolveResult, doesFQDNResolveError := doesFQDNResolve(fqdn, dnsServer)
		if doesFQDNResolveError == nil && doesFQDNResolveResult {
			resolvedFQDNs = append(resolvedFQDNs, fqdn)
		}
	}
	return resolvedFQDNs
}

func doesFQDNResolve(fqdn string, dnsServer string) (bool, error) {

	resolveFQDNARecordsResult, _, resolveFQDNARecordsError := resolveFQDNARecords(fqdn, dnsServer)

	if resolveFQDNARecordsError != nil {
		return false, errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->engines->dns_engine->doesFQDNResolve->resolveFQDNARecords:", resolveFQDNARecordsError.Error()))
	}

	if len(resolveFQDNARecordsResult.Answer) == 0 {
		return false, errors.New(fmt.Sprintf("|%s|", "GOLEM->engines->dns_engine->doesFQDNResolve->resolveFQDNARecords: No records returned."))
	}

	for _, ans := range resolveFQDNARecordsResult.Answer {
		switch ans.(type) {
		case *dns.A:
			return true, nil
		default:
			continue
		}
	}

	return false, nil

}

func triageDNSServersAndUpdateLocalCache(testFQDN string, dnsServerIPs []string, privateIPsNets []*net.IPNet, channel chan []string) {

	//resolutionMap := make(map[string]interface{}, 0)
	functionalDNSServers := make(map[string]bool, 0)

	for _, dnsIP := range dnsServerIPs {

		resolveFQDNARecordsResult, _, resolveFQDNARecordsError := resolveFQDNARecords(testFQDN, dnsIP)

		if resolveFQDNARecordsError != nil {
			fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->triageDNSServersAndUpdateLocalCache->resolveFQDNARecords: Resolve: %s Using: %s Result: %s Reason: %s", testFQDN, dnsIP, "[FAIL]", resolveFQDNARecordsError.Error()))
			continue
		}

		if len(resolveFQDNARecordsResult.Answer) == 0 {
			fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->triageDNSServersAndUpdateLocalCache->resolveFQDNARecords: Resolve: %s Using: %s Result: %s Reason: %s", testFQDN, dnsIP, "[FAIL]", "No records returned"))
			continue
		}

	Loop:
		for _, ans := range resolveFQDNARecordsResult.Answer {
			switch ans.(type) {
			case *dns.A:
				if !util.IsIPReserved(ans.(*dns.A).A.String(), privateIPsNets) {
					fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->triageDNSServersAndUpdateLocalCache: Resolve: %s Using: %s Result: %s Status: %s Reason: %s", testFQDN, dnsIP, ans.(*dns.A).A, "[OK]", "Valid IP"))
					functionalDNSServers[dnsIP] = true
				} else {
					fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->triageDNSServersAndUpdateLocalCache: Resolve: %s Using: %s Result: %s (reserved) Status: %s Reason: %s", testFQDN, dnsIP, ans.(*dns.A).A, "[FAILED]", "Reserved IP"))
				}
				break Loop
			default:
				continue
			}
		}
	}

	sliceOfDNSServers := make([]string, 0)
	for dnsIP, _ := range functionalDNSServers {
		sliceOfDNSServers = append(sliceOfDNSServers, dnsIP)
	}

	channel <- sliceOfDNSServers
}

func TriageDNSServersAndUpdateLocalCacheStub(testDomain string, dnsServerIPsMatrix []interface{}, privateIPsNets []*net.IPNet) ([]string, error) {

	dnsServersSlice := make([]string, 0)

	channel := make(chan []string, len(dnsServerIPsMatrix))
	for i := 0; i < len(dnsServerIPsMatrix); i++ {
		go triageDNSServersAndUpdateLocalCache(testDomain, dnsServerIPsMatrix[i].([]string), privateIPsNets, channel)
	}

	for i := 0; i < len(dnsServerIPsMatrix); i++ {
		result := <-channel
		for _, dnsServerIP := range result {
			dnsServersSlice = append(dnsServersSlice, dnsServerIP)
		}
	}

	return dnsServersSlice, nil
}

func bruteDomains(fqdns []string, dnsServers []string, privateIPsNets []*net.IPNet, channel chan map[string]bool) {

	resolutionMap := make(map[string]bool, 0)
	var chosenDNS string

	for _, fqdn := range fqdns {

		/*
			if len(excludedDNSServersMap) == len(dnsServers) {
				fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->bruteDomains: Out of servers to resolve: %s", fqdn))
				continue //continuing so i can tell which hostnames were not resolved...
			}*/

		excludedDNSServersMap := make(map[string]bool, 0)

		for true {

			if len(excludedDNSServersMap) == len(dnsServers) {
				fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->bruteDomains: Out of servers to resolve: %s", fqdn))
				break
			}

			for true {
				chosenDNS = dnsServers[rand.Intn(len(dnsServers)-0)]
				if _, ok := excludedDNSServersMap[chosenDNS]; !ok {
					break
				}
			}

			resolveFQDNARecordsResult, _, resolveFQDNARecordsError := resolveFQDNARecords(fqdn, chosenDNS)

			if resolveFQDNARecordsError != nil {
				fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->bruteDomains->resolveFQDNARecords: Resolve: %s Using: %s Result: %s Reason: %s", fqdn, chosenDNS, "[FAIL]", resolveFQDNARecordsError.Error()))
				excludedDNSServersMap[chosenDNS] = true
				continue
			}

			if len(resolveFQDNARecordsResult.Answer) == 0 {
				fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->bruteDomains->resolveFQDNARecords: Resolve: %s Using: %s Result: %s Reason: %s", fqdn, chosenDNS, "[FAIL]", "No records returned"))
				break //??? if it does not return anything should i break???
			}

		Loop:
			for _, ans := range resolveFQDNARecordsResult.Answer {
				switch ans.(type) {
				case *dns.A:
					if !util.IsIPReserved(ans.(*dns.A).A.String(), privateIPsNets) {
						fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->bruteDomains: Resolve: %s Using: %s Result: %s Status: %s Reason: %s", fqdn, chosenDNS, ans.(*dns.A).A, "[OK]", "Valid IP"))
						resolutionMap[fqdn] = true
					} else {
						fmt.Println(fmt.Sprintf("GOLEM->engines->dns_engine->bruteDomains: Resolve: %s Using: %s Result: %s (reserved) Status: %s Reason: %s", fqdn, chosenDNS, ans.(*dns.A).A, "[FAILED]", "Reserved IP"))
						continue
					}
					break Loop
				default:
					continue
				}
			}

			if _, ok := resolutionMap[fqdn]; !ok {
				continue
			}
			break
		}
	}

	channel <- resolutionMap

}

func BruteFQDNsStub(matrixOfFQDNs []interface{}, dnsServers []string, privateIPsNets []*net.IPNet) ([]string, error) {

	sliceOfFoundFQDNs := make([]string, 0)

	channel := make(chan map[string]bool, len(matrixOfFQDNs))
	for i := 0; i < len(matrixOfFQDNs); i++ {
		go bruteDomains(matrixOfFQDNs[i].([]string), dnsServers, privateIPsNets, channel)
	}

	for i := 0; i < len(matrixOfFQDNs); i++ {
		result := <-channel
		for fqdn := range result {
			sliceOfFoundFQDNs = append(sliceOfFoundFQDNs, fqdn)
		}
	}

	return sliceOfFoundFQDNs, nil
}
