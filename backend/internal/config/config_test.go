package config

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type mockSSMClient struct {
	getParametersFunc func(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error)
	callCount         int
	mu                sync.Mutex
}

func (m *mockSSMClient) GetParameters(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
	m.mu.Lock()
	m.callCount++
	m.mu.Unlock()
	return m.getParametersFunc(ctx, params, optFns...)
}

func (m *mockSSMClient) getCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.callCount
}

func TestGetEnvVariable_FallbackToEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "from_env")
	defer os.Unsetenv("TEST_VAR")

	globalCache.initialized = false
	globalCache.initErr = nil
	globalCache.initOnce = sync.Once{}

	val := GetEnvVariable("TEST_VAR")
	if val != "from_env" {
		t.Errorf("Expected 'from_env', got %s", val)
	}
}

func TestGetEnvVariable_SSMSuccess(t *testing.T) {
	mockClient := &mockSSMClient{
		getParametersFunc: func(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
			return &ssm.GetParametersOutput{
				Parameters: []types.Parameter{
					{
						Name:  aws.String("/gestione-caselo/test/TEST_VAR"),
						Value: aws.String("from_ssm"),
					},
				},
			}, nil
		},
	}

	globalCache.mu.Lock()
	globalCache.client = mockClient
	globalCache.environment = "test"
	globalCache.initialized = true
	globalCache.initErr = nil
	globalCache.parameters = make(map[string]string)
	globalCache.lastFetch = time.Time{}
	globalCache.mu.Unlock()

	val := GetEnvVariable("TEST_VAR")
	if val != "from_ssm" {
		t.Errorf("Expected 'from_ssm', got %s", val)
	}

	if mockClient.getCallCount() != 1 {
		t.Errorf("Expected 1 SSM call, got %d", mockClient.getCallCount())
	}
}

func TestGetEnvVariable_CachingWorks(t *testing.T) {
	mockClient := &mockSSMClient{
		getParametersFunc: func(ctx context.Context, params *ssm.GetParametersInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersOutput, error) {
			return &ssm.GetParametersOutput{
				Parameters: []types.Parameter{
					{
						Name:  aws.String("/gestione-caselo/test/CACHED_VAR"),
						Value: aws.String("cached_value"),
					},
				},
			}, nil
		},
	}

	globalCache.mu.Lock()
	globalCache.client = mockClient
	globalCache.environment = "test"
	globalCache.initialized = true
	globalCache.initErr = nil
	globalCache.parameters = make(map[string]string)
	globalCache.lastFetch = time.Time{}
	globalCache.mu.Unlock()

	val1 := GetEnvVariable("CACHED_VAR")
	val2 := GetEnvVariable("CACHED_VAR")

	if val1 != "cached_value" || val2 != "cached_value" {
		t.Errorf("Expected both values to be 'cached_value', got %s and %s", val1, val2)
	}

	if mockClient.getCallCount() != 1 {
		t.Errorf("Expected 1 SSM call (cached on second access), got %d", mockClient.getCallCount())
	}
}
