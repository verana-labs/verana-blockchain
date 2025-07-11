package v6

import "time"

// TrustRegistryData holds the trust registry objects
type TrustRegistryData struct {
	TrustRegistries              []TrustRegistry               `json:"trust_registries"`
	GovernanceFrameworkVersions  []GovernanceFrameworkVersion  `json:"governance_framework_versions"`
	GovernanceFrameworkDocuments []GovernanceFrameworkDocument `json:"governance_framework_documents"`
	Counters                     []Counter                     `json:"counters"`
}

type TrustRegistry struct {
	ID            string    `json:"id"`
	DID           string    `json:"did"`
	Controller    string    `json:"controller"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
	Deposit       string    `json:"deposit"`
	AKA           string    `json:"aka"`
	ActiveVersion int32     `json:"active_version"`
	Language      string    `json:"language"`
}

type GovernanceFrameworkVersion struct {
	ID          string    `json:"id"`
	TRID        string    `json:"tr_id"`
	Created     time.Time `json:"created"`
	Version     int32     `json:"version"`
	ActiveSince time.Time `json:"active_since"`
}

type GovernanceFrameworkDocument struct {
	ID        string    `json:"id"`
	GFVID     string    `json:"gfv_id"`
	Created   time.Time `json:"created"`
	Language  string    `json:"language"`
	URL       string    `json:"url"`
	DigestSRI string    `json:"digest_sri"`
}

type Counter struct {
	EntityType string `json:"entity_type"`
	Value      string `json:"value"`
}

// TrustDepositData holds trust deposit accounts
type TrustDepositData struct {
	TrustDeposits []TrustDeposit `json:"trust_deposits"`
}

type TrustDeposit struct {
	Account   string `json:"account"`
	Share     string `json:"share"`
	Amount    string `json:"amount"`
	Claimable string `json:"claimable"`
}

// CredentialSchemaData holds credential schemas
type CredentialSchemaData struct {
	CredentialSchemas []CredentialSchema `json:"credential_schemas"`
	SchemaCounter     string             `json:"schema_counter"`
}

type CredentialSchema struct {
	ID                                      string    `json:"id"`
	TRID                                    string    `json:"tr_id"`
	Created                                 time.Time `json:"created"`
	Modified                                time.Time `json:"modified"`
	Deposit                                 string    `json:"deposit"`
	JSONSchema                              string    `json:"json_schema"`
	IssuerGrantorValidationValidityPeriod   uint32    `json:"issuer_grantor_validation_validity_period"`
	VerifierGrantorValidationValidityPeriod uint32    `json:"verifier_grantor_validation_validity_period"`
	IssuerValidationValidityPeriod          uint32    `json:"issuer_validation_validity_period"`
	VerifierValidationValidityPeriod        uint32    `json:"verifier_validation_validity_period"`
	HolderValidationValidityPeriod          uint32    `json:"holder_validation_validity_period"`
	IssuerPermManagementMode                string    `json:"issuer_perm_management_mode"`
	VerifierPermManagementMode              string    `json:"verifier_perm_management_mode"`
}

// PermissionData holds permissions
type PermissionData struct {
	Permissions      []Permission `json:"permissions"`
	NextPermissionID string       `json:"next_permission_id"`
}

type Permission struct {
	ID                 string    `json:"id"`
	SchemaID           string    `json:"schema_id"`
	Type               string    `json:"type"`
	DID                string    `json:"did"`
	Grantee            string    `json:"grantee"`
	Created            time.Time `json:"created"`
	CreatedBy          string    `json:"created_by"`
	EffectiveFrom      time.Time `json:"effective_from"`
	EffectiveUntil     time.Time `json:"effective_until"`
	Modified           time.Time `json:"modified"`
	ValidationFees     string    `json:"validation_fees"`
	IssuanceFees       string    `json:"issuance_fees"`
	VerificationFees   string    `json:"verification_fees"`
	Deposit            string    `json:"deposit"`
	Country            string    `json:"country"`
	ValidatorPermID    string    `json:"validator_perm_id"`
	VPState            string    `json:"vp_state"`
	VPValidatorDeposit string    `json:"vp_validator_deposit"`
	VPCurrentFees      string    `json:"vp_current_fees"`
	VPCurrentDeposit   string    `json:"vp_current_deposit"`
	VPSummaryDigestSRI string    `json:"vp_summary_digest_sri"`
}

// DIDDirectoryData holds DID directories
type DIDDirectoryData struct {
	DIDDirectories []DIDDirectory `json:"did_directories"`
}

type DIDDirectory struct {
	DID        string    `json:"did"`
	Controller string    `json:"controller"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
	Exp        time.Time `json:"exp"`
	Deposit    string    `json:"deposit"`
}

// GetExtractedData returns all the extracted data from genesis
func GetExtractedData() (TrustRegistryData, TrustDepositData, CredentialSchemaData, PermissionData, DIDDirectoryData) {
	trustRegistryData := TrustRegistryData{
		TrustRegistries: []TrustRegistry{
			{
				ID:            "1",
				DID:           "did:example:184a2fddab1b3d505d477adbf0643446",
				Controller:    "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:       parseTime("2025-06-18T16:27:13.531941769Z"),
				Modified:      parseTime("2025-06-18T16:27:13.531941769Z"),
				Deposit:       "10000000",
				AKA:           "http://example-aka.com",
				ActiveVersion: 1,
				Language:      "en",
			},
			{
				ID:            "2",
				DID:           "did:example:184a3017e0a19d4018d0d621f7b7f9ee",
				Controller:    "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:       parseTime("2025-06-18T16:31:23.515164834Z"),
				Modified:      parseTime("2025-06-18T16:31:23.515164834Z"),
				Deposit:       "10000000",
				AKA:           "http://example-aka.com",
				ActiveVersion: 1,
				Language:      "en",
			},
			{
				ID:            "3",
				DID:           "did:example:184a305eceb9dfb09334f7673d0c2208",
				Controller:    "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:       parseTime("2025-06-18T16:36:26.350087938Z"),
				Modified:      parseTime("2025-06-18T16:36:26.350087938Z"),
				Deposit:       "10000000",
				AKA:           "http://example-aka.com",
				ActiveVersion: 1,
				Language:      "en",
			},
		},
		GovernanceFrameworkVersions: []GovernanceFrameworkVersion{
			{
				ID:          "1",
				TRID:        "1",
				Created:     parseTime("2025-06-18T16:27:13.531941769Z"),
				Version:     1,
				ActiveSince: parseTime("2025-06-18T16:27:13.531941769Z"),
			},
			{
				ID:          "2",
				TRID:        "2",
				Created:     parseTime("2025-06-18T16:31:23.515164834Z"),
				Version:     1,
				ActiveSince: parseTime("2025-06-18T16:31:23.515164834Z"),
			},
			{
				ID:          "3",
				TRID:        "3",
				Created:     parseTime("2025-06-18T16:36:26.350087938Z"),
				Version:     1,
				ActiveSince: parseTime("2025-06-18T16:36:26.350087938Z"),
			},
		},
		GovernanceFrameworkDocuments: []GovernanceFrameworkDocument{
			{
				ID:        "1",
				GFVID:     "1",
				Created:   parseTime("2025-06-18T16:27:13.531941769Z"),
				Language:  "en",
				URL:       "https://example.com/governance-framework.pdf",
				DigestSRI: "sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26",
			},
			{
				ID:        "2",
				GFVID:     "2",
				Created:   parseTime("2025-06-18T16:31:23.515164834Z"),
				Language:  "en",
				URL:       "https://example.com/governance-framework.pdf",
				DigestSRI: "sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26",
			},
			{
				ID:        "3",
				GFVID:     "3",
				Created:   parseTime("2025-06-18T16:36:26.350087938Z"),
				Language:  "en",
				URL:       "https://example.com/governance-framework.pdf",
				DigestSRI: "sha384-MzNNbQTWCSUSi0bbz7dbua+RcENv7C6FvlmYJ1Y+I727HsPOHdzwELMYO9Mz68M26",
			},
		},
		Counters: []Counter{
			{EntityType: "tr", Value: "3"},
			{EntityType: "gfv", Value: "3"},
			{EntityType: "gfd", Value: "3"},
		},
	}

	trustDepositData := TrustDepositData{
		TrustDeposits: []TrustDeposit{
			{
				Account:   "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Share:     "70000000",
				Amount:    "70000000",
				Claimable: "0",
			},
			{
				Account:   "verana1k6exwj6644xy028vxtzxs2fhf9nt8hymeuqkz7",
				Share:     "85000000",
				Amount:    "85000000",
				Claimable: "0",
			},
			{
				Account:   "verana1sxau0xyttphpck7vhlvt8s82ez70nlzw2mhya0",
				Share:     "5000000",
				Amount:    "5000000",
				Claimable: "0",
			},
		},
	}

	credentialSchemaData := CredentialSchemaData{
		CredentialSchemas: []CredentialSchema{
			{
				ID:                                      "1",
				TRID:                                    "1",
				Created:                                 parseTime("2025-06-18T16:27:18.757210803Z"),
				Modified:                                parseTime("2025-06-18T16:27:18.757210803Z"),
				Deposit:                                 "10000000",
				JSONSchema:                              "{\n\t\t\"$schema\": \"https://json-schema.org/draft/2020-12/schema\",\n\t\t\"$id\": \"/vpr/v1/cs/js/1\",\n\t\t\"type\": \"object\",\n\t\t\"$defs\": {},\n\t\t\"properties\": {\n\t\t\t\"name\": {\n\t\t\t\t\"type\": \"string\"\n\t\t\t}\n\t\t},\n\t\t\"required\": [\"name\"],\n\t\t\"additionalProperties\": false\n\t}",
				IssuerGrantorValidationValidityPeriod:   360,
				VerifierGrantorValidationValidityPeriod: 0,
				IssuerValidationValidityPeriod:          360,
				VerifierValidationValidityPeriod:        0,
				HolderValidationValidityPeriod:          0,
				IssuerPermManagementMode:                "GRANTOR_VALIDATION",
				VerifierPermManagementMode:              "OPEN",
			},
			{
				ID:                                      "2",
				TRID:                                    "2",
				Created:                                 parseTime("2025-06-18T16:31:28.865040149Z"),
				Modified:                                parseTime("2025-06-18T16:31:28.865040149Z"),
				Deposit:                                 "10000000",
				JSONSchema:                              "{\n\t\t\"$schema\": \"https://json-schema.org/draft/2020-12/schema\",\n\t\t\"$id\": \"/vpr/v1/cs/js/2\",\n\t\t\"type\": \"object\",\n\t\t\"$defs\": {},\n\t\t\"properties\": {\n\t\t\t\"name\": {\n\t\t\t\t\"type\": \"string\"\n\t\t\t}\n\t\t},\n\t\t\"required\": [\"name\"],\n\t\t\"additionalProperties\": false\n\t}",
				IssuerGrantorValidationValidityPeriod:   360,
				VerifierGrantorValidationValidityPeriod: 0,
				IssuerValidationValidityPeriod:          360,
				VerifierValidationValidityPeriod:        0,
				HolderValidationValidityPeriod:          0,
				IssuerPermManagementMode:                "GRANTOR_VALIDATION",
				VerifierPermManagementMode:              "OPEN",
			},
		},
		SchemaCounter: "2",
	}

	permissionData := PermissionData{
		Permissions: []Permission{
			{
				ID:                 "2",
				SchemaID:           "1",
				Type:               "PERMISSION_TYPE_ECOSYSTEM",
				DID:                "did:example:184a2fddab1b3d505d477adbf0643446",
				Grantee:            "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:            parseTime("2025-06-18T16:27:24.170051412Z"),
				CreatedBy:          "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				EffectiveFrom:      parseTime("2025-06-18T16:27:34.859255Z"),
				EffectiveUntil:     parseTime("2026-06-13T16:27:34.859255Z"),
				Modified:           parseTime("2025-06-18T16:27:24.170051412Z"),
				ValidationFees:     "0",
				IssuanceFees:       "0",
				VerificationFees:   "0",
				Deposit:            "0",
				Country:            "",
				ValidatorPermID:    "0",
				VPState:            "VALIDATION_STATE_UNSPECIFIED",
				VPValidatorDeposit: "0",
				VPCurrentFees:      "0",
				VPCurrentDeposit:   "0",
				VPSummaryDigestSRI: "",
			},
			{
				ID:                 "3",
				SchemaID:           "2",
				Type:               "PERMISSION_TYPE_ECOSYSTEM",
				DID:                "did:example:184a3017e0a19d4018d0d621f7b7f9ee",
				Grantee:            "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:            parseTime("2025-06-18T16:31:34.090851283Z"),
				CreatedBy:          "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				EffectiveFrom:      parseTime("2025-06-18T16:31:44.854252Z"),
				EffectiveUntil:     parseTime("2026-06-13T16:31:44.854252Z"),
				Modified:           parseTime("2025-06-18T16:31:34.090851283Z"),
				ValidationFees:     "0",
				IssuanceFees:       "0",
				VerificationFees:   "0",
				Deposit:            "0",
				Country:            "",
				ValidatorPermID:    "0",
				VPState:            "VALIDATION_STATE_UNSPECIFIED",
				VPValidatorDeposit: "0",
				VPCurrentFees:      "0",
				VPCurrentDeposit:   "0",
				VPSummaryDigestSRI: "",
			},
		},
		NextPermissionID: "3",
	}

	didDirectoryData := DIDDirectoryData{
		DIDDirectories: []DIDDirectory{
			{
				DID:        "did:example:101",
				Controller: "verana1k6exwj6644xy028vxtzxs2fhf9nt8hymeuqkz7",
				Created:    parseTime("2025-06-28T13:09:40.828360752Z"),
				Modified:   parseTime("2025-06-30T17:36:20.241498148Z"),
				Exp:        parseTime("2028-06-28T13:09:40.828360752Z"),
				Deposit:    "15000000",
			},
			{
				DID:        "did:example:110",
				Controller: "verana1k6exwj6644xy028vxtzxs2fhf9nt8hymeuqkz7",
				Created:    parseTime("2025-07-03T20:00:23.521858818Z"),
				Modified:   parseTime("2025-07-03T20:17:30.295110933Z"),
				Exp:        parseTime("2026-07-03T20:00:23.521858818Z"),
				Deposit:    "5000000",
			},
			{
				DID:        "did:example:184a2fdd",
				Controller: "verana1k6exwj6644xy028vxtzxs2fhf9nt8hymeuqkz7",
				Created:    parseTime("2025-06-28T16:29:05.385753071Z"),
				Modified:   parseTime("2025-06-30T20:57:51.134728245Z"),
				Exp:        parseTime("2028-06-28T16:29:05.385753071Z"),
				Deposit:    "15000000",
			},
			{
				DID:        "did:example:184a2fddab1b3d505d4123345565645",
				Controller: "verana1sxau0xyttphpck7vhlvt8s82ez70nlzw2mhya0",
				Created:    parseTime("2025-07-03T20:08:17.270454236Z"),
				Modified:   parseTime("2025-07-03T20:08:17.270454236Z"),
				Exp:        parseTime("2026-07-03T20:08:17.270454236Z"),
				Deposit:    "5000000",
			},
			{
				DID:        "did:example:184a2fddab1b3d505d477adbf064",
				Controller: "verana1k6exwj6644xy028vxtzxs2fhf9nt8hymeuqkz7",
				Created:    parseTime("2025-06-28T02:51:40.197491115Z"),
				Modified:   parseTime("2025-06-28T02:51:40.197491115Z"),
				Exp:        parseTime("2026-06-28T02:51:40.197491115Z"),
				Deposit:    "5000000",
			},
			{
				DID:        "did:example:184a2fddab1b3d505d477adbf0643446",
				Controller: "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:    parseTime("2025-06-18T16:27:29.384488056Z"),
				Modified:   parseTime("2025-06-28T02:31:15.539025824Z"),
				Exp:        parseTime("2026-06-18T16:27:29.384488056Z"),
				Deposit:    "5000000",
			},
			{
				DID:        "did:example:184a2fddab1b3d505d477adbf06453",
				Controller: "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:    parseTime("2025-06-27T14:59:55.698025636Z"),
				Modified:   parseTime("2025-07-02T21:29:28.647039705Z"),
				Exp:        parseTime("2026-06-27T14:59:55.698025636Z"),
				Deposit:    "5000000",
			},
			{
				DID:        "did:example:184a2fddab1b3d505d477adbf06543",
				Controller: "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:    parseTime("2025-06-30T15:38:16.181703297Z"),
				Modified:   parseTime("2025-06-30T15:38:16.181703297Z"),
				Exp:        parseTime("2026-06-30T15:38:16.181703297Z"),
				Deposit:    "5000000",
			},
			{
				DID:        "did:example:184a3017e0a19d4018d0d621f7b7f9ee",
				Controller: "verana12dyk649yce4dvdppehsyraxe6p6jemzg2qwutf",
				Created:    parseTime("2025-06-18T16:31:39.510648750Z"),
				Modified:   parseTime("2025-06-30T21:31:33.749398814Z"),
				Exp:        parseTime("2026-06-18T16:31:39.510648750Z"),
				Deposit:    "5000000",
			},
		},
	}

	return trustRegistryData, trustDepositData, credentialSchemaData, permissionData, didDirectoryData
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}
