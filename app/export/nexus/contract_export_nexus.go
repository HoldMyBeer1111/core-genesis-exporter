package nexus

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	terra "github.com/terra-money/core/app"
	"github.com/terra-money/core/app/export/util"
	wasmtypes "github.com/terra-money/core/x/wasm/types"
)

var (
	AddressBATOMVault           = "terra1lh3h7l5vsul2pxlevraucwev42ar6kyx33u4c8"
	AddressNATOM                = "terra1jtdc6zpf95tvh9peuaxwp3v0yqszcnwl8j5ade"
	AddressWASAVAXVault         = "terra1hn9rzu66s422rl9kg0a7j2yxdjef0szkqvy7ws"
	AddressNAVAX                = "terra13k62n0285wj8ug0ngcgpf7dgnkzqeu279tz636"
	AddressCNLUNA               = "terra1u553zk43jd4rwzc53qrdrq4jc2p8rextyq09dj"
	AddressNLUNA                = "terra10f2mt82kjnkxqj2gepgwl637u2w4ue2z5nhz5j"
	AddressCNLUNAAutoCompounder = "terra1au4h305fn4w3zpka2ql59e0t70jnqzu4mj2txx"
	AddressAnchorOverseer       = "terra1tmnqgvg567ypvsvk6rwsga3srp7e3lg6u0elp8"
)

// func ExportNexus(app *terra.TerraApp, fromLP util.SnapshotBalanceAggregateMap, bl util.Blacklist) (util.SnapshotBalanceAggregateMap, error) {
// 	ctx := util.PrepCtx(app)
// 	qs := util.PrepWasmQueryServer(app)

// 	keeper := app.WasmKeeper

// 	// get all cnLuna holders, unwrap to nLuna
// 	var cnLunaHolderMap = make(util.BalanceMap)
// 	if err := util.GetCW20AccountsAndBalances(ctx, keeper, AddressCNLUNA, cnLunaHolderMap); err != nil {
// 		return nil, fmt.Errorf("failed to fetch cnLUNA holders: %v", err)
// 	}

// 	// get nLUNA balance of cnLuna Autocompounder
// 	nLunaInAutocompounder, err := util.GetCW20Balance(ctx, qs, AddressNLUNA, AddressCNLUNAAutoCompounder)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch nLUNA balance in autocompounder: %v", err)
// 	}

// 	// get total cnLUNA supply
// 	var cnLunaSupply struct {
// 		TotalSupply sdk.Int `json:"total_supply"`
// 	}
// 	if err := util.ContractQuery(ctx, qs, &wasmtypes.QueryContractStoreRequest{
// 		ContractAddress: AddressCNLUNA,
// 		QueryMsg:        []byte("{\"token_info\":{}}"),
// 	}, &cnLunaSupply); err != nil {
// 		return nil, fmt.Errorf("failed to fetch cnLUNA supply")
// 	}

// 	// calc nLUNA <> cnLUNA ratio
// 	ratio := sdk.NewDecFromInt(cnLunaSupply.TotalSupply).QuoInt(nLunaInAutocompounder)

// 	// iterate over cnLuna holders, convert it to nLUNA
// 	var nLunaHolderMap = make(util.BalanceMap)
// 	for userAddr, cnLunaHolding := range cnLunaHolderMap {
// 		nLunaHolderMap[userAddr] = ratio.MulInt(cnLunaHolding).TruncateInt()
// 	}

// 	// iterate over nLuna holders, add it to nLunaHolderMap
// 	// (bar pairs from dexes)
// 	var nLunaHolderMapFlat = make(util.BalanceMap)
// 	if err := util.GetCW20AccountsAndBalances(ctx, keeper, AddressNLUNA, nLunaHolderMapFlat); err != nil {
// 		return nil, fmt.Errorf("failed to fetch nLUNA holder")
// 	}

// 	// merge holder maps + nLUNA holdings from LP
// 	blacklist := bl.GetAddressesByDenomMap(util.DenomNLUNA)
// 	nLunaHoldingsFromLP := fromLP.PickDenomIntoBalanceMap(util.DenomNLUNA)
// 	mergednLunaHolderMap := util.MergeMaps(nLunaHolderMap, nLunaHolderMapFlat, nLunaHoldingsFromLP)

// 	nAssetTobAssetRatio, err := getnAssetTobAssetRatio(ctx, qs)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// iterate over merged nLUNA holder map, apply nLUNA -> bLUNA ratio
// 	var finalBalance = make(util.SnapshotBalanceAggregateMap)
// 	for userAddr, nLunaHolding := range mergednLunaHolderMap {

// 		// bar blacklisted addresses (pairs, ...)
// 		if _, exists := blacklist[userAddr]; exists {
// 			continue
// 		}

// 		bLunaAmount := nAssetTobAssetRatio.MulInt(nLunaHolding)

// 		// there can't be more than 1 holding -- this is fine
// 		finalBalance[userAddr] = []util.SnapshotBalance{
// 			{
// 				Denom:   util.DenomBLUNA,
// 				Balance: bLunaAmount.TruncateInt(),
// 			},
// 		}
// 	}

// 	return finalBalance, nil
// }

// func getnAssetTobAssetRatio(ctx context.Context, qs wasmtypes.QueryServer) (sdk.Dec, error) {
// 	// nLUNA -> bLUNA ratio
// 	// get bLUNA held in collateral
// 	var collaterals struct {
// 		Collaterals [][2]string `json:"collaterals"`
// 	}
// 	if err := util.ContractQuery(ctx, qs, &wasmtypes.QueryContractStoreRequest{
// 		ContractAddress: AddressAnchorOverseer,
// 		QueryMsg:        []byte(fmt.Sprintf("{\"collaterals\":{\"borrower\":\"%s\"}}", AddressBLUNAVault)),
// 	}, &collaterals); err != nil {
// 		return sdk.Dec{}, fmt.Errorf("failed to fetch Nexus bLUNA vault collateral: %v", err)
// 	}

// 	bLUNAProvision, _ := sdk.NewIntFromString(collaterals.Collaterals[0][1])

// 	// calc nAsset->bAsset ratio
// 	var nLunaSupply struct {
// 		TotalSupply sdk.Int `json:"total_supply"`
// 	}
// 	if err := util.ContractQuery(ctx, qs, &wasmtypes.QueryContractStoreRequest{
// 		ContractAddress: AddressNLUNA,
// 		QueryMsg:        []byte("{\"token_info\":{}}"),
// 	}, &nLunaSupply); err != nil {
// 		return sdk.Dec{}, fmt.Errorf("failed to fetch nLUNA total supply: %v", err)
// 	}
// 	nAssetTobAssetRatio := sdk.NewDecFromInt(bLUNAProvision).QuoInt(nLunaSupply.TotalSupply)
// 	return nAssetTobAssetRatio, nil
// }

// func ResolveToBLuna(app *terra.TerraApp, snapshot util.SnapshotBalanceAggregateMap, bl util.Blacklist) error {
// 	ctx := util.PrepCtx(app)
// 	qs := util.PrepWasmQueryServer(app)

// 	nAssetTobAssetRatio, err := getnAssetTobAssetRatio(ctx, qs)
// 	if err != nil {
// 		return err
// 	}

// 	for _, sbs := range snapshot {
// 		for i, sb := range sbs {
// 			if sb.Denom == util.DenomNLUNA {
// 				sbs[i] = util.SnapshotBalance{
// 					Denom:   util.DenomBLUNA,
// 					Balance: nAssetTobAssetRatio.MulInt(sb.Balance).TruncateInt(),
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

func get_b_to_n_asset_ratio(ctx context.Context, qs wasmtypes.QueryServer, nasset_token string, basset_vault_addr string) (sdk.Dec, error) {
	// nAsset -> bAsset ratio
	// get bAsset held in collateral
	var collaterals struct {
		Collaterals [][2]string `json:"collaterals"`
	}
	if err := util.ContractQuery(ctx, qs, &wasmtypes.QueryContractStoreRequest{
		ContractAddress: AddressAnchorOverseer,
		QueryMsg:        []byte(fmt.Sprintf("{\"collaterals\":{\"borrower\":\"%s\"}}", basset_vault_addr)),
	}, &collaterals); err != nil {
		return sdk.Dec{}, fmt.Errorf("failed to fetch Nexus bLUNA vault collateral: %v", err)
	}

	bAssetProvision, _ := sdk.NewIntFromString(collaterals.Collaterals[0][1])

	// calc nAsset->bAsset ratio
	var nAssetSupply struct {
		TotalSupply sdk.Int `json:"total_supply"`
	}
	if err := util.ContractQuery(ctx, qs, &wasmtypes.QueryContractStoreRequest{
		ContractAddress: nasset_token,
		QueryMsg:        []byte("{\"token_info\":{}}"),
	}, &nAssetSupply); err != nil {
		return sdk.Dec{}, fmt.Errorf("failed to fetch nAsset total supply: %v", err)
	}
	nAssetTobAssetRatio := sdk.NewDecFromInt(bAssetProvision).QuoInt(nAssetSupply.TotalSupply)
	return nAssetTobAssetRatio, nil
}

func FindLiquidation(app *terra.TerraApp, height int64) (bool, error) {
	ctx := util.PrepCtxByHeight(app, height)
	qs := util.PrepWasmQueryServer(app)

	atom_found := false
	avax_found := false

	if ratio, err := get_b_to_n_asset_ratio(ctx, qs, AddressNATOM, AddressBATOMVault); err == nil && !ratio.Equal(sdk.OneDec()) {
		fmt.Printf("----> N to B ratio: %d\n", ratio)
		fmt.Printf("-----> bAtom vault liquidation height: %d\n", height)
		atom_found = true
	}
	if ratio, err := get_b_to_n_asset_ratio(ctx, qs, AddressNAVAX, AddressWASAVAXVault); err == nil && !ratio.Equal(sdk.OneDec()) {
		fmt.Printf("----> N to B ratio: %d\n", ratio)
		fmt.Printf("-----> wasAvax vault liquidation height: %d\n", height)
		avax_found = true
	}

	if atom_found && avax_found {
		return true, nil
	}
	return false, nil
}
