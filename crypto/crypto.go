package crypto

/* sha3 code was adapted from @dowlandaiello */

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

const (
	// hashLength is the standardized length of a hash.
	hashLength = 32

	// tokenSize is the size (in bytes) of the token for the GenRandomToken()
	// method
	tokenSize = 64
)

// Hash represents the streamlined hash type to be used.
type Hash [hashLength]byte

// NewHash constructs a new hash given a hash, API so it returns an error.
func NewHash(b []byte) (Hash, error) {
	var hash Hash // Setup the hash
	bCropped := b // Setup the cropped buffer

	// Check the crop side
	if len(b) > len(hash) {
		bCropped = bCropped[len(bCropped)-hashLength:] // Crop the hash
	}

	// Copy the source
	copy(
		hash[hashLength-len(bCropped):],
		bCropped,
	)

	return hash, nil
}

// newHash constructs a new hash given a hash, returns no error
func newHash(b []byte) Hash {
	var hash Hash // Setup the hash
	bCropped := b // Setup the cropped buffer

	// Check the crop side
	if len(b) > len(hash) {
		bCropped = bCropped[len(bCropped)-hashLength:] // Crop the hash
	}

	// Copy the source
	copy(
		hash[hashLength-len(bCropped):],
		bCropped,
	)

	return hash
}

// Sha3 hashes a []byte using sha3.
func Sha3(b []byte) Hash {
	hash := sha3.New256()
	hash.Write(b)
	return newHash(hash.Sum(nil))
}

// Sha3String hashes a given message via sha3 and encodes the hashed
// message to a hex string.
func Sha3String(s string) string {
	return Sha3([]byte(s)).String()
}

// String returns the hash as a hex string.
func (hash Hash) String() string {
	b := hash.Bytes()
	return hex.EncodeToString(b) // Convert to a hex string
}

// Bytes returns the bytes of the hash.
func (hash Hash) Bytes() []byte {
	return hash[:]
}

// GenRandomToken generates a random token to be stored in the cookie's state.
func GenRandomToken() string {
	b := make([]byte, tokenSize) // Or 48
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
