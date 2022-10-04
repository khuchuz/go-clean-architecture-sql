package utils

import (
	"testing"
)

func TestHashThis(t *testing.T) {
	password := "Password"
	salt := "salt"
	ans := HashThis(password, salt)
	expected := "38a8fde622c0cf723934ba7138a72beaccfc69d4"
	if ans != expected {
		t.Errorf("HashThis(%s, %s) = %s; want %s", password, salt, ans, expected)
	}
}

func BenchmarkHashThis(b *testing.B) {
	password := "Password"
	salt := "salt"
	for i := 0; i < b.N; i++ {
		HashThis(password, salt)
	}
}
