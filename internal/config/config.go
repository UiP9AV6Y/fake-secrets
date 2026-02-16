package config

import (
	"flag"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"net/url"
	"strconv"
	"time"

	env "github.com/caarlos0/env/v11"
)

const (
	EnvPrefix = "FAKE_SECRETS_"

	LogFormatJSON = "JSON"

	CommandHealthCheck = "healthcheck"
	CommandVersion     = "version"
	CommandServe       = "serve"
)

var (
	rootURL, _ = url.Parse("/")
	startTime  = time.Now()
)

type Config struct {
	LogLevel      string        `env:"LOG_LEVEL"`
	LogFormat     string        `env:"LOG_FORMAT"`
	StorageDir    string        `env:"STORAGE_DIR"`
	PathPrefix    string        `env:"PATH_PREFIX"`
	ReadTimeout   time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout  time.Duration `env:"WRITE_TIMEOUT"`
	ListenAddress string        `env:"LISTEN_ADDRESS"`
	ListenPort    int           `env:"LISTEN_PORT"`
	RandomSeed    int64         `env:"RANDOM_SEED"`

	Command string `env:"-"`

	name string
}

func New(name string) *Config {
	result := &Config{
		name:         name,
		LogLevel:     slog.LevelInfo.String(),
		LogFormat:    LogFormatJSON,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 10 * time.Second,
		ListenPort:   8080,
	}

	return result
}

func (c *Config) String() string {
	return c.name
}

func (c *Config) LoadEnv() error {
	opts := env.Options{
		Prefix: EnvPrefix,
	}

	return env.ParseWithOptions(c, opts)
}

func (c *Config) LoadArgs(args []string) (func(io.Writer), error) {
	fs := flag.NewFlagSet(c.name, flag.ContinueOnError)
	help := fs.Bool("help", false, "show this message and exit")

	fs.StringVar(&c.LogLevel, "log.level", c.LogLevel, "log verbosity")
	fs.StringVar(&c.LogFormat, "log.format", c.LogFormat, "log output format")
	fs.StringVar(&c.StorageDir, "storage.dir", c.StorageDir, "base directory to serve secrets from")

	fs.StringVar(&c.PathPrefix, "http.path-prefix", c.PathPrefix, "URL prefix under which to serve API requests")
	fs.DurationVar(&c.ReadTimeout, "http.read-timeout", c.ReadTimeout, "maximum duration for reading the entire HTTP request")
	fs.DurationVar(&c.WriteTimeout, "http.write-timeout", c.WriteTimeout, "maximum duration before timing out writes of the response")
	fs.StringVar(&c.ListenAddress, "listen.address", c.ListenAddress, "interface address to bind to")
	fs.IntVar(&c.ListenPort, "listen.port", c.ListenPort, "tcp port to bind to")
	fs.Int64Var(&c.RandomSeed, "random.seed", c.RandomSeed, "seed for the pseudo-random generator")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *help {
		usage := func(w io.Writer) {
			_, _ = fmt.Fprintf(w, "Usage of %s:\n", c.name)
			fs.SetOutput(w)
			fs.PrintDefaults()
		}

		return usage, nil
	}

	if arg := fs.Arg(0); arg == "" {
		c.Command = CommandServe
	} else {
		c.Command = arg
	}

	return nil, nil
}

func (c *Config) LogVerbosity() slog.Level {
	result := slog.LevelInfo

	if c.LogLevel == "" {
		return result
	} else if err := result.UnmarshalText([]byte(c.LogLevel)); err != nil {
		return result
	}

	return result
}

func (c *Config) LogHandler(w io.Writer) slog.Handler {
	opts := &slog.HandlerOptions{
		Level: c.LogVerbosity(),
	}

	if c.LogFormat == "" || c.LogFormat == "json" || c.LogFormat == LogFormatJSON {
		return slog.NewJSONHandler(w, opts)
	}

	return slog.NewTextHandler(w, opts)
}

func (c *Config) Listen() string {
	port := strconv.FormatInt(int64(c.ListenPort), 10)

	return net.JoinHostPort(c.ListenAddress, port)
}

func (c *Config) HandlerPattern(p ...string) string {
	elem := append([]string{c.PathPrefix}, p...)

	return rootURL.JoinPath(elem...).Path
}

func (c *Config) RandomSeedTime() time.Time {
	if c.RandomSeed == 0 {
		return startTime
	}

	return time.Unix(c.RandomSeed, 0)
}

func (c *Config) RandomSource() rand.Source {
	return rand.NewSource(c.RandomSeedTime().UnixNano())
}

func (c *Config) RandomGenerator() *rand.Rand {
	return rand.New(c.RandomSource()) //nolint:gosec
}

func (c *Config) SelfURL(api string) (*url.URL, error) {
	port := strconv.FormatInt(int64(c.ListenPort), 10)
	addr := net.JoinHostPort("127.0.0.1", port)
	path := c.HandlerPattern(api)

	return url.Parse("http://" + addr + path)
}
