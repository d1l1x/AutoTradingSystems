package main

import (
	gf "github.com/d1l1x/gofin"
	"github.com/d1l1x/gofin/indicators"
	"github.com/d1l1x/gofin/providers"
	"github.com/d1l1x/gofin/utils"
	"go.uber.org/zap"

	"github.com/d1l1x/gofin/brokers"
)

var log = utils.NewZapLogger("MyAwesomeTradingSystem", utils.Debug)

func PrepareAssets(broker *brokers.AlpacaBroker) *[]utils.Asset {

	brokerAssets, err := broker.GetListOfAssets("active", "us_equity", "")
	if err != nil {
		log.Fatal("Get list of assets", zap.Error(err))
	}
	log.Info("Number of assets", zap.Int("Number of assets", len(brokerAssets)))

	var assets []utils.Asset

	log.Debug("Transform list of assets")
	for _, asset := range brokerAssets {
		assets = append(assets, utils.Asset{
			Symbol: asset.Symbol,
			Name:   asset.Name,
			Id:     asset.ID},
		)
	}
	return &assets
}

func SetupFilters() *[]utils.Filter {
	log.Info("Setup filters")
	var filters []utils.Filter
	roc3 := utils.NewFilter(indicators.ROC([]float64{}, 3), utils.LT, 1.9)

	filters = append(filters, *roc3)

	return &filters
}

func SetupRanking() *utils.Ranking {
	log.Info("Setup ranking")
	ranking := utils.Ranking{Indicator: indicators.ROC([]float64{}, 3), Order: utils.Descending}
	return &ranking
}

func SetupCalendar() *utils.TradingCalendar {
	log.Info("Setup trading calendar")
	cal, err := utils.NewTradingCalendarUS()
	if err != nil {
		log.Fatal("New trading calendar", zap.Error(err))
	}
	return cal
}

func main() {

	alpaca := brokers.Alpaca(nil, "")

	provider := providers.FMP("", 300)

	cal := SetupCalendar()
	assets := PrepareAssets(alpaca)
	filters := SetupFilters()
	ranking := SetupRanking()

	watchlist := utils.NewWatchlist(*assets, *filters, ranking)

	wl := utils.Watchlist{
		Assets:  watchlist.Assets[:100],
		Filters: watchlist.Filters,
		Ranking: watchlist.Ranking,
	}
	ts := gf.NewTradingSystem(alpaca, &wl, cal, provider)

	//ts := gf.NewTradingSystem(alpaca, watchlist, cal, provider)

	ts.Init()
	ts.Run()

}
