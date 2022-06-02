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
	wasmtypes "github.com/terra-money/core/x/wasm/types"
)

const (
	DragonSB = "terra14xg04lntty04vgqdcl8jnkclrkg8m532q6w2ga"
	Bitlocus = "terra1fvt7pfnxqrc8fx45nc5weg2hjwfz3ayevt3gha"
	Defiato  = "terra13kq9rqxn0k252qfs4ww4zxzye83ypuye6w2hhm"
	Credefi  = "terra14czgy66f5wf9vxgvvmt5ajrv5urxvhm0359lrk"
)

func ExportTfmFarms(app *terra.TerraApp, bl util.Blacklist) (util.SnapshotBalanceAggregateMap, error) {
	ctx := util.PrepCtx(app)
	q := util.PrepWasmQueryServer(app)
	snapshot := make(util.SnapshotBalanceAggregateMap)
	logger := app.Logger()
	keeper := app.WasmKeeper
	logger.Info("Exporting TFM farms")

	totalUST := sdk.NewInt(0)

	var endowments struct {
		Endowments []struct {
			Address string `json:"address"`
		} `json:"endowments"`
	}
	prefix := util.GeneratePrefix("reward")
	delegatorAddr, err := sdk.AccAddressFromBech32(DragonSB)
	if err != nil {
		return nil, err
	}

	app.WasmKeeper.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), delegatorAddr, prefix, func(key, value []byte) bool {
		if strings.Contains(string(key), "uusd") {
			stakingHoldings, err := getStakingHoldings(ctx, keeper)
			if err != nil {
				return nil, err
			}
			for lp, staking := range stakingHoldings {
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

	if err := util.ContractQuery(ctx, q, &wasmtypes.QueryContractStoreRequest{
		ContractAddress: Endowments,
		QueryMsg:        []byte("{\"endowment_list\":{}}"),
	}, &endowments); err != nil {
		return nil, err
	}

	for _, endowment := range endowments.Endowments {
		var apANCBalance struct {
			LockedCW20 []struct {
				Address string  `json:"address"`
				Amount  sdk.Int `json:"amount"`
			} `json:"locked_cw20"`
			LiquidCW20 []struct {
				Address string  `json:"address"`
				Amount  sdk.Int `json:"amount"`
			} `json:"liquid_cw20"`
		}

		if err := util.ContractQuery(ctx, q, &wasmtypes.QueryContractStoreRequest{
			ContractAddress: APANC,
			QueryMsg:        []byte(fmt.Sprintf("{\"balance\":{\"address\":\"%s\"}}", endowment.Address)),
		}, &apANCBalance); err != nil {
			return nil, err
		}

		USTBalance := sdk.NewInt(0)

		if len(apANCBalance.LiquidCW20) != 0 {
			USTBalance = USTBalance.Add(apANCBalance.LiquidCW20[0].Amount)
		}

		if len(apANCBalance.LockedCW20) != 0 {
			USTBalance = USTBalance.Add(apANCBalance.LockedCW20[0].Amount)
		}

		if USTBalance.IsZero() {
			continue
		}

		totalUST = totalUST.Add(USTBalance)

		// Fetch endowment owner from InitMsg.
		var initMsg struct {
			Owner string `json:"owner_sc"`
		}

		if err := util.ContractInitMsg(ctx, q, &wasmtypes.QueryContractInfoRequest{
			ContractAddress: endowment.Address,
		}, &initMsg); err != nil {
			return nil, err
		}

		snapshot.AppendOrAddBalance(initMsg.Owner, util.SnapshotBalance{
			Denom:   util.DenomAUST,
			Balance: USTBalance,
		})
	}

	logger.Info(fmt.Sprintf("total UST indexed: %d", totalUST.Int64()))

	// These balances are counted using apANC tokens above.
	bl.RegisterAddress(util.DenomUST, DANO)
	bl.RegisterAddress(util.DenomAUST, DANO)
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
