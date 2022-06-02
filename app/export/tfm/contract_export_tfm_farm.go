package tfm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	terra "github.com/terra-money/core/app"
	"github.com/terra-money/core/app/export/util"
	wasmtypes "github.com/terra-money/core/x/wasm/types"
)

func ExportTfmFarms(app *terra.TerraApp, bl util.Blacklist) (util.SnapshotBalanceAggregateMap, error) {
	ctx := util.PrepCtx(app)
	snapshot := make(util.SnapshotBalanceAggregateMap)
	logger := app.Logger()
	qs := util.PrepWasmQueryServer(app)
	logger.Info("Exporting TFM farms")

	prefix := util.GeneratePrefix("reward")
	for _, staking := range StakingContracts {
		delegatorAddr, err := sdk.AccAddressFromBech32(staking)
		if err != nil {
			return nil, err
		}

		app.WasmKeeper.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), delegatorAddr, prefix, func(key, value []byte) bool {
			var info stakerInfo
			USTBalance := sdk.NewInt(0)
			if err := util.ContractQuery(ctx, qs, &wasmtypes.QueryContractStoreRequest{
				ContractAddress: staking,
				QueryMsg:        []byte(fmt.Sprintf("{\"staker_info\": {\"owner\": \"%s\"}}}", key))}, &info); err != nil {
				return false
			}

			if !info.bondAmount.IsZero() {
				USTBalance = USTBalance.Add(info.bondAmount.Quo(sdk.NewInt(2)))
			}
			if !info.pendingRewards.IsZero() {
				USTBalance = USTBalance.Add(info.pendingRewards.Quo(sdk.NewInt(2)))
			}

			snapshot.AppendOrAddBalance(string(key), util.SnapshotBalance{
				Denom:   util.DenomUST,
				Balance: USTBalance,
			})

			return false
		})
	}
	return snapshot, nil
}
