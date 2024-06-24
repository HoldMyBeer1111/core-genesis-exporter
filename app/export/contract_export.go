package app

import (
	"github.com/terra-money/core/app/export/nexus"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
	terra "github.com/terra-money/core/app"
)

func ExportContracts(app *terra.TerraApp, height int64) []types.Balance {
	// var err error
	// var snapshotType util.Snapshot
	// if app.LastBlockHeight() == 7544910 {
	// 	snapshotType = util.Snapshot(util.PreAttack)
	// } else {
	// 	snapshotType = util.Snapshot(util.PostAttack)
	// }

	// logger := app.Logger()
	// logger.Info(fmt.Sprintf("Exporting Contracts @ %d - %s", app.LastBlockHeight(), snapshotType))

	nexus.FindLiquidation(app, height)

	// check(err)

	// snapshot := util.MergeSnapshots(nexusSs)

	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func LiquidationSearch(app *terra.TerraApp, height int64) (bool, error) {
	return nexus.FindLiquidation(app, height)
}
