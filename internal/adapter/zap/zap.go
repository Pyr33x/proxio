package zap

import "go.uber.org/zap"

type Adapter struct {
	Logger *zap.Logger
}

func New(environment string) *Adapter {
	logger, err := createLogger(environment)
	if err != nil {
		if fallback, fallbackErr := zap.NewDevelopment(); fallbackErr == nil {
			fallback.Error("failed to create configured logger, using development fallback",
				zap.Error(err))
			return &Adapter{Logger: fallback}
		}

		println("CRITICAL: failed to create any logger:", err.Error())
		return nil
	}

	return &Adapter{Logger: logger}
}

func createLogger(environment string) (*zap.Logger, error) {
	switch environment {
	case "prod", "production":
		return zap.NewProduction()
	case "dev", "development", "local":
		return zap.NewDevelopment()
	case "test":
		return zap.NewNop(), nil
	default:
		return zap.NewDevelopment()
	}
}
