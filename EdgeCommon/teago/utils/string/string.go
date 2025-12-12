package stringutil

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Md5 returns a lowercase MD5 hex digest of input.
func Md5(input string) string {
	sum := md5.Sum([]byte(input))
	return hex.EncodeToString(sum[:])
}

// Rand returns a random alphanumeric string with the given length.
func Rand(length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RegexpCompile wraps regexp.Compile.
func RegexpCompile(pattern string) (*regexp.Regexp, error) {
	return regexp.Compile(pattern)
}

// VersionCompare compares semantic-like version strings.
// Returns -1 if v1<v2, 0 if equal, 1 if v1>v2.
func VersionCompare(v1, v2 string) int {
	s1 := strings.Split(v1, ".")
	s2 := strings.Split(v2, ".")
	max := len(s1)
	if len(s2) > max {
		max = len(s2)
	}
	for i := 0; i < max; i++ {
		var n1, n2 int
		if i < len(s1) {
			n1 = parseInt(s1[i])
		}
		if i < len(s2) {
			n2 = parseInt(s2[i])
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}

func parseInt(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return 0
}
