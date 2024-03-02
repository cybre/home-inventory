package helpers

import (
	"context"
	"time"
)

func Poll[T any](ctx context.Context, fn func(context.Context) (T, error), interval, timeout time.Duration) (T, error) {
	lastErr := context.DeadlineExceeded

	value, err := fn(ctx)
	if err == nil {
		return value, nil
	}
	lastErr = err

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return *new(T), lastErr
		case <-ticker.C:
			value, err := fn(ctx)
			if err == nil {
				return value, nil
			}

			lastErr = err
		}
	}
}
