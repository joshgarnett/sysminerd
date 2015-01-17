package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fzzy/radix/redis"
)

const RedisModuleName = "redis"

type RedisInputModule struct {
	Host   string
	Port   int
	client *redis.Client
}

func (m *RedisInputModule) Name() string {
	return RedisModuleName
}

func (m *RedisInputModule) Init(config *Config, moduleConfig *ModuleConfig) error {
	// parse redis settings
	redisHost, err := moduleConfig.SettingsString("host")
	if err != nil || redisHost == "" {
		log.Fatalf("host must be specified: %v", err)
	}

	redisPort, err := moduleConfig.SettingsInt("port")
	if err != nil {
		log.Fatalf("Unable to parse port: %v", err)
	} else if redisPort < 1 || redisPort > 65535 {
		log.Fatalf("invalid port number: %d", redisPort)
	}

	// save config data
	m.Host = redisHost
	m.Port = redisPort

	// connect to redis
	m.client, err = connectToRedis(m.Host, m.Port)

	return err
}

func (m *RedisInputModule) TearDown() error {
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}

func (m *RedisInputModule) GetMetrics() (*ModuleMetrics, error) {
	metrics := make([]Metric, 0, 48)

	// attempt to reconnect to redis
	if m.client == nil {
		client, err := connectToRedis(m.Host, m.Port)
		if err != nil {
			return nil, err
		}
		log.Print("Reconnected to redis")
		m.client = client
	}

	reply := m.client.Cmd("INFO")
	if reply.Err != nil {
		log.Printf("Error collecting metrics from redis: %v", reply.Err)

		// close the existing connection
		m.client.Close()
		m.client = nil

		return nil, reply.Err
	}

	values, err := reply.Str()
	if err != nil {
		log.Printf("Problem processing redis reply: %v", err)
		return nil, err
	}

	lines := strings.Split(values, "\n")

	processStats := false
	processKeyspace := false

	for _, line := range lines {
		if strings.Contains(line, "# Clients") {
			processStats = true
		}
		if strings.Contains(line, "# Keyspace") {
			processKeyspace = true
		}

		// Don't handle stats until we see the # Clients header
		if !processStats {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}

		key := strings.TrimSpace(fields[0])
		svalue := strings.TrimSpace(fields[1])

		// throwaway human formatted values
		if strings.Contains(key, "human") {
			continue
		}

		if processKeyspace {
			// format db0:keys=1,expires=0,avg_ttl=0
			dbValues := strings.Split(svalue, ",")
			for _, dbLine := range dbValues {
				fields = strings.Split(dbLine, "=")
				dbKey := fields[0]
				svalue = fields[1]
				dbValue, err := strconv.ParseFloat(svalue, 64)
				if err != nil {
					continue
				}
				dbKey = fmt.Sprintf("%s.%s", key, dbKey)

				metrics = append(metrics, NewMetric(dbKey, dbValue))
			}
		} else {
			// format used_cpu_user_children:0.00
			// throwaway non numeric values
			value, err := strconv.ParseFloat(svalue, 64)
			if err != nil {
				continue
			}

			metrics = append(metrics, NewMetric(key, value))
		}
	}

	return &ModuleMetrics{Module: m.Name(), Metrics: metrics}, nil
}

func connectToRedis(host string, port int) (*redis.Client, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	client, err := redis.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Printf("Failed to connect to redis: %v", err)
	}

	return client, err
}
