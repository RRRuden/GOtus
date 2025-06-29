package logger

//go:generate mockgen -source=interface.go -destination=../mocks/logger_mock.go -package=mocks
type Logger interface {
	Log(key string, message string, ttlSeconds int) error
}
