package tfm

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	terra "github.com/terra-money/core/app"
	"github.com/terra-money/core/app/export/util"
	"github.com/terra-money/core/x/wasm/types"
	wasmtypes "github.com/terra-money/core/x/wasm/types"
)

func ExportTfmFarms(app *terra.TerraApp, bl util.Blacklist) (util.SnapshotBalanceAggregateMap, error) {
	ctx := util.PrepCtx(app)
	snapshot := make(util.SnapshotBalanceAggregateMap)
	logger := app.Logger()
	qs := util.PrepWasmQueryServer(app)
	logger.Info("Exporting TFM farms")

	for _, farm := range FarmContracts {
		farmAddr, err := sdk.AccAddressFromBech32(farm)
		if err != nil {
			return nil, err
		}
		prefix := util.GeneratePrefix("reward")
		app.WasmKeeper.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), farmAddr, prefix, func(key, value []byte) bool {

			userAddr := string(key)
			stakerUstBal := getStakerUstBal(ctx, qs, userAddr, farmAddr)

			snapshot.AppendOrAddBalance(userAddr, util.SnapshotBalance{
				Denom:   util.DenomUST,
				Balance: stakerUstBal,
			})

			return false
		})
	}
	return snapshot, nil
}

func getStakerUstBal(ctx context.Context, q types.QueryServer, userAddr string, farmAddr sdk.AccAddress) sdk.Int {
	USTBalance := sdk.NewInt(0)
	var info stakerInfo
	if err := util.ContractQuery(ctx, q, &wasmtypes.QueryContractStoreRequest{
		ContractAddress: string(farmAddr),
		QueryMsg:        []byte(fmt.Sprintf("{\"staker_info\": {\"owner\": \"%s\"}}}", userAddr))}, &info); err != nil {
		fmt.Println(err)
	}
	if !info.bondAmount.IsZero() {
		USTBalance = USTBalance.Add(info.bondAmount.Quo(sdk.NewInt(2)))
	}
	if !info.pendingRewards.IsZero() {
		USTBalance = USTBalance.Add(info.pendingRewards.Quo(sdk.NewInt(2)))
	}
	return USTBalance
}
