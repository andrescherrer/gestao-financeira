package valueobjects

import (
	"testing"
)

func TestGenerateSlugFromName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "Alimentação", "alimentacao"},
		{"name with spaces", "Transporte Público", "transporte-publico"},
		{"name with numbers", "Conta 123", "conta-123"},
		{"name with special chars", "Alimentação & Bebidas", "alimentacao-bebidas"},
		{"accented characters", "Educação", "educacao"},
		{"cedilha (ç)", "Correção", "correcao"},
		{"cedilha (ç) in middle", "Ação", "acao"},
		{"cedilha (ç) multiple", "Preço e Correção", "preco-e-correcao"},
		{"mixed case", "AlimentAção", "alimentacao"},
		{"multiple spaces", "Alimentação   Pública", "alimentacao-publica"},
		{"leading/trailing spaces", "  Alimentação  ", "alimentacao"},
		{"only special chars", "!!!", "categoria"},
		{"empty string", "", "categoria"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := GenerateSlugFromName(tt.input)
			if slug.Value() != tt.expected {
				t.Errorf("GenerateSlugFromName(%q) = %q, want %q", tt.input, slug.Value(), tt.expected)
			}
		})
	}
}

func TestNewCategorySlug(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid slug", "alimentacao", false},
		{"valid slug with hyphens", "transporte-publico", false},
		{"valid slug with numbers", "conta-123", false},
		{"empty slug", "", true},
		{"too short", "a", true},
		{"invalid characters", "alimentação", true},
		{"uppercase", "ALIMENTACAO", true},
		{"spaces", "alimentacao publica", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := NewCategorySlug(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("NewCategorySlug() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && slug.IsEmpty() {
				t.Error("NewCategorySlug() returned empty slug for valid input")
			}
		})
	}
}
