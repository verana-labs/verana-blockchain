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
