package persistence

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"scheduling/internal/domain/entities"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserMySQLRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, *entities.User)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "usuário encontrado com sucesso",
			userID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "created_at"}).
					AddRow(1, "João Silva", "joao@email.com", "client", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
				mock.ExpectQuery("SELECT id, name, email, role, created_at FROM users WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, user *entities.User) {
				if user.ID() != 1 {
					t.Errorf("ID esperado 1, obtido %d", user.ID())
				}
				if user.Name() != "João Silva" {
					t.Errorf("nome esperado 'João Silva', obtido '%s'", user.Name())
				}
				if user.Email() != "joao@email.com" {
					t.Errorf("email esperado 'joao@email.com', obtido '%s'", user.Email())
				}
				if user.Role() != "client" {
					t.Errorf("role esperado 'client', obtido '%s'", user.Role())
				}
				expectedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
				if !user.CreatedAt().Equal(expectedTime) {
					t.Errorf("createdAt esperado %v, obtido %v", expectedTime, user.CreatedAt())
				}
			},
			wantErr: false,
		},
		{
			name:   "usuário encontrado sem created_at",
			userID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "created_at"}).
					AddRow(2, "Maria Santos", "maria@email.com", "admin", nil)
				mock.ExpectQuery("SELECT id, name, email, role, created_at FROM users WHERE id = ?").
					WithArgs(2).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, user *entities.User) {
				if user.ID() != 2 {
					t.Errorf("ID esperado 2, obtido %d", user.ID())
				}
				if user.Name() != "Maria Santos" {
					t.Errorf("nome esperado 'Maria Santos', obtido '%s'", user.Name())
				}
				if user.Email() != "maria@email.com" {
					t.Errorf("email esperado 'maria@email.com', obtido '%s'", user.Email())
				}
				if user.Role() != "admin" {
					t.Errorf("role esperado 'admin', obtido '%s'", user.Role())
				}
			},
			wantErr: false,
		},
		{
			name:   "erro no banco de dados",
			userID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, role, created_at FROM users WHERE id = ?").
					WithArgs(3).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:   "usuário não encontrado",
			userID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, role, created_at FROM users WHERE id = ?").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  sql.ErrNoRows.Error(),
		},
		{
			name:   "email inválido retornado do banco",
			userID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "created_at"}).
					AddRow(4, "Pedro Lima", "email-invalido", "client", time.Now())
				mock.ExpectQuery("SELECT id, name, email, role, created_at FROM users WHERE id = ?").
					WithArgs(4).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("erro ao criar mock do banco: %v", err)
			}
			defer db.Close()

			tt.mockFn(mock)

			repo := NewUserMySQLRepository(db)
			got, err := repo.FindByID(tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("erro esperado, mas não obtive nenhum")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("erro esperado '%s', obtido '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if tt.want != nil {
				tt.want(t, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestUserMySQLRepository_Exists(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		mockFn  func(sqlmock.Sqlmock)
		want    bool
		wantErr bool
		errMsg  string
	}{
		{
			name:   "usuário existe",
			userID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "usuário não existe",
			userID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = ?").
					WithArgs(999).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:   "múltiplos usuários com mesmo ID (caso improvável mas testado)",
			userID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = ?").
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:   "erro no banco de dados",
			userID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = ?").
					WithArgs(3).
					WillReturnError(errors.New("database connection error"))
			},
			want:    false,
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:   "erro no scan",
			userID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow("invalid_number")
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = ?").
					WithArgs(4).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("erro ao criar mock do banco: %v", err)
			}
			defer db.Close()

			tt.mockFn(mock)

			repo := NewUserMySQLRepository(db)
			got, err := repo.Exists(tt.userID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("erro esperado, mas não obtive nenhum")
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("erro esperado '%s', obtido '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if got != tt.want {
				t.Errorf("resultado esperado %v, obtido %v", tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestNewUserMySQLRepository(t *testing.T) {
	tests := []struct {
		name string
		db   *sql.DB
		want func(*testing.T, *UserMySQLRepository)
	}{
		{
			name: "criação bem-sucedida do repositório",
			db:   &sql.DB{},
			want: func(t *testing.T, repo *UserMySQLRepository) {
				if repo == nil {
					t.Error("repositório não deve ser nil")
				}
				if repo.db == nil {
					t.Error("db do repositório não deve ser nil")
				}
			},
		},
		{
			name: "criação com db nil",
			db:   nil,
			want: func(t *testing.T, repo *UserMySQLRepository) {
				if repo == nil {
					t.Error("repositório não deve ser nil")
				}
				if repo.db != nil {
					t.Error("db do repositório deve ser nil quando passado nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserMySQLRepository(tt.db)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}