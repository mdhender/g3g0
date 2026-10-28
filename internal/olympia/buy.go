// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import (
	"fmt"
	"slices"

	"github.com/mdhender/g3g0/internal/clib"
)

// Ported from olympia/buy.c: the per-game trade-route buyer secret (G3
// issue #46), which load_db() sets up with init_trade_routes() and
// save_db() keeps with save_trade_routes(), and
// investigate_possible_trade(), which add_item() calls, and complete trade
// settlement/dispatch. Monthly market restock is in buy_month.go.
//
// The FIND BUY buyer test picks which cities buy a tradegood with
// md5_int(city_sold, where, item, SECRET) & 1. The secret is 0xb05c0e
// unless the game has a seed of its own, from -G or lib/trade-route-seed
// when it is created, which is kept in lib/trade_routes; a kept seed always
// wins, so the secret cannot change mid-game.

// tradeRouteDefaultSecret is TRADE_ROUTE_DEFAULT_SECRET.
const tradeRouteDefaultSecret = 0xb05c0e

// splitmix64 is splitmix64().
func splitmix64(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}

// applyTradeRouteSeed is apply_trade_route_seed(): seed becomes the game's
// seed, and the secret is splitmix64 of it, folded from 64 bits to 32.
func (w *world) applyTradeRouteSeed(seed uint64) {
	s := splitmix64(seed)

	w.tradeRouteHaveSeed = true
	w.tradeRouteSeed = seed
	w.tradeRouteSecret = uint32(s ^ (s >> 32)) // fold 64 -> 32
}

// readTradeRouteSeed reads an unlimited decimal token as read_trade_route_seed
// does, saturating overflow even after a minus sign. ok reports fopen success.
func readTradeRouteSeed(path string) (seed uint64, scanned, ok bool) {
	data, err := clib.ReadFile(path)
	if err != nil {
		return 0, false, false
	}

	i := 0
	for i < len(data) && clib.IsSpace(data[i]) {
		i++
	}
	neg := false
	if i < len(data) && (data[i] == '+' || data[i] == '-') {
		neg = data[i] == '-'
		i++
	}
	over := false
	for ; i < len(data) && isdigit(data[i]); i++ {
		d := uint64(data[i] - '0')
		if seed > (^uint64(0)-d)/10 {
			over = true
		} else {
			seed = seed*10 + d
		}
		scanned = true
	}
	if over {
		seed = ^uint64(0)
	} else if neg {
		seed = -seed
	}
	return seed, scanned, true
}

// initTradeRoutes is init_trade_routes(): it sets the buyer secret. A
// readable libdir/trade_routes (a continuing game) always wins, with a
// warning if -G was given; then -G; then libdir/trade-route-seed. A file
// that does not start with a number leaves the default, and, for
// trade_routes, still wins.
func (w *world) initTradeRoutes() {
	fnam := sout("%s/trade_routes", w.libdir)
	if seed, scanned, ok := readTradeRouteSeed(fnam); ok {
		if scanned {
			w.applyTradeRouteSeed(seed)
		}

		if w.tradeRouteSeedPending {
			fmt.Fprintf(w.stderr, "trade-route seed: %s already exists; "+
				"ignoring -G/lib/trade-route-seed\n", fnam)
		}
		return
	}

	if w.tradeRouteSeedPending {
		w.applyTradeRouteSeed(w.tradeRouteSeedValue)
		return
	}

	fnam = sout("%s/trade-route-seed", w.libdir)
	if seed, scanned, ok := readTradeRouteSeed(fnam); ok && scanned {
		w.applyTradeRouteSeed(seed)
	}
}

// saveTradeRoutes is save_trade_routes(): it writes the game's seed to
// libdir/trade_routes, if it has one; the default secret writes nothing.
func (w *world) saveTradeRoutes() {
	if !w.tradeRouteHaveSeed {
		return
	}

	fnam := sout("%s/trade_routes", w.libdir)
	fp, err := clib.Fopen(fnam)
	if err != nil {
		fmt.Fprintf(w.stderr, "couldn't write %s: ", fnam)
		clib.Perror(w.stderr, "", err)
		return
	}

	fmt.Fprintf(fp, "%d\n", w.tradeRouteSeed)
	fp.Close()
}

// investigatePossibleTrade is investigate_possible_trade(): who has been
// given some item, and had oldHas of it before. If a pending trade of who's
// was not already satisfied at the old quantity (a BUY, for gold; a SELL of
// the item otherwise), who is queued on trades_to_check, for
// match_all_trades() to try soon.
//
// The gold test divides by the trade's cost, as C does; a cost of 0 would
// be SIGFPE in C and a runtime panic here.
func (w *world) investigatePossibleTrade(who, item, oldHas int) {
	check := false

	if item == itemGold {
		for t := range w.loopTrade(who) {
			if t.kind != tradeBuy {
				continue
			}

			if (int32(oldHas)-t.haveLeft)/t.cost < t.qty {
				check = true
				break
			}
		}
	} else {
		for t := range w.loopTrade(who) {
			if t.kind != tradeSell || int(t.item) != item {
				continue
			}

			if int32(oldHas)-t.haveLeft < t.qty {
				check = true
				break
			}
		}
	}

	if check {
		w.tradesToCheck.Append(who)
	}
}

// marketHere is buy.c:342–355. A ship cuts off the walk even in a city.
func (w *world) marketHere(who int) int {
	for who > 0 && w.subkind(who) != subCity {
		if w.isShip(who) {
			return 0
		}
		who = w.loc(who)
	}
	return who
}

// sellerList and buyerList are buy.c:386–464. Seller sort keys mutate the
// actual trades; city offers follow all sorted character offers.
func (w *world) sellerList(where, except int) tradesList {
	var l tradesList
	var count int32
	for _, i := range w.allCharHere(where) {
		if i == except || w.isPrisoner(i) || w.charMoving(i) != 0 {
			continue
		}
		for t := range w.loopTrade(i) {
			if t.kind == tradeSell {
				t.sort = count
				count++
				l.Append(t)
			}
		}
	}
	slices.SortFunc(l, func(a, b *trade) int {
		if a.cost == b.cost {
			return int(a.sort - b.sort)
		}
		return int(a.cost - b.cost)
	})
	for t := range w.loopTrade(where) {
		if t.kind == tradeSell {
			l.Append(t)
		}
	}
	return l
}

func (w *world) buyerList(where, except int) tradesList {
	var l tradesList
	for _, i := range w.allCharHere(where) {
		if i == except || w.isPrisoner(i) || w.charMoving(i) != 0 {
			continue
		}
		for t := range w.loopTrade(i) {
			if t.kind == tradeBuy {
				l.Append(t)
			}
		}
	}
	for t := range w.loopTrade(where) {
		if t.kind == tradeBuy {
			l.Append(t)
		}
	}
	return l
}

// scanTrades is buy.c:635–653.
func (w *world) scanTrades(t *trade, l tradesList) {
	for _, other := range l {
		if t.qty <= 0 {
			break
		}
		if other.item != t.item || other.who == t.who {
			continue
		}
		if t.kind == tradeBuy {
			w.attemptTrade(t, other)
		} else if t.kind == tradeSell {
			w.attemptTrade(other, t)
		}
	}
}

// reduceQty is buy.c:475–503. Locations generate stock/gold and do not
// apply inventory or reserve limits. All arithmetic remains C int width.
func (w *world) reduceQty(t *trade, cost int32) int32 {
	if w.kind(int(t.who)) == tLoc {
		return t.qty
	}
	if t.kind == tradeBuy {
		has := max(int32(w.hasItem(int(t.who), itemGold))-t.haveLeft, 0)
		// G3 #358: a loaded positive-quantity, zero-price BUY still divides
		// by zero, even with no available gold. Keep the original operation,
		// not a guard or suppressed row. C's failure is host-dependent; Go
		// panics. See https://github.com/playbymail/olympia-g3/issues/358.
		return min(t.qty, has/cost)
	}
	if t.kind == tradeSell {
		has := max(int32(w.hasItem(int(t.who), int(t.item)))-t.haveLeft, 0)
		return min(t.qty, has)
	}
	assert(false, "olympia/buy.c", 501, "FALSE")
	return 0
}

// listMarketItems is the complete buy.c:1027–1066 renderer.
func (w *world) listMarketItems(who int, l tradesList, first bool) bool {
	for _, t := range l {
		if t.cloak >= 2 {
			continue
		}
		qty := w.reduceQty(t, t.cost)
		if qty <= 0 {
			continue
		}
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "%5s %*s %7s %6s %9s   %-25s", "{<b>}trade", charField, "who", "price", "qty", "wt/ea", "item{</b>}")
			w.commandOut(who, "%5s %*s %7s %6s %9s   %-25s", "-----", charField, "---", "-----", "---", "-----", "----")
			first = false
		}
		kind := "sell"
		if t.kind == tradeBuy {
			kind = "buy"
		}
		trader := "?"
		if t.cloak == 0 {
			trader = boxCodeLess(int(t.who))
		}
		weight := 0
		if p := w.rpItem(int(t.item)); p != nil {
			weight = int(p.weight)
		}
		w.commandOut(who, "%5s %*s %7s %6s %9s   %-25s", kind, charField, trader,
			commaNum(int(t.cost)), commaNum(int(qty)), commaNum(weight), w.pluralItemBox(int(t.item), int(qty)))
	}
	return first
}

// marketReport is the complete buy.c:1069–1120 report, not settlement.
// The lists retain the actual buyer/seller traversal and saved sort effects.
func (w *world) marketReport(who, where int) {
	w.commandOut(who, "")
	w.commandOut(who, "Market report:")
	w.indent += 3
	flag := true
	for t := range w.loopTrade(where) {
		if t.kind == tradeProduce && t.monthProd != 0 {
			if flag {
				w.commandOut(who, "")
				flag = false
			}
			w.commandOut(who, "%s produces %s on month %d.", w.justName(where), w.pluralItemName(int(t.item), 2), t.monthProd)
		}
	}
	first := true
	l := w.buyerList(where, 0)
	if len(l) > 0 {
		first = w.listMarketItems(who, l, first)
	}
	l = w.sellerList(where, 0)
	if len(l) > 0 {
		first = w.listMarketItems(who, l, first)
	}
	if first {
		w.commandOut(who, "No goods offered for trade.")
	}
	w.indent -= 3
}

// attemptTrade is the complete buy.c:506–632 settlement, using the shared
// inventory/lore/token and output machinery (including its failure behavior).
func (w *world) attemptTrade(buyer, seller *trade) {
	item := int(buyer.item)
	assert(buyer != nil, "olympia/buy.c", 518, "buyer != NULL")
	assert(seller != nil, "olympia/buy.c", 519, "seller != NULL")
	assert(buyer.item == seller.item, "olympia/buy.c", 520, "buyer->item == seller->item")
	if buyer.cost < seller.cost {
		return
	}
	qty := min(w.reduceQty(buyer, seller.cost), w.reduceQty(seller, seller.cost))
	if qty <= 0 {
		return
	}
	cost := seller.cost * qty
	buyer.qty -= qty
	seller.qty -= qty
	assert(buyer.qty >= 0, "olympia/buy.c", 543, "buyer->qty >= 0")
	assert(seller.qty >= 0, "olympia/buy.c", 544, "seller->qty >= 0")
	b, s := int(buyer.who), int(seller.who)
	if w.kind(b) == tLoc {
		w.genItem(s, itemGold, int(cost))
		w.consumeItem(s, item, int(qty))
		if item == itemPot || item == itemBasket {
			w.goldPotBasket += cost
		} else if item == itemOpium {
			w.logWrite(11, "%s earned %s selling opium.", w.boxName(s), goldS(int(cost))) // LOG_SPECIAL
			w.goldOpium += cost
		} else {
			w.goldTrade += cost
		}
	} else if w.kind(s) == tLoc {
		w.consumeItem(b, itemGold, int(cost))
		if w.itemUnique(item) != 0 {
			w.moveItem(s, b, item, int(qty))
		} else {
			w.genItem(b, item, int(qty))
		}
	} else {
		w.moveItem(b, s, itemGold, int(cost))
		w.moveItem(s, b, item, int(qty))
	}
	sellerS, buyerS := "", ""
	if seller.cloak == 0 {
		sellerS = sout(" from %s", w.boxName(s))
	}
	if buyer.cloak == 0 {
		buyerS = sout(" to %s", w.boxName(b))
	}
	if w.kind(b) != tLoc {
		w.commandOut(b, "Bought %s%s for %s.", w.boxNameQty(item, int(qty)), sellerS, goldS(int(cost)))
	}
	if w.kind(s) != tLoc {
		w.commandOut(s, "Sold %s%s for %s.", w.boxNameQty(item, int(qty)), buyerS, goldS(int(cost)))
	}
	where := 0
	if w.kind(b) == tLoc {
		where = b
	} else if w.kind(s) == tLoc {
		where = s
	} else {
		where = w.subloc(b)
	}
	if seller.cloak == 0 && buyer.cloak == 0 {
		w.commandOut(where, "%s bought %s from %s for %s.", w.boxName(b), w.boxNameQty(item, int(qty)), w.boxName(s), goldS(int(cost)))
	} else if seller.cloak != 0 {
		w.commandOut(where, "%s bought %s for %s.", w.boxName(b), w.boxNameQty(item, int(qty)), goldS(int(cost)))
	} else if buyer.cloak != 0 {
		w.commandOut(where, "%s sold %s for %s.", w.boxName(s), w.boxNameQty(item, int(qty)), goldS(int(cost)))
	}
}

// matchTrades is the complete buy.c:656–695 city-arrival dispatcher.
func (w *world) matchTrades(who int) {
	where := w.subloc(who)
	firstBuy, firstSell := true, true
	var sellers, buyers tradesList
	if w.marketHere(who) == 0 {
		return
	}
	assert(w.validBox(who), "olympia/buy.c", 669, "valid_box(who)")
	for t := range w.loopTrade(who) {
		assert(int(t.who) == who, "olympia/buy.c", 671, "t->who == who")
		if t.kind == tradeBuy {
			if firstBuy {
				firstBuy = false
				sellers = w.sellerList(where, who)
			}
			w.scanTrades(t, sellers)
		} else if t.kind == tradeSell {
			if firstSell {
				firstSell = false
				buyers = w.buyerList(where, who)
			}
			w.scanTrades(t, buyers)
		}
	}
}

// matchAllTrades is buy.c:698–729 (phase instrumentation omitted).
func (w *world) matchAllTrades() {
	for where := range w.loopKind(tLoc) {
		if w.marketHere(where) == 0 {
			continue
		}
		sellers := w.sellerList(where, 0)
		buyers := w.buyerList(where, 0)
		if len(buyers) == 0 || len(sellers) == 0 {
			continue
		}
		for _, buyer := range buyers {
			w.scanTrades(buyer, sellers)
		}
	}
}

// checkValidatedTrades is buy.c:735–746. addItem during settlement can
// append more entries: iterate the live list, not a snapshot or fixed length.
func (w *world) checkValidatedTrades() {
	for i := 0; i < w.tradesToCheck.Len(); i++ {
		w.matchTrades(w.tradesToCheck[i])
	}
	w.tradesToCheck.Clear()
}

// listPendingTrades is complete buy.c:1123–1158. Unlike market_report,
// this displays the stored positive-quantity offers without affordability
// calculations, sorting, or updating the original trade records.
func (w *world) listPendingTrades(who, num int) {
	first := true
	assert(w.validBox(num), "olympia/buy.c", 1129, "valid_box(num)")
	for t := range w.loopTrade(num) {
		if t.kind != tradeBuy && t.kind != tradeSell {
			continue
		}
		if first {
			w.commandOut(who, "")
			w.commandOut(who, "Pending trades:")
			w.commandOut(who, "")
			w.indent += 3
			first = false
			w.commandOut(who, "%5s  %7s  %5s   %s", "trade", "price", "qty", "item")
			w.commandOut(who, "%5s  %7s  %5s   %s", "-----", "-----", "---", "----")
		}
		kind := "sell"
		if t.kind == tradeBuy {
			kind = "buy"
		}
		w.commandOut(who, "%5s  %7s  %5s   %s", kind, commaNum(int(t.cost)), commaNum(int(t.qty)), w.boxName(int(t.item)))
	}
	if !first {
		w.indent -= 3
	}
}
