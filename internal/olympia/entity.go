// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

import "github.com/mdhender/g3g0/internal/lists"

// The entity structures of olympia/oly.h. Each type is the C struct of the
// same name in camel case (struct entity_char is entityChar), with the same
// fields in the same order, also in camel case (x_loc_info is xLocInfo).
// Nothing else may be added to them: cSizeof (box.go) derives the C layout,
// and so sizeof(struct box), from these declarations, and the entity-compat
// tool holds every struct to the C layout field by field.
//
// Field types keep the C widths, because the engine shows their overflow and
// truncation (a (schar) atoi() in io.c, a short that wraps):
//
//   - schar is int8, uchar is uint8, short is int16 and int is int32 (the C
//     int of every supported host). Plain char is int8: it is signed on the
//     macOS reference builds (arm64 and x86_64); it would be unsigned on
//     Linux arm64.
//   - A char * is a *string, nil for NULL. io.c prints a string only when it
//     is non-NULL and non-empty, but zero_check() treats any non-NULL pointer
//     as data, so NULL and "" are kept apart.
//   - Every list is a lists.List, nil for NULL. A non-nil empty list is not
//     NULL, which zero_check() also sees. Keep one owner per list (see the
//     internal/lists package comment) and pass &field, as the C code passes
//     &list.
//   - A pointer to a sub-structure is a Go pointer, nil for NULL.
//
// Under those rules reflect.Value.IsZero on a structure is zero_check() on
// it, padding aside.

// sparse is a sparse ilist, the C typedef of the same name.
type sparse = lists.IList

// The template list types of lib/lists.h that the entity structures use.
type (
	acceptEntsList = lists.List[*acceptEnt]
	admitsList     = lists.List[*admit]
	itemEntsList   = lists.List[*itemEnt]
	ordersList     = lists.List[*orderList]
	reqEntsList    = lists.List[*reqEnt]
	skillEntsList  = lists.List[*skillEnt]
	tradesList     = lists.List[*trade]
	waitArgsList   = lists.List[*waitArg]

	// cstringsList is cstrings_list, a list of char *. Its elements are Go
	// strings. C's cstrings_lookup and cstrings_rem_value compare pointers,
	// not contents; the engine never calls either.
	cstringsList = lists.List[string]
)

// olytime is olytime, a date in game time.
type olytime struct {
	day            int16 // day of month
	turn           int16 // turn number
	daysSinceEpoch int32 // days since game begin
}

// locInfo is struct loc_info.
type locInfo struct {
	where    int32
	hereList lists.IList
}

// box is struct box, one entity.
type box struct {
	kind  int8
	skind int8
	name  *string

	xLocInfo locInfo
	xPlayer  *entityPlayer
	xChar    *entityChar
	xLoc     *entityLoc
	xSubloc  *entitySubloc
	xItem    *entityItem
	xSkill   *entitySkill
	xGate    *entityGate
	xMisc    *entityMisc
	xDisp    *attEnt

	cmd    *command
	items  itemEntsList // ilist of items held
	trades tradesList   // pending buys/sells

	temp        int32 // scratch space
	outputOrder int32 // for report ordering -- not saved

	xNextKind int32 // link to next entity of same type
	xNextSub  int32 // link to next of same subkind
}

// entityPlayer is struct entity_player.
type entityPlayer struct {
	fullName      *string
	email         *string
	visEmail      *string // address to put in player list
	lastEmail     *string // where did the last orders come from?
	password      *string
	firstTurn     int32      // which turn was their first?
	lastOrderTurn int32      // last turn orders were submitted
	orders        ordersList // ilist of orders for units in this faction
	known         sparse     // visited, lore seen, encountered

	units    lists.IList // what units are in our faction?
	admits   admitsList  // admit permissions list
	unformed lists.IList // nobles as yet unformed

	splitLines int32 // split mail at this many lines
	splitBytes int32 // split mail at this many bytes
	fastStudy  int16 // instant study days available

	noblePoints int16 // how many NP's the player has

	format       int8 // turn report formatting control
	notab        int8 // player can't tolerate tabs
	firstTower   int8 // has player built first tower yet?
	sentOrders   int8 // sent in orders this turn?
	dontRemind   int8 // don't send a reminder
	compuserve   int8 // get Times from CIS
	brokenMailer int8 // quote begin lines

	// not saved:

	publicTurn    int8        // turn is public view
	timesPaid     int8        // only 1 Times credit per player
	swearThisTurn int8        // have we used SWEAR this turn?
	cmdCount      int16       // count of cmds started this turn
	npGained      int16       // np's added this turn -- not saved
	npSpent       int16       // np's lost this turn -- not saved
	deliverLore   lists.IList // show these to player -- not saved
	weatherSeen   sparse      // locs we've viewed the weather
	output        sparse      // units with output -- not saved
	locs          sparse      // locs we touched -- not saved
}

// orderList is struct order_list.
type orderList struct {
	unit int32        // unit orders are for
	l    cstringsList // ilist of orders for unit
}

// acceptEnt is struct accept_ent.
type acceptEnt struct {
	item    int32 // 0 = any item
	fromWho int32 // 0 = anyone, else char or player
	qty     int32 // 0 = any qty
}

// attEnt is struct att_ent.
type attEnt struct {
	neutral lists.IList
	hostile lists.IList
	defend  lists.IList
}

// entityChar is struct entity_char.
type entityChar struct {
	unitItem int32 // unit is made of this kind of item

	health int8
	sick   int8 // 1=character is getting worse

	guard   int8  // character is guarding the loc
	loyKind int8  // LOY_xxx
	loyRate int32 // level with kind of loyalty

	deathTime olytime // when was character killed

	skills skillEntsList // ilist of skills known by char

	moving   int32 // daystamp of beginning of movement
	unitLord int32 // who is our owner?
	prevLord int32 // who was our previous owner?

	contact lists.IList // who have we contacted, also, who has found us

	xCharMagic *charMagic

	prisoner   int8 // is this character a prisoner?
	behind     int8 // are we behind in combat?
	timeFlying int8 // time airborne over ocean
	breakPoint int8 // break point when fighting
	rank       int8 // noble peerage status
	npcProg    int8 // npc program

	attack  int16 // fighter attack rating
	defense int16 // fighter defense rating
	missile int16 // capable of missile attacks?

	// The following are not saved by io.c:

	meltMe    int8           // in process of melting away
	freshHire int8           // don't erode loyalty
	newLord   int8           // got a new lord this turn
	studied   int8           // num days we studied
	accept    acceptEntsList // what we can be given
}

// charMagic is struct char_magic.
type charMagic struct {
	maxAura   int32 // maximum aura level for magician
	curAura   int32 // current aura level for magician
	auraculum int32 // char created an auraculum

	visions sparse // visions revealed
	pledge  int32  // lands are pledged to another
	token   int32  // we are controlled by this art
	fee     int32  // gold/100 wt. to board this ship

	projectCast   int32 // project next cast
	quickCast     int16 // speed next cast
	abilityShroud int16

	hideMage         int8 // hide magician status
	hinderMeditation int8
	magician         int8 // is a magician
	pray             int8 // have prayed
	auraReflect      int8 // reflect aura blast
	hideSelf         int8 // character is hidden
	swearOnRelease   int8 // swear to one who frees us
	knowsWeather     int8 // knows weather magic
	visProtect       int8 // vision protection level
	defaultGarr      int8 // default initial garrison

	mageWorked  int8        // worked this month -- not saved
	ferryFlag   int8        // ferry has tooted its horn -- ns
	pledgedToUs lists.IList // temp -- not saved
}

// skillEnt is struct skill_ent.
type skillEnt struct {
	skill       int32
	daysStudied int32 // days studied * TOUGH_NUM
	experience  int16 // experience level with skill
	know        int8  // SKILL_xxx

	// Not saved:

	expThisMonth int8 // flag for add_skill_experience()
}

// itemEnt is struct item_ent.
type itemEnt struct {
	item int32
	qty  int32
}

// entityLoc is struct entity_loc.
type entityLoc struct {
	provDest     lists.IList
	shroud       int16 // magical scry shroud
	barrier      int16 // magical barrier
	civ          int8  // civilization level (0 = wild)
	hidden       int8  // is location hidden?
	distFromGate int8
	seaLane      int8  // fast ocean travel here; also "tracks" for npc ferries
	next         int32 // temp loc list link
}

// entitySubloc is struct entity_subloc.
type entitySubloc struct {
	teaches   lists.IList // skills location offers
	opiumEcon int32       // addiction level of city
	defense   int32       // defense rating of structure

	loot       int8  // loot & pillage level
	recentLoot int8  // pillaged this month -- not saved
	damage     uint8 // 0=none, 100=fully destroyed
	galleyRam  int8  // galley is fitted with a ram
	shaftDepth int16 // depth of mine shaft
	castleLev  int8  // level of castle improvement

	buildMaterials int8  // fifths of materials we've used
	effortRequired int32 // not finished if nonzero
	effortGiven    int32

	moving   int32 // daystamp of beginning of movement
	capacity int32 // capacity of ship

	nearCities  lists.IList // cities rumored to be nearby
	safe        int8        // safe haven
	major       int8        // major city
	prominence  int8        // prominence of city
	uldimFlag   int8        // Uldim pass
	summerFlag  int8        // Summerbridge
	questLate   int8        // quest decay counter
	tunnelLevel int8        // depth of tunnel

	linkWhen    int8        // month link is open, -1 = never
	linkOpen    int8        // link is open now
	linkTo      lists.IList // where we are linked to
	linkFrom    lists.IList // where we are linked from
	boundStorms lists.IList // storms bound to this ship
}

// entityItem is struct entity_item.
type entityItem struct {
	weight  int16
	landCap int16
	rideCap int16
	flyCap  int16
	attack  int16 // fighter attack rating
	defense int16 // fighter defense rating
	missile int16 // capable of missile attacks?

	isManItem  int8 // unit is a character like thing
	animal     int8 // unit is or contains a horse or an ox
	prominent  int8 // big things that everyone sees
	capturable int8 // ni-char contents are capturable

	pluralName *string
	basePrice  int32 // base price of item for market seeding
	whoHas     int32 // who has this unique item

	xItemMagic *itemMagic
}

// itemMagic is struct item_magic.
type itemMagic struct {
	creator       int32
	regionCreated int32
	lore          int32 // deliver this lore for the item

	curseLoyalty int8 // curse noncreator loyalty
	cloakRegion  int8
	cloakCreator int8
	useKey       int8 // special use action

	mayUse   lists.IList // list of usable skills via this
	mayStudy lists.IList // list of skills studying from this

	projectCast int32 // stored projected cast
	tokenNi     int32 // ni for controlled npc units
	quickCast   int16 // stored quick cast

	auraBonus  int16
	aura       int16 // auraculum aura
	relicDecay int16 // countdown timer

	attackBonus  int8
	defenseBonus int8
	missileBonus int8

	tokenNum    int8 // how many token controlled units
	orbUseCount int8 // how many uses left in the orb

	// Not saved:

	oneTurnUse int8 // flag for one use per turn
}

// entitySkill is struct entity_skill.
type entitySkill struct {
	timeToLearn   int32       // days of study req'd to learn skill
	requiredSkill int32       // skill required to learn this skill
	npReq         int32       // noble points required to learn
	offered       lists.IList // skills learnable after this one
	research      lists.IList // skills researable with this one

	req      reqEntsList // ilist of items required for use or cast
	produced int32       // simple production skill result

	noExp int32 // this skill not rated for experience

	// not saved

	useCount   int32 // times skill used during turn
	lastUseWho int32 // who last used the skill (this turn)
}

// reqEnt is struct req_ent.
type reqEnt struct {
	item    int32 // item required to use
	qty     int32 // quantity required
	consume int8  // REQ_xx
}

// entityGate is struct entity_gate.
type entityGate struct {
	toLoc        int32 // destination of gate
	notifyJumps  int32 // whom to notify
	notifyUnseal int32 // whom to notify
	sealKey      int16 // numeric gate password
	roadHidden   int8  // this is a hidden road or passage
}

// entityMisc is struct entity_misc.
type entityMisc struct {
	display        *string // entity display banner
	npcCreated     int32   // turn peasant mob created
	npcHome        int32   // where npc was created
	npcCookie      int32   // allocation cookie item for us
	summonedBy     int32   // who summoned us?
	saveName       *string // orig name of noble for dead bodies
	oldLord        int32   // who did this dead body used to belong to
	npcMemory      sparse
	onlyVulnerable int32 // only defeatable with this rare artifact
	garrCastle     int32 // castle which owns this garrison
	bindStorm      int32 // ship this storm is bound to

	stormStr  int16 // storm strength
	npcDir    int8  // last direction npc moved
	mineDelay int8  // time until collapsed mine vanishes
	cmdAllow  int8  // unit under restricted control

	opiumDouble int8        // improved opium production -- not saved
	postTxt     []string    // text of posted sign -- not saved; char **, nil for NULL
	stormMove   int32       // next loc storm will move to -- not saved
	garrWatch   lists.IList // units garrison watches for -- not saved
	garrHost    lists.IList // units garrison will attack -- not saved
	garrTax     int32       // garrison taxes collected -- not saved
	garrForward int32       // garrison taxes forwarded -- not saved
}

// waitArg is struct wait_arg.
type waitArg struct {
	tag  int32
	a1   int32
	a2   int32
	flag *string
}

// command is struct command, an in-process command.
type command struct {
	who  int32 // entity this is under (redundant)
	wait int32 // time until completion
	cmd  int32 // index into cmd_tbl

	useSkill      int32 // skill we are using, if any
	useEnt        int32 // index into use_tbl[] for skill usage
	useExp        int32 // experience level at using this skill
	daysExecuting int32 // how long has this command been running

	a, b, c, d, e, f, g, h int32 // command arguments

	line       *string      // original command line
	parsedLine *string      // cut-up line, pointed to by parse
	parse      cstringsList // ilist of parsed arguments

	state         int8 // STATE_LOAD, STATE_RUN, STATE_ERROR, STATE_DONE
	status        int8 // success or failure
	poll          int8 // call finish routine each day?
	pri           int8 // command priority or precedence
	conditional   int8 // 0=none 1=last succeeded 2=last failed
	inhibitFinish int8 // don't call d_xxx

	fuzzy      int8         // command matched fuzzy -- not saved
	secondWait int8         // delay resulting from auto attacks -- ns
	waitParse  waitArgsList // not saved
	debug      int8         // debugging check -- not saved
}

// cmdTblEnt is struct cmd_tbl_ent, a row of cmd_tbl.
type cmdTblEnt struct {
	allow string // who may execute the command
	name  string // name of command

	start     func(*world, *command) int // initiator
	finish    func(*world, *command) int // conclusion
	interrupt func(*world, *command) int // interrupted order

	time int32 // how long command takes
	poll int32 // call finish each day, not just at end
	pri  int32 // command priority or precedence
}

// trade is struct trade.
type trade struct {
	kind      int32 // BUY or SELL
	item      int32
	qty       int32
	cost      int32
	cloak     int32 // don't reveal identity of trader
	haveLeft  int32
	monthProd int32 // month city produces item
	expire    int32 // countdown timer for tradegoods
	who       int32 // redundant -- not saved
	sort      int32 // temp key for sorting -- not saved
}

// admit is struct admit.
type admit struct {
	targ  int32 // char or loc admit is declared for
	sense int32 // 0=default no, 1=all but..
	l     lists.IList

	flag int32 // first time set this turn -- not saved
}
