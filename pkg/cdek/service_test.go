package cdek

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sony/gobreaker/v2"
)

func TestNewService(t *testing.T) {
	config := &AccountConfig{
		Name:         "test",
		ClientID:     "id",
		ClientSecret: "secret",
		BaseURL:      URLProduction,
	}

	client, err := NewAuthenticatedClient(config)
	if err != nil {
		t.Fatalf("NewAuthenticatedClient() error = %v", err)
	}

	service := NewService(client, nil)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}
	if service.client != client {
		t.Error("Service client doesn't match")
	}
	if service.breaker == nil {
		t.Error("Circuit breaker not initialized")
	}
	if service.costCalculator == nil {
		t.Error("Cost calculator not initialized")
	}
	if service.orderValidator == nil {
		t.Error("Order validator not initialized")
	}
	if service.parser == nil {
		t.Error("Parser not initialized")
	}
	if service.mapper == nil {
		t.Error("Mapper not initialized")
	}
}

func newTestService(t *testing.T) *Service {
	t.Helper()
	client, err := NewAuthenticatedClient(&AccountConfig{
		Name:         "test",
		ClientID:     "id",
		ClientSecret: "secret",
		BaseURL:      URLProduction,
	})
	if err != nil {
		t.Fatalf("NewAuthenticatedClient() error = %v", err)
	}
	return NewService(client, nil)
}

func TestServiceBreaker_ClientErrorsDoNotTrip(t *testing.T) {
	service := newTestService(t)

	// Множество ошибок валидации/использования не должны размыкать breaker,
	// т.к. они не свидетельствуют о нестабильности CDEK API.
	for i := 0; i < 10; i++ {
		_, _ = service.breaker.Execute(func() (any, error) {
			return nil, fmt.Errorf("bad input: %w", ErrInvalidRequest)
		})
	}
	if state := service.breaker.State(); state != gobreaker.StateClosed {
		t.Errorf("breaker state = %v, want Closed after client-side errors", state)
	}
}

func TestServiceBreaker_ServerErrorsTrip(t *testing.T) {
	service := newTestService(t)

	// Реальные сбои (сеть/сервер) должны по-прежнему размыкать breaker.
	for i := 0; i < 3; i++ {
		_, _ = service.breaker.Execute(func() (any, error) {
			return nil, fmt.Errorf("boom: %w", ErrServerError)
		})
	}
	if state := service.breaker.State(); state != gobreaker.StateOpen {
		t.Errorf("breaker state = %v, want Open after server errors", state)
	}
}

func TestServiceBreakerIsSuccessful(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, true},
		{"invalid request", ErrInvalidRequest, true},
		{"not found", ErrNotFound, true},
		{"unauthorized", ErrUnauthorized, true},
		{"rate limit", ErrRateLimitExceeded, false},
		{"server error", ErrServerError, false},
		{"generic/network error", errors.New("connection reset"), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := newTestService(t)
			before := service.breaker.Counts()
			_, _ = service.breaker.Execute(func() (any, error) {
				return nil, tc.err
			})
			after := service.breaker.Counts()

			gotSuccess := after.TotalSuccesses > before.TotalSuccesses
			if gotSuccess != tc.want {
				t.Errorf("IsSuccessful(%v) treated as success = %v, want %v", tc.err, gotSuccess, tc.want)
			}
		})
	}
}

func TestService_GetClient(t *testing.T) {
	config := &AccountConfig{
		Name:         "test",
		ClientID:     "id",
		ClientSecret: "secret",
		BaseURL:      URLProduction,
	}

	authClient, _ := NewAuthenticatedClient(config)
	service := NewService(authClient, nil)

	client := service.GetClient()
	if client == nil {
		t.Error("GetClient() returned nil")
	}
	if client != authClient {
		t.Error("GetClient() returned different client")
	}
}

func TestService_HealthCheck(t *testing.T) {
	t.Run("fails for invalid credentials", func(t *testing.T) {
		config := &AccountConfig{
			Name:         "test",
			ClientID:     "invalid",
			ClientSecret: "invalid",
			BaseURL:      URLProduction,
			Timeout:      2 * time.Second,
		}

		client, _ := NewAuthenticatedClient(config)
		service := NewService(client, nil)

		ctx := context.Background()
		err := service.HealthCheck(ctx)

		// Ожидаем ошибку т.к. credentials невалидны
		if err == nil {
			t.Error("HealthCheck() should return error for invalid credentials")
		}
	})

	t.Run("respects context", func(t *testing.T) {
		config := &AccountConfig{
			Name:         "test",
			ClientID:     "id",
			ClientSecret: "secret",
			BaseURL:      "http://localhost:9999", // Unreachable
			Timeout:      100 * time.Millisecond,
		}

		client, _ := NewAuthenticatedClient(config)
		service := NewService(client, nil)

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err := service.HealthCheck(ctx)
		if err == nil {
			t.Error("HealthCheck() should return error for timeout")
		}
	})
}
