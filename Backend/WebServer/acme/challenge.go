package acme

import (
	"net/http"
	"strings"
	tools "vServer/Backend/tools"
)

// HandleChallenge обрабатывает HTTP-01 ACME challenge
// Путь: /.well-known/acme-challenge/{token}
func (m *Manager) HandleChallenge(w http.ResponseWriter, r *http.Request) bool {
	path := r.URL.Path
	
	// Проверяем что это ACME challenge
	if !strings.HasPrefix(path, "/.well-known/acme-challenge/") {
		return false
	}
	
	// Извлекаем token из пути
	token := strings.TrimPrefix(path, "/.well-known/acme-challenge/")
	if token == "" {
		http.Error(w, "Token not found", http.StatusNotFound)
		return true
	}
	
	// Ищем challenge по token
	m.mu.RLock()
	challenge, exists := m.challenges[token]
	m.mu.RUnlock()
	
	if !exists {
		tools.Logs_file(1, "ACME", "⚠️ Challenge не найден для token: "+token, "logs_acme.log", false)
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return true
	}
	
	// Отдаём KeyAuth для подтверждения владения доменом
	tools.Logs_file(0, "ACME", "✅ Challenge ответ для домена: "+challenge.Domain, "logs_acme.log", true)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(challenge.KeyAuth))
	
	return true
}

// addChallenge добавляет challenge в хранилище
func (m *Manager) addChallenge(token, keyAuth, domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.challenges[token] = &ChallengeData{
		Token:   token,
		KeyAuth: keyAuth,
		Domain:  domain,
		Created: getCurrentTimestamp(),
	}
	
	tools.Logs_file(0, "ACME", "📝 Challenge добавлен для: "+domain, "logs_acme.log", false)
}

// removeChallenge удаляет challenge из хранилища
func (m *Manager) removeChallenge(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if challenge, exists := m.challenges[token]; exists {
		tools.Logs_file(0, "ACME", "🗑️ Challenge удалён для: "+challenge.Domain, "logs_acme.log", false)
		delete(m.challenges, token)
	}
}

// cleanupOldChallenges удаляет старые challenges (старше 10 минут)
func (m *Manager) cleanupOldChallenges() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	now := getCurrentTimestamp()
	maxAge := int64(600) // 10 минут
	
	for token, challenge := range m.challenges {
		if now-challenge.Created > maxAge {
			delete(m.challenges, token)
		}
	}
}
