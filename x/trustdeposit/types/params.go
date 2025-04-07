package types

import (
	"cosmossdk.io/math"
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DefaultTrustDepositReclaimBurnRate = "0.6" // 60%
	DefaultTrustDepositShareValue      = "1.0" // Initial value: 1
	DefaultTrustDepositRate            = "0.2" // 20%
	DefaultWalletUserAgentRewardRate   = "0.2" // 20%
	DefaultUserAgentRewardRate         = "0.2" // 20%
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	trustDepositReclaimBurnRate math.LegacyDec,
	trustDepositShareValue math.LegacyDec,
	trustDepositRate math.LegacyDec,
	walletUserAgentRewardRate math.LegacyDec,
	userAgentRewardRate math.LegacyDec,
) Params {
	return Params{
		TrustDepositReclaimBurnRate: trustDepositReclaimBurnRate,
		TrustDepositShareValue:      trustDepositShareValue,
		TrustDepositRate:            trustDepositRate,
		WalletUserAgentRewardRate:   walletUserAgentRewardRate,
		UserAgentRewardRate:         userAgentRewardRate,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	TrustDepositReclaimBurnRate, _ := math.LegacyNewDecFromStr(DefaultTrustDepositReclaimBurnRate)
	TrustDepositShareValue, _ := math.LegacyNewDecFromStr(DefaultTrustDepositShareValue)
	TrustDepositRate, _ := math.LegacyNewDecFromStr(DefaultTrustDepositRate)
	WalletUserAgentRewardRate, _ := math.LegacyNewDecFromStr(DefaultWalletUserAgentRewardRate)
	UserAgentRewardRate, _ := math.LegacyNewDecFromStr(DefaultUserAgentRewardRate)

	return NewParams(
		TrustDepositReclaimBurnRate,
		TrustDepositShareValue,
		TrustDepositRate,
		WalletUserAgentRewardRate,
		UserAgentRewardRate,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			[]byte("TrustDepositReclaimBurnRate"),
			&p.TrustDepositReclaimBurnRate,
			validatePercentage,
		),
		paramtypes.NewParamSetPair(
			[]byte("TrustDepositShareValue"),
			&p.TrustDepositShareValue,
			validatePositiveUint64,
		),
		paramtypes.NewParamSetPair(
			[]byte("TrustDepositRate"),
			&p.TrustDepositRate,
			validatePercentage,
		),
		paramtypes.NewParamSetPair(
			[]byte("WalletUserAgentRewardRate"),
			&p.WalletUserAgentRewardRate,
			validatePercentage,
		),
		paramtypes.NewParamSetPair(
			[]byte("UserAgentRewardRate"),
			&p.UserAgentRewardRate,
			validatePercentage,
		),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validatePercentage(p.TrustDepositReclaimBurnRate); err != nil {
		return err
	}
	if err := validatePositiveUint64(p.TrustDepositShareValue); err != nil {
		return err
	}
	if err := validatePercentage(p.TrustDepositRate); err != nil {
		return err
	}
	if err := validatePercentage(p.WalletUserAgentRewardRate); err != nil {
		return err
	}
	if err := validatePercentage(p.UserAgentRewardRate); err != nil {
		return err
	}
	return nil
}

func validatePercentage(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v > 100 {
		return fmt.Errorf("percentage value cannot be greater than 100: %d", v)
	}

	return nil
}

func validatePositiveUint64(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("value must be positive: %d", v)
	}

	return nil
}
