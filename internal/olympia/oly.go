// Copyright (c) 2026 Michael D Henderson.
// SPDX-License-Identifier: AGPL-3.0-or-later

package olympia

// The numeric constants #defined in olympia/oly.h, every one of them. Names
// are the C names in camel case (T_char is tChar, sub_ocean is subOcean,
// MAX_BOXES is maxBoxes); NEUTRAL, HOSTILE and DEFEND become attNeutral,
// attHostile and attDefend, and BUY, SELL, PRODUCE and CONSUME become
// tradeBuy, tradeSell, tradeProduce and tradeConsume, which would otherwise
// be too generic. The constants are untyped, as the C macros are. The phase
// coverage enum (PH_xxx) is not ported here: it belongs to --phase-coverage,
// not to the entity model.

// Sizes.
const (
	maxBoxes  = 150000
	monthDays = 30
	numMonths = 8
)

// olyMonth and olyYear are oly_month and oly_year. Promote the short turn
// before subtracting; C division and remainder truncate toward zero. In
// particular turn zero has month -1, year 0 (only the masthead clamps it).
func olyMonth(a olytime) int { return (int(a.turn) - 1) % numMonths }
func olyYear(a olytime) int  { return (int(a.turn) - 1) / numMonths }

// Entity kinds: T_xxx.
const (
	tDeleted  = 0 // forget on save
	tPlayer   = 1
	tChar     = 2
	tLoc      = 3
	tItem     = 4
	tSkill    = 5
	tGate     = 6
	tRoad     = 7
	tDeadchar = 8
	tShip     = 9
	tPost     = 10
	tStorm    = 11
	tUnform   = 12 // unformed noble
	tLore     = 13
	tMax      = 14 // one past highest T_xxx define
)

// Subkinds: sub_xxx.
const (
	subOcean            = 1
	subForest           = 2
	subPlain            = 3
	subMountain         = 4
	subDesert           = 5
	subSwamp            = 6
	subUnder            = 7  // underground
	subFaeryHill        = 8  // gateway to Faery
	subIsland           = 9  // island subloc
	subStoneCir         = 10 // ring of stones
	subMallornGrove     = 11
	subBog              = 12
	subCave             = 13
	subCity             = 14
	subLair             = 15 // dragon lair
	subGraveyard        = 16
	subRuins            = 17
	subBattlefield      = 18
	subEnchForest       = 19 // enchanted forest
	subRockyHill        = 20
	subTreeCircle       = 21
	subPits             = 22
	subPasture          = 23
	subOasis            = 24
	subYewGrove         = 25
	subSandPit          = 26
	subSacredGrove      = 27
	subPoppyField       = 28
	subTemple           = 29
	subGalley           = 30
	subRoundship        = 31
	subCastle           = 32
	subGalleyNotdone    = 33
	subRoundshipNotdone = 34
	subGhostShip        = 35
	subTempleNotdone    = 36
	subInn              = 37
	subInnNotdone       = 38
	subCastleNotdone    = 39
	subMine             = 40
	subMineNotdone      = 41
	subScroll           = 42 // item is a scroll
	subMagic            = 43 // this skill is magical
	subPalantir         = 44
	subAuraculum        = 45
	subTower            = 46
	subTowerNotdone     = 47
	subPlSystem         = 48 // system player
	subPlRegular        = 49 // regular player
	subRegion           = 50 // region wrapper loc
	subPlSavage         = 51 // Savage King
	subPlNpc            = 52
	subMineCollapsed    = 53
	subNi               = 54 // ni=noble_item
	subUndead           = 55 // undead lord
	subDeadBody         = 56 // dead noble's body
	subFog              = 57
	subWind             = 58
	subRain             = 59
	subHadesPit         = 60
	subArtifact         = 61
	subPlSilent         = 62
	subNpcToken         = 63 // npc group control art
	subGarrison         = 64 // npc group control art
	subCloud            = 65 // cloud terrain type
	subRaft             = 66 // raft made out of flotsam
	subRaftNotdone      = 67
	subSuffuseRing      = 68
	subRelic            = 69 // 400 series artifacts
	subTunnel           = 70
	subSewer            = 71
	subChamber          = 72
	subTradegood        = 73
)

const (
	subMax = 74 // one past highest sub_
)

// Items: item_xxx.
const (
	itemGold = 1
)

const (
	itemPeasant        = 10
	itemWorker         = 11
	itemSoldier        = 12
	itemArcher         = 13
	itemKnight         = 14
	itemEliteGuard     = 15
	itemPikeman        = 16
	itemBlessedSoldier = 17
	itemGhostWarrior   = 18
	itemSailor         = 19
	itemSwordsman      = 20
	itemCrossbowman    = 21
	itemEliteArch      = 22
	itemAngryPeasant   = 23
	itemPirate         = 24
	itemElf            = 25
	itemSpirit         = 26
)

const (
	itemCorpse    = 31
	itemSavage    = 32
	itemSkeleton  = 33
	itemBarbarian = 34
)

const (
	itemWildHorse   = 51
	itemRidingHorse = 52
	itemWarmount    = 53
	itemPegasus     = 54
	itemNazgul      = 55
)

const (
	itemFlotsam        = 59
	itemBatteringRam   = 60
	itemCatapult       = 61
	itemSiegeTower     = 62
	itemRatspiderVenom = 63
	itemLanaBark       = 64
	itemAviniaLeaf     = 65
	itemSpinyRoot      = 66
	itemFarrenstone    = 67
	itemYew            = 68
	itemElfstone       = 69
	itemMallornWood    = 70
	itemPretusBones    = 71
	itemLongbow        = 72
	itemPlate          = 73
	itemLongsword      = 74
	itemPike           = 75
	itemOx             = 76
	itemLumber         = 77
	itemStone          = 78
	itemIron           = 79
	itemLeather        = 80
	itemRatspider      = 81
	itemMithril        = 82
	itemGateCrystal    = 83
	itemBlankScroll    = 84
	itemCrossbow       = 85
	itemFish           = 87
	itemOpium          = 93
	itemBasket         = 94 // woven basket
	itemPot            = 95 // clay pot
	itemTaxCookie      = 96
	itemDrum           = 98
	itemHide           = 99
	itemMobCookie      = 101
	itemLead           = 102
)

const (
	itemGlue = 261
)

const (
	itemCentaur      = 271
	itemMinotaur     = 272
	itemUndeadCookie = 273
	itemFogCookie    = 274
	itemWindCookie   = 275
	itemRainCookie   = 276
	itemMageMenial   = 277 // mage menial labor cookie
	itemSpider       = 278 // giant spider
	itemRat          = 279 // horde of rats
	itemLion         = 280
	itemBird         = 281 // giant bird
	itemLizard       = 282
	itemBandit       = 283
	itemChimera      = 284
	itemHarpie       = 285
	itemDragon       = 286
	itemOrc          = 287
	itemGorgon       = 288
	itemWolf         = 289
	itemOrb          = 290
	itemCyclops      = 291
	itemGiant        = 292
	itemFaery        = 293
	itemPettyThief   = 294
	itemHound        = 295
)

// Lore sheets: lore_xxx.
const (
	loreSkeletonNpcToken  = 931
	loreOrcNpcToken       = 932
	loreUndeadNpcToken    = 933
	loreSavageNpcToken    = 934
	loreBarbarianNpcToken = 935
	loreOrb               = 936
	loreFaeryStone        = 937
	loreBarbarianKill     = 938
	loreSavageKill        = 939
	loreUndeadKill        = 940
	loreOrcKill           = 941
	loreSkeletonKill      = 942
)

// Skills: sk_xxx.
const (
	skShipcraft    = 600
	skPilotShip    = 601
	skShipbuilding = 602
	skFishing      = 603
)

const (
	skCombat       = 610
	skSurviveFatal = 611
	skFightToDeath = 612
	skMakeCatapult = 613
	skDefense      = 614
	skArchery      = 615
	skSwordplay    = 616
	skWeaponsmith  = 617
)

const (
	skStealth    = 630
	skPettyThief = 631
	skSpyInv     = 632 // determine char inventory
	skSpySkills  = 633 // determine char skill
	skSpyLord    = 634 // determine char's lord
	skHideLord   = 635
	skFindRich   = 636
	skTorture    = 637
	skHideSelf   = 638
	skSneakBuild = 639
)

const (
	skBeast         = 650
	skBirdSpy       = 651
	skCaptureBeasts = 652
	skUseBeasts     = 653 // use beasts in battle
	skBreedBeasts   = 654
	skCatchHorse    = 655
	skTrainWild     = 656
	skTrainWarmount = 657
	skSummonSavage  = 658
	skKeepSavage    = 659
	skBreedHound    = 661
)

const (
	skPersuasion   = 670
	skBribeNoble   = 671
	skPersuadeOath = 672
	skRaiseMob     = 673
	skRallyMob     = 674
	skInciteMob    = 675
	skTrainAngry   = 676
)

const (
	skConstruction = 680
	skMakeSiege    = 681
	skQuarryStone  = 682
)

const (
	skAlchemy      = 690
	skBrewHeal     = 691
	skRecordSkill  = 692
	skExtractVenom = 693 // from ratspider
	skBrewSlave    = 694 // potion of slavery
	skCollectElem  = 695
	skBrewDeath    = 696
	skLeadToGold   = 697
)

const (
	skForestry       = 700
	skMakeRam        = 701 // make battering ram
	skHarvestLumber  = 702
	skHarvestYew     = 703
	skCollectFoliage = 704
	skHarvestMallorn = 705
	skHarvestOpium   = 706
	skImproveOpium   = 707
)

const (
	skMining      = 720
	skMineIron    = 721
	skMineGold    = 722
	skMineMithril = 723
)

const (
	skTrade      = 730
	skCloakTrade = 731
	skFindSell   = 732
	skFindBuy    = 733
)

const (
	skReligion      = 750
	skRevealVision  = 751
	skLastRites     = 752
	skPray          = 753
	skResurrect     = 754
	skRemoveBless   = 755
	skVisionProtect = 756
)

const (
	skBasic        = 800
	skMeditate     = 801
	skMageMenial   = 802 // menial labor for mages
	skAppearCommon = 803
	skViewAura     = 804
	skHeal         = 805
	skWriteBasic   = 806
	skRevealMage   = 807 // reveal abilities of mage
	skTapHealth    = 808
	skShroudAbil   = 809 // ability shroud
	skDetectAbil   = 811 // detect ability scry
	skDispelAbil   = 812 // dispel ability shroud
	skAdvMed       = 813 // advanced meditation
	skHinderMed    = 814 // hinder meditation
)

const (
	skWeather      = 820
	skFierceWind   = 821
	skBindStorm    = 822
	skWriteWeather = 823
	skSummonWind   = 824
	skSummonRain   = 825
	skSummonFog    = 826
	skDirectStorm  = 827
	skDissipate    = 828
	skRenewStorm   = 829
	skLightning    = 831
	skSeizeStorm   = 832
	skDeathFog     = 833
)

const (
	skScry          = 840
	skScryRegion    = 841
	skWriteScry     = 842
	skShroudRegion  = 843
	skDispelRegion  = 844 // dispel region shroud
	skBarLoc        = 845 // create location barrier
	skUnbarLoc      = 846
	skLocateChar    = 847
	skDetectScry    = 848 // detect region scry
	skProjCast      = 849 // project next cast
	skSaveProj      = 851 // save projected cast
	skBanishCorpses = 852
)

const (
	skGate         = 860
	skDetectGates  = 861
	skJumpGate     = 862
	skWriteGate    = 863
	skSealGate     = 864
	skUnsealGate   = 865
	skNotifyUnseal = 866
	skRemSeal      = 867 // forcefully unseal gate
	skRevealKey    = 868
	skNotifyJump   = 869
	skTeleport     = 871
	skRevJump      = 872
)

const (
	skArtifact      = 880
	skForgeAura     = 881 // forge auraculum
	skWriteArt      = 882
	skForgeWeapon   = 883
	skForgeArmor    = 884
	skForgeBow      = 885
	skCurseNoncreat = 886 // curse noncreator loyalty
	skShowArtCreat  = 887 // learn who created art
	skShowArtReg    = 888 // learn where art created
	skDestroyArt    = 889
	skCloakCreat    = 891
	skCloakReg      = 892
	skRemArtCloak   = 893 // dispel artifact cloaks
	skForgePalantir = 894
)

const (
	skNecromancy     = 900
	skRaiseCorpses   = 901 // summon undead corpses
	skSummonGhost    = 902 // summon ghost warriors
	skWriteNecro     = 903
	skUndeadLord     = 904 // summon undead unit
	skRenewUndead    = 905
	skBanishUndead   = 906
	skEatDead        = 907
	skAuraBlast      = 908
	skAbsorbBlast    = 909
	skTranscendDeath = 911
)

const (
	skAdvSorcery   = 920
	skTrance       = 921
	skTeleportItem = 922
)

// Dead skills.
const (
	skQuickCast = 999 // speed next cast
	skSaveQuick = 998 // save speeded cast
	skAddRam    = 997 // add ram to galley
)

// NPC programs: PROG_xxx.
const (
	progBandit        = 1 // wilderness spice
	progSublocMonster = 2
	progNpcToken      = 3
	progFaeryBandit   = 4
	progHadesBandit   = 5
)

// Special item uses: use_xxx.
const (
	useDeathPotion   = 1
	useHealPotion    = 2
	useSlavePotion   = 3
	usePalantir      = 4
	useProjCast      = 5 // stored projected cast
	useQuickCast     = 6 // stored cast speedup
	useDrum          = 7 // beat savage's drum
	useFaeryStone    = 8 // Faery gate opener
	useOrb           = 9 // crystal orb
	useBarbarianKill = 10
	useSavageKill    = 11
	useCorpseKill    = 12
	useOrcKill       = 13
	useSkeletonKill  = 14
	useBtaSkull      = 15
	useKeyMax        = 15 // highest use key (issue #181)
)

// Directions: DIR_xxx.
const (
	dirN    = 1
	dirE    = 2
	dirS    = 3
	dirW    = 4
	dirUp   = 5
	dirDown = 6
	dirIn   = 7
	dirOut  = 8
	maxDir  = 9 // one past highest direction
)

// Location depths: LOC_xxx.
const (
	locRegion   = 1 // top most continent/island group
	locProvince = 2 // main location area
	locSubloc   = 3 // inner sublocation
	locBuild    = 4 // building, structure, etc.
)

// Loyalty kinds: LOY_xxx.
const (
	loyUnchanged = -1
	loyUnsworn   = 0
	loyContract  = 1
	loyOath      = 2
	loyFear      = 3
	loyNpc       = 4
	loySummon    = 5
)

// Experience levels: exp_xxx.
const (
	expNovice     = 1 // apprentice
	expJourneyman = 2
	expTeacher    = 3
	expMaster     = 4
	expGrand      = 5 // grand master
)

// Attitudes: ATT_NONE, NEUTRAL, HOSTILE and DEFEND. The last three get
// an att prefix here.
const (
	attNone    = 0 // no attitude -- default
	attNeutral = 1 // explicitly neutral
	attHostile = 2 // attack on sight
	attDefend  = 3 // defend if attacked
)

// Skill knowledge: SKILL_xxx.
const (
	skillDont     = 0 // don't know the skill
	skillLearning = 1 // in the process of learning it
	skillKnow     = 2 // know it
)

// Item requirements: REQ_xxx.
const (
	reqNo  = 0 // don't consume item
	reqYes = 1 // consume item
	reqOr  = 2 // or with next
)

// Command states: STATE_xxx.
const (
	stateDone  = 0
	stateLoad  = 1
	stateRun   = 2
	stateError = 3
)

// Trade kinds: BUY, SELL, PRODUCE and CONSUME, with a trade prefix here.
const (
	tradeBuy     = 1
	tradeSell    = 2
	tradeProduce = 3
	tradeConsume = 4
)

// do_command() outcomes counted by coverage_count(): COV_xxx.
const (
	covRejected = 0 // check_allow() failed, or unrecognized
	covUnimpl   = 1 // cmd_tbl start is NULL
	covOk       = 2 // start returned TRUE
	covFail     = 3 // start returned FALSE
	covMax      = 4
)

// v_use()/v_use_item() outcomes: UCOV_xxx.
const (
	ucovChecks = 0 // a shared check failed before dispatch
	ucovOk     = 1 // the skill or item use returned TRUE
	ucovFail   = 2 // the skill or item use returned FALSE
	ucovMax    = 3
)

// Field length for box_code_less: CHAR_FIELD.
const (
	charField = 6 // field length for box_code_less
)

// Maximum line length for posts and messages: MAX_POST.
const (
	maxPost = 60 // max line length for posts and messages
)

// style() tags: STYLE_xxx. The default style is 0 (regular).
const (
	styleText = 1
	styleHTML = 2
	stylePrev = -1
)

// Relics: RELIC_xxx.
const (
	relicThrone   = 401
	relicCrown    = 402
	relicBtaSkull = 403
)
