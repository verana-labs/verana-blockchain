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
	AttributeKeyAgentPermID                  = "agent_perm_id"
	AttributeKeyIssuerPermID                 = "issuer_perm_id"
	AttributeKeyVerifierPermID               = "verifier_perm_id"
	AttributeKeyWalletAgentPermID            = "wallet_agent_perm_id"
)
