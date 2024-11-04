package recovery

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"sync"
	"time"
	"validation/config"

	"github.com/jordan-wright/email"
)

type Service struct {
	config config.RecoveryConfig
	mu     sync.Mutex
}

func NewService(config config.RecoveryConfig) *Service {
	return &Service{
		config: config,
		mu:     sync.Mutex{},
	}
}

func (s *Service) SendEmail(toEmail, hash string) error {
	e := email.NewEmail()
	e.From = s.config.SMTP.From
	e.To = []string{toEmail}
	e.Subject = "Восстановление пароля"

	resetLink := fmt.Sprintf("http://localhost:8081/verify/%s", hash)

	e.HTML = []byte(fmt.Sprintf(`
        <h2>Восстановление пароля</h2>
        <p>Вы запросили восстановление пароля.</p>
        <p>Для восстановления перейдите по ссылке: <a href="%s">%s</a></p>
        <p>Если это были не вы, проигнорируйте это письмо.</p>
    `, resetLink, resetLink))

	auth := smtp.PlainAuth("",
		s.config.SMTP.Username,
		s.config.SMTP.Password,
		s.config.SMTP.Host,
	)

	smtpAddr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	fmt.Printf("Connecting to SMTP server: %s\n", smtpAddr)

	err := e.Send(smtpAddr, auth)
	if err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
		return err
	}

	fmt.Printf("Email sent successfully to: %s\n", toEmail)
	return nil
}

func (s *Service) SaveToFile(hash, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var data []RecoveryData

	if _, err := os.Stat(s.config.DataFile); err == nil {
		fileData, err := os.ReadFile(s.config.DataFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(fileData, &data); err != nil {
			return err
		}
	}

	data = append(data, RecoveryData{
		Hash:      hash,
		Email:     email,
		CreatedAt: time.Now(),
	})

	fileData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.config.DataFile, fileData, 0644)
}

func (s *Service) GetRecoveryData(hash string) (*RecoveryData, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var data []RecoveryData

	fileData, err := os.ReadFile(s.config.DataFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, err
	}

	for _, record := range data {
		if record.Hash == hash {
			if time.Since(record.CreatedAt) > 24*time.Hour {
				return nil, fmt.Errorf("recovery link expired")
			}
			return &record, nil
		}
	}

	return nil, fmt.Errorf("recovery hash not found")
}

func (s *Service) RemoveRecoveryData(hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var data []RecoveryData

	fileData, err := os.ReadFile(s.config.DataFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		return err
	}

	newData := make([]RecoveryData, 0, len(data))
	for _, record := range data {
		if record.Hash != hash {
			newData = append(newData, record)
		}
	}

	fileData, err = json.MarshalIndent(newData, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.config.DataFile, fileData, 0644)
}
