package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewTransactionManager(t *testing.T) {
	tests := []struct {
		name string
		db   *sql.DB
		want func(*testing.T, TransactionManager)
	}{
		{
			name: "criação bem-sucedida do transaction manager",
			db:   &sql.DB{},
			want: func(t *testing.T, tm TransactionManager) {
				if tm == nil {
					t.Error("transaction manager não deve ser nil")
				}

				if _, ok := tm.(*transactionManager); !ok {
					t.Error("transaction manager deve ser do tipo *transactionManager")
				}
			},
		},
		{
			name: "criação com db nil",
			db:   nil,
			want: func(t *testing.T, tm TransactionManager) {
				if tm == nil {
					t.Error("transaction manager não deve ser nil")
				}
				if _, ok := tm.(*transactionManager); !ok {
					t.Error("transaction manager deve ser do tipo *transactionManager")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTransactionManager(tt.db)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestTransactionManager_StartTransaction(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(sqlmock.Sqlmock)
		ctx     context.Context
		wantErr bool
		errMsg  string
	}{
		{
			name: "transação iniciada com sucesso",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			ctx:     context.Background(),
			wantErr: false,
		},
		{
			name: "transação iniciada com sucesso - contexto com timeout",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
			},
			ctx:     context.WithValue(context.Background(), "timeout", "5s"),
			wantErr: false,
		},
		{
			name: "erro ao iniciar transação",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("database connection error"))
			},
			ctx:     context.Background(),
			wantErr: true,
			errMsg:  "failed to start transaction: database connection error",
		},
		{
			name: "erro de contexto cancelado",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(context.Canceled)
			},
			ctx:     context.Background(),
			wantErr: true,
			errMsg:  "failed to start transaction: context canceled",
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

			tm := NewTransactionManager(db)
			tx, err := tm.StartTransaction(tt.ctx)

			if tt.wantErr {
				if err == nil {
					t.Fatal("erro esperado, mas não obtive nenhum")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("erro esperado '%s', obtido '%s'", tt.errMsg, err.Error())
				}
				if tx != nil {
					t.Error("transação deve ser nil quando há erro")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if tx == nil {
				t.Error("transação não deve ser nil em caso de sucesso")
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("expectativas do mock não foram atendidas: %v", err)
			}
		})
	}
}

func TestTransactionManager_Commit(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(sqlmock.Sqlmock) *sql.Tx
		wantErr bool
		errMsg  string
	}{
		{
			name: "commit realizado com sucesso",
			mockFn: func(mock sqlmock.Sqlmock) *sql.Tx {
				mock.ExpectBegin()
				mock.ExpectCommit()
				db, _ := sql.Open("sqlmock", "")
				tx, _ := db.Begin()
				return tx
			},
			wantErr: false,
		},
		{
			name: "erro durante commit",
			mockFn: func(mock sqlmock.Sqlmock) *sql.Tx {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
				db, _ := sql.Open("sqlmock", "")
				tx, _ := db.Begin()
				return tx
			},
			wantErr: true,
			errMsg:  "failed to commit transaction: commit failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("erro ao criar mock do banco: %v", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			if tt.wantErr {
				mock.ExpectCommit().WillReturnError(errors.New("commit failed"))
			} else {
				mock.ExpectCommit()
			}

			tm := NewTransactionManager(db)
			tx, err := tm.StartTransaction(context.Background())
			if err != nil {
				t.Fatalf("erro inesperado ao iniciar transação: %v", err)
			}

			err = tm.Commit(tx)

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

func TestTransactionManager_Rollback(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
		wantErr bool
		errMsg  string
	}{
		{
			name:    "rollback realizado com sucesso",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "rollback com ErrTxDone - não deve retornar erro",
			mockErr: sql.ErrTxDone,
			wantErr: false,
		},
		{
			name:    "erro durante rollback",
			mockErr: errors.New("rollback failed"),
			wantErr: true,
			errMsg:  "failed to rollback transaction: rollback failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("erro ao criar mock do banco: %v", err)
			}
			defer db.Close()

			mock.ExpectBegin()
			if tt.mockErr != nil {
				mock.ExpectRollback().WillReturnError(tt.mockErr)
			} else {
				mock.ExpectRollback()
			}

			tm := NewTransactionManager(db)
			tx, err := tm.StartTransaction(context.Background())
			if err != nil {
				t.Fatalf("erro inesperado ao iniciar transação: %v", err)
			}

			err = tm.Rollback(tx)

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

func TestTransactionManager_WithTransaction(t *testing.T) {
	tests := []struct {
		name         string
		mockFn       func(sqlmock.Sqlmock)
		txFunc       func(tx *sql.Tx) error
		wantErr      bool
		errMsg       string
		expectCommit bool
	}{
		{
			name: "transação executada com sucesso - commit automático",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO test").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			txFunc: func(tx *sql.Tx) error {
				_, err := tx.Exec("INSERT INTO test (id) VALUES (1)")
				return err
			},
			wantErr:      false,
			expectCommit: true,
		},
		{
			name: "erro na função - rollback automático",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			txFunc: func(tx *sql.Tx) error {
				return errors.New("business logic error")
			},
			wantErr:      true,
			errMsg:       "business logic error",
			expectCommit: false,
		},
		{
			name: "erro ao iniciar transação",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			txFunc: func(tx *sql.Tx) error {
				return nil
			},
			wantErr:      true,
			errMsg:       "failed to start transaction: connection error",
			expectCommit: false,
		},
		{
			name: "erro no commit - rollback automático",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO test").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
				mock.ExpectRollback()
			},
			txFunc: func(tx *sql.Tx) error {
				_, err := tx.Exec("INSERT INTO test (id) VALUES (1)")
				return err
			},
			wantErr:      true,
			errMsg:       "failed to commit transaction: commit error",
			expectCommit: false,
		},
		{
			name: "função com operações múltiplas - sucesso",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec("UPDATE users").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			txFunc: func(tx *sql.Tx) error {
				if _, err := tx.Exec("INSERT INTO users (name) VALUES ('test')"); err != nil {
					return err
				}
				if _, err := tx.Exec("UPDATE users SET active = true WHERE id = 1"); err != nil {
					return err
				}
				return nil
			},
			wantErr:      false,
			expectCommit: true,
		},
		{
			name: "função que retorna nil - commit automático",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			txFunc: func(tx *sql.Tx) error {
				return nil
			},
			wantErr:      false,
			expectCommit: true,
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

			tm := NewTransactionManager(db)
			err = tm.WithTransaction(context.Background(), tt.txFunc)

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

func TestTransactionManager_WithTransaction_ContextCancellation(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		mockFn func(sqlmock.Sqlmock)
		txFunc func(tx *sql.Tx) error
		wantErr bool
		errContains string
	}{
		{
			name: "contexto cancelado durante início da transação",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(context.Canceled)
			},
			txFunc: func(tx *sql.Tx) error {
				return nil
			},
			wantErr: true,
			errContains: "context canceled",
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

			tm := NewTransactionManager(db)
			err = tm.WithTransaction(tt.ctx, tt.txFunc)

			if tt.wantErr {
				if err == nil {
					t.Fatal("erro esperado, mas não obtive nenhum")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("erro deve conter '%s', obtido '%s'", tt.errContains, err.Error())
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

func TestTransactionManager_WithTransaction_PanicRecovery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock do banco: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	tm := NewTransactionManager(db)
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("esperava panic, mas não ocorreu")
		}
	}()

	err = tm.WithTransaction(context.Background(), func(tx *sql.Tx) error {
		panic("test panic")
	})

	t.Error("não deveria chegar aqui após panic")
}

func TestTransactionManager_IntegrationScenarios(t *testing.T) {
	tests := []struct {
		name        string
		mockFn      func(sqlmock.Sqlmock)
		operations  func(tm TransactionManager) error
		wantErr     bool
		errContains string
	}{
		{
			name: "múltiplas transações sequenciais",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			operations: func(tm TransactionManager) error {
				err := tm.WithTransaction(context.Background(), func(tx *sql.Tx) error {
					_, err := tx.Exec("INSERT INTO test (id) VALUES (1)")
					return err
				})
				if err != nil {
					return err
				}

				return tm.WithTransaction(context.Background(), func(tx *sql.Tx) error {
					_, err := tx.Exec("UPDATE test SET name = 'updated' WHERE id = 1")
					return err
				})
			},
			wantErr: false,
		},
		{
			name: "transação manual seguida de WithTransaction",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
				
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			operations: func(tm TransactionManager) error {
				tx, err := tm.StartTransaction(context.Background())
				if err != nil {
					return err
				}

				_, err = tx.Exec("INSERT INTO test (id) VALUES (1)")
				if err != nil {
					tm.Rollback(tx)
					return err
				}

				err = tm.Commit(tx)
				if err != nil {
					return err
				}

				return tm.WithTransaction(context.Background(), func(tx *sql.Tx) error {
					_, err := tx.Exec("UPDATE test SET name = 'updated' WHERE id = 1")
					return err
				})
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

			tm := NewTransactionManager(db)
			err = tt.operations(tm)

			if tt.wantErr {
				if err == nil {
					t.Fatal("erro esperado, mas não obtive nenhum")
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("erro deve conter '%s', obtido '%s'", tt.errContains, err.Error())
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

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
			s[len(s)-len(substr):] == substr || 
			containsAt(s, substr))))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}