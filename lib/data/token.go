package data

// Token .
type Token struct {
	UID       int64
	IsStore   bool
	Identity  string
	ExpiredIn int64
}
