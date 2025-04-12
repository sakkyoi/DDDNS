package main

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/miekg/dns"
	"github.com/sakkyoi/DDDNS/internal/config"
	"github.com/sakkyoi/DDDNS/internal/store"
	"github.com/sakkyoi/DDDNS/internal/store/memory"
	"github.com/sakkyoi/DDDNS/internal/store/redis"
	"golang.org/x/sync/errgroup"
)

var s store.Store
var cfg *config.Config

func init() {
	// initialize DNS server
	dns.HandleFunc(".", dnsRequestHandler)

	// load configuration
	cfg = config.Load()

	// initialize logger
	logLevel, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal("❌ Invalid log level", "error", err)
	}

	log.SetLevel(logLevel)

	// initialize store
	if cfg.RedisHost != "" { // replace with actual condition to check if Redis is available
		if s, err = redis.New(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB, cfg.RedisUser, cfg.RedisPass); err != nil {
			log.Fatal("❌ Failed to connect to Redis", "error", err)
		}
	} else {
		s = memory.New()
	}
}

func main() {
	log.Debug("CFG", "cfg", cfg)
	log.Infof("🚀 Starting DDDNS...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(startDnsServer)
	g.Go(startApiServer)

	if err := g.Wait(); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
		cancel()
	}
}
