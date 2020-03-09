# GOLEM

<p align="center">
  <img src="images/golem.jpeg" width=40% heigh=40%>
</p>


# Description
GOLEM is a multi-threaded FQDN bruteforcer with some twists:
- Spreads requests across multiple DNS servers
- Resolution timeouts are solved by choosing another DNS server from the list
- At the end of the bruteforce, a trusted resolver provided by the user is used to check if the FQDNs 


# Dependencies
- Miekg DNS Go Package (https://github.com/miekg/dns)



## The Current Features
- Request spreading across DNS servers
- Multi-threaded

## The Project Structure
My programming projects tend to follow the same structure: 
- Engines are the files where you find the low-level interction with the necessary SDKs.
- The util folder contains classes to write files, process net masks, etc
- Commandline parsers and processors: classes that generate command line processors (out-of-the-box Go flags), process them and call the appropriate engines. 

## Command Line Exampes

Updating local cache:
```
go run main.go -c u -tf www.google.com -t 30 -cf ./dnsServerCacheFile.txt -duf ./dnsURLs.txt 
```

Bruteforcing domain:
```
go run main.go -c b -bd youtube.com -t 40 -rf ./foundHostnames.txt -hw ./dnsbrute-names-large.txt -cf ./dnsServerCacheFile.txt  -vs 8.8.8.8 -vf ./foundHostnamesTriage.txt -duf ./dnsURLs.txt
```