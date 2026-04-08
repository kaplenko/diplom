package models

import "testing"

func TestPaginationParams_Normalize(t *testing.T) {
	tests := []struct {
		name         string
		input        PaginationParams
		wantPage     int
		wantPageSize int
	}{
		{
			name:         "zero values default",
			input:        PaginationParams{Page: 0, PageSize: 0},
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "negative values default",
			input:        PaginationParams{Page: -5, PageSize: -1},
			wantPage:     1,
			wantPageSize: 10,
		},
		{
			name:         "page_size exceeds max",
			input:        PaginationParams{Page: 2, PageSize: 500},
			wantPage:     2,
			wantPageSize: 100,
		},
		{
			name:         "valid values unchanged",
			input:        PaginationParams{Page: 3, PageSize: 25},
			wantPage:     3,
			wantPageSize: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", tt.input.Page, tt.wantPage)
			}
			if tt.input.PageSize != tt.wantPageSize {
				t.Errorf("PageSize = %d, want %d", tt.input.PageSize, tt.wantPageSize)
			}
		})
	}
}

func TestPaginationParams_Offset(t *testing.T) {
	p := PaginationParams{Page: 3, PageSize: 10}
	if got := p.Offset(); got != 20 {
		t.Errorf("Offset() = %d, want 20", got)
	}
}
