package pow

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	hcVersion      = 1         // Hashcash version 1
	numOfZeroes    = 4         // 4 first zeroes in hash
	zeroAsByte     = 48        // Used to check each rune in hash
	randUpperBound = 1_000_000 // Random nonce
)

// NewHashcash generates a challenge for a client to solve
func NewHashcash(resource string) Hashcash {
	return Hashcash{
		Version:  hcVersion,
		Bits:     numOfZeroes,
		Date:     time.Now().Unix(),
		Resource: resource, // ip:port of client
		RandB64:  base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", rand.Intn(randUpperBound)))),
		Counter:  0,
	}
}

// Hashcash is implementation of Adam Back's (Hashcash PoW system)[https://en.wikipedia.org/wiki/Hashcash]
type Hashcash struct {
	Version  uint8
	Bits     uint8
	Date     int64
	Resource string
	RandB64  string
	Counter  int
}

// String implements Stringer interface
func (h Hashcash) String() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Version, h.Bits, h.Date, h.Resource, h.RandB64, h.Counter)
}

// GetHash returns hash of Hashcash instance
func (h Hashcash) GetHash() string {
	hash := sha256.New()
	hash.Write([]byte(h.String()))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Solve searches solution for challenge via brute force
// MaxInt64 is used for demonstration purposes, may be used another upper bound value
// Each unsuccessful iteration increments counter if leading zeroes
func (h Hashcash) Solve() (Hashcash, error) {
	for i := 0; i < math.MaxInt64; i++ {
		hash := h.GetHash()
		if solutionFound(hash) {
			return h, nil
		}
		h.Counter++
	}
	return h, fmt.Errorf("max iterations exceeded")
}

// solutionFound checks if number of leading zeroes in hash
func solutionFound(hash string) bool {
	if len(hash) < numOfZeroes {
		return false
	}

	for _, r := range hash[:numOfZeroes] {
		if r != zeroAsByte {
			return false
		}
	}

	return true
}
