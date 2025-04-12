package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/miekg/dns"
	"net"
	"strings"
)

func dnsRequestHandler(w dns.ResponseWriter, r *dns.Msg) {
	for _, q := range r.Question {
		log.Info("Got New DNS Request", "from", w.RemoteAddr().String(), "id", r.Id, "qtype", dns.TypeToString[q.Qtype], "qname", q.Name)
		log.Debug("DNS Request", "body", r)
	}

	var sourceIp string

	for _, extra := range r.Extra {
		if opt, ok := extra.(*dns.OPT); ok {
			for _, subOpt := range opt.Option {
				if ecs, ok := subOpt.(*dns.EDNS0_SUBNET); ok {
					log.Debug("ECS Info", "ip", ecs.Address, "mask", ecs.SourceNetmask)
					sourceIp = ecs.Address.String()
				}
			}
		}
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		// Lookup from store
		destIp, err := s.Lookup(fmt.Sprintf("%s.", strings.ToLower(q.Name)), sourceIp)
		if err != nil {
			destIp = "127.0.0.1" // TODO this is a placeholder, need to set to a fallback ip
		}

		if q.Qtype == dns.TypeA {
			a := &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP(destIp),
			}

			m.Answer = append(m.Answer, a)
		}
	}

	if err := w.WriteMsg(m); err != nil {
		log.Errorf("Failed to write DNS response: %v", err)
	}
}
