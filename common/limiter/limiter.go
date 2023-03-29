package limiter

type Limiter interface {
	TryAcquire() bool
}
