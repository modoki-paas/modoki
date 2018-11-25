package constants

// ContextKeyType is for keys in context.WithValue
type ContextKeyType int

const (
	// UIDContextKeyForAPIKey is for uid identified from API Key
	UIDContextKeyForAPIKey ContextKeyType = iota
)
