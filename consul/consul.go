package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog"
)

type Config struct {
	StartTime           time.Time
	ConsulAddr          string        `mapstructure:"CONSUL_ADDRESS"`
	ServerName          string        `mapstructure:"SERVER_NAME"`
	ServiceName         string        `mapstructure:"SERVICE_NAME"`
	HealthCheckInterval time.Duration `mapstructure:"HEALTH_CHECK_INTERVAL"`
	HealthCheckTimeout  time.Duration `mapstructure:"HEALTH_CHECK_TIMEOUT"`
	RestPort            int           `mapstructure:"REST_PORT"`
	GrpcPort            int           `mapstructure:"GRPC_PORT"`
}

type Client struct {
	logr   zerolog.Logger
	client *api.Client
	config Config
}

func NewClient(logr zerolog.Logger, config Config) (*Client, error) {
	client, err := api.NewClient(&api.Config{Address: config.ConsulAddr})
	if err != nil {
		return nil, err
	}
	return &Client{client: client, config: config, logr: logr}, nil
}

func (c *Client) Register() (err error) {
	timeout := c.config.HealthCheckTimeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	interval := c.config.HealthCheckInterval
	if interval == 0 {
		interval = 60 * time.Second
	}

	grpcAddress := fmt.Sprintf("%s:%d", c.config.ServerName, c.config.GrpcPort)

	registration := &api.AgentServiceRegistration{
		Name: c.config.ServiceName,
		ID:   c.config.ServerName,
		Port: c.config.RestPort,
		Tags: []string{c.config.ServiceName},
		Meta: map[string]string{
			"grpc_port":    fmt.Sprintf("%d", c.config.GrpcPort),
			"rest_port":    fmt.Sprintf("%d", c.config.RestPort),
			"server_name":  c.config.ServerName,
			"service_name": c.config.ServiceName,
			"start_time":    c.config.StartTime.String(),
		},
		Check: &api.AgentServiceCheck{
			Interval: interval.String(),
			Timeout:  timeout.String(),
			GRPC:     grpcAddress,
		},
	}

	err = c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	return
}

func (c *Client) Deregister() {
	if err := c.client.Agent().ServiceDeregister(c.config.ServerName); err != nil {
		c.logr.Error().Err(err).Msg("failed to deregister service")
	}
}
