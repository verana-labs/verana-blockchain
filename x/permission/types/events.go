package types

const (
	EventTypeCreateRootPermission            = "create_root_permission"
	AttributeKeyRootPermissionID             = "root_permission_id"
	AttributeKeySchemaID                     = "schema_id"
	AttributeKeyTimestamp                    = "timestamp"
	EventTypeStartPermissionVP               = "start_permission_vp"
	AttributeKeyPermissionID                 = "permission_id"
	AttributeKeyCreator                      = "creator"
	AttributeKeyFees                         = "fees"
	AttributeKeyDeposit                      = "deposit"
	EventTypeCreateOrUpdatePermissionSession = "create_update_csps"
	AttributeKeySessionID                    = "session_id"
	AttributeKeyExecutorPermID               = "executor_perm_id"
	AttributeKeyBeneficiaryPermID            = "beneficiary_perm_id"
	AttributeKeyUserAgentDID                 = "user_agent_did"
	AttributeKeyTotalFees                    = "total_fees"
	AttributeKeyAgentPermID                  = "agent_perm_id"
)
