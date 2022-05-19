package web

import (
	"Waffle/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandleOsuOsz2GetScores(ctx *gin.Context) {
	skipScores := ctx.Query("s")
	beatmapChecksum := ctx.Query("c")
	osuFilename := ctx.Query("f")
	queryUserId := ctx.Query("u")
	queryPlaymode := ctx.Query("m")
	//beatmapsetId := ctx.Query("i")
	//osz2hash := ctx.Query("h")

	userId, parseErr := strconv.ParseInt(queryUserId, 10, 64)
	//playmode, parseErr := strconv.ParseInt(queryPlaymode, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	leaderboardBeatmapQueryResult, leaderboardBeatmap := database.BeatmapsGetByFilename(osuFilename)

	if leaderboardBeatmapQueryResult == -2 {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	if leaderboardBeatmapQueryResult == -1 {
		ctx.String(http.StatusOK, "-1|false")
		return
	}

	if beatmapChecksum != leaderboardBeatmap.BeatmapMd5 {
		ctx.String(http.StatusOK, "1|false")
		return
	}

	beatmapsetQueryResult, beatmapset := database.BeatmapsetsGetBeatmapsetById(leaderboardBeatmap.BeatmapsetId)

	if beatmapsetQueryResult == -2 {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	returnString := ""

	returnRankedStatus := "0"

	switch leaderboardBeatmap.RankingStatus {
	case 0:
		ctx.String(http.StatusOK, "0|false")
		return
	case 1:
		returnRankedStatus = "2"
		break
	case 2:
		returnRankedStatus = "3"
		break
	}
	//Ranked Status|Server has osz2 of map
	returnString += returnRankedStatus + "|false\n"
	//Online Offset, currently we don't store any so eh, TODO
	returnString += "0\n"
	//Display Title
	returnString += fmt.Sprintf("[bold:0,size:20]%s|%s\n", beatmapset.Artist, beatmapset.Title)
	//Online Rating, currently rating doesnt exist, so TODO
	returnString += "0\n"

	if skipScores == "1" {
		ctx.String(http.StatusOK, returnString)
		return
	}

	playmode, parseErr := strconv.ParseInt(queryPlaymode, 10, 64)

	if parseErr != nil {
		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	userBestScoreQueryResult, userBestScore, userUsername, userOnlineRank := database.ScoresGetUserLeaderboardBest(leaderboardBeatmap.BeatmapId, uint64(userId), int8(playmode))

	if userBestScoreQueryResult == -1 || userBestScore.Passed == 0 {
		returnString += "\n"
	} else {
		returnString += userBestScore.ScoresFormatLeaderboardScore(userUsername, int32(userOnlineRank))
	}

	leaderboardQuery, leaderboardQueryErr := database.Database.Query("SELECT ROW_NUMBER() OVER (ORDER BY score DESC) AS 'online_rank', users.username, scores.* FROM waffle.scores LEFT JOIN waffle.users ON scores.user_id = users.user_id WHERE beatmap_id = ? AND leaderboard_best = 1 AND passed = 1 AND playmode = ? ORDER BY score DESC", leaderboardBeatmap.BeatmapId, int8(playmode))

	if leaderboardQueryErr != nil {
		if leaderboardQuery != nil {
			leaderboardQuery.Close()
		}

		ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
		return
	}

	for leaderboardQuery.Next() {
		returnScore := database.Score{}

		var username string
		var onlineRank int64

		scanErr := leaderboardQuery.Scan(&onlineRank, &username, &returnScore.ScoreId, &returnScore.BeatmapId, &returnScore.BeatmapsetId, &returnScore.UserId, &returnScore.Playmode, &returnScore.Score, &returnScore.MaxCombo, &returnScore.Ranking, &returnScore.Hit300, &returnScore.Hit100, &returnScore.Hit50, &returnScore.HitMiss, &returnScore.HitGeki, &returnScore.HitKatu, &returnScore.EnabledMods, &returnScore.Perfect, &returnScore.Passed, &returnScore.Date, &returnScore.LeaderboardBest, &returnScore.MapsetBest, &returnScore.ScoreHash)

		if scanErr != nil {
			leaderboardQuery.Close()
			ctx.String(http.StatusOK, "deliberatly fucked string to give out an error client side cuz pain\n")
			return
		}

		returnString += returnScore.ScoresFormatLeaderboardScore(username, int32(onlineRank))
	}

	leaderboardQuery.Close()

	ctx.String(http.StatusOK, returnString)
}
