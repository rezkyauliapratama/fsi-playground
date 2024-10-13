package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	keto "github.com/ory/keto-client-go"
	"go.uber.org/zap"
)

// KetoRepository handles interaction with the Keto API
type KetoRepository struct {
	client *keto.APIClient
	logger *zap.Logger
}

// NewKetoRepository initializes a new KetoRepository with connection pooling and logger
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
func (r *KetoRepository) FetchUserRoles(ctx context.Context, userID string) (map[string]struct{}, error) {
	r.logger.Debug("Fetching roles from Keto", zap.String("userID", userID))
	response, _, err := r.client.RelationshipApi.GetRelationships(ctx).
		Namespace("users").
		Object(fmt.Sprintf("user:%s", userID)).
		Relation("role").
		Execute()

	if err != nil {
		r.logger.Error("Failed to fetch user roles", zap.String("userID", userID), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}

	roles := make(map[string]struct{})
	for _, tuple := range response.RelationTuples {
		roles[*tuple.SubjectId] = struct{}{}
	}

	r.logger.Debug("Fetched roles for user", zap.String("userID", userID), zap.Int("roleCount", len(roles)))
	return roles, nil
}

// FetchUserUnits fetches all units associated with a user based on their roles
func (r *KetoRepository) FetchUserUnits(ctx context.Context, userID string, roles map[string]struct{}) ([]string, error) {
	var units []string
	r.logger.Debug("Fetching units for user", zap.String("userID", userID))

	for role := range roles {
		r.logger.Debug("Fetching units for role", zap.String("role", role))

		response, _, err := r.client.RelationshipApi.GetRelationships(ctx).
			Namespace("units").
			Relation("member").
			SubjectSetNamespace("users").
			SubjectSetObject(fmt.Sprintf("user:%s", userID)).
			SubjectSetRelation(role).
			Execute()

		if err != nil {
			r.logger.Error("Failed to fetch user units", zap.String("userID", userID), zap.String("role", role), zap.Error(err))
			return nil, fmt.Errorf("failed to fetch user units: %w", err)
		}

		for _, tuple := range response.RelationTuples {
			units = append(units, tuple.Object)
		}
	}

	r.logger.Debug("Fetched units for user", zap.String("userID", userID), zap.Strings("units", units))
	return units, nil
}

// FetchAccessibleModules fetches all modules accessible by the user's unit and role, using dynamic actions
func (r *KetoRepository) FetchAccessibleModules(ctx context.Context, units []string, roles map[string]struct{}) ([]map[string]string, error) {
	r.logger.Debug("Fetching accessible modules")
	var modules []map[string]string

	for _, unit := range units {
		for role := range roles {
			// Fetch all dynamic actions for each role-unit combination
			response, _, err := r.client.RelationshipApi.GetRelationships(ctx).
				Namespace("modules").
				SubjectSetNamespace("units").
				SubjectSetObject(unit).
				SubjectSetRelation(role).
				Execute()

			if err != nil {
				r.logger.Error("Failed to fetch modules for unit and role", zap.String("unit", unit), zap.String("role", role), zap.Error(err))
				return nil, fmt.Errorf("failed to fetch modules: %w", err)
			}

			// Collect actions and modules
			for _, tuple := range response.RelationTuples {
				modules = append(modules, map[string]string{
					"action": tuple.Relation, // Dynamic action (e.g., "manage", "view", etc.)
					"module": tuple.Object,
				})
			}
		}
	}

	r.logger.Debug("Fetched accessible modules", zap.Int("moduleCount", len(modules)))
	return modules, nil
}
