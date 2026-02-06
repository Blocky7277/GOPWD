module util

go 1.24.2

require golang.org/x/term v0.33.0

require (
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
)

require internal/cryptoutil v1.0.0

replace internal/cryptoutil => ../crypto
