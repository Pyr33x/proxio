package err

import "errors"

var (
	ErrEmptyCacheKey = errors.New("cache key cannot be empty")
)
