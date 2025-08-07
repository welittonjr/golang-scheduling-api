package entities

import (
	"scheduling/internal/domain/valueobject"
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		userName    string
		email       string
		password    string
		role        string
		wantErr     bool
		errMsg      string
		checkFields func(*testing.T, *User)
	}{
		{
			name:     "criação de cliente bem-sucedida",
			id:       1,
			userName: "Well Ju",
			email:    "wellju@gmail.com",
			password: "senha123",
			role:     RoleClient,
			wantErr:  false,
			checkFields: func(t *testing.T, u *User) {
				if u.ID() != 1 {
					t.Errorf("ID esperado 1, obtido %d", u.ID())
				}
				if u.Name() != "Well Ju" {
					t.Errorf("nome esperado 'Well Ju', obtido'%s'", u.Name())
				}
				if u.Email() != "wellju@gmail.com" {
					t.Errorf("e-mail esperado 'wellju@gmail.com', obtido '%s'", u.Email())
				}
				if u.Role() != RoleClient {
					t.Errorf("função esperada 'cliente', obtida'%s'", u.Role())
				}
				if !u.CheckPassword("senha123") {
					t.Error("falha na verificação da senha para senha correta")
				}
				if u.CheckPassword("wrong") {
					t.Error("verificação de senha aprovada para senha errada")
				}
				if u.CreatedAt().IsZero() {
					t.Error("createdAt não deve ser zero")
				}
			},
		},
		{
			name:     "Nome vazio",
			id:       2,
			userName: "",
			email:    "wellju@gmail.com",
			password: "senha123",
			role:     RoleClient,
			wantErr:  true,
			errMsg:   "nome é obrigatório",
		},
		{
			name:     "senha curta",
			id:       3,
			userName: "Well ju",
			email:    "wellju@gmail.com",
			password: "curto",
			role:     RoleClient,
			wantErr:  true,
			errMsg:   "a senha deve ter ao menos 6 caracteres",
		},
		{
			name:     "função de usuario invalido",
			id:       4,
			userName: "funcao invalida",
			email:    "wellju@gmail.com",
			password: "senha123",
			role:     "invalid_role",
			wantErr:  true,
			errMsg:   "papel inválido",
		},
		{
			name:     "email invalido",
			id:       5,
			userName: "email invalido",
			email:    "email-invalido",
			password: "senha123",
			role:     RoleClient,
			wantErr:  true,
			errMsg:   valueobject.ErrInvalidEmail.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.id, tt.userName, tt.email, tt.password, tt.role)

			if tt.wantErr {
				if err == nil {
					t.Fatal("erro esperado, não obtive resultado algum")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("erro esperado '%s', obtido '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if tt.checkFields != nil {
				tt.checkFields(t, got)
			}
		})
	}
}

func TestRebuildUser(t *testing.T) {
	email, err := valueobject.NewEmail("wellju@gmail.com")
	if err != nil {
		t.Fatalf("falha ao criar e-mail: %v", err)
	}

	tests := []struct {
		name     string
		id       int
		userName string
		email    valueobject.Email
		role     string
		check    func(*testing.T, *User)
	}{
		{
			name:     "reconstruir cliente",
			id:       1,
			userName: "Cliente Reconstruído",
			email:    email,
			role:     RoleClient,
			check: func(t *testing.T, u *User) {
				if u.ID() != 1 {
					t.Errorf("esperado ID 1, obtido %d", u.ID())
				}
				if u.Name() != "Cliente Reconstruído" {
					t.Errorf("esperado nome 'Cliente Reconstruído', obtido '%s'", u.Name())
				}
				if u.Email() != "wellju@gmail.com" {
					t.Errorf("esperado e-mail 'wellju@gmail.com', obtido '%s'", u.Email())
				}
				if u.Role() != RoleClient {
					t.Errorf("esperado perfil 'client', obtido '%s'", u.Role())
				}
				if !u.CreatedAt().IsZero() {
					t.Error("esperado createdAt zerado")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := RebuildUser(tt.id, tt.userName, tt.email, tt.role)
			tt.check(t, user)
		})
	}
}

func TestUserRoles(t *testing.T) {
	email, err := valueobject.NewEmail("wellju@gmail.com")
	if err != nil {
		t.Fatalf("falha ao criar e-mail: %v", err)
	}

	tests := []struct {
		name       string
		role       string
		wantAdmin  bool
		wantClient bool
	}{
		{
			name:       "perfil de administrador",
			role:       RoleAdmin,
			wantAdmin:  true,
			wantClient: false,
		},
		{
			name:       "perfil de cliente",
			role:       RoleClient,
			wantAdmin:  false,
			wantClient: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := RebuildUser(1, "Usuário Teste", email, tt.role)

			if got := user.CanAccessAdminPanel(); got != tt.wantAdmin {
				t.Errorf("CanAccessAdminPanel() = %v, esperado %v", got, tt.wantAdmin)
			}
			if got := user.IsClient(); got != tt.wantClient {
				t.Errorf("IsClient() = %v, esperado %v", got, tt.wantClient)
			}
		})
	}
}

func TestUserMethods(t *testing.T) {
	email, err := valueobject.NewEmail("wellju@gmail.com")
	if err != nil {
		t.Fatalf("falha ao criar e-mail: %v", err)
	}

	t.Run("SetCreatedAt", func(t *testing.T) {
		user := RebuildUser(1, "Usuário Teste", email, RoleClient)
		testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

		user.SetCreatedAt(testTime)
		if got := user.CreatedAt(); got != testTime {
			t.Errorf("CreatedAt() = %v, esperado %v", got, testTime)
		}
	})

	t.Run("CheckPassword", func(t *testing.T) {
		user, err := NewUser(1, "Usuário Teste", "wellju@gmail.com", "senha_correta", RoleClient)
		if err != nil {
			t.Fatalf("falha ao criar usuário: %v", err)
		}

		if !user.CheckPassword("senha_correta") {
			t.Error("CheckPassword() = false para senha correta, esperado true")
		}
		if user.CheckPassword("senha_errada") {
			t.Error("CheckPassword() = true para senha errada, esperado false")
		}
	})
}

func TestUserGetters(t *testing.T) {
	emailStr := "wellju@gmail.com"
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		t.Fatalf("falha ao criar e-mail: %v", err)
	}
	createdAt := time.Now()

	user := &User{
		id:        1,
		name:      "Usuário Teste",
		email:     email,
		password:  "senha123",
		role:      RoleClient,
		createdAt: createdAt,
	}

	t.Run("ID", func(t *testing.T) {
		if got := user.ID(); got != 1 {
			t.Errorf("ID() = %d, esperado 1", got)
		}
	})

	t.Run("Name", func(t *testing.T) {
		if got := user.Name(); got != "Usuário Teste" {
			t.Errorf("Name() = %s, esperado 'Usuário Teste'", got)
		}
	})

	t.Run("Email", func(t *testing.T) {
		if got := user.Email(); got != emailStr {
			t.Errorf("Email() = %s, esperado '%s'", got, emailStr)
		}
	})

	t.Run("Role", func(t *testing.T) {
		if got := user.Role(); got != RoleClient {
			t.Errorf("Role() = %s, esperado 'client'", got)
		}
	})

	t.Run("CreatedAt", func(t *testing.T) {
		if got := user.CreatedAt(); got != createdAt {
			t.Errorf("CreatedAt() = %v, esperado %v", got, createdAt)
		}
	})
}
