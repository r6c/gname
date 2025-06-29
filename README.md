# GNAME for `libdns`

[![godoc reference](https://pkg.go.dev/badge/github.com/r6c/gname.svg)](https://pkg.go.dev/github.com/r6c/gname)

This package implements the [libdns](https://github.com/libdns/libdns) interfaces for [GNAME](https://www.gname.com), allowing you to manage DNS records with GNAME's API.

## Features

- ✅ **Complete libdns implementation**: Supports all CRUD operations for DNS records
- ✅ **Custom HTTP Client**: Configure timeouts and other HTTP client settings
- ✅ **Comprehensive error handling**: Detailed error messages for troubleshooting
- ✅ **Thread-safe**: Safe for concurrent use
- ✅ **Automatic authentication**: Handles GNAME API authentication automatically

## Authenticating

To use this package, you need to obtain API credentials from GNAME:

### Requirements

1. **Create a GNAME account**: [Register here](https://www.gname.com)
2. **Upgrade to reseller account**: Required for API access ([Reseller Plan](https://www.gname.com/domain/api))
3. **Apply for API access**: Submit API service agreement with GNAME
4. **Get your credentials**: You'll receive `APPID` and `APPKey` 

### Important Security Notes

- 🔐 **Keep your APPKey secure**: It should be treated as a password
- 🚫 **Never commit credentials**: Use environment variables or secure configuration
- ✅ **Use least privilege**: Only grant necessary permissions

## Installation

```bash
go get github.com/r6c/gname
```

## Usage Examples

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    
    "github.com/libdns/libdns"
    "github.com/r6c/gname"
)

func main() {
    // Configure the DNS provider
    provider := gname.Provider{
        APPID:  "your_app_id",
        APPKey: "your_app_key", // Keep this secure!
    }
    
    ctx := context.Background()
    zone := "example.com." // Note the trailing dot
    
    // Get all records
    records, err := provider.GetRecords(ctx, zone)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d records\n", len(records))
}
```

### Creating Records

```go
// Create new A record
newRecords, err := provider.AppendRecords(ctx, zone, []libdns.Record{
    libdns.RR{
        Name: "www",        // Creates www.example.com
        Type: "A",
        Data: "192.168.1.1",
        TTL:  300 * time.Second,
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created %d records\n", len(newRecords))
```

### Updating Records

```go
// Update or create records (SetRecords will update if exists, create if not)
updatedRecords, err := provider.SetRecords(ctx, zone, []libdns.Record{
    libdns.RR{
        Name: "www",
        Type: "A", 
        Data: "192.168.1.100", // New IP address
        TTL:  600 * time.Second,
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Updated %d records\n", len(updatedRecords))
```

### Deleting Records

```go
// Delete specific records
deletedRecords, err := provider.DeleteRecords(ctx, zone, []libdns.Record{
    libdns.RR{
        Name: "www",
        Type: "A",
        // Data is not required for deletion, only Name and Type
    },
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Deleted %d records\n", len(deletedRecords))
```

### Advanced Configuration

```go
import "net/http"

// With custom HTTP client (e.g., custom timeout)
provider := gname.Provider{
    APPID:  "your_app_id",
    APPKey: "your_app_key",
    HTTPClient: &http.Client{
        Timeout: 60 * time.Second,
        // Add other custom settings as needed
    },
}
```

### Environment Variables

For better security, use environment variables:

```bash
export GNAME_APP_ID="your_app_id"
export GNAME_APP_KEY="your_app_key"
```

```go
import "os"

provider := gname.Provider{
    APPID:  os.Getenv("GNAME_APP_ID"),
    APPKey: os.Getenv("GNAME_APP_KEY"),
}
```

## Supported Record Types

GNAME supports all standard DNS record types through their API:

- **A** - IPv4 address
- **AAAA** - IPv6 address  
- **CNAME** - Canonical name
- **MX** - Mail exchange
- **TXT** - Text
- **NS** - Name server
- **SRV** - Service
- **And more** - Check GNAME's API documentation for the complete list

## Error Handling

The package provides detailed error messages to help with troubleshooting:

```go
records, err := provider.GetRecords(ctx, zone)
if err != nil {
    // Errors are wrapped and provide context
    fmt.Printf("Failed to get records: %v\n", err)
    // Example output: "failed to get records for zone example.com.: API error: Invalid APPID (code: 0)"
}
```

## Integration Examples

### With Caddy

This provider works seamlessly with [Caddy](https://caddyserver.com/) for automatic HTTPS:

```json
{
    "apps": {
        "tls": {
            "automation": {
                "policies": [{
                    "subjects": ["*.example.com"],
                    "issuers": [{
                        "module": "acme",
                        "challenges": {
                            "dns": {
                                "provider": {
                                    "name": "gname",
                                    "app_id": "{env.GNAME_APP_ID}",
                                    "app_key": "{env.GNAME_APP_KEY}"
                                }
                            }
                        }
                    }]
                }]
            }
        }
    }
}
```

### With cert-manager

Use with Kubernetes [cert-manager](https://cert-manager.io/) for automatic certificate management.

## Rate Limits & Best Practices

- **Respect API limits**: GNAME may have rate limits on their API
- **Cache when possible**: Don't fetch records unnecessarily
- **Use appropriate TTLs**: Set reasonable TTL values for your use case
- **Handle errors gracefully**: Implement proper error handling and retries

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues.

### Development Setup

```bash
git clone https://github.com/r6c/gname
cd gname
go mod download
go test -v
```

## Troubleshooting

### Common Issues

1. **Invalid APPID error**: 
   - Verify your APPID is correct
   - Ensure your account has API access enabled

2. **Authentication failed**:
   - Check that your APPKey is correct and not expired
   - Verify you have a reseller account

3. **Zone not found**:
   - Ensure the domain is managed by your GNAME account
   - Check that the zone name includes the trailing dot (e.g., `example.com.`)

### Getting Help

- 📚 [GNAME API Documentation](https://www.gname.com/domain/api)
- 🐛 [Report Issues](https://github.com/r6c/gname/issues)
- 💬 [libdns Community](https://github.com/libdns/libdns)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the [libdns](https://github.com/libdns/libdns) project for providing the standard interfaces
- Inspired by other libdns providers, especially [libdns/cloudflare](https://github.com/libdns/cloudflare)
