// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// findTrade/newTrade are buy.c:809–826 and 851–870. A zero-quantity entry is not found:
// loop_trade copies the pointer list and filters positive quantities.
func (w *world) findTrade(who, kind, item int) *trade {
	assert(w.validBox(who), "olympia/buy.c", 815, "valid_box(who)")
	for t := range w.loopTrade(who) {
		if int(t.kind) == kind && int(t.item) == item {
			return t
		}
	}
	return nil
}

func (w *world) newTrade(who, kind, item int) *trade {
	if t := w.findTrade(who, kind, item); t != nil {
		return t
	}
	t := &trade{who: int32(who), kind: int32(kind), item: int32(item)}
	w.bx[who].trades.Append(t)
	return t
}

// add_city_trade, buy.c:1161–1172, also used by INIT's market seeding.
func (w *world) addCityTrade(where, kind, item, qty, cost, month int) *trade {
	t := w.newTrade(where, kind, item)
	t.qty, t.cost, t.monthProd = int32(qty), int32(cost), int32(month)
	return t
}

// Actual opium_data, buy.c:1195–1212 (not the older prose price table).
var opiumData = [...]struct{ qty, cost int32 }{
	{15, 17}, {28, 17}, {37, 18}, {50, 18},
	{55, 19}, {66, 19}, {70, 20}, {80, 20},
}

// opium_market_delta, buy.c:1215–1269. Deliberately raw: the first BUY
// includes sold-out entries excluded by find_trade (upstream issue #252).
func (w *world) opiumMarketDelta(where int) {
	assert(w.subkind(where) == subCity, "olympia/buy.c", 1222, "subkind(where) == sub_city")
	var t *trade
	for _, e := range w.bx[where].trades {
		if e.kind == tradeBuy && e.item == itemOpium {
			t = e
			break
		}
	}
	p := w.pSubloc(where)
	if t != nil {
		if t.qty < 1 {
			p.opiumEcon++
		} else if t.qty == opiumData[p.opiumEcon].qty {
			p.opiumEcon--
		}
	}
	if p.opiumEcon > 7 {
		p.opiumEcon = 7
	}
	if p.opiumEcon < 0 {
		p.opiumEcon = 0
	}
	t = w.newTrade(where, tradeConsume, itemOpium)
	t.qty, t.cost = opiumData[p.opiumEcon].qty, opiumData[p.opiumEcon].cost
	t.cloak = 1
	if p.opiumEcon == 0 {
		t.cloak = 2
	}
}

// expire_trades, buy.c:1271–1309. Raw, live deletion loop: even zero-qty
// offers age. BUY survival checks character inventory, not other markets.
func (w *world) expireTrades(where int) {
	for i := 0; i < len(w.bx[where].trades); i++ {
		t := w.bx[where].trades[i]
		if t.expire > 0 {
			t.expire--
		}
		if w.subkind(int(t.item)) == subTradegood && t.expire <= 0 {
			done := true
			if t.kind == tradeBuy {
				for j := range w.loopKind(tChar) {
					for e := range w.loopInv(j) {
						if e.item == t.item {
							done = false
						}
					}
				}
			}
			if done {
				w.bx[where].trades.Delete(i)
				i--
			}
		}
	}
}

// loc_trade_sup, buy.c:1318–1364. Annual gating applies only to PRODUCE;
// CONSUME always restocks. Existing surplus is a floor, never added to.
func (w *world) locTradeSup(where int, override bool) {
	w.expireTrades(where)
	assert(w.validBox(where), "olympia/buy.c", 1326, "valid_box(where)")
	for t := range w.loopTrade(where) {
		okay := true
		if t.monthProd != 0 && !override {
			thisMonth := olyMonth(w.sysclock) - 1
			nextMonth := (thisMonth + 1) % numMonths
			if nextMonth != int(t.monthProd)-1 {
				okay = false
			}
		}
		var n *trade
		if t.kind == tradeProduce && okay {
			n = w.newTrade(where, tradeSell, int(t.item))
		} else if t.kind == tradeConsume {
			n = w.newTrade(where, tradeBuy, int(t.item))
		}
		if n != nil {
			if n.qty < t.qty {
				n.qty = t.qty
			}
			n.cost, n.cloak, n.expire = t.cost, t.cloak, t.expire
		}
	}
}

// trade_suffuse_ring, buy.c:1367–1395, on the existing world's econ stream.
func (w *world) tradeSuffuseRing(where int) {
	found := false
	assert(w.validBox(where), "olympia/buy.c", 1375, "valid_box(where)")
	for t := range w.loopTrade(where) {
		if w.subkind(int(t.item)) == subSuffuseRing && t.kind == tradeSell && t.qty > 0 {
			found = true
		}
	}
	w.streams.BeginEconomy(where)
	if found || w.streams.EconStock(subSuffuseRing, where, 1, 3) < 3 {
		return
	}
	item := w.newSuffuseRing(where)
	n := w.newTrade(where, tradeSell, item)
	n.qty = 1
	n.cost = int32(450 + w.streams.EconCost(subSuffuseRing, where, 0, 12)*50)
	n.cloak = 0
}

// location_trades, buy.c:1398–1414. No INIT city/route reseeding here.
func (w *world) locationTrades() {
	w.stage("location_trades()")
	for where := range w.loopSubkind(subCity) {
		w.opiumMarketDelta(where)
		w.locTradeSup(where, false)
		if r := w.region(where); r == int(w.faeryRegion) || r == int(w.cloudRegion) {
			w.tradeSuffuseRing(where)
		}
	}
}
