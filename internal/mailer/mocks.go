package mailer

type MockMailer struct{}

func (m *MockMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	return nil
}
