package server

// ServerConfig configuration for the server
type ServerConfig struct {
	host                   string
	port                   int
	defaultCachePartitions int // default number of cache partition to use when create new cache
	githubAPIURL           string
	githubAPIUser          string
}

// NewServerConfig return new configuration instance for server
func NewServerConfig(host string, port int, githubAPIURL string, githubAPIUser string) *ServerConfig {
	return &ServerConfig{
		host:                   host,
		port:                   port,
		defaultCachePartitions: 7,
		githubAPIURL:           githubAPIURL,
		githubAPIUser:          githubAPIUser,
	}
}
