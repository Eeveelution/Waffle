package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func UserStatsFromDatabase(id uint64, mode int8) (int8, UserStats) {
	returnStats := UserStats{}

	queryResult, queryErr := database.Query("SELECT * FROM (SELECT user_id, mode, ROW_NUMBER() OVER (ORDER BY ranked_score DESC) AS 'rank', ranked_score, total_score, user_level, accuracy, playcount, count_ssh, count_ss, count_sh, count_s, count_a, count_b, count_c, count_d, hit300, hit100, hit50, hitMiss, hitGeki, hitKatu, replays_watched FROM waffle.stats WHERE mode = ?) t WHERE user_id = ?", mode, id)
	defer queryResult.Close()

	if queryErr != nil {
		fmt.Printf("Failed to Fetch User Stats from Database, MySQL query failed.\n")

		return -2, returnStats
	}

	if queryResult.Next() {
		scanErr := queryResult.Scan(&returnStats.UserID, &returnStats.Mode, &returnStats.Rank, &returnStats.RankedScore, &returnStats.TotalScore, &returnStats.Level, &returnStats.Accuracy, &returnStats.Playcount, &returnStats.CountSSH, &returnStats.CountSS, &returnStats.CountSH, &returnStats.CountS, &returnStats.CountA, &returnStats.CountB, &returnStats.CountC, &returnStats.CountD, &returnStats.Hit300, &returnStats.Hit100, &returnStats.Hit50, &returnStats.HitMiss, &returnStats.HitGeki, &returnStats.HitKatu, &returnStats.ReplaysWatched)

		if scanErr != nil {
			fmt.Printf("Failed to Scan database results onto UserStats object.\n")

			return -2, returnStats
		}

		return 0, returnStats
	}

	return -1, returnStats
}
