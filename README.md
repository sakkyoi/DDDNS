# Dynamic Destination DNS (DDDNS)

> **⚠️ Note**: This project is experimental and not recommended for production use.

## The Idea

The goal of DDDNS is to create a DNS server that responds with different private IP addresses for the same domain, depending on the requester's IP address.

## The Challenges

1. The IP address seen by the DNS server is usually that of the DNS resolver, not the actual client.
2. To address this, we use EDNS Client Subnet (ECS) to obtain the original client’s subnet. 
3. However, not all public DNS resolvers forward ECS data, often for privacy reasons. 
4. ECS only provides a subnet, not the full client IP. 
5. DDDNS matches all stored IPs that fall within the same subnet, and returns all of them. 
6. It's up to the client’s Happy Eyeballs algorithm (or equivalent logic) to decide which IP to connect to first.

## The Implementation

### Requirements
- Go 1.24.2 (used for development)

### Optional Dependencies
- Redis

### Building

```bash
$ go build -o dddns ./cmd/dddns
```

### Configuration

```yaml
listen_host: "0.0.0.0"
dns_port: 5355
api_port: 80

ttl: "0"  # seconds, 0 means no expiration (currently not implemented)
domain: "example.com"
mode: "ip"  # ip or ecs

fallback: "127.0.0.1"
fallback_type: "A" # "A" or "CNAME"

# Redis (optional)
# redis_host: "localhost"
# redis_port: 6379
# redis_db: 0
# redis_user: ""
# redis_pass: ""

log_level: "debug"  # debug, info, warning, error, fatal
```

### DNS Setup

You must configure your domain with:

1. An **A** record pointing to your DDDNS server's public IP: <br>
→ e.g. ns.example.com A 1.2.3.4
2. An **NS** record for the target domain pointing to that nameserver: <br>
→ e.g. example.com NS ns.example.com

### API

The API is an HTTP server that allows registering and unregistering destination IPs based on the client's IP address.

#### Register IP

Endpoint: `POST /api/register`

Request Body (JSON):

```json
{
  "domain": ".",
  "dest_ip": "192.168.1.1"
}
```

#### Unregister IP

Endpoint: `DELETE /api/unregister`

Request Body (JSON):

```json
{
  "domain": "."
}
```

> `.` means the root domain. You can also specify a subdomain like `sub` means `sub.example.com`.

## License

This project is licensed under the LGPLv3 License. See the [LICENSE](LICENSE) file for details.
