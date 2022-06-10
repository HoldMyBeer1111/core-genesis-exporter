package nexus

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

var (
	AddressPSIToken = "terra12897djskt9rge8dtmm86w654g7kzckkd698608"
	AddressNexusGov = "terra1xrk6v2tfjrhjz2dsfecj40ps7ayanjx970gy0j"

	AddressAstroUSTPair      = "terra1v5ct2tuhfqd0tf8z0wwengh4fg77kaczgf6gtx"
	AddressAstroUSTLPToken   = "terra1cspx9menzglmn7xt3tcn8v8lg6gu9r50d7lnve"
	AddressAstroNLUNAPair    = "terra10lv5wz84kpwxys7jeqkfxx299drs3vnw0lj8mz"
	AddressAstroNLUNALPToken = "terra1t53c8p0zwvj5xx7sxh3qtse0fq5765dltjrg33"
	AddressAstroNETHPair     = "terra18hjdxnnkv8ewqlaqj3zpn0vsfpzdt3d0y2ufdz"
	AddressAstroNETHLPToken  = "terra1pjfqacx7k6dg63v2h5q96zjg7w5q25093wnkjc"
	AddressAstroNAVAXPair    = "terra10usmg35qsa92fagh49np7phmhhr4ryhyl27749"
	AddressAstroNAVAXLPToken = "terra1p3zj8tkzufw9szmm97taj7x6kkd0cy7k2mpdws"
	AddressAstroNATOMPair    = "terra1spcf4486jjn8678hstwqzeeudu98yp7pyyltnl"
	AddressAstroNATOMLPToken = "terra1pyavxxun3vuakqq0wyqft69l3zjns0q76wut7z"

	AddressTerraUSTPair      = "terra163pkeeuwxzr0yhndf8xd2jprm9hrtk59xf7nqf"
	AddressTerraUSTLPToken   = "terra1q6r8hfdl203htfvpsmyh8x689lp2g0m7856fwd"
	AddressTerraNLUNAPair    = "terra1zvn8z6y8u2ndwvsjhtpsjsghk6pa6ugwzxp6vx"
	AddressTerraNLUNALPToken = "terra1tuw46dwfvahpcwf3ulempzsn9a0vhazut87zec"
	AddressTerraNETHPair     = "terra14zhkur7l7ut7tx6kvj28fp5q982lrqns59mnp3"
	AddressTerraNETHLPToken  = "terra1y8kxhfg22px5er32ctsgjvayaj8q36tr590qtp"

	AddressAstroGenerator      = "terra1zgrx9jjqrfye8swykfgmd6hpde60j0nszzupp9"
	AddressAstroUSTLPStaking   = "terra1fmu29xhg5nk8jr0p603y5qugpk2r0ywcyxyv7k"
	AddressAstroNLUNALPStaking = "terra1sxzggeujnxrd7hsx7uf2l6axh2uuv4zz5jadyg"
	AddressAstroNETHLPStaking  = "terra13n2sqaj25ugkt79k3evhvua30ut9qt8q0268zc"

	AddressTerraUSTLPStaking   = "terra12kzewegufqprmzl20nhsuwjjq6xu8t8ppzt30a"
	AddressTerraNLUNALPStaking = "terra1hs4ev0ghwn4wr888jwm56eztfpau6rjcd8mczc"
	AddressTerraNETHLPStaking  = "terra1lws09x0slx892ux526d6atwwgdxnjg58uan8ph"

	AddressApolloFarm = "terra1kn5kc4n9zu5mectfhsxukxd9qxhlc5q28pstyv"

	AddressSpectrumUSTLPStrategy   = "terra1jxh7hahwxlsy5cckkyhuz50a60mpn5tr0px6tq"
	AddressSpectrumNLUNALPStrategy = "terra19kzel57gvx42e628k6frh624x5vm2kpck9cr9c"
	AddressSpectrumNETHStrategy    = "terra1xw3jzqwrql5fvddchzxycd2ygrep5kudsden5c"
)

func ExportNexus(app *terra.TerraApp) (util.SnapshotBalanceAggregateMap, error) {
	logger := app.Logger()
	logger.Info("Exporting Nexus START")

	ctx := util.PrepCtx(app)
	qs := util.PrepWasmQueryServer(app)
	keeper := app.WasmKeeper

	//================================

	logger.Info("Fetching PSI holders")
	var psiHolderMap = make(util.BalanceMap)
	if err := getCW20Balances(ctx, keeper, AddressPSIToken, psiHolderMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI stakers")
	var psiStakerMap = make(util.BalanceMap)
	if err := getPSIGovBalances(ctx, keeper, qs, psiStakerMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI-UST LP holders on Astroport")
	var astroUSTLPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressAstroUSTLPToken, AddressAstroUSTPair, astroUSTLPHolderMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NLUNA LP holders on Astroport")
	var astroNLUNALPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressAstroNLUNALPToken, AddressAstroNLUNAPair, astroNLUNALPHolderMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NETH LP holders on Astroport")
	var astroNETHLPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressAstroNETHLPToken, AddressAstroNETHPair, astroNETHLPHolderMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NAVAX LP holders on Astroport")
	var astroNAVAXLPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressAstroNAVAXLPToken, AddressAstroNAVAXPair, astroNAVAXLPHolderMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NATOM LP holders on Astroport")
	var astroNATOMLPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressAstroNATOMLPToken, AddressAstroNATOMPair, astroNATOMLPHolderMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI-UST LP holders on Terraswap")
	var terraUSTLPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressTerraUSTLPToken, AddressTerraUSTPair, terraUSTLPHolderMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NLUNA LP holders on Terraswap")
	var terraNLUNALPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressTerraNLUNALPToken, AddressTerraNLUNAPair, terraNLUNALPHolderMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NETH LP holders on Terraswap")
	var terraNETHLPHolderMap = make(util.BalanceMap)
	if err := getLPBalances(ctx, keeper, qs, AddressTerraNETHLPToken, AddressTerraNETHPair, terraNETHLPHolderMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI-UST LP stakers on Astroport")
	var astroUSTLPStakerMap = make(util.BalanceMap)
	if err := getLPStakerBalances(ctx, keeper, qs, AddressAstroUSTLPStaking, AddressAstroUSTLPToken, AddressAstroUSTPair, astroUSTLPStakerMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NLUNA LP stakers on Astroport")
	var astroNLUNALPStakerMap = make(util.BalanceMap)
	if err := getLPStakerBalances(ctx, keeper, qs, AddressAstroNLUNALPStaking, AddressAstroNLUNALPToken, AddressAstroNLUNAPair, astroNLUNALPStakerMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NETH LP stakers on Astroport")
	var astroNETHLPStakerMap = make(util.BalanceMap)
	if err := getLPStakerBalances(ctx, keeper, qs, AddressAstroNETHLPStaking, AddressAstroNETHLPToken, AddressAstroNETHPair, astroNETHLPStakerMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI-UST LP stakers on Terraswap")
	var terraUSTLPStakerMap = make(util.BalanceMap)
	if err := getLPStakerBalances(ctx, keeper, qs, AddressTerraUSTLPStaking, AddressTerraUSTLPToken, AddressTerraUSTPair, terraUSTLPStakerMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NLUNA LP stakers on Terraswap")
	var terraNLUNALPStakerMap = make(util.BalanceMap)
	if err := getLPStakerBalances(ctx, keeper, qs, AddressTerraNLUNALPStaking, AddressTerraNLUNALPToken, AddressTerraNLUNAPair, terraNLUNALPStakerMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NETH LP stakers on Terraswap")
	var terraNETHLPStakerMap = make(util.BalanceMap)
	if err := getLPStakerBalances(ctx, keeper, qs, AddressTerraNETHLPStaking, AddressTerraNETHLPToken, AddressTerraNETHPair, terraNETHLPStakerMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching LP stakers on Astroport")
	var astroLPStakerMap = make(util.BalanceMap)
	pairs := []pair{
		{ContractAddr: AddressAstroUSTPair, LPTokenAddr: AddressAstroUSTLPToken, LPTokenSupply: sdk.ZeroInt(), PSIBalance: sdk.ZeroInt()},
		{ContractAddr: AddressAstroNLUNAPair, LPTokenAddr: AddressAstroNLUNALPToken, LPTokenSupply: sdk.ZeroInt(), PSIBalance: sdk.ZeroInt()},
		{ContractAddr: AddressAstroNETHPair, LPTokenAddr: AddressAstroNETHLPToken, LPTokenSupply: sdk.ZeroInt(), PSIBalance: sdk.ZeroInt()},
	}
	if err := getAstroLPStakerBalances(ctx, keeper, qs, AddressAstroGenerator, pairs, astroLPStakerMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI-UST Apollo farmers")
	var apolloFarmerMap = make(util.BalanceMap)
	if err := getApolloBalances(ctx, keeper, qs, apolloFarmerMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Fetching PSI-UST Spectrum depositors")
	var spectrumUSTLPDepositorMap = make(util.BalanceMap)
	if err := getSpectrumBalances(ctx, keeper, qs, AddressSpectrumUSTLPStrategy, AddressAstroUSTPair, AddressAstroUSTLPToken, spectrumUSTLPDepositorMap); err != nil {
		return nil, err
	}

	logger.Info("Fetching PSI-NLUNA Spectrum depositors")
	var spectrumNLUNALPDepositorMap = make(util.BalanceMap)
	if err := getSpectrumBalances(ctx, keeper, qs, AddressSpectrumNLUNALPStrategy, AddressAstroNLUNAPair, AddressAstroNLUNALPToken, spectrumNLUNALPDepositorMap); err != nil {
		return nil, err
	}

	//================================

	logger.Info("Merging")
	merged := util.MergeMaps(
		psiHolderMap,
		psiStakerMap,
		astroUSTLPHolderMap,
		astroNLUNALPHolderMap,
		astroNETHLPHolderMap,
		astroNAVAXLPHolderMap,
		astroNATOMLPHolderMap,
		terraUSTLPHolderMap,
		terraNLUNALPHolderMap,
		terraNETHLPHolderMap,
		astroLPStakerMap,
		apolloFarmerMap,
		spectrumUSTLPDepositorMap,
		spectrumNLUNALPDepositorMap,
	)

	logger.Info("Finalizing")
	var finalBalance = make(util.SnapshotBalanceAggregateMap)
	for userAddr, psiBalance := range merged {
		finalBalance[userAddr] = []util.SnapshotBalance{
			{
				Denom:   "PSI",
				Balance: psiBalance,
			},
		}
	}

	logger.Info("Exporting Nexus END")
	return finalBalance, nil
}

func isContract(ctx context.Context, keeper wasmkeeper.Keeper, address string) (bool, error) {
	contractAddr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return true, err
	}

	if _, err := keeper.GetContractInfo(sdk.UnwrapSDKContext(ctx), contractAddr); err != nil {
		return false, nil
	}

	return true, nil
}

func getCW20Balances(ctx context.Context, k wasmkeeper.Keeper, tokenAddress string, balanceMap map[string]sdk.Int) error {
	prefix := util.GeneratePrefix("balance")
	tokenAddr, err := sdk.AccAddressFromBech32(tokenAddress)
	if err != nil {
		return err
	}
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), tokenAddr, prefix, func(key, value []byte) bool {
		if contract, _ := isContract(ctx, k, string(key)); contract {
			return false
		}
		balance, ok := sdk.NewIntFromString(string(value[1 : len(value)-1]))
		// fmt.Printf("%s %s\n", string(key), balance.String())
		if ok {
			if balance.IsZero() {
				return false
			}
			if strings.Contains(string(key), "terra") {
				balanceMap[string(key)] = balance
			} else {
				addr := sdk.AccAddress(key)
				balanceMap[addr.String()] = balance
			}
		}
		return false
	})
	return nil
}

func getPSIGovBalances(ctx context.Context, k wasmkeeper.Keeper, q wasmtypes.QueryServer, balanceMap map[string]sdk.Int) error {
	var govState struct {
		TotalShare   sdk.Int `json:"total_share"`
		TotalDeposit sdk.Int `json:"total_deposit"`
	}
	if err := util.ContractQuery(ctx, q, &wasmtypes.QueryContractStoreRequest{
		ContractAddress: AddressNexusGov,
		QueryMsg:        []byte("{\"state\":{}}"),
	}, &govState); err != nil {
		return err
	}

	govAddr, err := sdk.AccAddressFromBech32(AddressNexusGov)
	if err != nil {
		return err
	}

	psiGovBalance, err := util.GetCW20Balance(ctx, q, AddressPSIToken, AddressNexusGov)
	if err != nil {
		return err
	}
	totalGovPSIBalance := psiGovBalance.Sub(govState.TotalDeposit)

	prefix := util.GeneratePrefix("bank")
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), govAddr, prefix, func(key, value []byte) bool {
		if contract, _ := isContract(ctx, k, string(key)); contract {
			return false
		}

		var staker struct {
			Share sdk.Int `json:"share"`
		}
		err := json.Unmarshal(value, &staker)
		if err != nil {
			panic(err)
		}

		if staker.Share.IsZero() {
			return false
		}

		balance := staker.Share.Mul(totalGovPSIBalance).Quo(govState.TotalShare)

		if strings.Contains(string(key), "terra") {
			balanceMap[string(key)] = balance
		} else {
			addr := sdk.AccAddress(key)
			balanceMap[addr.String()] = balance
		}
		return false
	})
	return nil
}

func getLPBalances(ctx context.Context, k wasmkeeper.Keeper, q wasmtypes.QueryServer, lpTokenAddress string, pairAddress string, balanceMap map[string]sdk.Int) error {
	lpSupply, err := util.GetCW20TotalSupply(ctx, q, lpTokenAddress)
	if err != nil {
		return err
	}

	pairPSIBalance, err := util.GetCW20Balance(ctx, q, AddressPSIToken, pairAddress)
	if err != nil {
		return err
	}

	prefix := util.GeneratePrefix("balance")
	lpTokenAddr, err := sdk.AccAddressFromBech32(lpTokenAddress)
	if err != nil {
		return err
	}
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), lpTokenAddr, prefix, func(key, value []byte) bool {
		if contract, _ := isContract(ctx, k, string(key)); contract {
			return false
		}
		lpBalance, ok := sdk.NewIntFromString(string(value[1 : len(value)-1]))
		if ok {
			if lpBalance.IsZero() {
				return false
			}
			balance := lpBalance.Mul(pairPSIBalance).Quo(lpSupply)
			if strings.Contains(string(key), "terra") {
				balanceMap[string(key)] = balance
			} else {
				addr := sdk.AccAddress(key)
				balanceMap[addr.String()] = balance
			}
		}
		return false
	})
	return nil
}

func getLPStakerBalances(ctx context.Context, k wasmkeeper.Keeper, q wasmtypes.QueryServer, contractAddress string, lpTokenAddress string, pairAddress string, balanceMap map[string]sdk.Int) error {
	lpSupply, err := util.GetCW20TotalSupply(ctx, q, lpTokenAddress)
	if err != nil {
		return err
	}

	pairPSIBalance, err := util.GetCW20Balance(ctx, q, AddressPSIToken, pairAddress)
	if err != nil {
		return err
	}

	contractAddr, err := sdk.AccAddressFromBech32(contractAddress)
	if err != nil {
		return err
	}

	prefix := util.GeneratePrefix("reward")
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), contractAddr, prefix, func(key, value []byte) bool {
		if contract, _ := isContract(ctx, k, string(key)); contract {
			return false
		}

		var staker struct {
			BondAmount    sdk.Int `json:"bond_amount"`
			PendingReward sdk.Int `json:"pending_reward"`
		}
		err := json.Unmarshal(value, &staker)
		if err != nil {
			panic(err)
		}

		if staker.BondAmount.IsZero() && staker.PendingReward.IsZero() {
			return false
		}

		balance := staker.BondAmount.Mul(pairPSIBalance).Quo(lpSupply).Add(staker.PendingReward)

		if strings.Contains(string(key), "terra") {
			balanceMap[string(key)] = balance
		} else {
			addr := sdk.AccAddress(key)
			balanceMap[addr.String()] = balance
		}
		return false
	})
	return nil
}

type pair struct {
	ContractAddr  string
	LPTokenAddr   string
	LPTokenSupply sdk.Int
	PSIBalance    sdk.Int
}

func getAstroLPStakerBalances(ctx context.Context, k wasmkeeper.Keeper, q wasmtypes.QueryServer, generatorAddress string, pairs []pair, balanceMap map[string]sdk.Int) error {
	for i := range pairs {
		supply, err := util.GetCW20TotalSupply(ctx, q, pairs[i].LPTokenAddr)
		if err != nil {
			return err
		}

		psiBalance, err := util.GetCW20Balance(ctx, q, AddressPSIToken, pairs[i].ContractAddr)
		if err != nil {
			return err
		}

		pairs[i].LPTokenSupply = supply
		pairs[i].PSIBalance = psiBalance
	}

	generatorAddr, err := sdk.AccAddressFromBech32(generatorAddress)
	if err != nil {
		return err
	}

	prefix := util.GeneratePrefix("user_info")
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), generatorAddr, prefix, func(key, value []byte) bool {
		lpTokenAddress := string(key[2:46])
		stakerAddress := string(key[46:90])

		var curPair *pair = nil
		for _, p := range pairs {
			if p.LPTokenAddr == lpTokenAddress {
				curPair = &p
				break
			}
		}
		if curPair == nil {
			return false
		}

		if contract, _ := isContract(ctx, k, stakerAddress); contract {
			return false
		}

		var staker struct {
			Amount sdk.Int `json:"amount"`
		}
		err := json.Unmarshal(value, &staker)
		if err != nil {
			panic(err)
		}

		if staker.Amount.IsZero() {
			return false
		}

		balance := staker.Amount.Mul(curPair.PSIBalance).Quo(curPair.LPTokenSupply)
		// fmt.Printf("%s %s %s %s\n", lpTokenAddress, stakerAddress, staker.Amount.String(), balance.String())

		if strings.Contains(stakerAddress, "terra") {
			if curBalance, ok := balanceMap[stakerAddress]; ok {
				balanceMap[stakerAddress] = curBalance.Add(balance)
			} else {
				balanceMap[stakerAddress] = balance
			}
		} else {
			addr := sdk.AccAddress(key[46:90])
			if curBalance, ok := balanceMap[addr.String()]; ok {
				balanceMap[addr.String()] = curBalance.Add(balance)
			} else {
				balanceMap[addr.String()] = balance
			}
		}
		return false
	})
	return nil
}

func getApolloBalances(ctx context.Context, k wasmkeeper.Keeper, q wasmtypes.QueryServer, balanceMap map[string]sdk.Int) error {
	var strategyInfo struct {
		TotalBondAmount sdk.Int `json:"total_bond_amount"`
		TotalShares     sdk.Int `json:"total_shares"`
	}

	if err := util.ContractQuery(ctx, q, &wasmtypes.QueryContractStoreRequest{
		ContractAddress: AddressApolloFarm,
		QueryMsg:        []byte("{\"strategy_info\":{}}"),
	}, &strategyInfo); err != nil {
		return err
	}

	supplyLPToken, err := util.GetCW20TotalSupply(ctx, q, AddressAstroUSTLPToken)
	if err != nil {
		return err
	}

	pairPSIBalance, err := util.GetCW20Balance(ctx, q, AddressPSIToken, AddressAstroUSTPair)
	if err != nil {
		return err
	}

	apolloFarmAddr, err := sdk.AccAddressFromBech32(AddressApolloFarm)
	if err != nil {
		return err
	}

	prefix := util.GeneratePrefix("user")
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), apolloFarmAddr, prefix, func(key, value []byte) bool {
		var userInfo struct {
			Shares sdk.Int `json:"shares"`
		}
		err := json.Unmarshal(value, &userInfo)
		if err != nil {
			panic(err)
		}

		if userInfo.Shares.IsZero() {
			return false
		}

		walletAddr := sdk.AccAddress(key)
		if contract, _ := isContract(ctx, k, walletAddr.String()); contract {
			return false
		}

		lpTokenAmount := userInfo.Shares.Mul(strategyInfo.TotalBondAmount).Quo(strategyInfo.TotalShares)
		balance := lpTokenAmount.Mul(pairPSIBalance).Quo(supplyLPToken)
		balanceMap[walletAddr.String()] = balance

		// fmt.Printf("APOLLO %s %s %s\n", walletAddr.String(), lpAmount.String(), balance.String())
		return false
	})

	return nil
}

func getSpectrumBalances(ctx context.Context, k wasmkeeper.Keeper, q wasmtypes.QueryServer, strategyAddress string, pairAddress string, lpTokenAddress string, balanceMap map[string]sdk.Int) error {
	supplyLPToken, err := util.GetCW20TotalSupply(ctx, q, lpTokenAddress)
	if err != nil {
		return err
	}

	pairPSIBalance, err := util.GetCW20Balance(ctx, q, AddressPSIToken, pairAddress)
	if err != nil {
		return err
	}

	strategyAddr, err := sdk.AccAddressFromBech32(strategyAddress)
	if err != nil {
		return err
	}

	prefix := util.GeneratePrefix("reward")
	k.IterateContractStateWithPrefix(sdk.UnwrapSDKContext(ctx), strategyAddr, prefix, func(key, value []byte) bool {
		walletAddress := sdk.AccAddress(key[2:22])
		if contract, _ := isContract(ctx, k, walletAddress.String()); contract {
			return false
		}

		var rewardInfo struct {
			RewardInfo []struct {
				LPTokenAmount sdk.Int `json:"bond_amount"`
			} `json:"reward_infos"`
		}
		err := util.ContractQuery(ctx, q, &wasmtypes.QueryContractStoreRequest{
			ContractAddress: strategyAddress,
			QueryMsg:        []byte(fmt.Sprintf("{\"reward_info\":{\"staker_addr\":\"%s\"}}", walletAddress.String())),
		}, &rewardInfo)
		if err != nil {
			panic(err)
		}

		lpTokenAmount := rewardInfo.RewardInfo[0].LPTokenAmount
		if lpTokenAmount.IsZero() {
			return false
		}
		balance := lpTokenAmount.Mul(pairPSIBalance).Quo(supplyLPToken)
		balanceMap[walletAddress.String()] = balance

		// fmt.Printf("SPECTRUM %s %s %s\n", walletAddress.String(), lpAmount.String(), balance.String())
		return false
	})

	return nil
}
