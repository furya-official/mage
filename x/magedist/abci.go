package Magedist

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/furya-official/mage/x/Magedist/keeper"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.MintPeriodInflation(ctx)
	if err != nil {
		panic(err)
	}
}
