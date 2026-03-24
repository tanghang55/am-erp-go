package seed

import "testing"

func TestGeneratePasswordReturnsRequestedLength(t *testing.T) {
	password, err := GeneratePassword(20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(password) != 20 {
		t.Fatalf("expected password length 20, got %d", len(password))
	}
}

func TestGeneratePasswordRejectsShortLength(t *testing.T) {
	if _, err := GeneratePassword(7); err == nil {
		t.Fatalf("expected error for short password length")
	}
}

func TestSplitStatementsSkipsBlankStatements(t *testing.T) {
	statements := SplitStatements("SET NAMES utf8mb4;;\n\nDELETE FROM menu;\n")
	if len(statements) != 2 {
		t.Fatalf("expected 2 statements, got %#v", statements)
	}
	if statements[0] != "SET NAMES utf8mb4" {
		t.Fatalf("unexpected first statement: %s", statements[0])
	}
	if statements[1] != "DELETE FROM menu" {
		t.Fatalf("unexpected second statement: %s", statements[1])
	}
}
