module github.com/Blocky7277/GOPWD

go 1.24.2

require internal/util v1.0.0

replace internal/util => ./internal/util

require internal/cryptoutil v1.0.0

replace internal/cryptoutil => ./internal/crypto

replace cmd => ./cmd

require golang.org/x/term v0.32.0

require (
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
