package persistence

import (
	"database/sql"
	"errors"
	"testing"

	"scheduling/internal/domain/entities"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestServiceMySQLRepository_FindByID(t *testing.T) {
	tests := []struct {
		name      string
		serviceID int
		mockFn    func(sqlmock.Sqlmock)
		want      func(*testing.T, *entities.Service)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "serviço encontrado com sucesso",
			serviceID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow(1, 101, "Corte de Cabelo", 30, 50.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, service *entities.Service) {
				if service.ID() != 1 {
					t.Errorf("ID esperado 1, obtido %d", service.ID())
				}
				if service.StaffID() != 101 {
					t.Errorf("StaffID esperado 101, obtido %d", service.StaffID())
				}
				if service.Name() != "Corte de Cabelo" {
					t.Errorf("Name esperado 'Corte de Cabelo', obtido '%s'", service.Name())
				}
				if service.DurationMinutes() != 30 {
					t.Errorf("Duration esperado 30, obtido %d", service.DurationMinutes())
				}
				if service.Price() != 50.0 {
					t.Errorf("Price esperado 50.0, obtido %.2f", service.Price())
				}
			},
			wantErr: false,
		},
		{
			name:      "erro no banco de dados",
			serviceID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE id = ?").
					WithArgs(2).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:      "serviço não encontrado",
			serviceID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE id = ?").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  sql.ErrNoRows.Error(),
		},
		{
			name:      "erro na criação da entidade - nome vazio",
			serviceID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow(3, 102, "", 45, 75.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE id = ?").
					WithArgs(3).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "nome do serviço é obrigatório",
		},
		{
			name:      "erro na criação da entidade - duração inválida",
			serviceID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow(4, 103, "Massagem", 0, 100.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE id = ?").
					WithArgs(4).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "a duração deve ser maior que zero",
		},
		{
			name:      "erro na criação da entidade - preço negativo",
			serviceID: 5,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow(5, 104, "Manicure", 30, -20.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE id = ?").
					WithArgs(5).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "preço não pode ser negativo",
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

			repo := NewServiceMySQLRepository(db)
			got, err := repo.FindByID(tt.serviceID)

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

func TestServiceMySQLRepository_FindAllByStaffID(t *testing.T) {
	tests := []struct {
		name    string
		staffID int
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, []*entities.Service)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "serviços encontrados com sucesso",
			staffID: 101,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow(1, 101, "Corte de Cabelo", 30, 50.0).
					AddRow(2, 101, "Barba", 15, 25.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE staff_id = ?").
					WithArgs(101).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, services []*entities.Service) {
				if len(services) != 2 {
					t.Errorf("esperado 2 serviços, obtido %d", len(services))
					return
				}
				
				service1 := services[0]
				if service1.ID() != 1 || service1.Name() != "Corte de Cabelo" {
					t.Errorf("primeiro serviço incorreto: ID=%d, Name=%s", service1.ID(), service1.Name())
				}
				
				service2 := services[1]
				if service2.ID() != 2 || service2.Name() != "Barba" {
					t.Errorf("segundo serviço incorreto: ID=%d, Name=%s", service2.ID(), service2.Name())
				}
			},
			wantErr: false,
		},
		{
			name:    "nenhum serviço encontrado",
			staffID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"})
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE staff_id = ?").
					WithArgs(999).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, services []*entities.Service) {
				if len(services) != 0 {
					t.Errorf("esperado 0 serviços, obtido %d", len(services))
				}
			},
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 102,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE staff_id = ?").
					WithArgs(102).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:    "erro no scan de uma linha",
			staffID: 103,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow("invalid_id", 103, "Massagem", 60, 120.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE staff_id = ?").
					WithArgs(103).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:    "erro na criação de entidade - nome vazio",
			staffID: 104,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "name", "duration", "price"}).
					AddRow(3, 104, "", 45, 75.0)
				mock.ExpectQuery("SELECT id, staff_id, name, duration, price FROM services WHERE staff_id = ?").
					WithArgs(104).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "nome do serviço é obrigatório",
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

			repo := NewServiceMySQLRepository(db)
			got, err := repo.FindAllByStaffID(tt.staffID)

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

func TestServiceMySQLRepository_Exists(t *testing.T) {
	tests := []struct {
		name      string
		serviceID int
		mockFn    func(sqlmock.Sqlmock)
		want      bool
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "serviço existe",
			serviceID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM services WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "serviço não existe",
			serviceID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM services WHERE id = ?").
					WithArgs(999).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:      "múltiplos serviços com mesmo ID (caso improvável mas testado)",
			serviceID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM services WHERE id = ?").
					WithArgs(2).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:      "erro no banco de dados",
			serviceID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM services WHERE id = ?").
					WithArgs(3).
					WillReturnError(errors.New("database connection error"))
			},
			want:    false,
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:      "erro no scan",
			serviceID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow("invalid_number")
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM services WHERE id = ?").
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

			repo := NewServiceMySQLRepository(db)
			got, err := repo.Exists(tt.serviceID)

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

func TestNewServiceMySQLRepository(t *testing.T) {
	tests := []struct {
		name string
		db   *sql.DB
		want func(*testing.T, *ServiceMySQLRepository)
	}{
		{
			name: "criação bem-sucedida do repositório",
			db:   &sql.DB{},
			want: func(t *testing.T, repo *ServiceMySQLRepository) {
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
			want: func(t *testing.T, repo *ServiceMySQLRepository) {
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
			got := NewServiceMySQLRepository(tt.db)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}