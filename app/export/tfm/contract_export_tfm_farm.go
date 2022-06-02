package tfm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	terra "github.com/terra-money/core/app"
	"github.com/terra-money/core/app/export/util"
	wasmkeeper "github.com/terra-money/core/x/wasm/keeper"
)

func ExportTfmFarms(app *terra.TerraApp, bl util.Blacklist) (util.SnapshotBalanceAggregateMap, error) {
	ctx := util.PrepCtx(app)
	snapshot := make(util.SnapshotBalanceAggregateMap)
	logger := app.Logger()
	keeper := app.WasmKeeper
	logger.Info("Exporting TFM farms")

	// totalUST := sdk.NewInt(0)

	prefix := util.GeneratePrefix("reward")
	for _, staking := range StakingContracts {
		delegatorAddr, err := sdk.AccAddressFromBech32(staking)
		if err != nil {
			return nil, err
		}
		app.WasmKeeper.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), delegatorAddr, prefix, func(key, value []byte) bool {
			if strings.Contains(string(key), "uusd") {
				stakingHoldings, err := getStakingHoldings(ctx, keeper)
				if err != nil {
					return false
				}
				for lp, staking := range stakingHoldings {
					fmt.Println(lp, staking)
				}
				// Filter out characters from start and end of the key.
				correctedAddress := string(key)[2:46]
				// Remove quotes from the value and convert to an Int.
				balance, ok := sdk.NewIntFromString(strings.Trim(string(value), "\""))
				if ok && !balance.IsZero() {
					snapshot.AppendOrAddBalance(correctedAddress, util.SnapshotBalance{
						Denom:   util.DenomLUNA,
						Balance: balance,
					})
				}
			}

			return false
		})
	}

	return snapshot, nil
}

func getStakingHoldings(ctx context.Context, k wasmkeeper.Keeper) (map[string]stakingHolders, error) {
	holdings := make(map[string]stakingHolders)
	for _, staking := range StakingContracts {
		stakingAddr := util.ToAddress(staking)
		var initMsg stakingInitMsg
		info, err := k.GetContractInfo(sdk.UnwrapSDKContext(ctx), stakingAddr)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(info.InitMsg, &initMsg); err != nil {
			return nil, err
		}
		var lpAddress string
		if initMsg.StakingToken != "" {
			lpAddress = initMsg.StakingToken
		} else if initMsg.LpToken != "" {
			lpAddress = initMsg.LpToken
		} else {
			continue
		}

		prefix := util.GeneratePrefix("reward")
		balances := make(map[string]sdk.Int)
		k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), stakingAddr, prefix, func(key, value []byte) bool {
			var reward struct {
				Amount              sdk.Int `json:"bond_amount"`
				StakingTokenVersion int     `json:"staking_token_version"`
			}
			json.Unmarshal(value, &reward)
			holderAddr := sdk.AccAddress(key)
			// Handle staking contracts that have multiple staking tokens
			if reward.StakingTokenVersion == 0 {
				balances[holderAddr.String()] = reward.Amount
			}
			return false
		})
		holdings[lpAddress] = stakingHolders{
			StakingAddr: staking,
			Holdings:    balances,
		}
	}
	return holdings, nil
}
