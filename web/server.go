package web

import (
	"Waffle/web/waffle_api"

	"github.com/gin-gonic/gin"
)

func RunOsuWeb() {
	ginServer := gin.Default()
	ginServer.SetTrustedProxies(nil)

	// /web
	ginServer.POST("/web/osu-screenshot.php", HandleOsuScreenshot)
	ginServer.GET("/web/osu-title-image.php", HandleTitleImage)
	ginServer.POST("/web/osu-submit-modular.php", HandleOsuSubmit)
	ginServer.GET("/web/osu-osz2-getscores.php", HandleOsuGetLeaderboards)
	ginServer.GET("/web/osu-getscores6.php", HandleOsuGetLeaderboards)
	ginServer.GET("/web/osu-getreplay.php", HandleGetReplay)
	ginServer.GET("/web/osu-getfavourites.php", HandleOsuGetFavourites)
	ginServer.GET("/web/osu-addfavourite.php", HandleOsuAddFavourite)
	ginServer.POST("/web/osu-comment.php", HandleOsuComments)
	ginServer.GET("/rating/ingame-rate2.php", HandleOsuIngameRate2)
	ginServer.GET("/web/osu-search.php", HandleOsuDirectSearch)
	ginServer.GET("/web/maps/:filename", HandleOsuMapUpdate)

	// updater
	//ginServer.GET("/p/changelog", HandleUpdaterChangelog)
	//ginServer.GET("/release/update2.txt", HandleUpdaterUpdate2)
	//ginServer.GET("/release/update2.php", HandleOsuUpdate2)
	//ginServer.GET("/release/:filename", HandleUpdaterGetFile)

	//direct stuff
	ginServer.GET("/mt/:filename", HandleOsuGetDirectThumbnail)
	ginServer.GET("/mp3/preview/:filename", HandleOsuGetDirectMp3Preview)
	ginServer.GET("/d/:filename", HandleOsuDirectDownload)

	//avatars
	ginServer.GET("/a/:filename", HandleOsuGetAvatar)

	// screenshots
	ginServer.GET("/ss/:filename", HandleOsuGetScreenshot)

	//api
	ginServer.POST("/api/waffle-login", waffle_api.ApiHandleWaffleLogin)
	ginServer.POST("/api/waffle-site-register", waffle_api.ApiHandleWaffleRegister)

	ginServer.Run("127.0.0.1:80")
}
