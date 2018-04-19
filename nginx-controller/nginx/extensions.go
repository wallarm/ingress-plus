package nginx

import (
	"fmt"
	"strings"
)

// ParseLBMethod parses method and matches it to a corresponding load balancing method in NGINX. An error is returned if method is not valid
func ParseLBMethod(method string) (string, error) {
	method = strings.TrimSpace(method)
	if method == "round_robin" {
		return "", nil
	}
	if strings.HasPrefix(method, "hash") {
		method, err := validateHashLBMethod(method)
		return method, err
	}

	if method == "least_conn" || method == "ip_hash" {
		return method, nil
	}
	return "", fmt.Errorf("Invalid load balancing method: %q", method)
}

var nginxPlusLBValidInput = map[string]bool{
	"least_time":                    true,
	"last_byte":                     true,
	"least_conn":                    true,
	"ip_hash":                       true,
	"least_time header":             true,
	"least_time last_byte":          true,
	"least_time header inflight":    true,
	"least_time last_byte inflight": true,
}

// ParseLBMethodForPlus parses method and matches it to a corresponding load balancing method in NGINX Plus. An error is returned if method is not valid
func ParseLBMethodForPlus(method string) (string, error) {
	method = strings.TrimSpace(method)
	if method == "round_robin" {
		return "", nil
	}
	if strings.HasPrefix(method, "hash") {
		method, err := validateHashLBMethod(method)
		return method, err
	}

	if _, exists := nginxPlusLBValidInput[method]; exists {
		return method, nil
	}
	return "", fmt.Errorf("Invalid load balancing method: %q", method)
}

func validateHashLBMethod(method string) (string, error) {
	keyWords := strings.Split(method, " ")
	if keyWords[0] == "hash" {
		if len(keyWords) == 2 || len(keyWords) == 3 && keyWords[2] == "consistent" {
			return method, nil
		}
	}
	return "", fmt.Errorf("Invalid load balancing method: %q", method)
}
