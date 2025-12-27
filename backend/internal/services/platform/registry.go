package platform

import (
	"fmt"
	"sync"

	"github.com/osmanmertacar/sosyal/backend/internal/database/models"
)

// PlatformRegistry manages all registered platform services
// It provides a centralized way to access platform-specific implementations
type PlatformRegistry struct {
	services map[models.Platform]PlatformService
	mu       sync.RWMutex // Protect concurrent access
}

// NewPlatformRegistry creates a new platform registry
func NewPlatformRegistry() *PlatformRegistry {
	return &PlatformRegistry{
		services: make(map[models.Platform]PlatformService),
	}
}

// Register registers a platform service
// This should be called during application initialization for each supported platform
func (r *PlatformRegistry) Register(service PlatformService) {
	r.mu.Lock()
	defer r.mu.Unlock()

	platform := service.GetPlatformName()
	r.services[platform] = service
}

// Get retrieves a platform service by platform name
// Returns an error if the platform is not supported
// Returns interface{} to avoid circular dependencies
func (r *PlatformRegistry) Get(platform models.Platform) (interface{}, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, ok := r.services[platform]
	if !ok {
		return nil, fmt.Errorf("platform %s is not supported or not registered", platform)
	}
	return service, nil
}

// GetAll returns all registered platform services
func (r *PlatformRegistry) GetAll() []PlatformService {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]PlatformService, 0, len(r.services))
	for _, service := range r.services {
		services = append(services, service)
	}
	return services
}

// IsSupported checks if a platform is supported
func (r *PlatformRegistry) IsSupported(platform models.Platform) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.services[platform]
	return ok
}

// GetSupportedPlatforms returns a list of all supported platform names
func (r *PlatformRegistry) GetSupportedPlatforms() []models.Platform {
	r.mu.RLock()
	defer r.mu.RUnlock()

	platforms := make([]models.Platform, 0, len(r.services))
	for platform := range r.services {
		platforms = append(platforms, platform)
	}
	return platforms
}

// Count returns the number of registered platforms
func (r *PlatformRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.services)
}
