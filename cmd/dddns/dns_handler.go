package main

import (
	"github.com/charmbracelet/log"
	"github.com/miekg/dns"
	"github.com/sakkyoi/DDDNS/internal/config"
	"net"
	"strings"
)

func dnsRequestHandler(w dns.ResponseWriter, r *dns.Msg) {
	for _, q := range r.Question {
		log.Info("Got New DNS Request", "from", w.RemoteAddr().String(), "id", r.Id, "qtype", dns.TypeToString[q.Qtype], "qname", q.Name)
		log.Debug("DNS Request", "body", r)
	}

	// Try to get the ecs info from the request
	var ecsIp net.IP
	var ecsMask uint8
	for _, extra := range r.Extra {
		if opt, ok := extra.(*dns.OPT); ok {
			for _, subOpt := range opt.Option {
				if ecs, ok := subOpt.(*dns.EDNS0_SUBNET); ok {
					if ecs.Address != nil {
						log.Debug("ECS Info", "ip", ecs.Address, "mask", ecs.SourceNetmask)
						ecsIp = ecs.Address
						ecsMask = ecs.SourceNetmask
					}
				}
			}
		}
	}

	log.Debug("DNS Request", "ecsIp", ecsIp, "ecsMask", ecsMask)

	// Build the response
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		// Only process A records
		if q.Qtype != dns.TypeA {
			continue
		}

		var rr []dns.RR // Resource Records

		// Check if the ECS info is available
		var mask *uint8
		var ip string
		if cfg.Mode == config.EcsMode {
			if ecsIp == nil {
				log.Error("ECS IP is nil, but ECS mode is enabled")
				continue
			}

			mask = &ecsMask
			ip = ecsIp.String()
		} else {
			ip = strings.Split(w.RemoteAddr().String(), ":")[0]
		}

		// Lookup from store
		destIps, err := s.Lookup(strings.ToLower(q.Name), ip, mask)
		if err != nil {
			log.Debug("Failed to lookup IP", "error", err)
		}

		for _, destIp := range destIps {
			rr = append(rr, &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP(destIp),
			})
		}

		// fallback
		if cfg.FallbackType == "A" {
			rr = append(rr, &dns.A{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				A: net.ParseIP(cfg.Fallback),
			})
		} else if cfg.FallbackType == "CNAME" {
			rr = append(rr, &dns.CNAME{
				Hdr: dns.RR_Header{
					Name:   q.Name,
					Rrtype: dns.TypeCNAME,
					Class:  dns.ClassINET,
					Ttl:    60,
				},
				Target: cfg.Fallback,
			})
		}

		m.Answer = append(m.Answer, rr...)
	}

	if err := w.WriteMsg(m); err != nil {
		log.Errorf("Failed to write DNS response: %v", err)
	}
}
