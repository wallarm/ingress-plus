package nginx

// Config holds NGINX configuration parameters
type Config struct {
	ProxyConnectTimeout string
	ProxyReadTimeout    string
}

// NewDefaultConfig creates a Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		ProxyConnectTimeout: "60s",
		ProxyReadTimeout:    "60s",
	}
}
