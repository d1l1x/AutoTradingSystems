package main

import (
	gf "github.com/d1l1x/gofin"
	"github.com/d1l1x/gofin/indicators"
	"github.com/d1l1x/gofin/providers"
	"github.com/d1l1x/gofin/utils"
	"go.uber.org/zap"

	"github.com/d1l1x/gofin/brokers"
	"github.com/sirupsen/logrus"
)

var log = utils.NewZapLogger("MyAwesomeTradingSystem", utils.Debug)

func main() {

	alpaca := brokers.Alpaca(nil, "")

	// provider of historical data
	provider := providers.FMP("", 300)

	cal, err := utils.NewTradingCalendarUS()
	if err != nil {
		log.Fatal("New trading calendar", zap.Error(err))
	}

	assets, err := alpaca.GetListOfAssets("active", "us_equity", "")
	if err != nil {
		log.Fatal("Get list of assets", zap.Error(err))
	}

	watchlist := utils.NewWatchlist()
	log.Debug("Add Assets to watchlist")
	for _, asset := range assets {
		watchlist.AddAsset(utils.Asset{Symbol: asset.Symbol, Name: asset.Name, Id: asset.ID})
	}

	// Setup filters to be applied to every asset
	roc3 := utils.NewFilter(indicators.ROC([]float64{}, 3), utils.LT, 1.9)
	watchlist.AddFilter(roc3)

	ranking := utils.Ranking{Indicator: indicators.ROC([]float64{}, 3), Order: utils.Descending}
	watchlist.AddRanking(&ranking)

	wl := utils.Watchlist{
		Assets:  watchlist.Assets[:100],
		Filters: watchlist.Filters,
		Ranking: watchlist.Ranking,
	}
	ts := gf.NewTradingSystem(alpaca, &wl, cal, provider)

	//ts := gf.NewTradingSystem(alpaca, watchlist, cal, provider)

	ts.Init(logrus.DebugLevel)
	ts.Run()

}
