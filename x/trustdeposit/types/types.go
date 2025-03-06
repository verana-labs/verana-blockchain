package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic implements sdk.Msg
func (msg *MsgReclaimTrustDepositInterests) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address (%s)", err)
	}
	return nil
}

func (msg *MsgReclaimTrustDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return fmt.Errorf("invalid creator address (%s)", err)
	}
	if msg.Claimed == 0 {
		return fmt.Errorf("claimed amount must be greater than 0")
	}
	return nil
}
