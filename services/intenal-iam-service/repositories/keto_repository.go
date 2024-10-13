package repositories

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	keto "github.com/ory/keto-client-go"
	"go.uber.org/zap"
)

// KetoRepository interacts with the Keto API with logger
type KetoRepository struct {
	client *keto.APIClient
	logger *zap.Logger
}

// NewKetoRepository initializes a new KetoRepository with connection pooling, timeouts, and logger
func NewKetoRepository(logger *zap.Logger) *KetoRepository {
	ketoConfig := keto.NewConfiguration()
	ketoConfig.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	ketoConfig.Servers = keto.ServerConfigurations{
		{URL: "http://localhost:4466"},
	}
	return &KetoRepository{client: keto.NewAPIClient(ketoConfig), logger: logger}
}

// FetchUserRoles fetches all roles associated with a user from Keto
func (r *KetoRepository) FetchUserRoles(ctx context.Context, userID string) ([]string, error) {
	r.logger.Info("Fetching roles from Keto", zap.String("userID", userID))
	response, _, err := r.client.RelationshipApi.GetRelationships(ctx).
		Namespace("users").
		Object(fmt.Sprintf("user:%s", userID)).
		Relation("role").
		Execute()

	if err != nil {
		r.logger.Error("Failed to fetch user roles", zap.String("userID", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}

	var roles []string
	for _, tuple := range response.RelationTuples {
		roles = append(roles, *tuple.SubjectId)
	}

	r.logger.Info("Fetched roles for user", zap.String("userID", userID), zap.Strings("roles", roles))
	return roles, nil
}

// FetchUserUnits fetches all units associated with a user from Keto based on roles using subject_set
func (r *KetoRepository) FetchUserUnits(ctx context.Context, userID string, roles []string) ([]string, error) {
	var units []string
	r.logger.Info("Fetching units for user", zap.String("userID", userID))

	var wg sync.WaitGroup
	var mu sync.Mutex
	errorsChan := make(chan error, len(roles))

	for _, role := range roles {
		wg.Add(1)
		go func(role string) {
			defer wg.Done()
			r.logger.Info("Fetching units for role", zap.String("role", role))

			response, _, err := r.client.RelationshipApi.GetRelationships(ctx).
				Namespace("units").
				Relation("member").
				SubjectSetNamespace("users").
				SubjectSetObject(fmt.Sprintf("user:%s", userID)).
				SubjectSetRelation(role).
				Execute()

			if err != nil {
				r.logger.Error("Failed to fetch user units", zap.String("userID", userID), zap.String("role", role), zap.Error(err))
				errorsChan <- err
				return
			}

			mu.Lock()
			for _, tuple := range response.RelationTuples {
				units = append(units, tuple.Object)
			}
			mu.Unlock()
		}(role)
	}

	wg.Wait()
	close(errorsChan)

	if len(errorsChan) > 0 {
		return nil, fmt.Errorf("errors occurred while fetching units for user %s", userID)
	}

	r.logger.Info("Fetched units for user", zap.String("userID", userID), zap.Strings("units", units))
	return units, nil
}

// FetchAccessibleModules fetches all modules accessible by the user's unit and role using subject_set
func (r *KetoRepository) FetchAccessibleModules(ctx context.Context, unit, role string) ([]string, error) {
	r.logger.Info("Fetching accessible modules", zap.String("unit", unit), zap.String("role", role))

	response, _, err := r.client.RelationshipApi.GetRelationships(ctx).
		Namespace("modules").
		SubjectSetNamespace("units").
		SubjectSetObject(unit).
		SubjectSetRelation("member").
		Execute()

	if err != nil {
		r.logger.Error("Failed to fetch accessible modules", zap.String("unit", unit), zap.String("role", role), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch accessible modules: %w", err)
	}

	var modules []string
	for _, tuple := range response.RelationTuples {
		modules = append(modules, tuple.Object)
	}

	r.logger.Info("Fetched modules for unit", zap.String("unit", unit), zap.String("role", role), zap.Strings("modules", modules))
	return modules, nil
}
