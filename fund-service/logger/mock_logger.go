package logger

type mockLogger struct{}

func NewMockLogger() Logger {
	return &mockLogger{}
}

func (m *mockLogger) Info(msg string, fields ...Field)  {}
func (m *mockLogger) Error(msg string, fields ...Field) {}
func (m *mockLogger) Debug(msg string, fields ...Field) {}
func (m *mockLogger) With(fields ...Field) Logger       { return m }
func (m *mockLogger) Sync() error                       { return nil }
