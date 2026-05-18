package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client encapsula o cliente Redis com helpers tipados
type Client struct {
	rdb *redis.Client
}

// NewClient cria um novo cliente Redis
// Se Redis não estiver disponível, retorna nil sem crashar a aplicação
func NewClient(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Warn("cache: Redis não disponível — continuando sem cache", "error", err)
		return &Client{rdb: nil}
	}

	slog.Info("cache: Redis conectado", "addr", addr)
	return &Client{rdb: rdb}
}

// IsAvailable retorna true se o Redis está conectado
func (c *Client) IsAvailable() bool {
	return c.rdb != nil
}

// Set armazena um valor JSON no cache com TTL
func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if !c.IsAvailable() {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, ttl).Err()
}

// Get recupera e desserializa um valor do cache
// Retorna false se não encontrado ou expirado
func (c *Client) Get(ctx context.Context, key string, dest any) (bool, error) {
	if !c.IsAvailable() {
		return false, nil
	}
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(data, dest)
}

// Delete remove uma chave do cache
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if !c.IsAvailable() {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

// Invalidate remove todas as chaves com um prefixo
func (c *Client) Invalidate(ctx context.Context, pattern string) error {
	if !c.IsAvailable() {
		return nil
	}
	keys, err := c.rdb.Keys(ctx, pattern).Result()
	if err != nil || len(keys) == 0 {
		return err
	}
	return c.rdb.Del(ctx, keys...).Err()
}
