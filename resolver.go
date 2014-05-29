package main

import (
		"net"
		"fmt"
		"os"
		"bufio"
		"log"
		"flag"
	)
	
type AddressPair struct {
	Address string
	IP string
}

func main() {
	var domain_file = flag.String("domains", "domains.txt", "Text file to read domain names from")
	var output_file = flag.String("output", "hosts", "Text file to write hosts to")
	flag.Parse()
	domain_array, domain_count := load_domains_from_file(*domain_file)
	address_channel := make(chan AddressPair)
	complete_channel := make(chan int)
	for _, domain := range domain_array {
		go resolve_to_channel(domain, address_channel, complete_channel)
	}
	go write_chan_to_file(address_channel, *output_file)
	for i := 0; i < domain_count; i++ {
		<- complete_channel
	}
	close(address_channel)
}

func load_domains_from_file(filename string) ([]string, int) {
	domain_count := 0
	domains := make([]string, 0)
	file, error := os.Open(filename)
	if error != nil {
		log.Fatal(error)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
		domain_count = domain_count + 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	return domains, domain_count
}

func resolve_to_channel(domain string, output chan AddressPair, done chan int) {
	ip_addr, error := net.ResolveIPAddr("ip4", domain)
	if error == nil {
		output <- AddressPair{domain,ip_addr.String()}
	} else {
		fmt.Println(error)
	}
	done <- 0
}

func write_chan_to_file(address_stream chan AddressPair, filename string) {
	file, error := os.Create(filename)
	if error != nil {
		log.Fatal(error)
	}
	for next_address := range address_stream {
		file.WriteString(fmt.Sprintf("%s\t%s\n", next_address.Address, next_address.IP))
	}
	file.Close()
}
