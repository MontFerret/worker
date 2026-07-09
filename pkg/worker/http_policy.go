package worker

import (
	"fmt"
	"strings"
	"time"

	ferretnet "github.com/MontFerret/ferret/v2/pkg/net"
	ferrethttp "github.com/MontFerret/ferret/v2/pkg/net/http"
)

const (
	defaultHTTPTimeout         = 10 * time.Second
	defaultHTTPMaxRequestSize  = 1024 * 1024
	defaultHTTPMaxResponseSize = 10 * 1024 * 1024
	defaultHTTPMaxRedirects    = 3
)

// HTTPPolicy configures the Ferret network client used by Worker.
// Start from DefaultHTTPPolicy and override the fields needed by a deployment.
type HTTPPolicy struct {
	AllowedSchemes        []string
	AllowedHosts          []string
	BlockedHosts          []string
	BlockedRequestHeaders []string
	Timeout               time.Duration
	MaxRequestSize        int64
	MaxResponseSize       int64
	MaxRedirects          int
	FollowRedirects       bool
	AllowAllHosts         bool
	AllowLocalhost        bool
	AllowPrivateNetworks  bool
}

// DefaultHTTPPolicy returns Worker's restrictive default Ferret HTTP policy.
func DefaultHTTPPolicy() HTTPPolicy {
	return HTTPPolicy{
		AllowedSchemes:        []string{"http", "https"},
		BlockedRequestHeaders: []string{"Authorization", "Cookie", "Proxy-Authorization"},
		Timeout:               defaultHTTPTimeout,
		MaxRequestSize:        defaultHTTPMaxRequestSize,
		MaxResponseSize:       defaultHTTPMaxResponseSize,
		MaxRedirects:          defaultHTTPMaxRedirects,
		FollowRedirects:       true,
	}
}

func validateHTTPPolicy(policy HTTPPolicy) error {
	policy = normalizeHTTPPolicy(policy)
	if policy.AllowAllHosts && len(policy.AllowedHosts) > 0 {
		return fmt.Errorf("http allowed hosts and allow-all hosts cannot both be set")
	}

	return nil
}

func newNetwork(policy HTTPPolicy) ferretnet.Network {
	policy = normalizeHTTPPolicy(policy)

	client := ferrethttp.Client(disabledHTTPClient{})
	if policy.AllowAllHosts || len(policy.AllowedHosts) > 0 {
		client = ferrethttp.New(httpPolicyOptions(policy)...)
	}

	return ferretnet.New(ferretnet.WithHTTPClient(client))
}

func httpPolicyOptions(policy HTTPPolicy) []ferrethttp.Policy {
	opts := []ferrethttp.Policy{
		ferrethttp.WithAllowedSchemes(policy.AllowedSchemes...),
		ferrethttp.WithBlockedHosts(policy.BlockedHosts...),
		ferrethttp.WithBlockedRequestHeaders(policy.BlockedRequestHeaders...),
		ferrethttp.WithTimeout(policy.Timeout),
		ferrethttp.WithMaxRequestSize(policy.MaxRequestSize),
		ferrethttp.WithMaxResponseSize(policy.MaxResponseSize),
		ferrethttp.WithMaxRedirects(policy.MaxRedirects),
		ferrethttp.WithFollowRedirects(policy.FollowRedirects),
		ferrethttp.WithAllowLocalhost(policy.AllowLocalhost),
		ferrethttp.WithAllowPrivateNetworks(policy.AllowPrivateNetworks),
	}

	if !policy.AllowAllHosts {
		opts = append(opts, ferrethttp.WithAllowedHosts(policy.AllowedHosts...))
	}

	return opts
}

func normalizeHTTPPolicy(policy HTTPPolicy) HTTPPolicy {
	policy.AllowedSchemes = cleanList(policy.AllowedSchemes)
	policy.AllowedHosts = cleanList(policy.AllowedHosts)
	policy.BlockedHosts = cleanList(policy.BlockedHosts)
	policy.BlockedRequestHeaders = cleanList(policy.BlockedRequestHeaders)

	return policy
}

func cleanList(values []string) []string {
	if values == nil {
		return nil
	}

	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		cleaned = append(cleaned, value)
	}

	return cleaned
}
