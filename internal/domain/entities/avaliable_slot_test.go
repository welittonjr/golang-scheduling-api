package entities

import (
	"testing"
	"time"
)

func TestNewAvailableSlot(t *testing.T) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)

	tests := []struct {
		name      string
		staffID   int
		weekday   Weekday
		startTime time.Time
		endTime   time.Time
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "criar slot válido",
			staffID:   1,
			weekday:   Monday,
			startTime: start,
			endTime:   end,
			wantErr:   false,
		},
		{
			name:      "staffID zero deve retornar erro",
			staffID:   0,
			weekday:   Tuesday,
			startTime: start,
			endTime:   end,
			wantErr:   true,
			errMsg:    "staffID é obrigatório",
		},
		{
			name:      "dia da semana inválido deve retornar erro",
			staffID:   1,
			weekday:   Weekday("invalid"),
			startTime: start,
			endTime:   end,
			wantErr:   true,
			errMsg:    "dia da semana inválido: invalid",
		},
		{
			name:      "horário final antes do inicial deve retornar erro",
			staffID:   1,
			weekday:   Wednesday,
			startTime: end,
			endTime:   start,
			wantErr:   true,
			errMsg:    "o horário inicial deve ser antes do final",
		},
		{
			name:      "horários iguais devem retornar erro",
			staffID:   1,
			weekday:   Thursday,
			startTime: start,
			endTime:   start,
			wantErr:   true,
			errMsg:    "o horário inicial deve ser antes do final",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot, err := NewAvailableSlot(tt.staffID, tt.weekday, tt.startTime, tt.endTime)

			if tt.wantErr {
				if err == nil {
					t.Error("esperado erro, mas nenhum foi retornado")
				} else if err.Error() != tt.errMsg {
					t.Errorf("mensagem de erro incorreta, esperado: '%s', obtido: '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("não esperava erro, mas obteve: %v", err)
			}

			if slot.StaffID() != tt.staffID {
				t.Errorf("StaffID esperado %d, obtido %d", tt.staffID, slot.StaffID())
			}
			if slot.Weekday() != tt.weekday {
				t.Errorf("Weekday esperado %s, obtido %s", tt.weekday, slot.Weekday())
			}
			if !slot.StartTime().Equal(tt.startTime) {
				t.Errorf("StartTime esperado %v, obtido %v", tt.startTime, slot.StartTime())
			}
			if !slot.EndTime().Equal(tt.endTime) {
				t.Errorf("EndTime esperado %v, obtido %v", tt.endTime, slot.EndTime())
			}
		})
	}
}

func TestAvailableSlotGetters(t *testing.T) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	slot := &AvailableSlot{
		id:        1,
		staffID:   101,
		weekday:   Friday,
		startTime: start,
		endTime:   end,
	}

	t.Run("Testar getters", func(t *testing.T) {
		if got := slot.ID(); got != 1 {
			t.Errorf("ID() = %d, esperado 1", got)
		}
		if got := slot.StaffID(); got != 101 {
			t.Errorf("StaffID() = %d, esperado 101", got)
		}
		if got := slot.Weekday(); got != Friday {
			t.Errorf("Weekday() = %s, esperado 'friday'", got)
		}
		if !slot.StartTime().Equal(start) {
			t.Errorf("StartTime() = %v, esperado %v", slot.StartTime(), start)
		}
		if !slot.EndTime().Equal(end) {
			t.Errorf("EndTime() = %v, esperado %v", slot.EndTime(), end)
		}
	})
}

func TestSetID(t *testing.T) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)
	end := start.Add(2 * time.Hour)
	slot := &AvailableSlot{
		staffID:   1,
		weekday:   Saturday,
		startTime: start,
		endTime:   end,
	}

	newID := 42
	slot.SetID(newID)

	if got := slot.ID(); got != newID {
		t.Errorf("ID() = %d, esperado %d após SetID", got, newID)
	}
}

func TestIsValidWeekday(t *testing.T) {
	tests := []struct {
		name    string
		weekday Weekday
		want    bool
	}{
		{"domingo válido", Sunday, true},
		{"segunda válido", Monday, true},
		{"terça válido", Tuesday, true},
		{"quarta válido", Wednesday, true},
		{"quinta válido", Thursday, true},
		{"sexta válido", Friday, true},
		{"sábado válido", Saturday, true},
		{"dia inválido", Weekday("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidWeekday(tt.weekday); got != tt.want {
				t.Errorf("isValidWeekday(%s) = %v, esperado %v", tt.weekday, got, tt.want)
			}
		})
	}
}

func TestFromTimeWeekday(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Weekday
		want      Weekday
		wantError bool
	}{
		{"domingo", time.Sunday, Sunday, false},
		{"segunda", time.Monday, Monday, false},
		{"terça", time.Tuesday, Tuesday, false},
		{"quarta", time.Wednesday, Wednesday, false},
		{"quinta", time.Thursday, Thursday, false},
		{"sexta", time.Friday, Friday, false},
		{"sábado", time.Saturday, Saturday, false},
		{"valor inválido", time.Weekday(7), Weekday(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromTimeWeekday(tt.input)

			if tt.wantError {
				if got != Weekday("") {
					t.Errorf("FromTimeWeekday(%d) = %s, esperado vazio", tt.input, got)
				}
			} else {
				if got != tt.want {
					t.Errorf("FromTimeWeekday(%d) = %s, esperado %s", tt.input, got, tt.want)
				}
			}
		})
	}
}
