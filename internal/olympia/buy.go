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
// investigate_possible_trade(), which add_item() calls. City-arrival matching
// below stops at the unported settlement dependency; the remaining economy
// work belongs to #81 and later tasks.
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
// be SIGFPE in C and a runtime panic here. No caller ported so far adds
// gold.
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

// scanTrades is buy.c:635–653. Reaching a candidate pair is a real
// settlement dependency, not permission to silently skip matching (#81).
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

// Settlement (buy.c:506–632) is stage task #81, not a successful stub.
func (w *world) attemptTrade(buyer, seller *trade) {
	panic(notImplemented("attempt_trade"))
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
