package keeper_test

import (
	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
)

func TestCreateAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)
	rewardDuration := app.AllianceKeeper.RewardDelayTime(ctx)

	// WHEN
	createErr := app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:        "",
		Description:  "",
		Denom:        "uluna",
		RewardWeight: sdk.OneDec(),
		TakeRate:     sdk.OneDec(),
	})
	alliancesRes, alliancesErr := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, createErr)
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         sdk.NewDec(1),
				TakeRate:             sdk.NewDec(1),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardStartTime:      ctx.BlockTime().Add(rewardDuration),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
				LastRewardChangeTime: ctx.BlockTime().Add(rewardDuration),
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	})
}

func TestCreateAllianceFailWithDuplicatedDenom(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			types.NewAllianceAsset("uluna", sdk.NewDec(1), sdk.NewDec(0), startTime),
		},
	})

	// WHEN
	createErr := app.AllianceKeeper.CreateAlliance(ctx, &types.MsgCreateAllianceProposal{
		Title:        "",
		Description:  "",
		Denom:        "uluna",
		RewardWeight: sdk.OneDec(),
		TakeRate:     sdk.OneDec(),
	})

	// THEN
	require.Error(t, createErr)
}

func TestUpdateAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         sdk.NewDec(2),
				TakeRate:             sdk.OneDec(),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN
	updateErr := app.AllianceKeeper.UpdateAlliance(ctx, &types.MsgUpdateAllianceProposal{
		Title:                "",
		Description:          "",
		Denom:                "uluna",
		RewardWeight:         sdk.NewDec(6),
		TakeRate:             sdk.NewDec(7),
		RewardChangeInterval: 0,
		RewardChangeRate:     sdk.ZeroDec(),
	})
	alliancesRes, alliancesErr := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, updateErr)
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: []types.AllianceAsset{
			{
				Denom:                "uluna",
				RewardWeight:         sdk.NewDec(6),
				TakeRate:             sdk.NewDec(7),
				TotalTokens:          sdk.ZeroInt(),
				TotalValidatorShares: sdk.NewDec(0),
				RewardChangeRate:     sdk.NewDec(0),
				RewardChangeInterval: 0,
			},
		},
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   1,
		},
	})
}

func TestDeleteAlliance(t *testing.T) {
	// GIVEN
	app, ctx := createTestContext(t)
	startTime := time.Now()
	ctx.WithBlockTime(startTime).WithBlockHeight(1)
	app.AllianceKeeper.InitGenesis(ctx, &types.GenesisState{
		Params: types.DefaultParams(),
		Assets: []types.AllianceAsset{
			{
				Denom:        "uluna",
				RewardWeight: sdk.NewDec(2),
				TakeRate:     sdk.OneDec(),
				TotalTokens:  sdk.ZeroInt(),
			},
		},
	})
	queryServer := keeper.NewQueryServerImpl(app.AllianceKeeper)

	// WHEN
	deleteErr := app.AllianceKeeper.DeleteAlliance(ctx, &types.MsgDeleteAllianceProposal{
		Denom: "uluna",
	})
	alliancesRes, alliancesErr := queryServer.Alliances(ctx, &types.QueryAlliancesRequest{})

	// THEN
	require.Nil(t, deleteErr)
	require.Nil(t, alliancesErr)
	require.Equal(t, alliancesRes, &types.QueryAlliancesResponse{
		Alliances: nil,
		Pagination: &query.PageResponse{
			NextKey: nil,
			Total:   0,
		},
	})
}
