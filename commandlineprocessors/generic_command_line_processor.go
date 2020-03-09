package commandlineprocessors

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/golem/engines"
	"github.com/golem/util"
)

func updateLocalCache(commandLineMap map[string]interface{}) error {
	//First we get all the command line parameters so we fail fast if something is missing

	var useURLsFileAsSource bool

	var dnsServersURLsSourcesInputFileTempString string
	if dnsServersURLsSourcesInputFile, ok := commandLineMap["dns_servers_url_sources_input_file"].(string); ok && dnsServersURLsSourcesInputFile != "" {
		dnsServersURLsSourcesInputFileTempString = dnsServersURLsSourcesInputFile
		useURLsFileAsSource = true
	} else {
		useURLsFileAsSource = false
	}

	var cacheInputOutputFileTempString string
	if cacheInputOutputFile, ok := commandLineMap["cache_input_output_file"].(string); ok && cacheInputOutputFile != "" {
		cacheInputOutputFileTempString = cacheInputOutputFile
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache: -cf argument missing"))
	}

	var numberOfThreadsTempInt int
	if numberOfThreads, ok := commandLineMap["number_of_threads"].(int); ok {
		numberOfThreadsTempInt = numberOfThreads
	} else {
		//does not happen. argument always there by default.
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache: -t argument missing"))
	}

	var testFQDNTempString string
	if testFQDN, ok := commandLineMap["test_fqdn"].(string); ok && testFQDN != "" {
		testFQDNTempString = testFQDN
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache: -th argument missing"))
	}

	getListOfReservedNetCIDRsResult, getListOfReservedNetCIDRsError := util.GetListOfReservedNetCIDRs()
	if getListOfReservedNetCIDRsError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->main->util.GetListOfReservedNetCIDRs:", getListOfReservedNetCIDRsError.Error()))
	}

	//Now we perform the calls

	var getDNSIPsResult []string
	var getDNSIPsError error
	if useURLsFileAsSource {
		readFileResult, readFileError := ioutil.ReadFile(dnsServersURLsSourcesInputFileTempString)
		if readFileError != nil {
			return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache->ioutil.ReadFile:", readFileError.Error()))
		}

		getDNSIPsResult, getDNSIPsError = util.GetDNSIPs(false, cacheInputOutputFileTempString, strings.Split(strings.Trim(string(readFileResult), "\n"), "\n"))
	} else {
		getDNSIPsResult, getDNSIPsError = util.GetDNSIPs(true, cacheInputOutputFileTempString, []string{})
	}

	if getDNSIPsError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache->util.GetDNSIPs:", getDNSIPsError.Error()))
	}

	splitArrayResult := util.SplitSliceIntoSubslices(getDNSIPsResult, numberOfThreadsTempInt)

	triageDNSServersAndUpdateLocalCacheResult, triageDNSServersAndUpdateLocalCacheError := engines.TriageDNSServersAndUpdateLocalCacheStub(testFQDNTempString, splitArrayResult, getListOfReservedNetCIDRsResult)
	if triageDNSServersAndUpdateLocalCacheError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache->engines.triageDNSServersAndUpdateLocalCache:", triageDNSServersAndUpdateLocalCacheError.Error()))
	}

	writeArrayOfStringsToFileError := util.WriteArrayOfStringsToFile(triageDNSServersAndUpdateLocalCacheResult, cacheInputOutputFileTempString)
	if writeArrayOfStringsToFileError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->updateLocalCache->util.WriteArrayOfStringsToFile:", writeArrayOfStringsToFileError.Error()))
	}

	return nil

}

func bruteforceDomains(commandLineMap map[string]interface{}) error {

	//First we get all the command line parameters so we fail fast if something is missing

	var cacheInputOutputFileTempString string
	if cacheInputOutputFile, ok := commandLineMap["cache_input_output_file"].(string); ok && cacheInputOutputFile != "" {
		cacheInputOutputFileTempString = cacheInputOutputFile
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->bruteforceDomains: -cf argument missing"))
	}

	var baseDomainTempString string
	if baseDomain, ok := commandLineMap["base_domain"].(string); ok && baseDomain != "" {
		baseDomainTempString = baseDomain
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->bruteforceDomains: -bd argument missing"))
	}

	var numberOfThreadsTempInt int
	if numberOfThreads, ok := commandLineMap["number_of_threads"].(int); ok {
		numberOfThreadsTempInt = numberOfThreads
	} else {
		//does not happen. argument always there by default.
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor: -t argument missing"))
	}

	var hostsWordlistFileTempString string
	if hostsWordlistFile, ok := commandLineMap["hosts_wordlist_wile"].(string); ok && hostsWordlistFile != "" {
		hostsWordlistFileTempString = hostsWordlistFile
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->bruteforceDomains: -hw argument missing"))
	}

	var resolvedFQDNsTempString string
	if resolvedFQDNs, ok := commandLineMap["resolve_fqdns"].(string); ok && resolvedFQDNs != "" {
		resolvedFQDNsTempString = resolvedFQDNs
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor: -rf argument missing"))
	}

	var validationDNSServerTempString string
	var validDNSOutputFileTempString string
	var validateFQDNs bool = false

	if validationDNSServer, ok := commandLineMap["validation_dns_server"].(string); ok && validationDNSServer != "" {
		validationDNSServerTempString = validationDNSServer

		if validDNSOutputFile, ok := commandLineMap["valid_dns_output_file"].(string); ok && validDNSOutputFile != "" {
			validDNSOutputFileTempString = validDNSOutputFile
			validateFQDNs = true
		} else {
			return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor: -vf argument missing"))
		}
	} else {
		return errors.New(fmt.Sprintf("|%s|", "GOLEM->commandlineprocessors->generic_command_line_processor: -vs argument missing"))
	}

	//Now we perform the calls
	getListOfReservedNetCIDRsResult, getListOfReservedNetCIDRsError := util.GetListOfReservedNetCIDRs()
	if getListOfReservedNetCIDRsError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->bruteforceDomains->util.GetListOfReservedNetCIDRs:", getListOfReservedNetCIDRsError.Error()))
	}

	getDNSIPsResult, getDNSIPsError := util.GetDNSIPs(true, cacheInputOutputFileTempString, []string{})
	if getDNSIPsError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->bruteforceDomains->util.GetDNSIPs:", getDNSIPsError.Error()))
	}

	generateListOfFQDNsResult, generateListOfFQDNsError := engines.GenerateSliceOfFQDNs(baseDomainTempString, hostsWordlistFileTempString)
	if generateListOfFQDNsError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->engines.GenerateListOfFQDNs:", generateListOfFQDNsError.Error()))
	}

	splitArrayResult := util.SplitSliceIntoSubslices(generateListOfFQDNsResult, numberOfThreadsTempInt)

	bruteFQDNsStubResult, bruteFQDNsStubError := engines.BruteFQDNsStub(splitArrayResult, getDNSIPsResult, getListOfReservedNetCIDRsResult)
	if bruteFQDNsStubError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->bruteDomainsStub->engines.BruteFQDNsStub:", bruteFQDNsStubError.Error()))
	}

	writeArrayOfStringsToFileError := util.WriteArrayOfStringsToFile(bruteFQDNsStubResult, resolvedFQDNsTempString)
	if writeArrayOfStringsToFileError != nil {
		return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->util.WriteArrayOfStringsToFile:", writeArrayOfStringsToFileError.Error()))
	}

	if validateFQDNs {
		checkIfHostnamesResolveResult := engines.CheckIfFQDNsResolve(bruteFQDNsStubResult, validationDNSServerTempString)

		writeArrayOfStringsToFileError = util.WriteArrayOfStringsToFile(checkIfHostnamesResolveResult, validDNSOutputFileTempString)
		if writeArrayOfStringsToFileError != nil {
			return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->commandlineprocessors->generic_command_line_processor->util.WriteArrayOfStringsToFile:", writeArrayOfStringsToFileError.Error()))
		}
	}

	return nil
}

func ProcessCommandLine(commandLineMap map[string]interface{}) error {

	if opt, ok := commandLineMap["option"]; ok {
		switch opt.(string) {
		case "u":
			updateLocalCacheError := updateLocalCache(commandLineMap)
			if updateLocalCacheError != nil {
				return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->main->ProcessCommandLine->updateLocalCache:", updateLocalCacheError.Error()))
			}
		case "b":
			bruteforceDomainsError := bruteforceDomains(commandLineMap)
			if bruteforceDomainsError != nil {
				return errors.New(fmt.Sprintf("|%s:%s|", "GOLEM->main->ProcessCommandLine->bruteforceDomains:", bruteforceDomainsError.Error()))
			}
		default:
			return errors.New("Invalid option")

		}

	}

	return nil

}
