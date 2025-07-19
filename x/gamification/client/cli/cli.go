/*
Copyright 2024 DeshChain Foundation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group gamification queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "gamification",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "gamification"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryProfile())
	cmd.AddCommand(CmdQueryProfileByUsername())
	cmd.AddCommand(CmdQueryAchievements())
	cmd.AddCommand(CmdQueryAchievement())
	cmd.AddCommand(CmdQueryLeaderboard())
	cmd.AddCommand(CmdQueryDailyChallenge())
	cmd.AddCommand(CmdQueryTeamBattles())
	cmd.AddCommand(CmdQueryTeamBattle())
	cmd.AddCommand(CmdQuerySocialPosts())
	cmd.AddCommand(CmdQueryHumorQuotes())
	cmd.AddCommand(CmdQueryLevelConfig())

	return cmd
}

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "gamification",
		Short:                      fmt.Sprintf("%s transactions subcommands", "gamification"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateProfile())
	cmd.AddCommand(CmdUpdateProfile())
	cmd.AddCommand(CmdSelectAvatar())
	cmd.AddCommand(CmdClaimAchievement())
	cmd.AddCommand(CmdRecordAction())
	cmd.AddCommand(CmdShareAchievement())
	cmd.AddCommand(CmdJoinTeamBattle())
	cmd.AddCommand(CmdCreateTeamBattle())
	cmd.AddCommand(CmdCompleteDailyChallenge())

	return cmd
}