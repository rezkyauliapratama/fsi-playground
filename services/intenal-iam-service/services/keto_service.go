package services

import (
	"context"
	"internal-iam-service/repositories"
	"sync"
	"time"

	"go.uber.org/zap"
)

// KetoService manages role and module retrieval with a logger
type KetoService struct {
	repo   *repositories.KetoRepository
	logger *zap.Logger
}

// NewKetoService initializes a new KetoService with a logger
func NewKetoService(repo *repositories.KetoRepository, logger *zap.Logger) *KetoService {
	return &KetoService{repo: repo, logger: logger}
}

// retryLogic retries a function that returns a slice of strings with a delay between retries
func retryLogic(attempts int, delay time.Duration, fn func() ([]string, error), logger *zap.Logger) ([]string, error) {
	var result []string
	var err error
	for i := 0; i < attempts; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		logger.Warn("Retrying after failure", zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(delay)
	}
	return nil, err
}

// GetUserModules fetches the list of modules categorized by actions (view, manage) for a user
func (s *KetoService) GetUserModules(ctx context.Context, userID string) (map[string][]string, error) {
	var wg sync.WaitGroup
	modulesByAction := make(map[string][]string)
	modulesChan := make(chan map[string][]string, 1)
	errorChan := make(chan error, 1)

	s.logger.Info("Fetching roles for user", zap.String("userID", userID))

	// Step 1: Fetch all roles for the user
	roles, roleErr := retryLogic(3, 1*time.Second, func() ([]string, error) {
		return s.repo.FetchUserRoles(ctx, userID)
	}, s.logger)
	if roleErr != nil {
		return nil, roleErr
	}

	s.logger.Info("Fetching units for user", zap.String("userID", userID))

	// Step 2: Fetch all units associated with the user based on roles
	units, unitErr := retryLogic(3, 1*time.Second, func() ([]string, error) {
		return s.repo.FetchUserUnits(ctx, userID, roles)
	}, s.logger)
	if unitErr != nil {
		return nil, unitErr
	}

	s.logger.Info("Fetching modules for each role-unit combination", zap.String("userID", userID))

	// Step 3: Fetch modules for each combination of role and unit
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, unit := range units {
			for _, role := range roles {
				modulesForUnitRole, moduleErr := retryLogic(3, 1*time.Second, func() ([]string, error) {
					return s.repo.FetchAccessibleModules(ctx, unit, role)
				}, s.logger)
				if moduleErr != nil {
					errorChan <- moduleErr
					return
				}
				// Aggregate the modules by their action (view/manage)
				for action, modules := range groupModulesByAction(modulesForUnitRole) {
					modulesByAction[action] = append(modulesByAction[action], modules...)
				}
			}
		}
		modulesChan <- modulesByAction
	}()

	// Wait for goroutines to finish
	go func() {
		wg.Wait()
		close(errorChan)
		close(modulesChan)
	}()

	// Handle errors and results
	select {
	case err := <-errorChan:
		s.logger.Error("Error fetching modules", zap.String("userID", userID), zap.Error(err))
		return nil, err
	case modules := <-modulesChan:
		s.logger.Info("Successfully fetched modules for user", zap.String("userID", userID), zap.Any("modules", modules))
		return modules, nil
	}
}

// groupModulesByAction helps to categorize modules by their action (view, manage)
func groupModulesByAction(modules []string) map[string][]string {
	modulesByAction := make(map[string][]string)
	// Assuming you have a way to distinguish between view and manage actions
	// For example, you can infer this from the module name or a different field
	for _, module := range modules {
		// Example of categorization logic (you can modify this based on actual module details)
		if isManageAction(module) {
			modulesByAction["manage"] = append(modulesByAction["manage"], module)
		} else {
			modulesByAction["view"] = append(modulesByAction["view"], module)
		}
	}
	return modulesByAction
}

// Example function that determines if a module is "manage" or "view"
func isManageAction(module string) bool {
	// Logic to determine if it's a manage action based on the module (this is just an example)
	return true // Replace with real logic
}
