package main

import (
	"strings"
	"time"

	"github.com/simonks2016/order_flow_analysis"
	"github.com/simonks2016/sliding_window"
)

type MicrostructureOvershootFactor struct {
	slidingWindow *sliding_window.SlidingWindow
	orderFlow     *order_flow_analysis.OrderFlowEngine
	dirScale      float64
	momentumScale float64
}

func NewMicrostructureOvershootFactor(emaAlpha float64, slidingWinDuration time.Duration, slidingWinCap int) *MicrostructureOvershootFactor {
	return &MicrostructureOvershootFactor{
		slidingWindow: sliding_window.NewSlidingWindow(slidingWinDuration, slidingWinCap, emaAlpha),
		orderFlow:     order_flow_analysis.NewOrderFlowEngine(emaAlpha, slidingWinCap),
		dirScale:      0.00005920,
		momentumScale: 0.00033251,
	}
}

func (o *MicrostructureOvershootFactor) SetScale(dirScale, momentumScale float64) *MicrostructureOvershootFactor {
	o.dirScale = dirScale
	o.momentumScale = momentumScale
	return o
}

// AddTrade
// Add Transaction to order flow analysis
// 添加交易信息到订单流分析器
func (o *MicrostructureOvershootFactor) AddTrade(isBuy bool, currentPrice float64, currentVol float64, currentTimestamp time.Time) {
	o.orderFlow.UpdateTrade(order_flow_analysis.Trade{
		Ts:     currentTimestamp,
		Price:  currentPrice,
		Volume: currentVol,
		Side: func() order_flow_analysis.Side {
			if isBuy {
				return order_flow_analysis.SideBuy
			}
			return order_flow_analysis.SideSell
		}(),
	})
}

// AddTick
// Add market data to a sliding window
// 将行情当中数据添加到滑动窗口当中
func (o *MicrostructureOvershootFactor) AddTick(price, vol float64, currentTimeStamp time.Time) {
	o.slidingWindow.Add(sliding_window.WindowPoint{
		Ts:     currentTimeStamp,
		Price:  price,
		Volume: vol,
	})
}

// AddOrderBook
// Add order book data to the order flow analyzer
// 添加盘口信息到订单流分析器当中
func (o *MicrostructureOvershootFactor) AddOrderBook(currentTime time.Time, level ...Level) {

	var bidsLevel []order_flow_analysis.Level
	var asksLevel []order_flow_analysis.Level

	for _, l := range level {

		if strings.ToLower(l.Type()) == "bids" {
			bidsLevel = append(bidsLevel, order_flow_analysis.Level{
				Price: l.Price(),
				Size:  l.Size(),
			})
		} else {
			asksLevel = append(asksLevel, order_flow_analysis.Level{
				Price: l.Price(),
				Size:  l.Size(),
			})
		}
	}
	o.orderFlow.UpdateOrderBook(order_flow_analysis.OrderBook{
		Ts:   currentTime,
		Bids: bidsLevel,
		Asks: asksLevel,
	}, 5)
}

func (o *MicrostructureOvershootFactor) Momentum() (float64, bool) {
	return o.slidingWindow.Momentum()
}
func (o *MicrostructureOvershootFactor) TotalVol() float64 {
	return o.slidingWindow.TotalVolume()
}
func (o *MicrostructureOvershootFactor) MidPrice() (float64, bool) {
	return o.slidingWindow.MedianPrice()
}
func (o *MicrostructureOvershootFactor) Score(currentMomentum float64) (float64, bool) {
	s, err := o.slidingWindow.ScoreWithMomentum(
		currentMomentum,
		o.dirScale,
		o.momentumScale,
		o.orderFlow.GetConfidence())
	if err != nil {
		return 0.0, false
	} else {
		return s, true
	}
}
