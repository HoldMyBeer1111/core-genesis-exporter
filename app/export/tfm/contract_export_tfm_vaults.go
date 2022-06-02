package tfm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	terra "github.com/terra-money/core/app"
	"github.com/terra-money/core/app/export/util"
	wasmtypes "github.com/terra-money/core/x/wasm/types"
)

const (
	DragonSB = "terra14xg04lntty04vgqdcl8jnkclrkg8m532q6w2ga"
	Bitlocus = "terra1fvt7pfnxqrc8fx45nc5weg2hjwfz3ayevt3gha"
	Defiato  = "terra13kq9rqxn0k252qfs4ww4zxzye83ypuye6w2hhm"
	Credefi  = "terra14czgy66f5wf9vxgvvmt5ajrv5urxvhm0359lrk"
)

func ExportVaults(app *terra.TerraApp, bl util.Blacklist) (util.SnapshotBalanceAggregateMap, error) {
	ctx := util.PrepCtx(app)
	q := util.PrepWasmQueryServer(app)
	snapshot := make(util.SnapshotBalanceAggregateMap)
	logger := app.Logger()
	logger.Info("Exporting Angel Protocol endowments")

	totalaUST := sdk.NewInt(0)

	var endowments struct {
		Endowments []struct {
			Address string `json:"address"`
		} `json:"endowments"`
	}

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

		aUSTBalance := sdk.NewInt(0)

		if len(apANCBalance.LiquidCW20) != 0 {
			aUSTBalance = aUSTBalance.Add(apANCBalance.LiquidCW20[0].Amount)
		}

		if len(apANCBalance.LockedCW20) != 0 {
			aUSTBalance = aUSTBalance.Add(apANCBalance.LockedCW20[0].Amount)
		}

		if aUSTBalance.IsZero() {
			continue
		}

		totalaUST = totalaUST.Add(aUSTBalance)

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
			Balance: aUSTBalance,
		})
	}

	logger.Info(fmt.Sprintf("total aUST indexed: %d", totalaUST.Int64()))

	// These balances are counted using apANC tokens above.
	bl.RegisterAddress(util.DenomUST, DANO)
	bl.RegisterAddress(util.DenomAUST, DANO)

	return snapshot, nil
}
