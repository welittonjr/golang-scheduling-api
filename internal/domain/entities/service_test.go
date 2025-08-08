package entities

import (
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		staffID         int
		serviceName     string
		durationMinutes int
		price           float64
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "criar serviço válido",
			id:              1,
			staffID:         101,
			serviceName:     "Corte de Cabelo",
			durationMinutes: 30,
			price:           50.0,
			wantErr:         false,
		},
		{
			name:            "nome vazio deve retornar erro",
			id:              2,
			staffID:         102,
			serviceName:     "",
			durationMinutes: 45,
			price:           75.0,
			wantErr:         true,
			errMsg:          "nome do serviço é obrigatório",
		},
		{
			name:            "duração inválida deve retornar erro",
			id:              3,
			staffID:         103,
			serviceName:     "Manicure",
			durationMinutes: 0,
			price:           40.0,
			wantErr:         true,
			errMsg:          "a duração deve ser maior que zero",
		},
		{
			name:            "preço negativo deve retornar erro",
			id:              4,
			staffID:         104,
			serviceName:     "Pedicure",
			durationMinutes: 60,
			price:           -10.0,
			wantErr:         true,
			errMsg:          "preço não pode ser negativo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewService(tt.id, tt.staffID, tt.serviceName, tt.durationMinutes, tt.price)

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

			if service.ID() != tt.id {
				t.Errorf("ID esperado %d, obtido %d", tt.id, service.ID())
			}
			if service.StaffID() != tt.staffID {
				t.Errorf("StaffID esperado %d, obtido %d", tt.staffID, service.StaffID())
			}
			if service.Name() != tt.serviceName {
				t.Errorf("Name esperado '%s', obtido '%s'", tt.serviceName, service.Name())
			}
			if service.DurationMinutes() != tt.durationMinutes {
				t.Errorf("DurationMinutes esperado %d, obtido %d", tt.durationMinutes, service.DurationMinutes())
			}
			if service.Price() != tt.price {
				t.Errorf("Price esperado %.2f, obtido %.2f", tt.price, service.Price())
			}
			if service.CreatedAt().IsZero() {
				t.Error("CreatedAt não deve ser zero")
			}
		})
	}
}

func TestServiceGetters(t *testing.T) {
	now := time.Now()
	service := &Service{
		id:              1,
		staffID:         101,
		name:            "Massagem",
		durationMinutes: 60,
		price:           120.0,
		createdAt:       now,
	}

	t.Run("Testar getters", func(t *testing.T) {
		if got := service.ID(); got != 1 {
			t.Errorf("ID() = %d, esperado 1", got)
		}
		if got := service.StaffID(); got != 101 {
			t.Errorf("StaffID() = %d, esperado 101", got)
		}
		if got := service.Name(); got != "Massagem" {
			t.Errorf("Name() = %s, esperado 'Massagem'", got)
		}
		if got := service.DurationMinutes(); got != 60 {
			t.Errorf("DurationMinutes() = %d, esperado 60", got)
		}
		if got := service.Price(); got != 120.0 {
			t.Errorf("Price() = %.2f, esperado 120.00", got)
		}
		if got := service.CreatedAt(); got != now {
			t.Errorf("CreatedAt() = %v, esperado %v", got, now)
		}
	})
}

func TestChangePrice(t *testing.T) {
	service, _ := NewService(1, 101, "Design de Sobrancelhas", 30, 45.0)

	tests := []struct {
		name     string
		newPrice float64
		wantErr  bool
	}{
		{
			name:     "alterar para preço válido",
			newPrice: 50.0,
			wantErr:  false,
		},
		{
			name:     "alterar para preço negativo deve retornar erro",
			newPrice: -10.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ChangePrice(tt.newPrice)

			if tt.wantErr {
				if err == nil {
					t.Error("esperado erro, mas nenhum foi retornado")
				} else if err.Error() != "preço não pode ser negativo" {
					t.Errorf("mensagem de erro incorreta, obtido: '%s'", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("não esperava erro, mas obteve: %v", err)
				}
				if service.Price() != tt.newPrice {
					t.Errorf("Price esperado %.2f, obtido %.2f", tt.newPrice, service.Price())
				}
			}
		})
	}
}

func TestChangeDuration(t *testing.T) {
	service, _ := NewService(1, 101, "Depilação", 45, 80.0)

	tests := []struct {
		name        string
		newDuration int
		wantErr     bool
	}{
		{
			name:        "alterar para duração válida",
			newDuration: 60,
			wantErr:     false,
		},
		{
			name:        "alterar para duração inválida deve retornar erro",
			newDuration: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ChangeDuration(tt.newDuration)

			if tt.wantErr {
				if err == nil {
					t.Error("esperado erro, mas nenhum foi retornado")
				} else if err.Error() != "duração inválida" {
					t.Errorf("mensagem de erro incorreta, obtido: '%s'", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("não esperava erro, mas obteve: %v", err)
				}
				if service.DurationMinutes() != tt.newDuration {
					t.Errorf("DurationMinutes esperado %d, obtido %d", tt.newDuration, service.DurationMinutes())
				}
			}
		})
	}
}
