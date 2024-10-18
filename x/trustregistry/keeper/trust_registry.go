package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/verana-labs/verana-blockchain/x/trustregistry/types"
	"net/url"
	"regexp"
	"time"
)

func (ms msgServer) validateCreateTrustRegistryParams(ctx sdk.Context, msg *types.MsgCreateTrustRegistry) error {
	if msg.Did == "" || msg.Language == "" || msg.DocUrl == "" || msg.DocHash == "" {
		return errors.New("missing mandatory parameter")
	}

	if !isValidDID(msg.Did) {
		return errors.New("invalid DID syntax")
	}

	// Check if a trust registry with this DID already exists
	_, err := ms.TrustRegistry.Get(ctx, msg.Did)
	if err == nil {
		return errors.New("trust registry with this DID already exists")
	}

	// Validate AKA URI if present
	if msg.Aka != "" && !isValidURI(msg.Aka) {
		return errors.New("invalid AKA URI")
	}

	// Validate language tag (rfc1766)
	if !isValidLanguageTag(msg.Language) {
		return errors.New("invalid language tag (must conform to rfc1766)")
	}

	if !isValidURL(msg.DocUrl) {
		return errors.New("invalid document URL")
	}

	if !isValidHash(msg.DocHash) {
		return errors.New("invalid document hash")
	}

	return nil
}

func isValidLanguageTag(lang string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9]{1,8}(-[a-zA-Z0-9]{1,8})*$`, lang)
	return match && len(lang) <= 17
}

// TODO: Remove comment before testing on real environment
func (ms msgServer) checkSufficientFees(ctx sdk.Context, creator string) error {
	//creatorAddr, err := sdk.AccAddressFromBech32(creator)
	//if err != nil {
	//	return fmt.Errorf("invalid creator address: %w", err)
	//}
	//
	//// Use the first denomination from minimum gas prices
	//minGasPrices := ctx.MinGasPrices()
	//if len(minGasPrices) == 0 {
	//	return fmt.Errorf("no minimum gas price set")
	//}
	//feeDenom := minGasPrices[0].Denom
	//
	//// Estimate fee (using a fixed gas amount for simplicity)
	//estimatedGas := uint64(200000)
	//estimatedFee := minGasPrices.AmountOf(feeDenom).MulInt64(int64(estimatedGas))
	//
	//// Check if the account has enough balance
	//balance := ms.k.bankKeeper.GetBalance(ctx, creatorAddr, feeDenom)
	//if balance.Amount.LT(estimatedFee.TruncateInt()) {
	//	return fmt.Errorf("insufficient funds to cover estimated transaction fees")
	//}

	return nil
}

func (ms msgServer) createTrustRegistryEntries(_ sdk.Context, msg *types.MsgCreateTrustRegistry, now time.Time) (types.TrustRegistry, types.GovernanceFrameworkVersion, types.GovernanceFrameworkDocument, error) {
	tr := types.TrustRegistry{
		Did:           msg.Did,
		Controller:    msg.Creator,
		Created:       now,
		Modified:      now,
		Deposit:       0,
		Aka:           msg.Aka,
		ActiveVersion: 1,
		Language:      msg.Language,
	}

	gfv := types.GovernanceFrameworkVersion{
		Id:          uuid.New().String(),
		TrDid:       msg.Did,
		Created:     now,
		Version:     1,
		ActiveSince: now,
	}

	gfd := types.GovernanceFrameworkDocument{
		Id:       uuid.New().String(),
		GfvId:    gfv.Id,
		Created:  now,
		Language: msg.Language,
		Url:      msg.DocUrl,
		Hash:     msg.DocHash,
	}

	return tr, gfv, gfd, nil
}

func (ms msgServer) persistEntries(ctx sdk.Context, tr types.TrustRegistry, gfv types.GovernanceFrameworkVersion, gfd types.GovernanceFrameworkDocument) error {
	if err := ms.TrustRegistry.Set(ctx, tr.Did, tr); err != nil {
		return fmt.Errorf("failed to persist TrustRegistry: %w", err)
	}

	if err := ms.GFVersion.Set(ctx, gfv.Id, gfv); err != nil {
		return fmt.Errorf("failed to persist GovernanceFrameworkVersion: %w", err)
	}

	if err := ms.GFDocument.Set(ctx, gfd.Id, gfd); err != nil {
		return fmt.Errorf("failed to persist GovernanceFrameworkDocument: %w", err)
	}

	return nil
}

// Helper functions

func isValidDID(did string) bool {
	// Basic DID validation regex
	// This is a simplified version and may need to be expanded based on specific DID method requirements
	didRegex := regexp.MustCompile(`^did:[a-zA-Z0-9]+:[a-zA-Z0-9._-]+$`)
	return didRegex.MatchString(did)
}

func isValidURI(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil
}

func isValidURL(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}

func isValidHash(hash string) bool {
	// This is a basic check for a SHA-256 hash (64 hexadecimal characters)
	// Adjust this based on your specific hash requirements
	hashRegex := regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	return hashRegex.MatchString(hash)
}
