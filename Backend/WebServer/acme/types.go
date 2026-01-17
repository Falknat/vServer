package acme

import (
	"crypto/ecdsa"
	"sync"
)

// ChallengeData хранит данные для HTTP-01 challenge
type ChallengeData struct {
	Token    string
	KeyAuth  string
	Domain   string
	Created  int64
}

// Manager управляет ACME сертификатами
type Manager struct {
	mu              sync.RWMutex
	accountKey      *ecdsa.PrivateKey
	challenges      map[string]*ChallengeData // token -> challenge data
	certDir         string
	acmeDir         string
	isProduction    bool
}

// CertInfo информация о сертификате
type CertInfo struct {
	Domain      string `json:"domain"`
	Issuer      string `json:"issuer"`
	NotBefore   string `json:"not_before"`
	NotAfter    string `json:"not_after"`
	DaysLeft    int    `json:"days_left"`
	AutoCreated bool   `json:"auto_created"`
}

// ObtainResult результат получения сертификата
type ObtainResult struct {
	Success bool   `json:"success"`
	Domain  string `json:"domain"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
