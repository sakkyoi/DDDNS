package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/miekg/dns"
)

func startDnsServer() error {
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost, cfg.DNSPort)

	server := &dns.Server{Addr: addr, Net: "udp"}

	log.Infof("🚀 DNS server listening on %s (UDP)", addr)

	return server.ListenAndServe()
}
