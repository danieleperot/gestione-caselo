package config

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient interface {
	GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
}

type cache struct {
	mu            sync.RWMutex
	client        SSMClient
	environment   string
	cacheTTL      time.Duration
	lastFetch     time.Time
	parameters    map[string]string
	initialized   bool
	initOnce      sync.Once
	initErr       error
}

var globalCache = &cache{
	cacheTTL:   5 * time.Minute,
	parameters: make(map[string]string),
}

func GetEnvVariable(key string) string {
	globalCache.initOnce.Do(func() {
		ctx := context.Background()

		environment := os.Getenv("ENVIRONMENT")
		if environment == "" {
			environment = "local"
		}
		globalCache.environment = environment

		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			globalCache.initErr = fmt.Errorf("failed to load AWS config: %w", err)
			return
		}

		globalCache.client = ssm.NewFromConfig(cfg)
		globalCache.initialized = true
	})

	if globalCache.initErr != nil || !globalCache.initialized {
		return os.Getenv(key)
	}

	ctx := context.Background()

	globalCache.mu.RLock()
	if time.Since(globalCache.lastFetch) < globalCache.cacheTTL {
		if val, ok := globalCache.parameters[key]; ok {
			globalCache.mu.RUnlock()
			return val
		}
	}
	globalCache.mu.RUnlock()

	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	if time.Since(globalCache.lastFetch) < globalCache.cacheTTL {
		if val, ok := globalCache.parameters[key]; ok {
			return val
		}
	}

	paramName := fmt.Sprintf("/gestione-caselo/%s/%s", globalCache.environment, key)
	result, err := globalCache.client.GetParameters(ctx, &ssm.GetParametersInput{
		Names:          []string{paramName},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return os.Getenv(key)
	}

	if len(result.Parameters) > 0 {
		globalCache.parameters[key] = *result.Parameters[0].Value
		globalCache.lastFetch = time.Now()
		return *result.Parameters[0].Value
	}

	return os.Getenv(key)
}
