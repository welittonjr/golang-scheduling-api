package persistence

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"scheduling/internal/domain/entities"
	"scheduling/internal/domain/valueobject"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserMySQLRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		user    func() *entities.User
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "usuário admin criado com sucesso",
			user: func() *entities.User {
				email, _ := valueobject.NewEmail("joao@gmail.com")
				user, _ := entities.NewUser(1, "João da Silva", email.String(), "123456", "admin")
				return user
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users \\(name, email, password, role, created_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("João da Silva", "joao@gmail.com", "123456", "admin", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "erro ao executar query",
			user: func() *entities.User {
				email, _ := valueobject.NewEmail("joao@gmail.com")
				user, _ := entities.NewUser(0, "João da Silva", email.String(), "123456", "admin")
				return user
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users \\(name, email, password, role, created_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("João da Silva", "joao@gmail.com", "123456", "admin", sqlmock.AnyArg()).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name: "erro ao obter último ID inserido",
			user: func() *entities.User {
				email, _ := valueobject.NewEmail("joao@gmail.com")
				user, _ := entities.NewUser(0, "João da Silva", email.String(), "123456", "admin")
				return user
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users \\(name, email, password, role, created_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)").
					WithArgs("João da Silva", "joao@gmail.com", "123456", "admin", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))
			},
			wantErr: true,
			errMsg:  "last insert id error",
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
			user := tt.user()
			err = repo.Create(context.Background(), user)

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

			if user.ID() != 1 {
				t.Errorf("ID esperado 1, obtido %d", user.ID())
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestUserMySQLRepository_FindByID(t *testing.T) {
    expectedQuery := "SELECT id, name, email, password, role, created_at FROM users WHERE id = \\?"

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
                rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at"}).
                    AddRow(1, "João Silva", "joao@email.com", "123456", "client", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))

                mock.ExpectQuery(expectedQuery).
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
                rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at"}).
                    AddRow(2, "Maria Santos", "maria@email.com", "123456", "admin", nil)

                mock.ExpectQuery(expectedQuery).
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
                if !user.CreatedAt().IsZero() {
                    t.Errorf("createdAt esperado zero value, obtido %v", user.CreatedAt())
                }
            },
            wantErr: false,
        },
        {
            name:   "erro no banco de dados",
            userID: 3,
            mockFn: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery(expectedQuery).
                    WithArgs(3).
                    WillReturnError(errors.New("database connection error"))
            },
            wantErr: true,
            errMsg:  "database connection error",
        },
        {
            name:   "email inválido retornado do banco",
            userID: 4,
            mockFn: func(mock sqlmock.Sqlmock) {
                rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at"}).
                    AddRow(4, "Pedro Lima", "email-invalido", "", "client", time.Now())
                mock.ExpectQuery(expectedQuery).
                    WithArgs(4).
                    WillReturnRows(rows)
            },
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
            got, err := repo.FindByID(context.Background(), tt.userID)

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
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = \\?").
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
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = \\?").
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
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = \\?").
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
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = \\?").
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
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE id = \\?").
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
			got, err := repo.Exists(context.Background(), tt.userID)

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

func TestUserMySQLRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		user    func() *entities.User
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "usuário atualizado com sucesso",
			user: func() *entities.User {
				email, _ := valueobject.NewEmail("joao.novo@email.com")
				user := entities.RebuildUser(1, "João Silva Atualizado", email, "admin")
				return user
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\?, email = \\?, role = \\? WHERE id = \\?").
					WithArgs("João Silva Atualizado", "joao.novo@email.com", "admin", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "erro ao executar update",
			user: func() *entities.User {
				email, _ := valueobject.NewEmail("joao@email.com")
				user := entities.RebuildUser(1, "João Silva", email, "admin")
				return user
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\?, email = \\?, role = \\? WHERE id = \\?").
					WithArgs("João Silva", "joao@email.com", "admin", 1).
					WillReturnError(errors.New("update error"))
			},
			wantErr: true,
			errMsg:  "update error",
		},
		{
			name: "usuário não encontrado para update",
			user: func() *entities.User {
				email, _ := valueobject.NewEmail("joao@email.com")
				user := entities.RebuildUser(999, "João Silva", email, "admin")
				return user
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET name = \\?, email = \\?, role = \\? WHERE id = \\?").
					WithArgs("João Silva", "joao@email.com", "admin", 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			errMsg:  sql.ErrNoRows.Error(),
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
			user := tt.user()
			err = repo.Update(context.Background(), user)

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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestUserMySQLRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "usuário deletado com sucesso",
			userID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET deleted_at = \\?, updated_at = \\? WHERE id = \\? AND deleted_at IS NULL").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name:   "erro ao deletar usuário",
			userID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET deleted_at = \\?, updated_at = \\? WHERE id = \\? AND deleted_at IS NULL").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 2).
					WillReturnError(errors.New("delete error"))
			},
			wantErr: true,
			errMsg:  "delete error",
		},
		{
			name:   "usuário não encontrado para deletar",
			userID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users SET deleted_at = \\?, updated_at = \\? WHERE id = \\? AND deleted_at IS NULL").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
			errMsg:  sql.ErrNoRows.Error(),
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
			err = repo.Delete(context.Background(), tt.userID)

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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestUserMySQLRepository_List(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		offset  int
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, []*entities.User)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "lista usuários com sucesso",
			limit:  10,
			offset: 0,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "created_at", "updated_at"}).
					AddRow(1, "João Silva", "joao@email.com", "admin", time.Now(), time.Now()).
					AddRow(2, "Maria Santos", "maria@email.com", "client", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, role, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \\? OFFSET \\?").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, users []*entities.User) {
				if len(users) != 2 {
					t.Errorf("esperado 2 usuários, obtido %d", len(users))
				}
				if users[0].ID() != 1 {
					t.Errorf("primeiro usuário ID esperado 1, obtido %d", users[0].ID())
				}
				if users[1].ID() != 2 {
					t.Errorf("segundo usuário ID esperado 2, obtido %d", users[1].ID())
				}
			},
			wantErr: false,
		},
		{
			name:   "lista vazia",
			limit:  10,
			offset: 0,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, name, email, role, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \\? OFFSET \\?").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, users []*entities.User) {
				if len(users) != 0 {
					t.Errorf("esperado 0 usuários, obtido %d", len(users))
				}
			},
			wantErr: false,
		},
		{
			name:   "erro na query",
			limit:  10,
			offset: 0,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, role, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \\? OFFSET \\?").
					WithArgs(10, 0).
					WillReturnError(errors.New("query error"))
			},
			wantErr: true,
			errMsg:  "query error",
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
			got, err := repo.List(context.Background(), tt.limit, tt.offset)

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

			if tt.want != nil {
				tt.want(t, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestUserMySQLRepository_Count(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(sqlmock.Sqlmock)
		want    int64
		wantErr bool
		errMsg  string
	}{
		{
			name: "contagem bem-sucedida",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE deleted_at IS NULL").
					WillReturnRows(rows)
			},
			want:    5,
			wantErr: false,
		},
		{
			name: "erro na contagem",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE deleted_at IS NULL").
					WillReturnError(errors.New("count error"))
			},
			want:    0,
			wantErr: true,
			errMsg:  "count error",
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
			got, err := repo.Count(context.Background())

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
				t.Errorf("contagem esperada %d, obtida %d", tt.want, got)
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