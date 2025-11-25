package validator

import "testing"

func TestURLValidator(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		allow   []string
		wantErr bool
		wantMsg string
	}{
		{
			name:    "Valid HTTP URL",
			url:     "http://example.com",
			allow:   []string{"http", "https"},
			wantErr: false,
			wantMsg: "",
		},
		{
			name:    "Valid HTTPS URL",
			url:     "https://example.com",
			allow:   []string{"http", "https"},
			wantErr: false,
			wantMsg: "",
		},
		{
			name:    "Valid relative URL",
			url:     "/relative/path",
			allow:   []string{"http", "https"},
			wantErr: false,
			wantMsg: "",
		},
		{
			name:    "Invalid Protocol",
			url:     "ftp://example.com",
			allow:   []string{"http", "https"},
			wantErr: true,
			wantMsg: "url protocol is not allowed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := URLValidator(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.wantMsg {
				t.Errorf("URLValidator() error message = %v, wantMsg %v", err.Error(), tt.wantMsg)
			}
		})
	}
}
