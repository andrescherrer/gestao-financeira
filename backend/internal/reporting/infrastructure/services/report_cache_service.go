package services

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"gestao-financeira/backend/pkg/cache"
)

// ReportCacheService provides caching for reports.
type ReportCacheService struct {
	cache *cache.CacheService
	ttl   time.Duration
}

// NewReportCacheService creates a new ReportCacheService instance.
func NewReportCacheService(cacheService *cache.CacheService, ttl time.Duration) *ReportCacheService {
	return &ReportCacheService{
		cache: cacheService,
		ttl:   ttl,
	}
}

// GenerateKey generates a cache key for a report based on its parameters.
func (s *ReportCacheService) GenerateKey(reportType string, params map[string]string) string {
	// Create a hash of the parameters for a consistent key
	hash := sha256.New()
	hash.Write([]byte(reportType))
	for key, value := range params {
		hash.Write([]byte(fmt.Sprintf("%s:%s", key, value)))
	}
	hashBytes := hash.Sum(nil)
	hashStr := hex.EncodeToString(hashBytes)
	return fmt.Sprintf("report:%s:%s", reportType, hashStr[:16])
}

// Get retrieves a cached report.
func (s *ReportCacheService) Get(key string, dest interface{}) (bool, error) {
	if s.cache == nil {
		return false, nil // Cache disabled
	}
	return s.cache.GetJSON(key, dest)
}

// Set stores a report in cache.
func (s *ReportCacheService) Set(key string, value interface{}) error {
	if s.cache == nil {
		return nil // Cache disabled
	}
	return s.cache.SetJSON(key, value, s.ttl)
}

// Invalidate invalidates cache for a specific report type or all reports.
func (s *ReportCacheService) Invalidate(reportType string) error {
	if s.cache == nil {
		return nil // Cache disabled
	}
	pattern := fmt.Sprintf("report:%s:*", reportType)
	return s.cache.DeletePattern(pattern)
}

// InvalidateAll invalidates all report caches.
func (s *ReportCacheService) InvalidateAll() error {
	if s.cache == nil {
		return nil // Cache disabled
	}
	return s.cache.DeletePattern("report:*")
}
