package entities

import (
	"testing"
	"time"
)

func TestNewAppointment(t *testing.T) {
	futureTime := time.Now().Add(2 * time.Hour)
	pastTime := time.Now().Add(-2 * time.Hour)

	tests := []struct {
		name        string
		clientID    int
		staffID     int
		serviceID   int
		scheduledAt time.Time
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "criar agendamento válido",
			clientID:    1,
			staffID:     101,
			serviceID:   1001,
			scheduledAt: futureTime,
			wantErr:     false,
		},
		{
			name:        "clientID zero deve retornar erro",
			clientID:    0,
			staffID:     101,
			serviceID:   1001,
			scheduledAt: futureTime,
			wantErr:     true,
			errMsg:      "cliente, profissional e serviço são obrigatórios",
		},
		{
			name:        "staffID zero deve retornar erro",
			clientID:    1,
			staffID:     0,
			serviceID:   1001,
			scheduledAt: futureTime,
			wantErr:     true,
			errMsg:      "cliente, profissional e serviço são obrigatórios",
		},
		{
			name:        "serviceID zero deve retornar erro",
			clientID:    1,
			staffID:     101,
			serviceID:   0,
			scheduledAt: futureTime,
			wantErr:     true,
			errMsg:      "cliente, profissional e serviço são obrigatórios",
		},
		{
			name:        "agendamento no passado deve retornar erro",
			clientID:    1,
			staffID:     101,
			serviceID:   1001,
			scheduledAt: pastTime,
			wantErr:     true,
			errMsg:      "não é possível agendar para o passado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment, err := NewAppointment(tt.clientID, tt.staffID, tt.serviceID, tt.scheduledAt)

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

			if appointment.ClientID() != tt.clientID {
				t.Errorf("ClientID esperado %d, obtido %d", tt.clientID, appointment.ClientID())
			}
			if appointment.StaffID() != tt.staffID {
				t.Errorf("StaffID esperado %d, obtido %d", tt.staffID, appointment.StaffID())
			}
			if appointment.ServiceID() != tt.serviceID {
				t.Errorf("ServiceID esperado %d, obtido %d", tt.serviceID, appointment.ServiceID())
			}
			if !appointment.ScheduledAt().Equal(tt.scheduledAt) {
				t.Errorf("ScheduledAt esperado %v, obtido %v", tt.scheduledAt, appointment.ScheduledAt())
			}
			if appointment.Status() != StatusScheduled {
				t.Errorf("Status esperado %s, obtido %s", StatusScheduled, appointment.Status())
			}
			if appointment.CreatedAt().IsZero() {
				t.Error("CreatedAt não deve ser zero")
			}
		})
	}
}

func TestAppointmentGetters(t *testing.T) {
	now := time.Now()
	scheduledAt := now.Add(1 * time.Hour)
	appointment := &Appointment{
		id:          1,
		clientID:    2,
		staffID:     3,
		serviceID:   4,
		scheduledAt: scheduledAt,
		status:      StatusScheduled,
		createdAt:   now,
	}

	t.Run("Testar getters", func(t *testing.T) {
		if got := appointment.ID(); got != 1 {
			t.Errorf("ID() = %d, esperado 1", got)
		}
		if got := appointment.ClientID(); got != 2 {
			t.Errorf("ClientID() = %d, esperado 2", got)
		}
		if got := appointment.StaffID(); got != 3 {
			t.Errorf("StaffID() = %d, esperado 3", got)
		}
		if got := appointment.ServiceID(); got != 4 {
			t.Errorf("ServiceID() = %d, esperado 4", got)
		}
		if !appointment.ScheduledAt().Equal(scheduledAt) {
			t.Errorf("ScheduledAt() = %v, esperado %v", appointment.ScheduledAt(), scheduledAt)
		}
		if appointment.Status() != StatusScheduled {
			t.Errorf("Status() = %s, esperado %s", appointment.Status(), StatusScheduled)
		}
		if !appointment.CreatedAt().Equal(now) {
			t.Errorf("CreatedAt() = %v, esperado %v", appointment.CreatedAt(), now)
		}
	})
}

func TestAppointmentStatusMethods(t *testing.T) {
	tests := []struct {
		name           string
		initialStatus  AppointmentStatus
		methodToCall   func(*Appointment)
		expectedStatus AppointmentStatus
		isScheduled    bool
		isCanceled     bool
	}{
		{
			name:           "cancelar agendamento marcado",
			initialStatus:  StatusScheduled,
			methodToCall:   (*Appointment).Cancel,
			expectedStatus: StatusCancelled,
			isScheduled:    false,
			isCanceled:     true,
		},
		{
			name:           "completar agendamento marcado",
			initialStatus:  StatusScheduled,
			methodToCall:   (*Appointment).Complete,
			expectedStatus: StatusCompleted,
			isScheduled:    false,
			isCanceled:     false,
		},
		{
			name:           "cancelar agendamento já cancelado",
			initialStatus:  StatusCancelled,
			methodToCall:   (*Appointment).Cancel,
			expectedStatus: StatusCancelled,
			isScheduled:    false,
			isCanceled:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appointment := &Appointment{
				status: tt.initialStatus,
			}

			tt.methodToCall(appointment)

			if appointment.Status() != tt.expectedStatus {
				t.Errorf("Status() = %s, esperado %s", appointment.Status(), tt.expectedStatus)
			}
			if got := appointment.IsScheduled(); got != tt.isScheduled {
				t.Errorf("IsScheduled() = %v, esperado %v", got, tt.isScheduled)
			}
			if got := appointment.IsCanceled(); got != tt.isCanceled {
				t.Errorf("IsCanceled() = %v, esperado %v", got, tt.isCanceled)
			}
		})
	}
}

func TestSetMethods(t *testing.T) {
	t.Run("Testar SetID", func(t *testing.T) {
		appointment := &Appointment{}
		newID := 42
		appointment.SetID(newID)

		if appointment.ID() != newID {
			t.Errorf("ID() = %d, esperado %d após SetID", appointment.ID(), newID)
		}
	})

	t.Run("Testar SetStatus", func(t *testing.T) {
		appointment := &Appointment{}

		tests := []struct {
			statusToSet    string
			expectedStatus AppointmentStatus
		}{
			{"scheduled", StatusScheduled},
			{"completed", StatusCompleted},
			{"cancelled", StatusCancelled},
			{"invalid", StatusScheduled},
		}

		for _, tt := range tests {
			appointment.SetStatus(tt.statusToSet)
			if appointment.Status() != tt.expectedStatus {
				t.Errorf("Após SetStatus(%s), Status() = %s, esperado %s",
					tt.statusToSet, appointment.Status(), tt.expectedStatus)
			}
		}
	})

	t.Run("Testar SetCreatedAt", func(t *testing.T) {
		appointment := &Appointment{}
		newTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		appointment.SetCreatedAt(newTime)

		if !appointment.CreatedAt().Equal(newTime) {
			t.Errorf("CreatedAt() = %v, esperado %v", appointment.CreatedAt(), newTime)
		}
	})
}
