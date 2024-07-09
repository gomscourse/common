package tools

import "context"

// HandleErrorWithContext runs handler and returns its error type result
// or done context error if context canceled or timed out
func HandleErrorWithContext(ctx context.Context, handler func() error) error {
	ch := make(chan error, 1)

	go func() {
		ch <- handler()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
