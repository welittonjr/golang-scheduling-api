package persistence

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"scheduling/internal/domain/entities"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewAppointmentMySQLRepository(t *testing.T) {
	tests := []struct {
		name string
		db   *sql.DB
		want func(*testing.T, *AppointmentMySQLRepository)
	}{
		{
			name: "criação bem-sucedida do repositório",
			db:   &sql.DB{},
			want: func(t *testing.T, repo *AppointmentMySQLRepository) {
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
			want: func(t *testing.T, repo *AppointmentMySQLRepository) {
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
			got := NewAppointmentMySQLRepository(tt.db)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestAppointmentMySQLRepository_FindByID(t *testing.T) {
	scheduledTime := time.Date(2025, 12, 25, 14, 30, 0, 0, time.UTC)
	createdTime := time.Date(2025, 12, 20, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		appointmentID  int
		mockFn         func(sqlmock.Sqlmock)
		want           func(*testing.T, *entities.Appointment)
		wantErr        bool
		errMsg         string
	}{
		{
			name:          "agendamento encontrado com sucesso",
			appointmentID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "staff_id", "service_id", "scheduled_at", "status", "created_at"}).
					AddRow(1, 2, 3, 4, scheduledTime, "scheduled", createdTime)
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, appointment *entities.Appointment) {
				if appointment.ID() != 1 {
					t.Errorf("ID esperado 1, obtido %d", appointment.ID())
				}
				if appointment.ClientID() != 2 {
					t.Errorf("ClientID esperado 2, obtido %d", appointment.ClientID())
				}
				if appointment.StaffID() != 3 {
					t.Errorf("StaffID esperado 3, obtido %d", appointment.StaffID())
				}
				if appointment.ServiceID() != 4 {
					t.Errorf("ServiceID esperado 4, obtido %d", appointment.ServiceID())
				}
				if !appointment.ScheduledAt().Equal(scheduledTime) {
					t.Errorf("ScheduledAt esperado %v, obtido %v", scheduledTime, appointment.ScheduledAt())
				}
				if appointment.Status() != "scheduled" {
					t.Errorf("Status esperado 'scheduled', obtido '%s'", appointment.Status())
				}
				if !appointment.CreatedAt().Equal(createdTime) {
					t.Errorf("CreatedAt esperado %v, obtido %v", createdTime, appointment.CreatedAt())
				}
			},
			wantErr: false,
		},
		{
			name:          "agendamento não encontrado",
			appointmentID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE id = ?").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  sql.ErrNoRows.Error(),
		},
		{
			name:          "erro no banco de dados",
			appointmentID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE id = ?").
					WithArgs(2).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:          "erro ao criar entidade appointment",
			appointmentID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "staff_id", "service_id", "scheduled_at", "status", "created_at"}).
					AddRow(3, 0, 3, 4, scheduledTime, "scheduled", createdTime)
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE id = ?").
					WithArgs(3).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "cliente, profissional e serviço são obrigatórios",
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

			repo := NewAppointmentMySQLRepository(db)
			got, err := repo.FindByID(tt.appointmentID)

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

func TestAppointmentMySQLRepository_FindAllByStaffID(t *testing.T) {
	scheduledTime1 := time.Date(2025, 12, 25, 14, 30, 0, 0, time.UTC)
	createdTime1 := time.Date(2025, 12, 20, 10, 0, 0, 0, time.UTC)
	scheduledTime2 := time.Date(2025, 12, 26, 16, 0, 0, 0, time.UTC)
	createdTime2 := time.Date(2025, 12, 21, 11, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		staffID int
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, []*entities.Appointment)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "agendamentos encontrados com sucesso",
			staffID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "staff_id", "service_id", "scheduled_at", "status", "created_at"}).
					AddRow(1, 2, 3, 4, scheduledTime1, "scheduled", createdTime1).
					AddRow(2, 5, 3, 6, scheduledTime2, "completed", createdTime2)
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE staff_id = ?").
					WithArgs(3).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, appointments []*entities.Appointment) {
				if len(appointments) != 2 {
					t.Errorf("esperado 2 agendamentos, obtido %d", len(appointments))
				}

				if appointments[0].ID() != 1 {
					t.Errorf("primeiro agendamento: ID esperado 1, obtido %d", appointments[0].ID())
				}
				if appointments[0].ClientID() != 2 {
					t.Errorf("primeiro agendamento: ClientID esperado 2, obtido %d", appointments[0].ClientID())
				}
				if appointments[0].Status() != "scheduled" {
					t.Errorf("primeiro agendamento: Status esperado 'scheduled', obtido '%s'", appointments[0].Status())
				}

				if appointments[1].ID() != 2 {
					t.Errorf("segundo agendamento: ID esperado 2, obtido %d", appointments[1].ID())
				}
				if appointments[1].ClientID() != 5 {
					t.Errorf("segundo agendamento: ClientID esperado 5, obtido %d", appointments[1].ClientID())
				}
				if appointments[1].Status() != "completed" {
					t.Errorf("segundo agendamento: Status esperado 'completed', obtido '%s'", appointments[1].Status())
				}
			},
			wantErr: false,
		},
		{
			name:    "nenhum agendamento encontrado",
			staffID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "staff_id", "service_id", "scheduled_at", "status", "created_at"})
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE staff_id = ?").
					WithArgs(999).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, appointments []*entities.Appointment) {
				if len(appointments) != 0 {
					t.Errorf("esperado 0 agendamentos, obtido %d", len(appointments))
				}
			},
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE staff_id = ?").
					WithArgs(4).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:    "erro ao fazer scan da linha",
			staffID: 5,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "client_id", "staff_id", "service_id", "scheduled_at", "status", "created_at"}).
					AddRow(1, 0, 5, 6, scheduledTime1, "scheduled", createdTime1)
				mock.ExpectQuery("SELECT id, client_id, staff_id, service_id, scheduled_at, status, created_at FROM appointments WHERE staff_id = ?").
					WithArgs(5).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "cliente, profissional e serviço são obrigatórios",
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

			repo := NewAppointmentMySQLRepository(db)
			got, err := repo.FindAllByStaffID(tt.staffID)

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

func TestAppointmentMySQLRepository_HasConflict(t *testing.T) {
	start := time.Date(2025, 12, 25, 14, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 25, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		staffID int
		start   time.Time
		end     time.Time
		mockFn  func(sqlmock.Sqlmock)
		want    bool
		wantErr bool
		errMsg  string
	}{
		{
			name:    "conflito encontrado",
			staffID: 1,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE staff_id = \? AND status = 'scheduled' AND \(\s*\(\s*scheduled_at BETWEEN \? AND \?\s*\)\s*OR \s*\(\s*\? BETWEEN scheduled_at AND DATE_ADD\s*\(\s*scheduled_at, INTERVAL 30 MINUTE\s*\)\s*\)\s*\)`).
					WithArgs(1, start, end, end).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "nenhum conflito encontrado",
			staffID: 2,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE staff_id = \? AND status = 'scheduled' AND \(\s*\(\s*scheduled_at BETWEEN \? AND \?\s*\)\s*OR \s*\(\s*\? BETWEEN scheduled_at AND DATE_ADD\s*\(\s*scheduled_at, INTERVAL 30 MINUTE\s*\)\s*\)\s*\)`).
					WithArgs(2, start, end, end).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "múltiplos conflitos",
			staffID: 3,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE staff_id = \? AND status = 'scheduled' AND \(\s*\(\s*scheduled_at BETWEEN \? AND \?\s*\)\s*OR \s*\(\s*\? BETWEEN scheduled_at AND DATE_ADD\s*\(\s*scheduled_at, INTERVAL 30 MINUTE\s*\)\s*\)\s*\)`).
					WithArgs(3, start, end, end).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 4,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM appointments WHERE staff_id = \? AND status = 'scheduled' AND \(\s*\(\s*scheduled_at BETWEEN \? AND \?\s*\)\s*OR \s*\(\s*\? BETWEEN scheduled_at AND DATE_ADD\s*\(\s*scheduled_at, INTERVAL 30 MINUTE\s*\)\s*\)\s*\)`).
					WithArgs(4, start, end, end).
					WillReturnError(errors.New("database connection error"))
			},
			want:    false,
			wantErr: true,
			errMsg:  "database connection error",
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

			repo := NewAppointmentMySQLRepository(db)
			got, err := repo.HasConflict(tt.staffID, tt.start, tt.end)

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

			if got != tt.want {
				t.Errorf("resultado esperado %v, obtido %v", tt.want, got)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestAppointmentMySQLRepository_Save(t *testing.T) {
	scheduledTime := time.Date(2025, 12, 15, 14, 30, 0, 0, time.UTC)
	
	tests := []struct {
		name        string
		appointment func() *entities.Appointment
		mockFn      func(sqlmock.Sqlmock)
		wantErr     bool
		errMsg      string
	}{
		{
			name: "agendamento salvo com sucesso",
			appointment: func() *entities.Appointment {
				apt, err := entities.NewAppointment(1, 2, 3, scheduledTime)
				if err != nil {
					panic("failed to create appointment: " + err.Error())
				}
				return apt
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO appointments \\(client_id, staff_id, service_id, scheduled_at, status, created_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
					WithArgs(1, 2, 3, scheduledTime, entities.StatusScheduled, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "erro no banco de dados durante inserção",
			appointment: func() *entities.Appointment {
				apt, err := entities.NewAppointment(4, 5, 6, scheduledTime)
				if err != nil {
					panic("failed to create appointment: " + err.Error())
				}
				return apt
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO appointments \\(client_id, staff_id, service_id, scheduled_at, status, created_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?\\)").
					WithArgs(4, 5, 6, scheduledTime, entities.StatusScheduled, sqlmock.AnyArg()).
					WillReturnError(errors.New("database insert error"))
			},
			wantErr: true,
			errMsg:  "database insert error",
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

			repo := NewAppointmentMySQLRepository(db)
			appointment := tt.appointment()
			err = repo.Save(appointment)

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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestAppointmentMySQLRepository_Update(t *testing.T) {
	scheduledTime := time.Date(2025, 12, 15, 14, 30, 0, 0, time.UTC)
	
	tests := []struct {
		name        string
		appointment func() *entities.Appointment
		mockFn      func(sqlmock.Sqlmock)
		wantErr     bool
		errMsg      string
	}{
		{
			name: "agendamento atualizado com sucesso",
			appointment: func() *entities.Appointment {
				apt, err := entities.NewAppointment(1, 2, 3, scheduledTime)
				if err != nil {
					panic("failed to create appointment: " + err.Error())
				}
				apt.SetID(1)
				apt.Complete()
				return apt
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE appointments SET status = \\? WHERE id = \\?").
					WithArgs(entities.StatusCompleted, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "erro no banco de dados durante atualização",
			appointment: func() *entities.Appointment {
				apt, err := entities.NewAppointment(4, 5, 6, scheduledTime)
				if err != nil {
					panic("failed to create appointment: " + err.Error())
				}
				apt.SetID(2)
				apt.Cancel()
				return apt
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE appointments SET status = \\? WHERE id = \\?").
					WithArgs(entities.StatusCancelled, 2).
					WillReturnError(errors.New("database update error"))
			},
			wantErr: true,
			errMsg:  "database update error",
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

			repo := NewAppointmentMySQLRepository(db)
			appointment := tt.appointment()
			err = repo.Update(appointment)

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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestAppointmentMySQLRepository_Delete(t *testing.T) {
	tests := []struct {
		name           string
		appointmentID  int
		mockFn         func(sqlmock.Sqlmock)
		wantErr        bool
		errMsg         string
	}{
		{
			name:          "agendamento deletado com sucesso",
			appointmentID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM appointments WHERE id = \\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:          "erro no banco de dados durante deleção",
			appointmentID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM appointments WHERE id = \\?").
					WithArgs(2).
					WillReturnError(errors.New("database delete error"))
			},
			wantErr: true,
			errMsg:  "database delete error",
		},
		{
			name:          "tentativa de deletar agendamento inexistente",
			appointmentID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM appointments WHERE id = \\?").
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(999, 0))
			},
			wantErr: false,
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

			repo := NewAppointmentMySQLRepository(db)
			err = repo.Delete(tt.appointmentID)

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

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}