package app

import (
	"fmt"

	"github.com/terra-money/core/app/export/nexus"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
	terra "github.com/terra-money/core/app"
	"github.com/terra-money/core/app/export/util"
)

func ExportContracts(app *terra.TerraApp) []types.Balance {
	// var err error
	var snapshotType util.Snapshot
	if app.LastBlockHeight() == 7544910 {
		snapshotType = util.Snapshot(util.PreAttack)
	} else {
		snapshotType = util.Snapshot(util.PostAttack)
	}

	logger := app.Logger()
	logger.Info(fmt.Sprintf("Exporting Contracts @ %d - %s", app.LastBlockHeight(), snapshotType))

	err := nexus.FindLiquidation(app, 7544910)
	check(err)

	// snapshot := util.MergeSnapshots(nexusSs)

	return nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
