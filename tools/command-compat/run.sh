#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-6a25f2d252ce45bd5adbe3baaccb3be3eed56279}"
WORK="${WORK:-${ROOT}/tmp/command-compat}"
CC="${CC:-cc}"
CASES="${ROOT}/internal/olympia/testdata/command-cases.txt"
EXPECT="${ROOT}/internal/olympia/testdata/command-c-output.txt"
update=0
case "${1:-}" in
"") ;;
--update) update=1 ;;
*) echo "usage: $0 [--update]" >&2; exit 2 ;;
esac
die() { echo "error: $*" >&2; exit 2; }
git -C "${G3_REPO}" cat-file -e "${G3_REV}^{commit}" || die "missing G3 revision"
rm -rf "${WORK}"
mkdir -p "${WORK}/g3" "${WORK}/run1" "${WORK}/run2"
git -C "${G3_REPO}" archive "${G3_REV}" lib olympia | tar -x -C "${WORK}/g3"
echo "G3 revision: ${G3_REV}"

python3 -I - "${WORK}/g3/olympia" "${WORK}" <<'PY'
import re, sys
src, dst = sys.argv[1:]
combat = open(f'{src}/combat.c').read()
with open(f'{dst}/combat-structs.h', 'w') as out:
    start = combat.index('struct fight {')
    end = combat.index('\n};', combat.index('struct wield {')) + 3
    out.write('#line 110 "olympia/combat.c"\n' + combat[start:end] + '\n')
functions = {
    'input.c': ['parse_line', 'remove_comment', 'remove_ctrl_chars', 'find_command',
                'oly_parse_cmd', 'parse_arg', 'oly_parse', 'cmd_shift', 'check_allow',
                'get_command', 'set_state', 'load_command', 'command_done',
                'finish_command', 'do_command', 'init_load_sup', 'exec_precedence',
                'exec_comp', 'sort_load_queue', 'sort_run_queue', 'initial_command_load',
                'min_pri_ready', 'init_wait_list', 'check_all_waits', 'start_phase', 'evening_phase',
                'daily_command_loop', 'process_player_orders', 'interrupt_order', 'process_interrupted_units', 'process_orders'],
    'order.c': ['p_order_head', 'rp_order_head', 'top_order', 'pop_order', 'is_stop_order', 'stop_order', 'flush_unit_orders', 'queue_order', 'prepend_order', 'queue', 'save_player_orders'],
    'swear.c': ['lord', 'player', 'set_loyal', 'reswear_all_sworn', 'set_lord', 'unit_deserts'],
    'move.c': ['check_captain_loses_sailors', 'restore_stack_actions', 'touch_loc_after_move', 'move_stack',
               'departure_message', 'arrival_message', 'discover_road', 'parse_exit_dir',
               'save_v_array', 'restore_v_array', 'move_exit_land', 'land_check',
               'can_move_here', 'can_move_at_outer_level', 'suspend_stack_actions',
               'clear_guard_flag', 'v_move', 'd_move',
               'init_ocean_chars', 'check_ocean_chars',
               'v_exit', 'v_enter', 'v_north', 'v_south', 'v_east', 'v_west'],
    'day.c': ['touch_loc_pl', 'touch_loc', 'storm_owner_touch_loc', 'init_locs_touched',
              'default_garrison_pay', 'maint_cost', 'storm_decay', 'storm_move',
              'near_rocky_coast', 'ship_coastal_damage', 'mine_calamity',
              'inn_calamity', 'random_loc_damage', 'heal_char_sup', 'heal_characters',
              'dogs_bark_at_hidden_chars', 'increment_stone_ring_aura',
              'noncreator_curse_erode', 'daily_events', 'compute_civ_levels',
              'post_production', 'move_city_gold', 'collect_taxes',
              'add_claim_gold', 'add_noble_points', 'add_unformed_sup', 'add_unformed',
              'special_locs_open', 'special_locs_close', 'clear_orders_sent',
              'increment_current_aura', 'decrement_ability_shroud', 'decrement_region_shroud',
              'decrement_meditation_hinder', 'decrement_loc_barrier', 'loyalty_decay',
              'pillage_decay', 'relic_decay', 'hide_mage_decay', 'inn_income', 'temple_income',
              'men_starve', 'unit_maint_cost', 'charge_maint_sup', 'charge_maint_costs', 'animal_deaths'],
    'scry.c': ['notify_loc_shroud'],
    'seed.c': ['int_comp', 'seed_taxes'],
    'stack.c': ['stack_parent', 'stack_leader', 'here_pos', 'stacked_beneath',
                'first_prisoner_pos', 'promote', 'unstack', 'leave_stack', 'stack', 'join_stack',
                'check_prisoner_escape', 'prisoner_escapes', 'prisoner_movement_escape_check',
                'weekly_prisoner_escape_check',
                'drop_stack', 'extract_stacked_unit', 'move_prisoner'],
    'dir.c': ['location_direction', 'exit_distance', 'is_port_city',
              'province_has_port_city', 'summer_uldim_open_now', 'add_province_exit',
              'extra_routes', 'province_exits', 'province_sub_exits', 'subloc_exits',
              'ship_exits', 'building_exits', 'exits_from_loc', 'exits_from_loc_nsew',
              'exits_from_loc_nsew_select', 'has_ocean_access'],
    'combat.c': ['fort_covers', 'auto_attack', 'check_auto_attack_sup', 'check_all_auto_attacks', 'clear_second_waits', 'find_wield', 'wield_s'],
    'loc.c': ['add_char_here', 'all_stack', 'subloc', 'first_character',
              'subloc_here', 'count_loc_structures',
              'in_safe_now', 'province', 'building_owner', 'viewloc', 'all_char_here',
              'in_here_list', 'add_to_here_list', 'remove_from_here_list', 'set_where',
              'region', 'mark_loc_stack_known', 'add_here', 'all_here'],
    'garr.c': ['garrison_here', 'province_admin', 'players_who_rule_here',
               'rank_s', 'garrison_notices', 'garrison_spot_check',
               'garrison_strength', 'ping_garrisons', 'garrison_gold'],
    'build.c': ['province_subloc'],
    'art.c': ['has_auraculum', 'max_eff_aura', 'token_player', 'swear_token_units',
              'melt_token_units', 'add_token_unit_sup', 'add_token_units', 'move_token', 'check_token_units', 'new_suffuse_ring'],
    'display.c': ['with_inventory_string', 'display_with', 'mage_s', 'priest_s',
                  'emperor', 'liner_desc_char', 'liner_desc', 'show_chars_below'],
    'u.c': ['stage', 'olytime_increment', 'olytime_turn_change', 'print_dot',
            'nprovinces', 'greater_region', 'deduct_aura', 'charge_aura', 'check_aura',
            'min', 'cap', 'set_bit', 'test_bit', 'count_stack_any', 'count_any', 'check_skill',
            'contacted', 'char_where', 'char_here', 'test_known', 'loc_hidden', 'has_item', 'loc_depth',
            'lookup', 'gold_s', 'comma_num', 'nice_num', 'loyal_s', 'add_np', 'deduct_np',
            'new_char', 'beast_capturable', 'create_unique_item', 'create_unique_item_alloc',
            'destroy_unique_item', 'add_item', 'sub_item', 'set_known', 'gen_item',
            'consume_item', 'move_item', 'drop_item',
            'stackmate_inheritor', 'take_unit_items', 'put_back_cookie',
            'survive_fatal', 'char_reclaim', 'dead_char_body', 'kill_char', 'hack_unique_item',
            'kill_stack_ocean', 'sink_ship', 'building_collapses',
            'get_rid_of_collapsed_mine', 'add_structure_damage', 'count_loc_char_item',
            'can_pay', 'charge', 'stack_has_item', 'has_use_key', 'stack_has_use_key',
            'stack_sub_item', 'autocharge', 'find_nearest_land', 'max',
            'add_item_weight', 'determine_unit_weights', 'determine_stack_weights',
            'count_stack_move_nobles', 'count_stack_figures', 'count_stack_units', 'bark_dogs', 'clear_temps'],
    'npc.c': ['controlled_humans_here', 'get_exit_dir', 'choose_npc_direction', 'npc_move',
              'create_hades_bandit', 'create_faery_bandit', 'hades_attack_check', 'faery_attack_check'],
    'necro.c': ['auto_undead'],
    'savage.c': ['create_savage', 'call_savage', 'savage_hates', 'savage_hate_here',
                 'auto_savage', 'init_savage_attacks'],
    'lore.c': ['queue_lore', 'scan_char_item_lore'],
    'buy.c': ['investigate_possible_trade', 'market_here', 'seller_comp',
              'seller_list', 'buyer_list', 'reduce_qty', 'attempt_trade',
              'scan_trades', 'match_trades', 'match_all_trades', 'check_validated_trades',
              'find_trade', 'new_trade', 'add_city_trade', 'opium_market_delta',
              'expire_trades', 'loc_trade_sup', 'trade_suffuse_ring', 'location_trades', 'is_tradegood'],
    'c1.c': ['flag_raised', 'clear_wait_parse', 'parse_wait_args', 'wait_leader', 'check_wait_conditions',
             'v_wait', 'd_wait', 'i_wait', 'next_np_turn'],
    'produce.c': ['init_collect_list', 'replenish', 'mine_production', 'location_production'],
    'perm.c': ['is_hostile', 'clear_all_att', 'rp_admit', 'will_admit'],
    'storm.c': ['weather_here', 'update_weather_view_loc_sup', 'init_weather_views',
                'update_weather_view_locs', 'new_storm', 'create_some_storms',
                'natural_weather', 'v_bind_storm', 'd_bind_storm',
                'move_storm', 'move_bound_storms', 'dissipate_storm'],
    'use.c': ['rp_skill_ent', 'exp_level', 'has_skill', 'char_np_total', 'forget_skill', 'skill_school'],
    'report.c': ['char_rep_location'],
    'c2.c': ['v_stop', 'v_format', 'v_notab'],
    'immed.c': ['v_xyzzy', 'v_plugh'],
    'stealth.c': ['v_hide', 'd_hide', 'clear_contacts'],
}
real = {n for names in functions.values() for n in names}
with open(f'{dst}/helpers.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <time.h>
#include "z.h"
#include "oly.h"
#include "rng.h"
#define MAX_PRI 5
#define MAX_BUFFER 50
#define BUF_LEN LEN
ilist load_q[MAX_PRI], run_q;
int cur_pri, cmd_wait = -1, fuzzy_find;
int evening, auto_attack_flag;
ilist wait_list, second_wait_list;
int pledge_backlinks;
int immed_see_all, month_done, show_combat_flag;
#define phase_ran(x) ((void)0)
#define phase_acted(x) ((void)0)
''')
    gm = open(f'{src}/gm.c').read()
    for name in ['gold_garrison', 'gold_combat', 'gold_combat_indep',
                 'gold_pot_basket', 'gold_trade', 'gold_opium', 'gold_taxes',
                 'gold_inn', 'gold_temple', 'gold_provinces']:
        m = re.search(rf'^int {name} = 0;\n', gm, re.M)
        if not m:
            sys.exit(f'no {name} definition in gm.c')
        out.write(f'#line {gm.count(chr(10), 0, m.start()) + 1} "olympia/gm.c"\n{m.group(0)}')
    for file, names in [('cloud.c', ['cloud_region']),
                        ('tunnel.c', ['tunnel_region', 'under_region']),
                        ('quest.c', ['nowhere_region', 'nowhere_loc'])]:
        source = open(f'{src}/{file}').read()
        for name in names:
            m = re.search(rf'^int {name} = 0;\n', source, re.M)
            if not m:
                sys.exit(f'no {name} definition in {file}')
            out.write(f'#line {source.count(chr(10), 0, m.start()) + 1} "olympia/{file}"\n{m.group(0)}')
    for file, names in functions.items():
        source = open(f'{src}/{file}').read()
        # Original data declarations too: no copied or modeled tag/number table.
        if file == 'dir.c':
            start = source.index('#define\tEXTAG_TURN')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/dir.c"\n')
            out.write(source[start:source.index('int max_map_row')])
        if file == 'combat.c':
            out.write('#include "combat-structs.h"\n')
        if file == 'buy.c':
            out.write(re.search(r'^static rng_stream economy_rng;\n', source, re.M).group(0))
            start, end = source.index('#define\tTAG_TURN'), source.index('/*\n *  Explanation of trade->cloak')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/buy.c"\n')
            out.write(source[start:end])
            start, end = source.index('struct\n{\n\tint qty;'), source.index('void\nopium_market_delta(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/buy.c"\n')
            out.write(source[start:end])
        if file == 'savage.c':
            start = source.index('#define\t\tMAX_SAVAGES')
            end = source.index('static int\ncreate_savage(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/savage.c"\n')
            out.write(source[start:end])
        if file == 'display.c':
            start = source.index('int show_display_string =')
            end = source.index('\n\n', start)
            out.write('#line 402 "olympia/display.c"\n' + source[start:end] + '\n')
            for name in ['liner_desc_ship', 'liner_desc_loc', 'liner_desc_road', 'liner_desc_storm']:
                out.write(f'static char *{name}(int n) {{ fprintf(stderr, "out-of-scope C dependency: {name}\\n"); abort(); }}\n')
        if file == 'c1.c':
            start, end = source.index('struct flag_ent {'), source.index('static int\nflag_raised(')
            out.write(source[start:end])
            start, end = source.index('static char *wait_tags[]'), source.index('void\nclear_wait_parse(')
            out.write(source[start:end])
        if file == 'u.c':
            out.write(re.search(r'^static int dot_count = 0;\n', source, re.M).group(0))
            # Whole original RNG prefix: shared persistent entity stream, tags,
            # turn guard and all keyed leaves, not replacement helper draws.
            start = source.index('#define\tETAG_TURN')
            out.write('#include "rng.h"\n')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/u.c"\n')
            out.write(source[start:source.index('void\nkill_stack_ocean(int who)')])
            start, end = source.index('static char *num_s[]'), source.index('char *\nnice_num(')
            out.write(source[start:end])
        if file == 'produce.c':
            # Original complete production tables/macros, not a driver copy.
            start, end = source.index('#define\t\tMOUNTAIN_STONE'), source.index('static void\nreplenish(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/produce.c"\n')
            out.write(source[start:end])
            out.write(re.search(r'^static ilist collectors = NULL;\n', source, re.M).group(0))
        if file == 'day.c':
            start, end = source.index('static rng_stream upkeep_rng;'), source.index('/*\n *  issue #25 endgame Unit B:')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/day.c"\n')
            out.write(source[start:end])
            start, end = source.index('static rng_stream calendar_rng;'), source.index('static int\nnear_rocky_coast(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/day.c"\n')
            out.write(source[start:end])
        if file == 'move.c':
            out.write('#line 8 "olympia/move.c"\n')
            out.write(re.search(r'^static ilist ocean_chars = NULL;\n', source, re.M).group(0))
        if file == 'storm.c':
            start = source.index('static rng_stream weather_rng;')
            end = source.index('int\nv_bind_storm(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/storm.c"\n')
            out.write(source[start:end])
        matches = []
        for name in names:
            pattern = rf'^(?:static )?[^\n(]*\n{name}\(.*?^\}}\n'
            if name == 'queue':
                pattern = r'^void queue\(.*?^\}\n'
            m = re.search(pattern, source, re.M | re.S)
            if not m:
                sys.exit(f'no definition of {name} in {file}')
            line = source.count('\n', 0, m.start()) + 1
            matches.append((m.start(), f'#line {line} "olympia/{file}"\n{m.group(0)}\n'))
        for _, text in sorted(matches):
            out.write(text)
    out.write('''
int drv_load(struct command *c) { return load_command(c); }
void drv_garrison_pay(void) { default_garrison_pay(); }
void drv_sink(int ship) { sink_ship(ship); }
void drv_collapse(int fort) { building_collapses(fort); }
int drv_coast(int where) { return near_rocky_coast(where); }
void drv_coastal(void) { ship_coastal_damage(); }
void drv_mine(int mine) { mine_calamity(mine); }
void drv_inn(int inn) { inn_calamity(inn); }
void drv_loc_damage(void) { random_loc_damage(); }
void drv_heal(void) { heal_characters(); }
void drv_dogs(void) { dogs_bark_at_hidden_chars(); }
void drv_ring(void) { increment_stone_ring_aura(); }
void drv_curse(void) { noncreator_curse_erode(); }
void drv_aura(void) { increment_current_aura(); }
void drv_ability(void) { decrement_ability_shroud(); }
void drv_shroud(void) { decrement_region_shroud(); }
void drv_hinder(void) { decrement_meditation_hinder(); }
void drv_barrier(void) { decrement_loc_barrier(); }
void drv_loyalty(void) { loyalty_decay(); }
void drv_pillage(void) { pillage_decay(); }
void drv_relic(void) { relic_decay(); }
void drv_hide_mage(void) { hide_mage_decay(); }
void drv_income_inn(void) { inn_income(); }
void drv_income_temple(void) { temple_income(); }
int drv_income_roll(int inn, int sub, int low, int high) { return up_income(inn, sub, low, high); }
void drv_maint(void) { charge_maint_costs(); }
void drv_starve(int who, int have) { men_starve(who, have); }
void drv_animal(void) { animal_deaths(); }
int drv_starve_roll(int who, int item) { return up_starve(who, item, 1, 2); }
int drv_animal_roll(int who, int item, int idx) { return up_animal(who, item, idx, 1, 1000); }
int drv_loyal_roll(int who, int sub) { return up_loyal(who, sub, 1, 2); }
void drv_city_gold(void) { move_city_gold(); }
void drv_taxes(void) { collect_taxes(); }
void drv_claim(void) { add_claim_gold(); }
void drv_np(void) { add_noble_points(); }
void drv_unformed(void) { add_unformed(); }
void drv_special_open(void) { special_locs_open(); }
void drv_special_close(void) { special_locs_close(); }
void drv_clear_sent(void) { clear_orders_sent(); }
void drv_replenish(int where, int item, int qty) { replenish(where, item, qty); }
int drv_bark_roll(int who, int hound) { return ent_bark(who, hound); }
void drv_event_days(void) {
 printf("curse day %d\\n", cal_day(0, 1, MONTH_DAYS));
 printf("faery day %d\\n", cal_day(1, MONTH_DAYS/2, MONTH_DAYS));
 printf("dog day %d\\n", cal_day(2, 1, MONTH_DAYS));
}
int drv_heal_roll(int who, int sub, int low, int high) { return up_heal(who, sub, low, high); }
void drv_state(struct command *c, int s) { set_state(c, s, c->pri); }
struct order_list *drv_orders(int pl, int who) { return p_order_head(pl, who); }
void drv_evening(void) { evening_phase(); }
void drv_daily(void) { daily_command_loop(); }
void drv_players(void) { process_player_orders(); }
void drv_interrupted(void) { process_interrupted_units(); }
int drv_sub_item(int who, int item, int qty) { return sub_item(who, item, qty); }
void drv_add_item(int who, int item, int qty) { add_item(who, item, qty); }
void drv_sort_load(int pri) { sort_load_queue(load_q[pri]); }
void drv_sort_run(void) { sort_run_queue(run_q); }
ilist drv_collectors(void) { return collectors; }
int drv_fort_covers(int n) { return fort_covers(n); }
int drv_bandit(int faery, int who, int where) {
 return faery ? create_faery_bandit(who, where) : create_hades_bandit(who, where);
}
int drv_market(int who) { return market_here(who); }
int drv_reduce(struct trade *t, int cost) { return reduce_qty(t, cost); }
void drv_attempt(struct trade *b, struct trade *s) { attempt_trade(b, s); }
void drv_storms(int num, int sk) { create_some_storms(num, sk); }
void drv_storm_decay(void) { storm_decay(); }
void drv_storm_move(void) { storm_move(); }
void drv_move_storm(int storm, int dest) { move_storm(storm, dest); }
void drv_storm_touch(void) { storm_owner_touch_loc(); }
int drv_savage_count(void) { return num_savages; }
int drv_savage_hates(int where) { return savage_hates(where); }
int drv_savage_here(int where) { return savage_hate_here(where); }
void drv_ocean(int who) { ilist_append(&ocean_chars, who); }
ilist drv_oceans(void) { return ocean_chars; }
void drv_entity_rolls(int who, int sub, int where) {
 printf("prison roll %d\\n", ent_prisoner(who, sub));
 printf("land dir99 %d\\n", ent_land_dir(where, 99));
 printf("land dir98 %d\\n", ent_land_dir(where, 98));
 printf("land pick %d\\n", ent_land_pick(where, 3));
}
void drv_inherit_rolls(int from, int item, int qty) {
 printf("take gate %d\\n", ent_take(from, item, 0, 1, 2));
 printf("take amount %d\\n", ent_take(from, item, 1, 0, qty));
 printf("drop gate %d\\n", ent_drop(from, item, 0, 1, 5));
 printf("drop message %d\\n", ent_drop(from, item, 1, 1, 3));
}
trades_list drv_offers(int seller, int where, int except) {
 return seller ? seller_list(where, except) : buyer_list(where, except);
}
char *drv_wait_check(struct command *c) { return check_wait_conditions(c); }
void drv_signal(int who, char *flag) {
 struct flag_ent *p = my_malloc(sizeof(*p));
 p->who = who; p->flag = str_save(flag); flag_ents_append(&flags, p);
}
''')
# Region prefixes retain the original tags, structures, turn guards and whole
# keyed leaf functions. No replaced draw or reset; the same master feeds mint.
for region in ['hades', 'faery']:
    source = open(f'{src}/{region}.c').read()
    with open(f'{dst}/{region}-rng.c', 'w') as out:
        out.write(f'#line 1 "olympia/{region}.c"\n')
        out.write(source[:source.index(f'void\ncreate_{region}(void)')])
        if region == 'hades':
            out.write('#define phase_ran(x) ((void)0)\n#define phase_acted(x) ((void)0)\n')
            for name in ['create_hades_nasty', 'auto_hades_sup', 'auto_hades', 'random_hades_loc']:
                m = re.search(rf'^(?:static )?[^\n(]*\n{name}\(.*?^\}}\n', source, re.M | re.S)
                if not m:
                    sys.exit(f'no definition of {name} in hades.c')
                line = source.count('\n', 0, m.start()) + 1
                out.write(f'#line {line} "olympia/hades.c"\n{m.group(0)}\n')
            out.write('void drv_nasty(int slot) { create_hades_nasty(slot); }\n')
            out.write('void drv_hades_unit(int who) { auto_hades_sup(who); }\n')
            out.write('int drv_nasty_roll(int slot, int sub, int low, int high) { return hads_nasty(slot, sub, low, high); }\n')
        else:
            out.write('#define phase_ran(x) ((void)0)\n#define phase_acted(x) ((void)0)\n')
            for name in ['create_elven_hunt', 'warn_human', 'auto_faery_sup', 'auto_faery']:
                m = re.search(rf'^(?:static )?[^\n(]*\n{name}\(.*?^\}}\n', source, re.M | re.S)
                if not m:
                    sys.exit(f'no definition of {name} in faery.c')
                line = source.count('\n', 0, m.start()) + 1
                out.write(f'#line {line} "olympia/faery.c"\n{m.group(0)}\n')
            out.write('void drv_hunt(int slot) { create_elven_hunt(slot); }\n')
            out.write('void drv_faery_unit(int who) { auto_faery_sup(who); }\n')
            out.write('int drv_hunt_loc(int slot, int retry, int size) { return faer_hunt_loc(slot, retry, 0, size-1); }\n')
            out.write('int drv_hunt_qty(int slot) { return faer_hunt_qty(slot); }\n')
source = open(f'{src}/npc.c').read()
with open(f'{dst}/npc-rng.c', 'w') as out:
    out.write('#line 1 "olympia/npc.c"\n')
    out.write(source[:source.index('int\ncontrolled_humans_here(')])
    out.write('#define phase_ran(x) ((void)0)\n#define phase_acted(x) ((void)0)\n')
    for name in ['auto_unsworn', 'auto_mob', 'auto_bandit', 'queue_npc_orders']:
        m = re.search(rf'^(?:static )?[^\n(]*\n{name}\(.*?^\}}\n', source, re.M | re.S)
        if not m:
            sys.exit(f'no definition of {name} in npc.c')
        line = source.count('\n', 0, m.start()) + 1
        out.write(f'#line {line} "olympia/npc.c"\n{m.group(0)}\n')
    out.write('void drv_unsworn(int who) { auto_unsworn(who); }\n')
    out.write('void drv_mob(int who) { auto_mob(who); }\n')
    out.write('void drv_auto_bandit(int who, int prog) { auto_bandit(who, prog); }\n')
glob = open(f'{src}/glob.c').read()
table = glob[glob.index('struct cmd_tbl_ent cmd_tbl[]'):glob.index('char *kind_s[]')]
routines = set()
for m in re.finditer(r'^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),', table, re.M):
    routines.update(n for n in m.groups() if n != 'NULL')
with open(f'{dst}/stubs.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include "z.h"
#include "oly.h"
#define STUB(n) { fprintf(stderr, "out-of-scope C dependency: %s\\n", n); abort(); }
''')
    for name in sorted(routines - real):
        out.write(f'int {name}(struct command *c) STUB("{name}")\n')
    out.write('''
/* Instrumentation only; coverage artifacts belong to stage 3. */
void coverage_count(int cmd, int result) { (void)cmd; (void)result; }
int immediate, flush_always = 1, pretty_data_files, time_self;
ilist trades_to_check;
''')
PY

case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -w -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib -I "${WORK}" \
 "${HERE}/driver.c" "${WORK}/helpers.c" "${WORK}/stubs.c" "${WORK}/hades-rng.c" "${WORK}/faery-rng.c" "${WORK}/npc-rng.c" \
 olympia/{code,glob,z,sout,io}.c \
 lib/{checked_alloc,ilist,exit_views,rng,rnd,cstrings,item_ents,skill_ents,orders,wait_args,flag_ents,admits,trades,req_ents}.c \
 -o "${WORK}/driver") >"${WORK}/build.log" 2>&1 || die "C build failed: ${WORK}/build.log"
sections=$(grep -c '^section ' "${CASES}")
section_names=() # Bash 3 (Darwin) has no mapfile.
while IFS= read -r line; do section_names+=("$line"); done < <(grep '^section ' "${CASES}")
for run in run1 run2; do
 for ((i=0; i<sections; i++)); do
  mkdir -p "${WORK}/${run}/${i}/log" "${WORK}/${run}/${i}/orders"
  status=0
  (run_limited "${RUN_TIMEOUT}" command-driver "${WORK}/driver" "${CASES}" "$i" "${WORK}/${run}/${i}") \
   >"${WORK}/${run}/stdout" 2>"${WORK}/${run}/stderr" || status=$?
  [ "$status" != 124 ] || die "C section $i timed out"
  ! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stderr" || die "section $i reached a C stub"
  if [[ "${section_names[$i]}" == 'section upkeep '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section upkeep invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected upkeep status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section income '* ]]; then
   expected_status=0
   if [[ "${section_names[$i]}" == 'section income invalid count preserves loop_here valid box assertion' ]]; then
    expected_status=134
   fi
   [ "$status" = "$expected_status" ] || die "unexpected income status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section monthly '* ]]; then
   [ "$status" = 0 ] || die "unexpected monthly aura/loyalty status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section accounts '* ]]; then
   [ "$status" = 0 ] || die "unexpected faction settlement status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section market '* ]]; then
   [ "$status" = 0 ] || die "unexpected market status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section production '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section production invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected production status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section economy '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section economy invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected economy status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section entities '* ]]; then
   [ "$status" = 0 ] || die "unexpected daily entity status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section healing '* || "${section_names[$i]}" == 'section ocean '* ]]; then
   [ "$status" = 0 ] || die "unexpected healing/ocean status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section damage '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section damage invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected daily damage status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section structure '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section structure invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected structure status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section faery '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section faery unlocated hunt preserves loop here assertion') expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected faery status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section weather '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section weather invalid '*|'section weather without province preserves loop assertion'|\
    'section weather mage unowned assertion') expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected weather status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section hades '* || "${section_names[$i]}" == 'section npc route '* ]]; then
   [ "$status" = 0 ] || die "unexpected Hades/NPC status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section savage '* ]]; then
   [ "$status" = 0 ] || die "unexpected savage status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section autonomous '* ]]; then
   [ "$status" = 0 ] || die "unexpected autonomous NPC status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section bandit '* || "${section_names[$i]}" == 'section dispatcher '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section dispatcher invalid program '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected bandit/dispatcher status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section move '* ]]; then
   [ "$status" = 0 ] || die "unexpected movement status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section garrison '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section garrison invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected garrison status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section inventory '* || "${section_names[$i]}" == 'section token '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section inventory asserts '*|'section inventory invalid '*|'section token invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected inventory/token status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section inheritance '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section inheritance invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected inheritance status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section death '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section death invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected death status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section reconcile '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section reconcile invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected reconciliation status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section arrival '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section arrival mismatched trade owner retains exact assertion'|\
    'section arrival move rejects noncharacter root'|\
    'section arrival city matching includes captive descendants'|\
    'section arrival unowned weather mage retains exact assertion') expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected arrival status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section escape '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section escape invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected escape status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section combat boundary '* ]]; then
   [ "$status" = 134 ] && grep -qx 'combat boundary: regular_combat' "${WORK}/${run}/stderr" \
    || die "section $i did not reach actual combat"
  else
   ! grep -q 'combat boundary:' "${WORK}/${run}/stderr" || die "unexpected combat in section $i"
  fi
  case "$status" in
   0|134|139) ;;
   *) die "unexpected C status $status in section $i; see ${WORK}/${run}" ;;
  esac
  {
   printf 'stdout:\n'; cat "${WORK}/${run}/stdout"
   printf 'stderr:\n'; cat "${WORK}/${run}/stderr"
   printf '\nexit %d\n' "$status"
  } >>"${WORK}/${run}/out.txt"
 done
done
cmp "${WORK}/run1/out.txt" "${WORK}/run2/out.txt" || die "C is not deterministic"
echo "C driver: ${sections} sections, $(wc -l <"${WORK}/run1/out.txt") lines, deterministic"
if [ "$update" = 1 ]; then
 cp "${WORK}/run1/out.txt" "${EXPECT}"
else
 diff -u "${EXPECT}" "${WORK}/run1/out.txt"
fi
(cd "${ROOT}" && RUN_TIMEOUT="${RUN_TIMEOUT}" go test -exec "bash ${HERE}/../lib/run-limited.sh" \
 -count=1 -run '^TestCommand' ./internal/olympia)
echo PASS
