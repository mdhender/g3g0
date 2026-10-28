#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-e733acb63bacb21167a1f2589d8dcdc023364896}"
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
    'order.c': ['p_order_head', 'rp_order_head', 'top_order', 'pop_order', 'is_stop_order', 'stop_order', 'flush_unit_orders', 'queue_order', 'prepend_order', 'queue', 'save_player_orders',
                'orders_template_sup', 'orders_other', 'orders_template', 'list_order_templates'],
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
              'men_starve', 'unit_maint_cost', 'charge_maint_sup', 'charge_maint_costs', 'animal_deaths',
              'collapsed_mine_decay', 'ghost_warrior_decay', 'corpse_decay', 'dead_body_rot',
              'auto_drop', 'link_decay', 'quest_decay', 'post_month'],
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
              'exits_from_loc_nsew_select', 'has_ocean_access', 'list_exit_extras',
              'list_exits_sup', 'list_road_sup', 'list_exits', 'list_sailable_routes'],
    'combat.c': ['fort_covers', 'auto_attack', 'check_auto_attack_sup', 'check_all_auto_attacks', 'clear_second_waits', 'find_wield', 'wield_s'],
    'loc.c': ['add_char_here', 'all_stack', 'subloc', 'first_character',
              'subloc_here', 'count_loc_structures',
              'in_safe_now', 'province', 'building_owner', 'viewloc', 'all_char_here',
              'in_here_list', 'add_to_here_list', 'remove_from_here_list', 'set_where',
              'region', 'mark_loc_stack_known', 'add_here', 'all_here'],
    'garr.c': ['garrison_here', 'province_admin', 'players_who_rule_here',
               'rank_s', 'garrison_notices', 'garrison_spot_check',
               'garrison_strength', 'ping_garrisons', 'garrison_gold', 'top_ruler',
               'nprovs_to_rank', 'find_kings', 'lower_noble_rank',
               'det_noble_rank_sup', 'determine_noble_ranks', 'may_rule_here',
               'garr_own_s', 'garrison_summary'],
    'build.c': ['province_subloc'],
    'art.c': ['has_auraculum', 'max_eff_aura', 'token_player', 'swear_token_units',
              'melt_token_units', 'add_token_unit_sup', 'add_token_units', 'move_token', 'check_token_units', 'new_suffuse_ring'],
    'display.c': ['loc_inside_string', 'show_loc_barrier', 'safe_haven_s', 'ship_cap_s',
                  'loc_civ_s', 'show_loc_header',
                  'display_owner', 'incomplete_string', 'liner_desc_ship', 'liner_desc_loc',
                  'liner_desc_road', 'liner_desc_storm',
                  'with_inventory_string', 'display_with', 'mage_s', 'priest_s',
                  'emperor', 'liner_desc_char', 'liner_desc', 'highlight_units',
                  'show_chars_below', 'show_chars_below_highlight', 'show_owner_stack',
                  'show_chars_here', 'show_inner_locs', 'show_sublocs_here', 'show_ships_here',
                  'show_nearby_cities', 'show_loc_skills', 'show_loc_posts',
                  'show_weather', 'show_loc_ruler', 'show_loc_stats', 'show_loc'],
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
            'ship_cap', 'ship_weight', 'first_char_here',
            'count_stack_move_nobles', 'count_stack_figures', 'count_stack_units', 'bark_dogs', 'clear_temps', 'ordinal', 'weeks', 'knum'],
    'npc.c': ['controlled_humans_here', 'get_exit_dir', 'choose_npc_direction', 'npc_move',
              'create_hades_bandit', 'create_faery_bandit', 'hades_attack_check', 'faery_attack_check'],
    'necro.c': ['auto_undead'],
    'savage.c': ['create_savage', 'call_savage', 'savage_hates', 'savage_hate_here',
                 'auto_savage', 'init_savage_attacks'],
    'lore.c': ['lore_function', 'deliver_lore_sheet', 'np_req_s', 'out_skill_line',
               'deliver_skill_lore', 'deliver_lore', 'queue_lore', 'lore_comp',
               'show_lore_sheets', 'gm_show_all_skills', 'scan_char_skill_lore', 'scan_char_item_lore'],
    'buy.c': ['investigate_possible_trade', 'market_here', 'seller_comp',
              'seller_list', 'buyer_list', 'reduce_qty', 'attempt_trade',
              'scan_trades', 'match_trades', 'match_all_trades', 'check_validated_trades',
              'find_trade', 'new_trade', 'add_city_trade', 'opium_market_delta',
              'expire_trades', 'loc_trade_sup', 'trade_suffuse_ring', 'location_trades', 'is_tradegood', 'list_pending_trades',
              'list_market_items', 'market_report'],
    'c1.c': ['flag_raised', 'clear_wait_parse', 'parse_wait_args', 'wait_leader', 'check_wait_conditions',
             'v_wait', 'd_wait', 'i_wait', 'next_np_turn', 'print_hiring_status', 'print_unformed'],
    'produce.c': ['init_collect_list', 'replenish', 'mine_production', 'location_production'],
    'perm.c': ['is_hostile', 'clear_all_att', 'rp_admit', 'will_admit',
               'admit_comp', 'print_admit_sup', 'print_admit', 'print_att_sup', 'print_att'],
    'storm.c': ['weather_here', 'can_see_weather_here', 'update_weather_view_loc_sup', 'init_weather_views',
                'update_weather_view_locs', 'new_storm', 'create_some_storms',
                'natural_weather', 'v_bind_storm', 'd_bind_storm',
                'move_storm', 'move_bound_storms', 'dissipate_storm', 'storm_report'],
    'use.c': ['rp_skill_ent', 'exp_level', 'has_skill', 'char_np_total', 'forget_skill', 'skill_school',
              'exp_s', 'rep_skill_comp', 'flat_skill_comp', 'list_skill_sup',
              'list_skills', 'fractional_skill_qualifier', 'list_partial_skills'],
    'report.c': ['char_rep_location', 'output_order_comp', 'sort_for_output',
                 'determine_output_order', 'inv_item_comp', 'extra_item_info',
                 'show_unclaimed', 'player_ent_info', 'show_char_inventory',
                 'show_carry_capacity', 'show_item_skills_sup', 'show_item_skills',
                 'char_rep_stack_info', 'collect_pledge_backlinks', 'show_pledged',
                 'prisoner_health', 'char_rep_prisoners', 'char_rep_health',
                 'char_rep_combat', 'char_rep_misc', 'char_rep_magic',
                 'char_rep_sup', 'character_report', 'turn_end_loc_reports', 'player_banner',
                 'sum_fighters', 'strip_leading_stupid_word', 'stupid_word', 'prev_word',
                 'summary_trunc_name', 'unit_summary_sup', 'unit_summary', 'loc_ind_s',
                 'loc_stack_catchup', 'loc_stack_rep_sup', 'loc_stack_report',
                 'player_report_sup', 'stack_capacity_report', 'player_report',
                 'rep_player', 'rep_char', 'rep_loc', 'inc', 'gen_include_sup',
                 'gen_include_section', 'report_account_sup', 'report_account'],
    'summary.c': ['rank_comp', 'ranking', 'top_rank', 'out_rank_mine_top',
                  'out_ranking', 'collect_game_totals', 'summary_units', 'summary_skills',
                  'summary_spells', 'summary_men', 'summary_gold', 'summary_land_owned',
                  'summary_provinces', 'summary_report'],
    'gm.c': ['pct', 'skill_use_comp', 'gm_show_skill_use_counts', 'skills_known_comp',
             'gm_show_skills_known', 'gm_show_interesting_attributes', 'gm_list_animate_items',
             'gm_show_gold', 'gm_show_control_arts', 'gm_count_stuff', 'gm_show_gate_stats',
             'gm_show_locs_visited', 'gm_loyalty_stats', 'region_occupy_comp', 'gm_land_stats',
             'wealth_list_comp', 'gm_faction_wealth', 'nobles_list_comp', 'gm_nobles_list',
             'gm_player_details', 'list_all_notices', 'list_all_items', 'gm_report'],
    'c2.c': ['v_stop', 'v_format', 'v_notab', 'open_times', 'close_times'],
    'immed.c': ['v_xyzzy', 'v_plugh'],
    'stealth.c': ['v_hide', 'd_hide', 'clear_contacts'],
    'main.c': ['write_totimes', 'write_email', 'list_a_player', 'write_player_list',
               'write_forward_sup', 'write_forwards', 'write_faction_sup', 'write_factions'],
    'add.c': ['fetch_inp', 'get_city_id', 'pick_starting_city', 'add_new_player',
              'make_new_players_sup', 'make_new_players', 'new_player_banners',
              'show_new_char_locs', 'new_player_report', 'new_char_report',
              'new_order_templates', 'new_player_list_sup', 'new_player_list', 'add_new_players'],
}
real = {n for names in functions.values() for n in names}
with open(f'{dst}/helpers.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <time.h>
#include <dirent.h>
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
                 'gold_inn', 'gold_temple', 'gold_provinces', 'gold_common_magic',
                 'gold_lead_to_gold', 'gold_times', 'gold_petty_thief', 'gold_pillage', 'gold_ferry']:
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
        if file == 'c2.c':
            start, end = source.index('#define DASH_LINE'), source.index('char *months[]')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/c2.c"\n')
            out.write(source[start:end])
        if file == 'add.c':
            for pattern in [r'^ilist new_players = NULL;\n', r'^static ilist new_chars = NULL;\n']:
                out.write(re.search(pattern, source, re.M).group(0))
        if file == 'dir.c':
            start = source.index('#define\tEXTAG_TURN')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/dir.c"\n')
            out.write(source[start:source.index('int max_map_row')])
        if file == 'combat.c':
            out.write('#include "combat-structs.h"\n')
        if file == 'summary.c':
            for pattern in [r'^#define\s+SHOW_TOP\s+FALSE\n', r'^static ilist ranks = NULL;\n',
                            r'^static int nplayers;\n', r'^static int ncontrolled;\n', r'^static int nother;\n']:
                out.write(re.search(pattern, source, re.M).group(0))
        if file == 'report.c':
            start, end = source.index('#define\tTRUNC_NAME'), source.index('static char *\nstrip_leading_stupid_word(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/report.c"\n')
            out.write(source[start:end])
            for pattern in [r'^static int sum_\w+;\n', r'^static char \*loyal_chars = .*;\n',
                            r'^static int loc_stack_explain;\n']:
                for m in re.finditer(pattern, source, re.M):
                    out.write(m.group(0))
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
            out.write(re.search(r'^int show_loc_no_header = FALSE;\n', source, re.M).group(0))
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
            if file == 'main.c' and name != 'list_a_player':
                pattern = rf'^void {name}\(.*?^\}}\n'
            if name == 'queue':
                pattern = r'^void queue\(.*?^\}\n'
            if name == 'rank_comp':
                pattern = r'^static\nint rank_comp\(.*?^\}\n'
            if name == 'get_city_id':
                pattern = r'^static int get_city_id\(.*?^\}\n'
            m = re.search(pattern, source, re.M | re.S)
            if name == 'list_skill_sup':
                # The first definition is disabled by #if 0; use the real #else.
                m = list(re.finditer(pattern, source, re.M | re.S))[-1]
            if not m:
                sys.exit(f'no definition of {name} in {file}')
            line = source.count('\n', 0, m.start()) + 1
            matches.append((m.start(), f'#line {line} "olympia/{file}"\n{m.group(0)}\n'))
        for _, text in sorted(matches):
            out.write(text)
    out.write('''
void drv_gm_economy(char *part, int pl) {
 if (!strcmp(part,"attributes")) gm_show_interesting_attributes(pl);
 else if (!strcmp(part,"animate")) gm_list_animate_items(pl);
 else if (!strcmp(part,"gold")) gm_show_gold(pl);
 else if (!strcmp(part,"control")) gm_show_control_arts(pl);
 else if (!strcmp(part,"structures")) gm_count_stuff(pl);
 else if (!strcmp(part,"items")) list_all_items(pl);
 else if (!strcmp(part,"use")) gm_show_skill_use_counts(pl);
 else if (!strcmp(part,"skills")) gm_show_skills_known(pl);
 else if (!strcmp(part,"gates")) gm_show_gate_stats(pl);
 else if (!strcmp(part,"locations")) gm_show_locs_visited(pl);
 else if (!strcmp(part,"loyalty")) gm_loyalty_stats(pl);
 else if (!strcmp(part,"regions")) gm_land_stats(pl);
 else abort();
}
int drv_gm_pct(int n, int d) { return pct(n,d); }
void drv_lore_sheet(int who, int num, int display_num) { deliver_lore_sheet(who, num, display_num); }
char *drv_summary_trunc_name(int who) { return summary_trunc_name(who); }
int drv_sum_fighters(int who) { return sum_fighters(who); }
void drv_unit_summary_sup(int pl, int who) { unit_summary_sup(pl, who); }
void drv_summary_totals(void) { printf("summary totals %d %d %d %d %d\\n", sum_gold, sum_peas, sum_work, sum_sail, sum_fight); }
char *drv_loc_ind(int where) { return loc_ind_s(where); }
void drv_loc_catchup(int pl, int where) { loc_stack_catchup(pl, where); }
void drv_loc_stack_sup(int pl, int where, int who) { loc_stack_rep_sup(pl, where, who); }
int drv_loc_barrier(int who, int where) { return show_loc_barrier(who, where); }
void drv_chars_here(int who, int where) { show_chars_here(who, where); }
void drv_inner_locs(int who, int where) { show_inner_locs(who, where); }
void drv_sublocs_here(int who, int where) { show_sublocs_here(who, where); }
void drv_ships_here(int who, int where) { show_ships_here(who, where); }
void drv_ancillary(int mode, int who, int where) {
 switch(mode) {
 case 0: show_nearby_cities(who, where); break;
 case 1: show_loc_skills(who, where); break;
 case 2: show_weather(who, where); break;
 case 3: show_loc_ruler(who, where); break;
 case 4: show_loc_stats(who, where); break;
 }
}
void drv_chars_below(int who, int n, int depth) { show_chars_below_highlight(who, n, depth); }
char *drv_highlight(int who, int n, int depth) { return highlight_units(who, n, depth); }
char *drv_fractional_skill_qualifier(struct skill_ent *p) { return fractional_skill_qualifier(p); }
void drv_char_detail(char *section, int who, int num) {
 if(!strcmp(section,"stack_info")) char_rep_stack_info(who,num);
 else if(!strcmp(section,"pledged")) show_pledged(who,num);
 else if(!strcmp(section,"prisoners")) char_rep_prisoners(who,num);
 else if(!strcmp(section,"health")) char_rep_health(who,num);
 else if(!strcmp(section,"combat")) char_rep_combat(who,num);
 else if(!strcmp(section,"detail_misc")) char_rep_misc(who,num);
 else if(!strcmp(section,"magic")) char_rep_magic(who,num);
 else abort();
}
char *drv_prisoner_health(int who) { return prisoner_health(who); }
void drv_route_view(int who, struct exit_view *v, int mode) {
 char first[LEN] = "Direct routes:";
 extern int second_indent;
 int old_indent = indent;
 if(mode == 1) first[0] = 0;
 if(v->road) list_road_sup(who, 0, v, mode == 0 ? NULL : first);
 else list_exits_sup(who, 0, v, mode == 0 ? NULL : first);
 printf("route first %d indent %d second %d\\n", first[0] != 0, indent, second_indent);
 indent = old_indent;
}
char *drv_fetch(FILE *fp) { return fetch_inp(fp); }
int drv_admission(char *acct, FILE *fp) { return make_new_players_sup(acct, fp); }
void drv_make_players(void) { make_new_players(); }
int drv_load(struct command *c) { return load_command(c); }
void drv_summary_provinces(void) { summary_provinces(); }
void drv_summary_parts(int mode, int total) {
 switch(mode) {
 case 0: collect_game_totals(); break;
 case 1: out_ranking("Test ranks:",total); break;
 case 2: summary_units(); break;
 case 3: summary_men(); break;
 case 4: summary_gold(); break;
 case 5: summary_land_owned(); break;
 case 6: summary_skills(); break;
 case 7: summary_spells(); break;
 }
}
void drv_summary_state(void) {
 int n;
 printf("game totals %d %d %d\\n",nplayers,ncontrolled,nother);
 printf("ranks allocated %d count %d\\n",ranks!=NULL,ilist_len(ranks));
 loop_player(n) { printf("rank temp %d %d\\n",n,bx[n]->temp); } next_player;
 loop_char(n) { printf("rank char %d %d\\n",n,bx[n]->temp); } next_char;
}
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
void drv_decay_mines(void) { collapsed_mine_decay(); }
void drv_decay_ghosts(void) { ghost_warrior_decay(); }
void drv_decay_corpses(void) { corpse_decay(); }
void drv_decay_bodies(void) { dead_body_rot(); }
void drv_auto_drop(void) { auto_drop(); }
void drv_link_decay(void) { link_decay(); }
void drv_quest_decay(void) { quest_decay(); }
int drv_rank_threshold(int n) { return nprovs_to_rank(n); }
int drv_rank_lower(int n) { return lower_noble_rank(n); }
int drv_rank_sup(int n) { return det_noble_rank_sup(n); }
void drv_kings(void) { find_kings(); }
int drv_corpse_roll(int who) { return up_corpse(who, 0, 2); }
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
int win_flag;
ilist trades_to_check;
''')
PY

# Expose only the existing static reader; all parsing and command loading
# remains in the complete unchanged io.c translation unit.
cat >"${WORK}/io-wrapper.c" <<'C'
#include "io.c"
void drv_read_boxes(char *path) { read_boxes(path); }
C
case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -w -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib -I "${WORK}" \
 "${HERE}/driver.c" "${WORK}/helpers.c" "${WORK}/stubs.c" "${WORK}/hades-rng.c" "${WORK}/faery-rng.c" "${WORK}/npc-rng.c" \
 olympia/{code,glob,z,sout}.c "${WORK}/io-wrapper.c" \
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
  ! grep -q 'out-of-scope C dependency:' "${WORK}/${run}/stdout" "${WORK}/${run}/stderr" || die "section $i reached a C stub"
  if [[ "${section_names[$i]}" == 'section GM economy '* || "${section_names[$i]}" == 'section GM stats '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section GM economy invalid '*|'section GM stats invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected GM economy status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section lore render '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section lore render invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected lore render status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section location root '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section location root invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected location root status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section ancillary '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section ancillary invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected ancillary status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section inner sightings '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section inner sightings invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected inner sightings status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section visible characters '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section visible characters invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected character display status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section market rendering '* ]]; then
   [ "$status" = 0 ] || die "unexpected market rendering status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section routes '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section routes invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected route status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section description '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section description invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected description status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section header '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section header invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected header status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section names '* ]]; then
   [ "$status" = 0 ] || die "unexpected names status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section sout '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section sout invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected sout status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section file observer '* ]]; then
   [ "$status" = 0 ] || die "unexpected file observer status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section lifecycle '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section lifecycle invalid '*) expected_status=134 ;;
    'section lifecycle fatal '*) expected_status=1 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected lifecycle status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section post writer '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section post writer invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected post writer status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section admission '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section admission invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected admission status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section skill lore '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section skill lore invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected lore status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section report '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section report invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected report state status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section closure '* ]]; then
   [ "$status" = 0 ] || die "unexpected monthly closure status $status in section $i; expected 0"
  fi
  if [[ "${section_names[$i]}" == 'section decay '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section decay invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected decay status $status in section $i; expected $expected_status"
  fi
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
   1) [[ "${section_names[$i]}" == 'section lifecycle fatal '* ]] || die "unexpected C exit 1 in section $i" ;;
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

# A saved multi-candidate diagnostic order would encode this filesystem's
# readdir order. Instead, reuse ONE untouched input tree for C, C, then Go.
# The real make_new_players definition above is unchanged, not a model oracle.
native="${WORK}/native"
mkdir -p "${native}/input/act" "${native}/input/log" "${native}/input/orders"
for acct in z0 a1 q2 b3 p4; do mkdir "${native}/input/act/${acct}"; done
printf '' >"${native}/input/act/z0/Join-g3"
printf 'faction' >"${native}/input/act/a1/Join-g3"
printf 'faction\n\ncharacter\n' >"${native}/input/act/q2/Join-g3"
printf '\nfaction\n \ncharacter\n\ncity' >"${native}/input/act/b3/Join-g3"
printf 'faction\ncharacter\ncity\nname\n \n' >"${native}/input/act/p4/Join-g3"
cat >"${native}/cases.txt" <<'EOF'
section admission live native traversal
lore_root
admission probe
admission state
save 1000
EOF
native_snapshot() {
 python3 -I - "${native}/input" <<'PY'
import json, os, pathlib, stat, sys
root = pathlib.Path(sys.argv[1])
def snapshot(path):
    mode = path.lstat().st_mode
    if stat.S_ISDIR(mode):
        # Keep native entry order too: even an unnoticed remove/recreate can
        # invalidate a traversal oracle despite preserving all file contents.
        return [mode, [[name, snapshot(path / name)] for name in os.listdir(path)]]
    if not stat.S_ISREG(mode):
        sys.exit(f'unexpected native fixture file type: {path}')
    return [mode, path.read_bytes().hex()]
print(json.dumps(snapshot(root)))
PY
}
native_snapshot >"${native}/before.json"
for run in run1 run2; do
 status=0
 run_limited "${RUN_TIMEOUT}" command-driver "${WORK}/driver" "${native}/cases.txt" 0 "${native}/input" \
  >"${native}/${run}.stdout" 2>"${native}/${run}.stderr" || status=$?
 ! grep -q 'out-of-scope C dependency:' "${native}/${run}.stdout" "${native}/${run}.stderr" || die "native traversal reached a C stub"
 [ "$status" = 0 ] || die "unexpected native traversal C status $status"
 printf '%d\n' "$status" >"${native}/${run}.status"
 native_snapshot >"${native}/after.json"
 cmp "${native}/before.json" "${native}/after.json" || die "C mutated native traversal input"
done
for stream in stdout stderr status; do
 cmp "${native}/run1.${stream}" "${native}/run2.${stream}" || die "native traversal C ${stream} is not deterministic"
done
{
 printf 'stdout:\n'; cat "${native}/run1.stdout"
 printf 'stderr:\n'; cat "${native}/run1.stderr"
 printf '\nexit %s\n' "$(cat "${native}/run1.status")"
} >"${native}/out.txt"
if [ "$update" = 1 ]; then
 cp "${WORK}/run1/out.txt" "${EXPECT}"
else
 diff -u "${EXPECT}" "${WORK}/run1/out.txt"
fi
(cd "${ROOT}" && COMMAND_NATIVE_DIR="${native}" RUN_TIMEOUT="${RUN_TIMEOUT}" go test -exec "bash ${HERE}/../lib/run-limited.sh" \
 -count=1 -run '^TestCommand' ./internal/olympia)
native_snapshot >"${native}/after.json"
cmp "${native}/before.json" "${native}/after.json" || die "Go mutated native traversal input"
echo "Live native traversal: raw C/C/Go match, shared input unchanged"
echo PASS
