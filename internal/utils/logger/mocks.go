package logger

func NewLoggerMock() Logger {
	return &MockLogger{}
}

type MockLogger struct{}

func (m *MockLogger) Debug(args ...interface{})                 {}
func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Info(args ...interface{})                  {}
func (m *MockLogger) Infof(format string, args ...interface{})  {}
func (m *MockLogger) Warn(args ...interface{})                  {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}
func (m *MockLogger) Error(args ...interface{})                 {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Fatal(args ...interface{})                 {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {}
func (m *MockLogger) Panic(args ...interface{})                 {}
func (m *MockLogger) Panicf(format string, args ...interface{}) {}
