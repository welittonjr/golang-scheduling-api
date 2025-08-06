package valueobject

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "email valido",
			address:     "user@example.com",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "email valido com subdominio",
			address:     "user@mail.example.com",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "email valido com caracter especial",
			address:     "name.lastname_tag@example.com.br",
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name:        "email invalido - ausente @",
			address:     "userexample.com",
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "emai invalido - ausente (dominio)[ex: gmail]",
			address:     "name@",
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "email invalido - ausente (parte local)[ex: br]",
			address:     "@gmail.com",
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "email invalido - invalido characters",
			address:     "name@g mail.com",
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "email invalido - tld muito curto",
			address:     "name@gmail.c",
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "email vazio",
			address:     "",
			wantErr:     true,
			expectedErr: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != tt.expectedErr {
				t.Errorf("NewEmail() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	emailStr := "user@gmail.com"
	email, err := NewEmail(emailStr)
	if err != nil {
		t.Fatalf("Falha ao criar e-mail: %v", err)
	}

	if got := email.String(); got != emailStr {
		t.Errorf("Email.String() = %v, quero %v", got, emailStr)
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := NewEmail("name@gmail.com")
	email2, _ := NewEmail("name@gmail.com")
	email3, _ := NewEmail("another@gmail.com")

	if !email1.Equals(email2) {
		t.Errorf("Espera-se que os e-mails sejam iguais")
	}

	if email1.Equals(email3) {
		t.Errorf("E-mails esperados para serem diferentes")
	}
}
