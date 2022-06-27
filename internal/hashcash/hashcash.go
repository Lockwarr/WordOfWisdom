package hashcash

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Lockwarr/WordOfWisdom/internal/repository"
)

const zeroByte = 48 // '0'

type Stamp struct {
	Version    int    `json:"version"`
	ZerosCount int    `json:"zerosCount"`
	Date       int64  `json:"date"`
	Resource   string `json:"resource"`
	Rand       string `json:"rand"`
	Counter    int    `json:"counter"`
}

func (s Stamp) ToString() string {
	return fmt.Sprintf("%d:%d:%d:%s::%s:%d", s.Version, s.ZerosCount, s.Date, s.Resource, s.Rand, s.Counter)
}

func HashToStamp(hash string) (Stamp, error) {
	splitHash := strings.Split(hash, ":")
	if len(splitHash) < 6 {
		return Stamp{}, fmt.Errorf("invalid hash")
	}
	version, err := strconv.Atoi(splitHash[0])
	if err != nil {
		return Stamp{}, fmt.Errorf("invalid version")
	}
	zerosCount, err := strconv.Atoi(splitHash[1])
	if err != nil {
		return Stamp{}, fmt.Errorf("invalid zeros count")
	}
	date, err := strconv.ParseInt(splitHash[2], 10, 64)
	if err != nil {
		return Stamp{}, fmt.Errorf("invalid date")
	}
	counter, err := strconv.Atoi(splitHash[6])
	if err != nil {
		return Stamp{}, fmt.Errorf("invalid counter")
	}

	return Stamp{
		Version:    version,
		ZerosCount: zerosCount,
		Date:       date,
		Resource:   splitHash[3],
		Rand:       splitHash[5],
		Counter:    counter,
	}, nil
}

// IsHashCorrect - checks that hash has leading <zerosCount> zeros
func (s *Stamp) IsHashSolved() bool {
	actualZerosCount := 0
	hashString := sha1Hash(s.ToString())
	if s.ZerosCount > len(hashString) {
		return false
	}
	for _, ch := range hashString[:s.ZerosCount] {
		if ch != zeroByte {
			return false
		}
		actualZerosCount++
	}
	return s.ZerosCount <= actualZerosCount
}

// ComputeHashcash - calculates correct hashcash by bruteforce
// until the resulting hash satisfies the condition of IsHashCorrect
// maxIterations to prevent endless computing (0 or -1 to disable it)
func (s Stamp) ComputeHashcash(maxIterations int) (Stamp, error) {
	for s.Counter <= maxIterations || maxIterations <= 0 {
		if s.IsHashSolved() {
			return s, nil
		}
		// if hash don't have needed count of leading zeros, we are increasing counter and try next hash
		s.Counter++
	}
	return s, fmt.Errorf("max iterations exceeded")
}

// ValidStamp - checks if the stamp is valid
func (p *Stamp) ValidStamp(ctx context.Context, stampForValidation Stamp, repo repository.Repository) bool {
	// TODO: add debug logs on errors
	if stampForValidation.Date > time.Now().Add(2*24*time.Hour).Unix() {
		return false // futuristic
	}
	if stampForValidation.Date < time.Now().Add(-28*(24*time.Hour)).Unix() {
		return false // expired
	}

	if !p.IsHashSolved() {
		return false // insufficient zeroes
	}

	v, err := strconv.ParseInt(stampForValidation.Rand, 10, 64)
	if err != nil {
		return false // invalid rand
	}

	_, err = repo.GetIndicator(ctx, v)
	return err == nil
}

// sha1Hash - calculates sha1 hash from given string
func sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
