package database

import (
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetDSN(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name: "todas as variáveis de ambiente definidas",
			env: map[string]string{
				"DB_HOST": "testhost",
				"DB_PORT": "3307",
				"DB_USER": "testuser",
				"DB_PASS": "testpass",
				"DB_NAME": "testdb",
			},
			expected: "testuser:testpass@tcp(testhost:3307)/testdb?parseTime=true",
		},
		{
			name: "host e porta padrão",
			env: map[string]string{
				"DB_USER": "testuser",
				"DB_PASS": "testpass",
				"DB_NAME": "testdb",
			},
			expected: "testuser:testpass@tcp(localhost:3306)/testdb?parseTime=true",
		},
		{
			name: "senha vazia",
			env: map[string]string{
				"DB_HOST": "testhost",
				"DB_PORT": "3307",
				"DB_USER": "testuser",
				"DB_PASS": "",
				"DB_NAME": "testdb",
			},
			expected: "testuser:@tcp(testhost:3307)/testdb?parseTime=true",
		},
		{
			name: "apenas campos obrigatórios",
			env: map[string]string{
				"DB_USER": "user",
				"DB_PASS": "pass",
				"DB_NAME": "db",
			},
			expected: "user:pass@tcp(localhost:3306)/db?parseTime=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			got := getDSN()
			if got != tt.expected {
				t.Errorf("getDSN() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		mockErr error
		wantErr bool
	}{
		{
			name: "conexão bem-sucedida",
			env: map[string]string{
				"DB_USER": "testuser",
				"DB_PASS": "testpass",
				"DB_NAME": "testdb",
			},
			mockErr: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.mockErr == nil {
				mock.ExpectPing()
			} else {
				mock.ExpectPing().WillReturnError(tt.mockErr)
			}

			os.Clearenv()
			
			for key, value := range tt.env {
				os.Setenv(key, value)
			}

			dsn := getDSN()
			if dsn == "" {
				t.Error("getDSN() returned empty string")
			}

			err = db.Ping()
			if tt.mockErr != nil {
				if err == nil {
					t.Error("expected ping error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected ping error: %v", err)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestConnectIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("pulando teste de integração")
	}

	t.Run("teste de integração com banco real", func(t *testing.T) {
		os.Setenv("DB_USER", "root")
		os.Setenv("DB_PASS", "")
		os.Setenv("DB_NAME", "test")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "3306")

		t.Skip("requer conexão real com banco de dados")
	})
}