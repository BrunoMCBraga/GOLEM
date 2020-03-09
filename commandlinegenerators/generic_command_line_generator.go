package commandlinegenerators

import (
	"flag"
	"strings"

	"github.com/golem/globalstringsproviders"
)

var option *string
var baseDomain *string
var testFQDN *string

var numberOfThreads *int

var hostsWordlistFile *string
var resolvedFQDNs *string
var cacheInputOutputFile *string
var validDNSOutputFile *string

var dnsServersURLSourcesInputFile *string

var validationDNSServer *string

func PrepareCommandLineProcessing() {

	optionHelp := globalstringsproviders.GetCommandsMenu()
	//command
	option = flag.String("c", "", strings.TrimLeft(optionHelp, "\n"))

	//input domains/fqdns
	baseDomain = flag.String("bd", "", "Base domain for bruteforce (used with -c b)")
	testFQDN = flag.String("tf", "", "FQDN used to test DNS server viability  (used with -c u)")

	//threading
	numberOfThreads = flag.Int("t", 1, "Number of threads  (used with -c b/u)")

	//input/output files
	hostsWordlistFile = flag.String("hw", "", "Hostnames Wordlist (used with -c b)")
	resolvedFQDNs = flag.String("rf", "", "Output file for resolved FQDNs (used with -c b)")
	cacheInputOutputFile = flag.String("cf", "", "Cache file  (used with -c b/u)")
	validDNSOutputFile = flag.String("vf", "", "Valid resolved FQDNs (used with -c b -ts). This is used in combination with -ts to perform a final resolution of found FQDNs using a trusted server supplied with -ts. The reson for this is because sometimes, some DNS servers return invalid IPs for some FQDNs that actually don't exist")

	dnsServersURLSourcesInputFile = flag.String("duf", "", "Input file containing URLs with lists of available DNS servers (used with -c u). If this flag is not provided with -c u, the -cf file will be used instead and updated")

	//other switches
	validationDNSServer = flag.String("vs", "", "FQDN trusted validation server (used with -c b and -ts)")

}

func ParseCommandLine() {
	flag.Parse()
}

func GetParametersDict() map[string]interface{} {

	parameters := make(map[string]interface{}, 0)
	parameters["option"] = *option
	parameters["base_domain"] = *baseDomain
	parameters["test_fqdn"] = *testFQDN
	parameters["number_of_threads"] = *numberOfThreads

	parameters["hosts_wordlist_wile"] = *hostsWordlistFile
	parameters["resolve_fqdns"] = *resolvedFQDNs
	parameters["cache_input_output_file"] = *cacheInputOutputFile
	parameters["valid_dns_output_file"] = *validDNSOutputFile

	parameters["dns_servers_url_sources_input_file"] = *dnsServersURLSourcesInputFile

	parameters["validation_dns_server"] = *validationDNSServer

	return parameters
}
