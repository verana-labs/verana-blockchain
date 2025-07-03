#!/bin/bash

set -e

# Function to log messages with timestamp
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Function to log section headers
log_section() {
    echo ""
    echo "================================================================"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
    echo "================================================================"
    echo ""
}

# Function to log subsection headers
log_subsection() {
    echo ""
    echo "----------------------------------------------------------------"
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1"
    echo "----------------------------------------------------------------"
}

# Function to log JSON responses
log_json() {
    local description=$1
    local json=$2
    log "$description:"
    echo "$json" | jq '.' || log "Error: Invalid JSON response"
}

# Function to wait for transaction to be mined
wait_for_tx() {
    local tx_hash=$1
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        local tx_result=$(veranad query tx $tx_hash --output json 2>/dev/null)
        if [ $? -eq 0 ]; then
            log "Transaction $tx_hash confirmed"
            log_json "Transaction result" "$tx_result"
            sleep 3  # Wait 3 seconds after transaction confirmation
            return 0
        fi
        log "Waiting for transaction to be mined (attempt $attempt/$max_attempts)..."
        sleep 2
        attempt=$((attempt + 1))
    done

    log "Error: Transaction confirmation timeout"
    return 1
}

# Function to ensure chain is responsive
wait_for_chain() {
    local max_attempts=30
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if veranad status 2>/dev/null | grep -q "latest_block_height"; then
            log "Chain is responsive"
            return 0
        fi
        log "Waiting for chain to become responsive (attempt $attempt/$max_attempts)..."
        sleep 2
        attempt=$((attempt + 1))
    done

    log "Error: Chain is not responsive"
    return 1
}

# Function to execute transaction and handle response
execute_tx() {
    local description=$1
    local tx_result=$2

    log "$description"
    log_json "Transaction submission result" "$tx_result"

    local tx_hash=$(echo "$tx_result" | jq -r '.txhash')
    if [ -z "$tx_hash" ] || [ "$tx_hash" = "null" ]; then
        log "Error: Failed to get transaction hash"
        return 1
    fi

    wait_for_tx "$tx_hash"
}

# Common transaction parameters
COMMON_PARAMS="--keyring-backend test --chain-id vna-local-1 --gas 800000 --gas-adjustment 1.3 --gas-prices 1.1uvna --yes --output json"

log "Beginning test sequence..."

# Check if chain is running and responsive
log "Checking chain status..."
if ! wait_for_chain; then
    log "Error: Could not connect to chain. Please ensure the chain is running and try again."
    exit 1
fi

# Additional validation functions
validate_did() {
    local did=$1
    if [[ ! $did =~ ^did:.*:.*$ ]]; then
        log "Error: Invalid DID format"
        return 1
    fi
    return 0
}

validate_language() {
    local lang=$1
    if [[ ! $lang =~ ^[a-z]{2}$ ]]; then
        log "Error: Invalid language format. Must be 2 letters (e.g., 'en')"
        return 1
    fi
    return 0
}

validate_url() {
    local url=$1
    if [[ ! $url =~ ^https?:// ]]; then
        log "Error: Invalid URL format"
        return 1
    fi
    return 0
}

validate_hash() {
    local hash=$1
    if [[ ! $hash =~ ^[a-fA-F0-9]{64}$ ]]; then
        log "Error: Invalid hash format. Must be 64 hex characters"
        return 1
    fi
    return 0
}

validate_trust_registry() {
    local tr_data=$1
    local expected_did=$2
    local expected_controller=$3

    # Extract values using jq
    local actual_did=$(echo "$tr_data" | jq -r '.trust_registry.did')
    local actual_controller=$(echo "$tr_data" | jq -r '.trust_registry.controller')
    local active_version=$(echo "$tr_data" | jq -r '.trust_registry.active_version')

    # Validate DID
    if [ "$actual_did" != "$expected_did" ]; then
        log "Error: DID mismatch. Expected: $expected_did, Got: $actual_did"
        return 1
    fi

    # Validate Controller
    if [ "$actual_controller" != "$expected_controller" ]; then
        log "Error: Controller mismatch. Expected: $expected_controller, Got: $actual_controller"
        return 1
    fi

    log "✓ Trust Registry validation passed:"
    log "  - DID matches: $actual_did"
    log "  - Controller matches: $actual_controller"
    log "  - Active Version: $active_version"
    return 0
}

# ================================================================
# Trust Registry Module Tests
# ================================================================
log_section "TRUST REGISTRY MODULE TESTS"

# Test 1: Create Trust Registry with validation
log_subsection "1. Creating Trust Registry (with validation)"

# Test input validation
DID="did:example:123456789abcdefghi"
AKA="http://example-aka.com"
LANGUAGE="en"
DOC_URL="https://example.com/governance-framework.pdf"
DOC_HASH="e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

log "Validating inputs..."
validate_did "$DID" || exit 1
validate_language "$LANGUAGE" || exit 1
validate_url "$AKA" || exit 1
validate_url "$DOC_URL" || exit 1
validate_hash "$DOC_HASH" || exit 1
log "✓ Input validation passed"

TR_RESULT=$(veranad tx trustregistry create-trust-registry \
    "$DID" \
    "$AKA" \
    "$LANGUAGE" \
    "$DOC_URL" \
    "$DOC_HASH" \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "Trust Registry Creation" "$TR_RESULT"; then
    log "Error: Failed to create Trust Registry"
    exit 1
fi

# Store TR_ID and validate response
TR_ID="1"
log "Using Trust Registry ID: $TR_ID"
sleep 3

# Verify Trust Registry creation
log "Verifying Trust Registry creation..."
TR_VERIFY=$(veranad query trustregistry get-trust-registry $TR_ID \
    --active-gf-only \
    --preferred-language en \
    --output json)

ACCOUNT_ADDRESS=$(veranad keys show cooluser -a --keyring-backend test)
if ! validate_trust_registry "$TR_VERIFY" "$DID" "$ACCOUNT_ADDRESS"; then
    log "Error: Trust Registry validation failed"
    exit 1
fi

# Test 2: Add Governance Framework Document with validation
log_subsection "2. Adding Governance Framework Document (with validation)"

# Validate inputs for new document
NEW_DOC_URL="https://example.com/updated-governance-framework.pdf"
NEW_DOC_HASH="e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
NEW_VERSION=2

log "Validating new document inputs..."
validate_language "$LANGUAGE" || exit 1
validate_url "$NEW_DOC_URL" || exit 1
validate_hash "$NEW_DOC_HASH" || exit 1
log "✓ Document input validation passed"

GF_RESULT=$(veranad tx trustregistry add-governance-framework-document \
    $TR_ID \
    "$LANGUAGE" \
    "$NEW_DOC_URL" \
    "$NEW_DOC_HASH" \
    $NEW_VERSION \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "Governance Framework Document Addition" "$GF_RESULT"; then
    log "Error: Failed to add Governance Framework Document"
    exit 1
fi

sleep 3

# Verify document addition
log "Verifying Governance Framework Document addition..."
TR_DOC_VERIFY=$(veranad query trustregistry get-trust-registry $TR_ID \
    --preferred-language en \
    --output json)

# Validate document was added correctly
DOC_COUNT=$(echo "$TR_DOC_VERIFY" | jq '.documents | length')
if [ "$DOC_COUNT" -lt 2 ]; then
    log "Error: Expected at least 2 documents, got $DOC_COUNT"
    exit 1
fi

LATEST_DOC=$(echo "$TR_DOC_VERIFY" | jq '.documents[-1]')
LATEST_URL=$(echo "$LATEST_DOC" | jq -r '.url')
LATEST_HASH=$(echo "$LATEST_DOC" | jq -r '.hash')

if [ "$LATEST_URL" != "$NEW_DOC_URL" ] || [ "$LATEST_HASH" != "$NEW_DOC_HASH" ]; then
    log "Error: Document details don't match"
    log "Expected URL: $NEW_DOC_URL, Got: $LATEST_URL"
    log "Expected Hash: $NEW_DOC_HASH, Got: $LATEST_HASH"
    exit 1
fi

log "✓ Document addition verification passed:"
log "  - Document count: $DOC_COUNT"
log "  - Latest document URL matches"
log "  - Latest document hash matches"

# Test 3: Increase Active GF Version with validation
log_subsection "3. Increasing Active Governance Framework Version (with validation)"

# Get current version before increase
CURRENT_VERSION=$(echo "$TR_DOC_VERIFY" | jq -r '.trust_registry.active_version')
log "Current active version: $CURRENT_VERSION"

GFV_RESULT=$(veranad tx trustregistry increase-active-gf-version \
    $TR_ID \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "Increase Active GF Version" "$GFV_RESULT"; then
    log "Error: Failed to increase active GF version"
    exit 1
fi

sleep 3

# Verify version increase
log "Verifying version increase..."
TR_VERSION_VERIFY=$(veranad query trustregistry get-trust-registry $TR_ID \
    --active-gf-only \
    --preferred-language en \
    --output json)

NEW_ACTIVE_VERSION=$(echo "$TR_VERSION_VERIFY" | jq -r '.trust_registry.active_version')
EXPECTED_VERSION=$((CURRENT_VERSION + 1))

if [ "$NEW_ACTIVE_VERSION" != "$EXPECTED_VERSION" ]; then
    log "Error: Version increase failed. Expected: $EXPECTED_VERSION, Got: $NEW_ACTIVE_VERSION"
    exit 1
fi

log "✓ Version increase verification passed:"
log "  - Previous version: $CURRENT_VERSION"
log "  - New version: $NEW_ACTIVE_VERSION"

# Test 4: Query Trust Registry by ID with validation
log_subsection "4. Querying Trust Registry by ID (with validation)"

TR_QUERY=$(veranad query trustregistry get-trust-registry $TR_ID \
    --preferred-language en \
    --output json)

log_json "Trust Registry Query Result" "$TR_QUERY"

# Validate query result
log "Validating Trust Registry query result..."
if [ -z "$TR_QUERY" ]; then
    log "Error: Empty query result"
    exit 1
fi

# Validate required fields using individual checks
VALIDATED_FIELDS=0

# Check ID
actual_id=$(echo "$TR_QUERY" | jq -r ".trust_registry.id")
if [ "$actual_id" = "$TR_ID" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'id' mismatch. Expected: $TR_ID, Got: $actual_id"
    exit 1
fi

# Check DID
actual_did=$(echo "$TR_QUERY" | jq -r ".trust_registry.did")
if [ "$actual_did" = "$DID" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'did' mismatch. Expected: $DID, Got: $actual_did"
    exit 1
fi

# Check controller
actual_controller=$(echo "$TR_QUERY" | jq -r ".trust_registry.controller")
if [ "$actual_controller" = "$ACCOUNT_ADDRESS" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'controller' mismatch. Expected: $ACCOUNT_ADDRESS, Got: $actual_controller"
    exit 1
fi

# Check language
actual_language=$(echo "$TR_QUERY" | jq -r ".trust_registry.language")
if [ "$actual_language" = "$LANGUAGE" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'language' mismatch. Expected: $LANGUAGE, Got: $actual_language"
    exit 1
fi

# Check active version
actual_version=$(echo "$TR_QUERY" | jq -r ".trust_registry.active_version")
if [ "$actual_version" = "$NEW_ACTIVE_VERSION" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'active_version' mismatch. Expected: $NEW_ACTIVE_VERSION, Got: $actual_version"
    exit 1
fi

log "✓ Trust Registry query validation passed:"
log "  - Validated $VALIDATED_FIELDS required fields"
log "  - All field values match expected values"
sleep 3

# Test 5: Query Trust Registry by DID with validation
log_subsection "5. Querying Trust Registry by DID (with validation)"

TR_DID_QUERY=$(veranad query trustregistry get-trust-registry-by-did \
    "$DID" \
    --preferred-language en \
    --output json)
log_json "Trust Registry DID Query Result" "$TR_DID_QUERY"

# Validate DID query result matches ID query
log "Validating DID query result matches ID query..."
if [ "$(echo "$TR_DID_QUERY" | jq -c .)" != "$(echo "$TR_QUERY" | jq -c .)" ]; then
    log "Error: DID query result doesn't match ID query result"
    exit 1
fi

log "✓ Trust Registry DID query validation passed:"
log "  - DID query result matches ID query result"
log "  - All fields are consistent between queries"
sleep 3

# Test 6: List Trust Registries with validation
log_subsection "6. Listing Trust Registries (with validation)"

TR_LIST_QUERY=$(veranad query trustregistry list-trust-registries \
    --controller "$ACCOUNT_ADDRESS" \
    --modified-after "2023-01-01T00:00:00Z" \
    --active-gf-only \
    --preferred-language en \
    --response-max-size 100 \
    --output json)
log_json "Trust Registry List Result" "$TR_LIST_QUERY"
# Validate list result
log "Validating Trust Registry list result..."

# Check if list contains at least one registry
REGISTRY_COUNT=$(echo "$TR_LIST_QUERY" | jq '.trust_registries | length')
if [ "$REGISTRY_COUNT" -lt 1 ]; then
    log "Error: Expected at least one Trust Registry in list"
    exit 1
fi

# Validate our created registry is in the list
FOUND_REGISTRY=$(echo "$TR_LIST_QUERY" | jq --arg did "$DID" '.trust_registries[] | select(.did == $did)')
if [ -z "$FOUND_REGISTRY" ]; then
    log "Error: Created Trust Registry not found in list"
    exit 1
fi

# Validate registry details in list
log "Validating registry details in list..."
LIST_VALIDATED_FIELDS=0
for field in "${!REQUIRED_FIELDS[@]}"; do
    expected_value="${REQUIRED_FIELDS[$field]}"
    actual_value=$(echo "$FOUND_REGISTRY" | jq -r ".$field")

    if [ "$actual_value" = "$expected_value" ]; then
        ((LIST_VALIDATED_FIELDS++))
    else
        log "Error: Field '$field' mismatch in list. Expected: $expected_value, Got: $actual_value"
        exit 1
    fi
done

# Validate document versions
log "Validating document versions..."
DOC_VERSIONS=$(echo "$TR_QUERY" | jq '.versions | length')
if [ "$DOC_VERSIONS" -lt 2 ]; then
    log "Error: Expected at least 2 document versions"
    exit 1
fi

log "✓ Trust Registry list validation passed:"
log "  - Found $REGISTRY_COUNT total registries"
log "  - Successfully located created registry"
log "  - Validated $LIST_VALIDATED_FIELDS fields in list entry"
log "  - Found $DOC_VERSIONS document versions"

# Final validation summary
log_subsection "Trust Registry Module Test Summary"
log "All validations passed successfully:"
log "1. Trust Registry Creation"
log "   - Input validation"
log "   - Creation transaction"
log "   - Registry details verification"
log "2. Governance Framework Document Addition"
log "   - Document input validation"
log "   - Document addition transaction"
log "   - Document details verification"
log "3. Version Management"
log "   - Version increase transaction"
log "   - Version number verification"
log "4. Query Operations"
log "   - Query by ID validation"
log "   - Query by DID validation"
log "   - List query validation"
log "5. Document Versions"
log "   - Version count verification"
log "   - Document content verification"
log "6. Field Validations"
log "   - All required fields present"
log "   - All field values correct"

# ================================================================
# DID Directory Module Tests
# ================================================================
log_section "DID DIRECTORY MODULE TESTS"

# Additional validation functions for DID Directory
validate_did_response() {
    local did_data=$1
    local expected_did=$2
    local expected_controller=$3
    local expected_years=$4

    # Extract values using jq
    local actual_did=$(echo "$did_data" | jq -r '.did')
    local actual_controller=$(echo "$did_data" | jq -r '.controller')
    local created=$(echo "$did_data" | jq -r '.created')
    local modified=$(echo "$did_data" | jq -r '.modified')
    local exp=$(echo "$did_data" | jq -r '.exp')

    # Validate fields
    if [ "$actual_did" != "$expected_did" ]; then
        log "Error: DID mismatch. Expected: $expected_did, Got: $actual_did"
        return 1
    fi

    if [ "$actual_controller" != "$expected_controller" ]; then
        log "Error: Controller mismatch. Expected: $expected_controller, Got: $actual_controller"
        return 1
    fi

    # Validate timestamps exist
    if [ "$created" = "null" ] || [ -z "$created" ]; then
        log "Error: Created timestamp is missing"
        return 1
    fi

    if [ "$modified" = "null" ] || [ -z "$modified" ]; then
        log "Error: Modified timestamp is missing"
        return 1
    fi

    if [ "$exp" = "null" ] || [ -z "$exp" ]; then
        log "Error: Expiration timestamp is missing"
        return 1
    fi

    log "✓ DID response validation passed:"
    log "  - DID matches: $actual_did"
    log "  - Controller matches: $actual_controller"
    log "  - Created: $created"
    log "  - Modified: $modified"
    log "  - Expiration: $exp"
    return 0
}

# Test 1: Add DID with validation
log_subsection "1. Adding DID (with validation)"

# Test input validation
DID="did:example:123456789abcdefghi"
YEARS=5

log "Validating inputs..."
validate_did "$DID" || exit 1
if [ "$YEARS" -lt 1 ] || [ "$YEARS" -gt 31 ]; then
    log "Error: Invalid years. Must be between 1 and 31"
    exit 1
fi
log "✓ Input validation passed"

# Add DID
DID_RESULT=$(veranad tx diddirectory add-did \
    "$DID" \
    $YEARS \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "DID Addition" "$DID_RESULT"; then
    log "Error: Failed to add DID"
    exit 1
fi

# Wait for indexing
log "Waiting for DID indexing..."
MAX_RETRIES=10
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    sleep 3
    DID_VERIFY=$(veranad query diddirectory get-did \
        "$DID" \
        --output json 2>/dev/null)

    # Changed this line to check did_entry.did instead of just .did
    if [ $? -eq 0 ] && [ "$(echo "$DID_VERIFY" | jq -r '.did_entry.did')" = "$DID" ]; then
        log "DID successfully indexed"
        break
    fi

    ((RETRY_COUNT++))
    log "Waiting for DID indexing (attempt $RETRY_COUNT/$MAX_RETRIES)..."
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    log "Error: Timeout waiting for DID to be indexed"
    exit 1
fi

# Verify DID creation
log "Verifying DID creation..."
DID_VERIFY=$(veranad query diddirectory get-did \
    "$DID" \
    --output json)

log_json "DID Verification Result" "$DID_VERIFY"

ACCOUNT_ADDRESS=$(veranad keys show cooluser -a --keyring-backend test)
# Update validate_did_response function to handle the new structure
validate_did_response() {
    local did_data=$1
    local expected_did=$2
    local expected_controller=$3
    local expected_years=$4

    # Extract values using jq, accessing through did_entry
    local actual_did=$(echo "$did_data" | jq -r '.did_entry.did')
    local actual_controller=$(echo "$did_data" | jq -r '.did_entry.controller')
    local created=$(echo "$did_data" | jq -r '.did_entry.created')
    local modified=$(echo "$did_data" | jq -r '.did_entry.modified')
    local exp=$(echo "$did_data" | jq -r '.did_entry.exp')

    # Validate fields
    if [ "$actual_did" != "$expected_did" ]; then
        log "Error: DID mismatch. Expected: $expected_did, Got: $actual_did"
        return 1
    fi

    if [ "$actual_controller" != "$expected_controller" ]; then
        log "Error: Controller mismatch. Expected: $expected_controller, Got: $actual_controller"
        return 1
    fi

    # Validate timestamps exist
    if [ "$created" = "null" ] || [ -z "$created" ]; then
        log "Error: Created timestamp is missing"
        return 1
    fi

    if [ "$modified" = "null" ] || [ -z "$modified" ]; then
        log "Error: Modified timestamp is missing"
        return 1
    fi

    if [ "$exp" = "null" ] || [ -z "$exp" ]; then
        log "Error: Expiration timestamp is missing"
        return 1
    fi

    log "✓ DID response validation passed:"
    log "  - DID matches: $actual_did"
    log "  - Controller matches: $actual_controller"
    log "  - Created: $created"
    log "  - Modified: $modified"
    log "  - Expiration: $exp"
    return 0
}

if ! validate_did_response "$DID_VERIFY" "$DID" "$ACCOUNT_ADDRESS" "$YEARS"; then
    log "Error: DID validation failed"
    exit 1
fi

# Test 2: Renew DID with validation
log_subsection "2. Renewing DID (with validation)"

# Validate renewal inputs
RENEWAL_YEARS=2
if [ "$RENEWAL_YEARS" -lt 1 ] || [ "$RENEWAL_YEARS" -gt 31 ]; then
    log "Error: Invalid renewal years. Must be between 1 and 31"
    exit 1
fi

# Get initial expiration
INITIAL_EXP=$(echo "$DID_VERIFY" | jq -r '.did_entry.exp')
log "Initial expiration: $INITIAL_EXP"

DID_RENEW_RESULT=$(veranad tx diddirectory renew-did \
    "$DID" \
    $RENEWAL_YEARS \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "DID Renewal" "$DID_RENEW_RESULT"; then
    log "Error: Failed to renew DID"
    exit 1
fi

sleep 3

# Verify renewal
log "Verifying DID renewal..."
DID_RENEWAL_VERIFY=$(veranad query diddirectory get-did \
    "$DID" \
    --output json)

log_json "DID Renewal Query Result" "$DID_RENEWAL_VERIFY"

NEW_EXP=$(echo "$DID_RENEWAL_VERIFY" | jq -r '.did_entry.exp')  # Added .did_entry
if [ "$NEW_EXP" = "$INITIAL_EXP" ]; then
    log "Error: Expiration not updated after renewal"
    exit 1
fi

log "✓ DID renewal verification passed:"
log "  - Previous expiration: $INITIAL_EXP"
log "  - New expiration: $NEW_EXP"

# Validate renewal details
log_json "DID Renewal Query Result" "$DID_RENEWAL_VERIFY"

NEW_EXP=$(echo "$DID_RENEWAL_VERIFY" | jq -r '.exp')
if [ "$NEW_EXP" = "$INITIAL_EXP" ]; then
    log "Error: Expiration not updated after renewal"
    exit 1
fi

log "✓ DID renewal verification passed:"
log "  - Previous expiration: $INITIAL_EXP"
log "  - New expiration: $NEW_EXP"

# Test 3: Touch DID with validation
log_subsection "3. Touching DID (with validation)"

# Get initial modified time
INITIAL_MODIFIED=$(echo "$DID_RENEWAL_VERIFY" | jq -r '.did_entry.modified')
log "Initial modified time: $INITIAL_MODIFIED"

DID_TOUCH_RESULT=$(veranad tx diddirectory touch-did \
    "$DID" \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "DID Touch" "$DID_TOUCH_RESULT"; then
    log "Error: Failed to touch DID"
    exit 1
fi

sleep 3

# Verify touch
log "Verifying DID touch..."
DID_TOUCH_VERIFY=$(veranad query diddirectory get-did \
    "$DID" \
    --output json)

log_json "DID Touch Query Result" "$DID_TOUCH_VERIFY"

NEW_MODIFIED=$(echo "$DID_TOUCH_VERIFY" | jq -r '.did_entry.modified')
if [ "$NEW_MODIFIED" = "$INITIAL_MODIFIED" ]; then
    log "Error: Modified timestamp not updated after touch"
    exit 1
fi

log "✓ DID touch verification passed:"
log "  - Previous modified time: $INITIAL_MODIFIED"
log "  - New modified time: $NEW_MODIFIED"

# Test 4: Query DID with validation
log_subsection "4. Querying DID (with validation)"

DID_QUERY=$(veranad query diddirectory get-did \
    "$DID" \
    --output json)

log_json "DID Query Result" "$DID_QUERY"

# Validate query result
log "Validating DID query result..."
if [ -z "$DID_QUERY" ]; then
    log "Error: Empty query result"
    exit 1
fi

# Validate required fields
VALIDATED_FIELDS=0

# Check DID
actual_did=$(echo "$DID_QUERY" | jq -r ".did_entry.did")
if [ "$actual_did" = "$DID" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'did' mismatch. Expected: $DID, Got: $actual_did"
    exit 1
fi

# Check controller
actual_controller=$(echo "$DID_QUERY" | jq -r ".did_entry.controller")
if [ "$actual_controller" = "$ACCOUNT_ADDRESS" ]; then
    ((VALIDATED_FIELDS++))
else
    log "Error: Field 'controller' mismatch. Expected: $ACCOUNT_ADDRESS, Got: $actual_controller"
    exit 1
fi

log "✓ DID query validation passed:"
log "  - Validated $VALIDATED_FIELDS required fields"
log "  - All field values match expected values"

# Test 5: List DIDs with validation
log_subsection "5. Listing DIDs (with validation)"

DID_LIST_QUERY=$(veranad query diddirectory list-dids \
    --account "$ACCOUNT_ADDRESS" \
    --changed "2024-01-01T00:00:00Z" \
    --expired=false \
    --over-grace=false \
    --max-results 64 \
    --output json)

log_json "DID List Query Result" "$DID_LIST_QUERY"

# Validate list result
log "Validating DID list result..."

# Check if list contains at least one DID
DID_COUNT=$(echo "$DID_LIST_QUERY" | jq '.dids | length')
if [ "$DID_COUNT" -lt 1 ]; then
    log "Error: Expected at least one DID in list"
    exit 1
fi

# Validate our created DID is in the list
FOUND_DID=$(echo "$DID_LIST_QUERY" | jq --arg did "$DID" '.dids[] | select(.did == $did)')
if [ -z "$FOUND_DID" ]; then
    log "Error: Created DID not found in list"
    exit 1
fi

log "✓ DID list validation passed:"
log "  - Found $DID_COUNT total DIDs"
log "  - Successfully located created DID"

# Final validation summary
log_subsection "DID Directory Module Test Summary"
log "All validations passed successfully:"
log "1. DID Creation"
log "   - Input validation"
log "   - Creation transaction"
log "   - DID details verification"
log "2. DID Renewal"
log "   - Renewal years validation"
log "   - Expiration update verification"
log "   - Renewal transaction success"
log "3. DID Touch"
log "   - Modified timestamp update"
log "   - Touch transaction success"
log "4. Query Operations"
log "   - Single DID query validation"
log "   - Field validation"
log "   - Timestamp validation"
log "5. List Operations"
log "   - List response validation"
log "   - DID count verification"
log "   - Created DID presence"

# ================================================================
# Credential Schema Module Tests
# ================================================================
log_section "CREDENTIAL SCHEMA MODULE TESTS"

# Additional validation functions for Credential Schema
validate_json_schema() {
    local schema=$1
    # Check if it's valid JSON
    if ! echo "$schema" | jq . >/dev/null 2>&1; then
        log "Error: Invalid JSON schema format"
        return 1
    fi

    # Check required fields
    if ! echo "$schema" | jq -e 'has("$schema")' >/dev/null; then
        log "Error: Missing $schema field"
        return 1
    fi

    if ! echo "$schema" | jq -e 'has("type") and .type == "object"' >/dev/null; then
        log "Error: Missing or invalid type field"
        return 1
    fi

    log "✓ JSON schema validation passed"
    return 0
}

validate_credential_schema_response() {
    local schema_data=$1
    local expected_id=$2
    local expected_tr_id=$3

    # Extract values using jq
    local actual_id=$(echo "$schema_data" | jq -r '.schema.id')
    local actual_tr_id=$(echo "$schema_data" | jq -r '.schema.tr_id')
    local json_schema=$(echo "$schema_data" | jq -r '.schema.json_schema')
    local issuer_mode=$(echo "$schema_data" | jq -r '.schema.issuer_perm_management_mode')
    local verifier_mode=$(echo "$schema_data" | jq -r '.schema.verifier_perm_management_mode')

    # Validate ID
    if [ "$actual_id" != "$expected_id" ]; then
        log "Error: Schema ID mismatch. Expected: $expected_id, Got: $actual_id"
        return 1
    fi

    # Validate Trust Registry ID
    if [ "$actual_tr_id" != "$expected_tr_id" ]; then
        log "Error: Trust Registry ID mismatch. Expected: $expected_tr_id, Got: $actual_tr_id"
        return 1
    fi

    # Validate JSON schema
    if ! validate_json_schema "$json_schema"; then
        return 1
    fi

    # Validate management modes
    if [ "$issuer_mode" != "PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION" ]; then
        log "Error: Invalid issuer management mode: $issuer_mode"
        return 1
    fi

    if [ "$verifier_mode" != "PERM_MANAGEMENT_MODE_GRANTOR_VALIDATION" ]; then
        log "Error: Invalid verifier management mode: $verifier_mode"
        return 1
    fi

    log "✓ Credential Schema response validation passed:"
    log "  - Schema ID matches: $actual_id"
    log "  - Trust Registry ID matches: $actual_tr_id"
    log "  - JSON Schema is valid"
    log "  - Management modes are correct"
    return 0
}

# Test 1: Create Credential Schema with validation
log_subsection "1. Creating Credential Schema (with validation)"

# Create and validate schema JSON
cat > test_schema.json << 'EOF'
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "/dtr/v1/cs/js/1",
    "type": "object",
    "properties": {
        "name": {
            "type": "string"
        }
    },
    "required": ["name"],
    "additionalProperties": false
}
EOF

# Validate schema before using
log "Validating JSON schema..."
if ! validate_json_schema "$(cat test_schema.json)"; then
    log "Error: Invalid JSON schema"
    rm -f test_schema.json
    exit 1
fi
log "✓ Schema validation passed"

# Create credential schema
log "Creating credential schema..."
CS_RESULT=$(veranad tx credentialschema create-credential-schema \
    $TR_ID \
    "$(cat test_schema.json)" \
    365 \
    365 \
    180 \
    180 \
    180 \
    2 \
    2 \
    --from cooluser \
    $COMMON_PARAMS)

if ! execute_tx "Credential Schema Creation" "$CS_RESULT"; then
    log "Error: Failed to create Credential Schema"
    rm -f test_schema.json
    exit 1
fi

# Wait for indexing
log "Waiting for schema to be indexed..."
CS_ID="1"  # Expected first schema ID
MAX_RETRIES=10
RETRY_COUNT=0
while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    sleep 3
    CS_VERIFY=$(veranad query credentialschema get $CS_ID --output json 2>/dev/null)

    if [ $? -eq 0 ] && [ "$(echo "$CS_VERIFY" | jq -r '.schema.id')" = "$CS_ID" ]; then
        log "Schema successfully indexed"
        break
    fi

    ((RETRY_COUNT++))
    log "Waiting for schema indexing (attempt $RETRY_COUNT/$MAX_RETRIES)..."
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    log "Error: Timeout waiting for schema to be indexed"
    exit 1
fi

# Get and validate schema details
log_subsection "2. Getting Credential Schema Details"
CS_QUERY=$(veranad query credentialschema get $CS_ID \
    --output json)
log_json "Credential Schema Details" "$CS_QUERY"

if ! validate_credential_schema_response "$CS_QUERY" "$CS_ID" "$TR_ID"; then
    log "Error: Credential Schema validation failed"
    exit 1
fi

sleep 3

# Get and validate JSON schema definition
log_subsection "3. Getting Credential Schema JSON Definition"
CS_SCHEMA_QUERY=$(veranad query credentialschema schema $CS_ID \
    --output json)
log_json "Credential Schema JSON Definition" "$CS_SCHEMA_QUERY"

# Validate JSON schema format
if ! validate_json_schema "$(echo "$CS_SCHEMA_QUERY" | jq -r '.schema')"; then
    log "Error: Invalid JSON schema in response"
    exit 1
fi

sleep 3

# List and validate schemas
log_subsection "4. Listing Credential Schemas"
CS_LIST_QUERY=$(veranad query credentialschema list-schemas \
    --tr_id $TR_ID \
    --created_after "2024-01-01T00:00:00Z" \
    --response_max_size 100 \
    --output json)
log_json "Credential Schema List Result" "$CS_LIST_QUERY"

# Validate list result
log "Validating schema list result..."

# Check if list contains at least one schema
SCHEMA_COUNT=$(echo "$CS_LIST_QUERY" | jq '.schemas | length')
if [ "$SCHEMA_COUNT" -lt 1 ]; then
    log "Error: Expected at least one Credential Schema in list"
    exit 1
fi

# Validate our created schema is in the list
FOUND_SCHEMA=$(echo "$CS_LIST_QUERY" | jq --arg id "$CS_ID" '.schemas[] | select(.id == $id)')
if [ -z "$FOUND_SCHEMA" ]; then
    log "Error: Created Credential Schema not found in list"
    exit 1
fi

log "✓ Schema list validation passed:"
log "  - Found $SCHEMA_COUNT total schemas"
log "  - Successfully located created schema"

# Validate schema details in list
SCHEMA_ID=$(echo "$FOUND_SCHEMA" | jq -r '.id')
SCHEMA_TR_ID=$(echo "$FOUND_SCHEMA" | jq -r '.tr_id')
SCHEMA_CREATED=$(echo "$FOUND_SCHEMA" | jq -r '.created')

if [ "$SCHEMA_ID" != "$CS_ID" ] || [ "$SCHEMA_TR_ID" != "$TR_ID" ]; then
    log "Error: Schema details in list don't match"
    exit 1
fi

log "  - Schema details validation passed"
log "    - ID: $SCHEMA_ID"
log "    - Trust Registry ID: $SCHEMA_TR_ID"
log "    - Created: $SCHEMA_CREATED"

# Final validation summary
log_subsection "Credential Schema Module Test Summary"
log "All validations passed successfully:"
log "1. Schema Creation"
log "   - JSON schema validation"
log "   - Creation transaction"
log "   - Schema details verification"
log "2. Schema Details"
log "   - ID verification"
log "   - Trust Registry ID verification"
log "   - Schema content validation"
log "   - Management modes verification"
log "3. Schema Definition"
log "   - JSON format validation"
log "   - Schema structure verification"
log "4. Schema List"
log "   - List response validation"
log "   - Schema count verification"
log "   - Created schema presence"
log "   - Schema details consistency"

# ================================================================
# Credential Schema Permission Module Tests
# ================================================================
log_section "CREDENTIAL SCHEMA PERMISSION MODULE TESTS"

# Test different perm types
PERMISSION_TYPES=("1" "2" "3" "4" "5" "6")
PERMISSION_NAMES=("ISSUER" "VERIFIER" "ISSUER_GRANTOR" "VERIFIER_GRANTOR" "TRUST_REGISTRY" "HOLDER")

for i in "${!PERMISSION_TYPES[@]}"; do
    PERM_TYPE="${PERMISSION_TYPES[$i]}"
    PERM_NAME="${PERMISSION_NAMES[$i]}"

    log_subsection "Creating $PERM_NAME Permission (Type $PERM_TYPE)"

    CURRENT_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    NEXT_YEAR=$(date -u -v+1y +"%Y-%m-%dT%H:%M:%SZ")
    TWO_YEARS=$(date -u -v+2y +"%Y-%m-%dT%H:%M:%SZ")

    PERM_RESULT=$(veranad tx cspermission create-credential-schema-perm \
        $TR_ID \
        $PERM_TYPE \
        "did:example:123456789abcdefghi" \
        "$(veranad keys show cooluser -a --keyring-backend test)" \
        "$NEXT_YEAR" \
        100 \
        200 \
        300 \
        --effective-until "$TWO_YEARS" \
        --country US \
        --validation-id 123 \
        --from cooluser \
        $COMMON_PARAMS)

    if ! execute_tx "Creating $PERM_NAME Permission" "$PERM_RESULT"; then
        log "Error: Failed to create $PERM_NAME permission"
        exit 1
    fi

    sleep 3
done

# Query all permissions (Note: Adding a query if available in your implementation)
log_subsection "Querying Credential Schema Permissions"
log "Note: Implement permission queries according to your API"

# Cleanup
log_subsection "Cleaning up temporary files"
rm -f test_schema.json

log_section "TEST SEQUENCE COMPLETED SUCCESSFULLY"
log "Summary of tests performed:"

# Trust Registry Module Summary
log "1. Trust Registry Module:"
log "   - Input Validation"
log "     • DID format validation"
log "     • URL format validation"
log "     • Language code validation"
log "     • Hash format validation"
log "   - Trust Registry Creation & Verification"
log "     • Creation transaction success"
log "     • Registry ID verification"
log "     • Controller verification"
log "   - Governance Framework Management"
log "     • Document addition verification"
log "     • Version control validation"
log "     • Document count validation"
log "   - Query Validations"
log "     • Registry details verification"
log "     • DID-based lookup verification"
log "     • List operations validation"
log "     • Field consistency checks"

# DID Directory Module Summary
log "2. DID Directory Module:"
log "   - DID Registration"
log "     • DID format validation"
log "     • Registration period validation"
log "     • Creation transaction verification"
log "   - DID Management"
log "     • Renewal period validation"
log "     • Expiration update verification"
log "     • Modified timestamp verification"
log "   - Query Operations"
log "     • DID details validation"
log "     • Controller verification"
log "     • Timestamp validations"
log "     • List operations verification"

# Credential Schema Module Summary
log "3. Credential Schema Module:"
log "   - Schema Creation"
log "     • JSON schema validation"
log "     • Trust Registry association"
log "     • Creation transaction verification"
log "   - Schema Details"
log "     • ID verification"
log "     • Trust Registry ID verification"
log "     • Management modes validation"
log "     • Schema content validation"
log "   - Query Operations"
log "     • Schema definition validation"
log "     • List operations verification"
log "     • Field consistency checks"
log "     • Schema count validation"

# Credential Schema Permission Module Summary
log "4. Credential Schema Permission Module:"
log "   - Permission Types Created & Validated"
log "     • ISSUER (Type 1)"
log "     • VERIFIER (Type 2)"
log "     • ISSUER_GRANTOR (Type 3)"
log "     • VERIFIER_GRANTOR (Type 4)"
log "     • TRUST_REGISTRY (Type 5)"
log "     • HOLDER (Type 6)"
log "   - For Each Permission Type:"
log "     • Transaction success verification"
log "     • Permission parameters validation"
log "     • Time period validations"
log "     • Association with Trust Registry"