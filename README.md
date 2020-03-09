# GOLEM

<p align="center">
  <img src="images/golem.jpeg" width=40% heigh=40%>
</p>


# Motivation
Tools don't seem to handle timeouts properly and assume the domains don't exist if so happens. On another note, it seems everyone goes crazy over threads and progress bars but no one thinks about the poor DNS server they are overloading with a bunch of requests. Why not (as we say in my country) "spreading the evil across the villages"? Also, requesting a recursive resolution hides the original IP of the request and let's you sleep better at night knowing Santa won't punish you next Christmas.

# Description
GOLEM is a multi-threaded FQDN bruteforcer with some twists:
- Spreads requests across multiple DNS servers
- Resolution timeouts are solved by choosing another DNS server from the list
- At the end of the bruteforce, a trusted resolver provided by the user is used to check if the FQDNs 

# Operating GOLEM
The command line tool is self-explanatory. However, how it all comes together requires a brief explanation. As explain above, GOLEM uses a list of DNS servers to savoid overloading the same DNS server and address problems with DNS resolution timeouts. The cache file is provided using the -cf flag. This file can be built in two ways:
1. Manually by you with a list of IPs separated by newline
2. Run
```
go run main.go -c u -tf www.google.com -t 30 -cf ./dnsServerCacheFile.txt -duf ./dnsURLs.txt 
```

-c allows you to choose what action is performed. -u means that the local cache is created (using a file containing URLs for files with DNS IPs separated by newlines) while -b assumes the existance of the cache file (passed with -cf), the hostnames wordlist (passed with -hw) and performs the bruteforce of FQDNs. You still need the -cf in this case since the resulting list of servers will be saved to the file you provide. The cache is created by pulling the DNSIPs, filtering out IPv6, and testing them with a A query for the FQDN you provide with -tf flag.

Once the cache is built you can run the bruteforce:

```
go run main.go -c b -bd youtube.com -t 40 -rf ./foundHostnames.txt -hw ./dnsbrute-names-large.txt -cf ./dnsServerCacheFile.txt  -vs 8.8.8.8 -vf ./foundHostnamesTriage.txt -duf ./dnsURLs.txt
```
 

The output is quite verbose so you can troubleshoot and make sure you don't miss any resolutions. GOLEM spreads the list of FQDNs to resolve across multiple GoO routines using the -t flag. Every thread performs the resolution by choosing a DNS server at random from the list of available DNS servers in the cache. Per testing i have observed three possible outcomes for each A resolution:
- No records found: GOLEM assumes the FQDN does not exist
- Timeout: this usually happens (or so it seems) when the server does not know the FQDN and performs the resolution. Later, if you perform the query again it is very likely that the record will return. It may also mean the DNS server is dead and/or not responding. 
- Invalid IP: sometimes, some DNS servers respond to some queries with invalid IPs (e.g. private IPs, IPs that don't belong to the Net range of the FQDN/domain owner). GOLEM performs a check on the IPs to make sure they don't fall in reserved ranges and if all a server has to offer is that, a new one is checked until a definitive empty response (i.e. no A records found) is returned.

The premisse here is: if the initial cache building and testing narrowed down the list of servers to a small test that can resolve the test FQDN you provided, then any records not found for a given FQDN in the bruteforcing process tells us that that domain does not exist. In the unfortunate case of a DNS server returning an IP that is not reserved and valid, GOLEM will assume that the FQDN exists. The last filter can be used by passing the flags -vs [TRUSTED_DNS_SERVER_IP] and -vf [FILTERED_LIST_OF_FQDNS]: all the FQDNs are again resolved against a trusted server (e.g. 8.8.8.8, 1.1.1.1). This last step is simple: perform A query for each FQDN. If an A record returns, it's assumed that the FQDN exists. Otherwise, it is filtered out.  




# Dependencies
- Miekg DNS Go Package (https://github.com/miekg/dns)


# The Current Features
- Request spreading across DNS servers
- Multi-threaded

# The Project Structure
My programming projects tend to follow the same structure: 
- Engines are the files where you find the low-level interction with the necessary SDKs.
- The util folder contains classes to write files, process net masks, etc
- Commandline parsers and processors: classes that generate command line processors (out-of-the-box Go flags), process them and call the appropriate engines. 

# Things to improve
- Support IPv6