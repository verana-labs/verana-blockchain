package types

import (
	"fmt"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:             DefaultParams(),
		Permissions:        []Permission{},
		PermissionSessions: []PermissionSession{},
		NextPermissionId:   1,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Check for duplicate permission IDs
	permissionIds := make(map[uint64]bool)
	maxPermId := uint64(0)

	for _, perm := range gs.Permissions {
		// Check if ID exists
		if perm.Id == 0 {
			return fmt.Errorf("permission ID cannot be 0")
		}

		// Check for duplicate IDs
		if _, exists := permissionIds[perm.Id]; exists {
			return fmt.Errorf("duplicate permission ID found: %d", perm.Id)
		}
		permissionIds[perm.Id] = true

		// Track highest permission ID
		if perm.Id > maxPermId {
			maxPermId = perm.Id
		}

		// Validate each permission
		if err := validatePermission(perm, gs.Permissions); err != nil {
			return err
		}

		// Validate timestamps are chronologically consistent
		if err := validatePermissionTimestamps(perm); err != nil {
			return err
		}
	}

	// Check for duplicate session IDs
	sessionIds := make(map[string]bool)
	for _, session := range gs.PermissionSessions {
		// Check if ID exists
		if session.Id == "" {
			return fmt.Errorf("permission session ID cannot be empty")
		}

		// Check for duplicate IDs
		if _, exists := sessionIds[session.Id]; exists {
			return fmt.Errorf("duplicate permission session ID found: %s", session.Id)
		}
		sessionIds[session.Id] = true

		// Validate permission references
		if err := validatePermissionSession(session, permissionIds); err != nil {
			return err
		}
	}

	// Validate next permission ID is greater than max permission ID
	if gs.NextPermissionId <= maxPermId {
		return fmt.Errorf("next_permission_id (%d) must be greater than the maximum permission ID (%d)",
			gs.NextPermissionId, maxPermId)
	}

	return nil
}

// validatePermission validates a single permission
func validatePermission(perm Permission, allPerms []Permission) error {
	// Check required fields
	if perm.Type == 0 {
		return fmt.Errorf("permission type cannot be 0 for permission ID %d", perm.Id)
	}

	if perm.Grantee == "" {
		return fmt.Errorf("grantee cannot be empty for permission ID %d", perm.Id)
	}

	// Validate validator permission reference
	if perm.ValidatorPermId != 0 {
		validatorFound := false

		// Check if validator permission exists
		for _, p := range allPerms {
			if p.Id == perm.ValidatorPermId {
				validatorFound = true
				break
			}
		}

		if !validatorFound {
			return fmt.Errorf("validator permission ID %d not found for permission ID %d",
				perm.ValidatorPermId, perm.Id)
		}
	}

	return nil
}

// validatePermissionTimestamps validates that timestamps are chronologically consistent
func validatePermissionTimestamps(perm Permission) error {
	// Check that modified time exists
	if perm.Modified == nil {
		return fmt.Errorf("modified timestamp is required for permission ID %d", perm.Id)
	}

	// Check that created time exists
	if perm.Created == nil {
		return fmt.Errorf("created timestamp is required for permission ID %d", perm.Id)
	}

	// If effective_from and effective_until both exist, ensure effective_from is before effective_until
	if perm.EffectiveFrom != nil && perm.EffectiveUntil != nil {
		if !perm.EffectiveFrom.Before(*perm.EffectiveUntil) {
			return fmt.Errorf("effective_from must be before effective_until for permission ID %d", perm.Id)
		}
	}

	// If extended time exists, it should be after created time
	if perm.Extended != nil && perm.Created != nil {
		if !perm.Created.Before(*perm.Extended) {
			return fmt.Errorf("extended timestamp must be after created timestamp for permission ID %d", perm.Id)
		}
	}

	return nil
}

// validatePermissionSession validates a single permission session
func validatePermissionSession(session PermissionSession, permissionIds map[uint64]bool) error {
	// Check that agent permission exists
	if session.AgentPermId == 0 {
		return fmt.Errorf("agent permission ID cannot be 0 for session ID %s", session.Id)
	}

	if !permissionIds[session.AgentPermId] {
		return fmt.Errorf("agent permission ID %d not found for session ID %s", session.AgentPermId, session.Id)
	}

	// Validate timestamps
	if session.Created == nil {
		return fmt.Errorf("created timestamp is required for session ID %s", session.Id)
	}

	if session.Modified == nil {
		return fmt.Errorf("modified timestamp is required for session ID %s", session.Id)
	}

	// Validate each authorization entry
	for i, authz := range session.Authz {
		// At least one of executor or beneficiary must be set
		if authz.ExecutorPermId == 0 && authz.BeneficiaryPermId == 0 {
			return fmt.Errorf("at least one of executor_perm_id or beneficiary_perm_id must be set for session ID %s, authz index %d",
				session.Id, i)
		}

		// Check that executor permission exists if set
		if authz.ExecutorPermId != 0 && !permissionIds[authz.ExecutorPermId] {
			return fmt.Errorf("executor permission ID %d not found for session ID %s, authz index %d",
				authz.ExecutorPermId, session.Id, i)
		}

		// Check that beneficiary permission exists if set
		if authz.BeneficiaryPermId != 0 && !permissionIds[authz.BeneficiaryPermId] {
			return fmt.Errorf("beneficiary permission ID %d not found for session ID %s, authz index %d",
				authz.BeneficiaryPermId, session.Id, i)
		}

		// Check that wallet agent permission exists if set
		if authz.WalletAgentPermId != 0 && !permissionIds[authz.WalletAgentPermId] {
			return fmt.Errorf("wallet agent permission ID %d not found for session ID %s, authz index %d",
				authz.WalletAgentPermId, session.Id, i)
		}
	}

	return nil
}
