package persistence

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"scheduling/internal/domain/entities"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewAvailableSlotMySQLRepository(t *testing.T) {
	tests := []struct {
		name string
		db   *sql.DB
		want func(*testing.T, *AvailableSlotMySQLRepository)
	}{
		{
			name: "criação bem-sucedida do repositório",
			db:   &sql.DB{},
			want: func(t *testing.T, repo *AvailableSlotMySQLRepository) {
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
			want: func(t *testing.T, repo *AvailableSlotMySQLRepository) {
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
			got := NewAvailableSlotMySQLRepository(tt.db)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestAvailableSlotMySQLRepository_FindByID(t *testing.T) {
	startTime := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 25, 17, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		slotID  int
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, *entities.AvailableSlot)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "slot encontrado com sucesso",
			slotID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(1, 2, "monday", startTime, endTime)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE id = \\?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, slot *entities.AvailableSlot) {
				if slot.ID() != 1 {
					t.Errorf("ID esperado 1, obtido %d", slot.ID())
				}
				if slot.StaffID() != 2 {
					t.Errorf("StaffID esperado 2, obtido %d", slot.StaffID())
				}
				if slot.Weekday() != entities.Monday {
					t.Errorf("Weekday esperado 'monday', obtido '%s'", slot.Weekday())
				}
				if !slot.StartTime().Equal(startTime) {
					t.Errorf("StartTime esperado %v, obtido %v", startTime, slot.StartTime())
				}
				if !slot.EndTime().Equal(endTime) {
					t.Errorf("EndTime esperado %v, obtido %v", endTime, slot.EndTime())
				}
			},
			wantErr: false,
		},
		{
			name:   "slot não encontrado",
			slotID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE id = \\?").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			errMsg:  sql.ErrNoRows.Error(),
		},
		{
			name:   "erro no banco de dados",
			slotID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE id = \\?").
					WithArgs(2).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:   "erro ao criar entidade available_slot - staffID inválido",
			slotID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(3, 0, "monday", startTime, endTime)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE id = \\?").
					WithArgs(3).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "staffID é obrigatório",
		},
		{
			name:   "erro ao criar entidade available_slot - weekday inválido",
			slotID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(4, 2, "invalid_day", startTime, endTime)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE id = \\?").
					WithArgs(4).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "dia da semana inválido: invalid_day",
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

			repo := NewAvailableSlotMySQLRepository(db)
			got, err := repo.FindByID(tt.slotID)

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

func TestAvailableSlotMySQLRepository_FindAllByStaffID(t *testing.T) {
	startTime1 := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	endTime1 := time.Date(2025, 12, 25, 12, 0, 0, 0, time.UTC)
	startTime2 := time.Date(2025, 12, 25, 13, 0, 0, 0, time.UTC)
	endTime2 := time.Date(2025, 12, 25, 17, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		staffID int
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, []*entities.AvailableSlot)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "slots encontrados com sucesso",
			staffID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(1, 2, "monday", startTime1, endTime1).
					AddRow(2, 2, "tuesday", startTime2, endTime2)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\?").
					WithArgs(2).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, slots []*entities.AvailableSlot) {
				if len(slots) != 2 {
					t.Errorf("esperado 2 slots, obtido %d", len(slots))
				}

				if slots[0].ID() != 1 {
					t.Errorf("primeiro slot: ID esperado 1, obtido %d", slots[0].ID())
				}
				if slots[0].StaffID() != 2 {
					t.Errorf("primeiro slot: StaffID esperado 2, obtido %d", slots[0].StaffID())
				}
				if slots[0].Weekday() != entities.Monday {
					t.Errorf("primeiro slot: Weekday esperado 'monday', obtido '%s'", slots[0].Weekday())
				}

				if slots[1].ID() != 2 {
					t.Errorf("segundo slot: ID esperado 2, obtido %d", slots[1].ID())
				}
				if slots[1].StaffID() != 2 {
					t.Errorf("segundo slot: StaffID esperado 2, obtido %d", slots[1].StaffID())
				}
				if slots[1].Weekday() != entities.Tuesday {
					t.Errorf("segundo slot: Weekday esperado 'tuesday', obtido '%s'", slots[1].Weekday())
				}
			},
			wantErr: false,
		},
		{
			name:    "nenhum slot encontrado",
			staffID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"})
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\?").
					WithArgs(999).
					WillReturnRows(rows)
			},
			want: func(t *testing.T, slots []*entities.AvailableSlot) {
				if len(slots) != 0 {
					t.Errorf("esperado 0 slots, obtido %d", len(slots))
				}
			},
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 3,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\?").
					WithArgs(3).
					WillReturnError(errors.New("database connection error"))
			},
			wantErr: true,
			errMsg:  "database connection error",
		},
		{
			name:    "erro ao fazer scan da linha",
			staffID: 4,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(1, 0, "monday", startTime1, endTime1)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\?").
					WithArgs(4).
					WillReturnRows(rows)
			},
			wantErr: true,
			errMsg:  "staffID é obrigatório",
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

			repo := NewAvailableSlotMySQLRepository(db)
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

func TestAvailableSlotMySQLRepository_FindByWeekday(t *testing.T) {
	startTime := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 25, 17, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		staffID int
		weekday entities.Weekday
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, []*entities.AvailableSlot)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "slots encontrados por weekday",
			staffID: 2,
			weekday: entities.Monday,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(1, 2, "monday", startTime, endTime)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\? AND weekday = \\?").
					WithArgs(2, "monday").
					WillReturnRows(rows)
			},
			want: func(t *testing.T, slots []*entities.AvailableSlot) {
				if len(slots) != 1 {
					t.Errorf("esperado 1 slot, obtido %d", len(slots))
				}
				if slots[0].Weekday() != entities.Monday {
					t.Errorf("Weekday esperado 'monday', obtido '%s'", slots[0].Weekday())
				}
			},
			wantErr: false,
		},
		{
			name:    "nenhum slot encontrado para o weekday",
			staffID: 2,
			weekday: entities.Sunday,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"})
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\? AND weekday = \\?").
					WithArgs(2, "sunday").
					WillReturnRows(rows)
			},
			want: func(t *testing.T, slots []*entities.AvailableSlot) {
				if len(slots) != 0 {
					t.Errorf("esperado 0 slots, obtido %d", len(slots))
				}
			},
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 2,
			weekday: entities.Monday,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\? AND weekday = \\?").
					WithArgs(2, "monday").
					WillReturnError(errors.New("database connection error"))
			},
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

			repo := NewAvailableSlotMySQLRepository(db)
			got, err := repo.FindByWeekday(tt.staffID, tt.weekday)

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

func TestAvailableSlotMySQLRepository_FindSlotsByStaffAndDate(t *testing.T) {
	startTime := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 25, 17, 0, 0, 0, time.UTC)
	testDate := time.Date(2025, 12, 22, 10, 0, 0, 0, time.UTC) // Monday

	tests := []struct {
		name    string
		staffID int
		date    time.Time
		mockFn  func(sqlmock.Sqlmock)
		want    func(*testing.T, []*entities.AvailableSlot)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "slots encontrados por data",
			staffID: 2,
			date:    testDate,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "staff_id", "weekday", "start_time", "end_time"}).
					AddRow(1, 2, "monday", startTime, endTime)
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\? AND weekday = \\?").
					WithArgs(2, "monday").
					WillReturnRows(rows)
			},
			want: func(t *testing.T, slots []*entities.AvailableSlot) {
				if len(slots) != 1 {
					t.Errorf("esperado 1 slot, obtido %d", len(slots))
				}
				if slots[0].Weekday() != entities.Monday {
					t.Errorf("Weekday esperado 'monday', obtido '%s'", slots[0].Weekday())
				}
			},
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 2,
			date:    testDate,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, staff_id, weekday, start_time, end_time FROM available_slots WHERE staff_id = \\? AND weekday = \\?").
					WithArgs(2, "monday").
					WillReturnError(errors.New("database connection error"))
			},
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

			repo := NewAvailableSlotMySQLRepository(db)
			got, err := repo.FindSlotsByStaffAndDate(tt.staffID, tt.date)

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

func TestAvailableSlotMySQLRepository_HasConflict(t *testing.T) {
	start := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 25, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		staffID int
		weekday entities.Weekday
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
			weekday: entities.Monday,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND \(\s*\(\s*start_time < \? AND end_time > \?\s*\)\s*OR\s*\(\s*start_time >= \? AND start_time < \?\s*\)\s*\)`).
					WithArgs(1, "monday", end, start, start, end).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "nenhum conflito encontrado",
			staffID: 2,
			weekday: entities.Monday,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND \(\s*\(\s*start_time < \? AND end_time > \?\s*\)\s*OR\s*\(\s*start_time >= \? AND start_time < \?\s*\)\s*\)`).
					WithArgs(2, "monday", end, start, start, end).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "múltiplos conflitos",
			staffID: 3,
			weekday: entities.Monday,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(3)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND \(\s*\(\s*start_time < \? AND end_time > \?\s*\)\s*OR\s*\(\s*start_time >= \? AND start_time < \?\s*\)\s*\)`).
					WithArgs(3, "monday", end, start, start, end).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 4,
			weekday: entities.Monday,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND \(\s*\(\s*start_time < \? AND end_time > \?\s*\)\s*OR\s*\(\s*start_time >= \? AND start_time < \?\s*\)\s*\)`).
					WithArgs(4, "monday", end, start, start, end).
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

			repo := NewAvailableSlotMySQLRepository(db)
			got, err := repo.HasConflict(tt.staffID, tt.weekday, tt.start, tt.end)

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

func TestAvailableSlotMySQLRepository_IsWithinAvailableSlot(t *testing.T) {
	start := time.Date(2025, 12, 22, 9, 0, 0, 0, time.UTC)  // Monday 9am
	end := time.Date(2025, 12, 22, 10, 0, 0, 0, time.UTC)   // Monday 10am

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
			name:    "dentro do slot disponível",
			staffID: 1,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND start_time <= \? AND end_time >= \?`).
					WithArgs(1, "monday", start, end).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "fora do slot disponível",
			staffID: 2,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND start_time <= \? AND end_time >= \?`).
					WithArgs(2, "monday", start, end).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: false,
		},
		{
			name:    "erro no banco de dados",
			staffID: 3,
			start:   start,
			end:     end,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM available_slots WHERE staff_id = \? AND weekday = \? AND start_time <= \? AND end_time >= \?`).
					WithArgs(3, "monday", start, end).
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

			repo := NewAvailableSlotMySQLRepository(db)
			got, err := repo.IsWithinAvailableSlot(tt.staffID, tt.start, tt.end)

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

func TestAvailableSlotMySQLRepository_Save(t *testing.T) {
	startTime := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 25, 17, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name    string
		slot    func() *entities.AvailableSlot
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "slot salvo com sucesso",
			slot: func() *entities.AvailableSlot {
				slot, err := entities.NewAvailableSlot(1, entities.Monday, startTime, endTime)
				if err != nil {
					panic("failed to create slot: " + err.Error())
				}
				return slot
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO available_slots \\(staff_id, weekday, start_time, end_time\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
					WithArgs(1, "monday", startTime, endTime).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "erro no banco de dados durante inserção",
			slot: func() *entities.AvailableSlot {
				slot, err := entities.NewAvailableSlot(2, entities.Tuesday, startTime, endTime)
				if err != nil {
					panic("failed to create slot: " + err.Error())
				}
				return slot
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO available_slots \\(staff_id, weekday, start_time, end_time\\) VALUES \\(\\?, \\?, \\?, \\?\\)").
					WithArgs(2, "tuesday", startTime, endTime).
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

			repo := NewAvailableSlotMySQLRepository(db)
			slot := tt.slot()
			err = repo.Save(slot)

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

func TestAvailableSlotMySQLRepository_Update(t *testing.T) {
	startTime := time.Date(2025, 12, 25, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 12, 25, 17, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name    string
		slot    func() *entities.AvailableSlot
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name: "slot atualizado com sucesso",
			slot: func() *entities.AvailableSlot {
				slot, err := entities.NewAvailableSlot(1, entities.Monday, startTime, endTime)
				if err != nil {
					panic("failed to create slot: " + err.Error())
				}
				slot.SetID(1)
				return slot
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE available_slots SET weekday = \\?, start_time = \\?, end_time = \\? WHERE id = \\?").
					WithArgs("monday", startTime, endTime, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "erro no banco de dados durante atualização",
			slot: func() *entities.AvailableSlot {
				slot, err := entities.NewAvailableSlot(2, entities.Tuesday, startTime, endTime)
				if err != nil {
					panic("failed to create slot: " + err.Error())
				}
				slot.SetID(2)
				return slot
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE available_slots SET weekday = \\?, start_time = \\?, end_time = \\? WHERE id = \\?").
					WithArgs("tuesday", startTime, endTime, 2).
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

			repo := NewAvailableSlotMySQLRepository(db)
			slot := tt.slot()
			err = repo.Update(slot)

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

func TestAvailableSlotMySQLRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		slotID  int
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errMsg  string
	}{
		{
			name:   "slot deletado com sucesso",
			slotID: 1,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM available_slots WHERE id = \\?").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name:   "erro no banco de dados durante deleção",
			slotID: 2,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM available_slots WHERE id = \\?").
					WithArgs(2).
					WillReturnError(errors.New("database delete error"))
			},
			wantErr: true,
			errMsg:  "database delete error",
		},
		{
			name:   "tentativa de deletar slot inexistente",
			slotID: 999,
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM available_slots WHERE id = \\?").
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

			repo := NewAvailableSlotMySQLRepository(db)
			err = repo.Delete(tt.slotID)

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