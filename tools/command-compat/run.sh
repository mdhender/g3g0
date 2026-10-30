#!/usr/bin/env bash
# Copyright (c) 2026 Michael D Henderson.
# SPDX-License-Identifier: AGPL-3.0-or-later
set -euo pipefail
HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "${HERE}/../.." && pwd)"
. "${HERE}/../lib/timeout.sh"
RUN_TIMEOUT="${RUN_TIMEOUT:-10}"
G3_REPO="${G3_REPO:-${HOME}/Software/playbymail/olympia-g3}"
G3_REV="${G3_REV:-9a03d126b329ed466c4af48425553c025eaca26f}"
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
    'swear.c': ['lord', 'player', 'set_loyal', 'reswear_all_sworn', 'set_lord', 'unit_deserts',
                'np_to_acquire', 'enough_np_to_acquire', 'v_honor', 'v_oath', 'v_swear',
                'v_raise', 'd_raise', 'v_rally', 'd_rally', 'v_bribe', 'thanks_for_gift', 'd_bribe',
                'terrorize_vassal', 'terrorize_prisoner', 'v_terrorize', 'd_terrorize', 'v_incite'],
    'move.c': ['check_captain_loses_sailors', 'restore_stack_actions', 'touch_loc_after_move', 'move_stack',
               'departure_message', 'arrival_message', 'discover_road', 'parse_exit_dir',
               'save_v_array', 'restore_v_array', 'move_exit_land', 'land_check',
               'can_move_here', 'can_move_at_outer_level', 'suspend_stack_actions',
               'clear_guard_flag', 'v_move', 'd_move',
               'move_exit_fly', 'fly_check', 'can_fly_here', 'can_fly_at_outer_level', 'v_fly', 'd_fly',
               'move_exit_water', 'sail_depart_message', 'sail_arrive_message', 'sail_check',
               'can_sail_here', 'can_sail_at_outer_level', 'v_sail', 'd_sail', 'i_sail',
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
    'scry.c': ['notify_loc_shroud', 'reset_cast_where'],
    'make.c': ['find_make', 'v_generic_make', 'd_generic_make', 'i_generic_make',
               'v_second_make', 'd_second_make', 'v_make', 'd_make', 'i_make',
               'v_use_train_riding', 'v_use_train_war'],
    'basic.c': ['v_hinder_med', 'hinder_med_omen', 'd_hinder_med'],
    'alchem.c': ['new_potion'],
    'seed.c': ['int_comp', 'seed_taxes'],
    'stack.c': ['stack_parent', 'stack_leader', 'here_pos', 'stacked_beneath',
                'first_prisoner_pos', 'promote', 'unstack', 'leave_stack', 'stack', 'join_stack',
                'check_prisoner_escape', 'prisoner_escapes', 'prisoner_movement_escape_check',
                'weekly_prisoner_escape_check',
                'drop_stack', 'extract_stacked_unit', 'move_prisoner', 'v_stack',
                'here_precedes', 'promote_stack', 'v_unstack', 'promote_after', 'v_promote',
                'take_prisoner', 'has_prisoner', 'prisoner_gives_today', 'give_prisoner'],
    'dir.c': ['los_province_distance', 'location_direction', 'exit_distance', 'is_port_city',
              'province_has_port_city', 'summer_uldim_open_now', 'add_province_exit',
              'extra_routes', 'province_exits', 'province_sub_exits', 'subloc_exits',
              'ship_exits', 'building_exits', 'exits_from_loc', 'exits_from_loc_nsew',
              'exits_from_loc_nsew_select', 'has_ocean_access', 'list_exit_extras',
              'list_exits_sup', 'list_road_sup', 'list_exits', 'list_sailable_routes'],
    'combat.c': ['fort_covers', 'auto_attack', 'check_auto_attack_sup', 'check_all_auto_attacks', 'clear_second_waits', 'find_wield', 'wield_s', 'v_guard'],
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
               'garr_own_s', 'garrison_summary', 'pledged_beneath', 'v_pledge'],
    'build.c': ['fort_default_defense', 'ship_loc_okay', 'temple_loc_okay',
                'tower_loc_okay', 'mine_loc_okay', 'inn_loc_okay', 'province_subloc',
                'castle_loc_okay', 'find_build', 'build_materials_check', 'create_new_building', 'start_build',
                'build_structure', 'unfinished_building', 'v_build', 'daily_build', 'd_build',
                'repair_points', 'v_repair', 'd_repair', 'i_repair', 'v_improve', 'd_improve',
                'v_raze', 'd_raze'],
    'art.c': ['has_auraculum', 'max_eff_aura', 'max_current_aura', 'limit_cur_aura', 'token_player', 'swear_token_units',
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
            'rest_name', 'nprovinces', 'greater_region', 'deduct_aura', 'charge_aura', 'check_aura',
            'min', 'cap', 'set_bit', 'test_bit', 'count_stack_any', 'count_any', 'check_skill',
            'contacted', 'char_where', 'char_here', 'check_char_where', 'check_char_here',
            'check_char_gone', 'check_still_here', 'add_char_damage', 'count_man_items',
            'diff_region', 'test_known', 'loc_hidden', 'has_item', 'loc_depth',
            'lookup', 'gold_s', 'comma_num', 'nice_num', 'loyal_s', 'add_np', 'deduct_np',
            'new_char', 'beast_capturable', 'my_prisoner', 'create_unique_item', 'create_unique_item_alloc',
            'destroy_unique_item', 'add_item', 'sub_item', 'set_known', 'gen_item',
            'consume_item', 'move_item', 'drop_item',
            'stackmate_inheritor', 'take_unit_items', 'put_back_cookie',
            'survive_fatal', 'char_reclaim', 'v_reclaim', 'dead_char_body', 'kill_char', 'hack_unique_item',
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
    'savage.c': ['create_savage', 'call_savage', 'v_use_drum', 'savage_hates', 'savage_hate_here',
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
             'v_wait', 'd_wait', 'i_wait', 'next_np_turn', 'print_hiring_status', 'print_unformed',
             'form_new_noble', 'v_form', 'd_form', 'v_accept', 'will_accept_sup', 'will_accept', 'how_many',
             'v_give', 'v_pay', 'may_take', 'v_get', 'may_name', 'v_name', 'v_times',
             'v_fullname', 'v_banner', 'v_flag', 'v_split', 'v_emote'],
    'produce.c': ['init_collect_list', 'replenish', 'mine_production', 'location_production',
                  'item_gen_here', 'find_harv', 'v_generic_harvest', 'mage_menial_how',
                  'i_generic_harvest', 'v_collect', 'd_generic_harvest', 'bump_other_collectors',
                  'd_collect', 'i_collect', 'v_quarry', 'v_recruit', 'v_raise_corpses',
                  'v_fish', 'v_wood', 'v_opium', 'v_mallorn', 'v_yew', 'v_catch', 'v_mage_menial'],
    'perm.c': ['is_hostile', 'is_defend', 'clear_all_att', 'rp_admit', 'p_admit', 'will_admit',
               'v_admit', 'set_att', 'v_set_att', 'v_hostile', 'v_defend', 'v_neutral', 'v_att_clear',
               'admit_comp', 'print_admit_sup', 'print_admit', 'print_att_sup', 'print_att'],
    'storm.c': ['weather_here', 'can_see_weather_here', 'update_weather_view_loc_sup', 'init_weather_views',
                'update_weather_view_locs', 'new_storm', 'create_some_storms',
                'natural_weather', 'v_bind_storm', 'd_bind_storm',
                'move_storm', 'move_bound_storms', 'dissipate_storm', 'storm_report'],
    'use.c': ['find_use_entry', 'use_tbl_skill', 'may_use_skill', 'magically_speed_casting',
              'correct_use_item', 'meets_requirements', 'consume_requirements', 'consume_scroll',
              'experience_use_speedup', 'add_skill_experience',
              'v_use', 'd_use', 'i_use', 'v_use_item', 'd_use_item',
              'rp_skill_ent', 'p_skill_ent', 'learn_skill', 'exp_level', 'has_skill', 'char_np_total', 'forget_skill', 'skill_school',
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
    'c2.c': ['v_stop', 'v_format', 'v_notab', 'open_times', 'times_masthead', 'close_times',
             'v_discard', 'v_claim', 'v_public', 'text_list_free', 'line_length_check',
             'parse_text_list', 'v_post', 'v_message', 'v_rumor', 'v_press',
             'drop_player', 'v_quit', 'v_die', 'v_fee', 'board_message', 'v_board',
             'unboard_message', 'v_unload', 'v_ferry'],
    'immed.c': ['v_xyzzy', 'v_plugh', 'v_credit', 'v_relore'],
    'stealth.c': ['v_hide', 'd_hide', 'clear_contacts', 'add_contact', 'v_contact', 'cloak_lord'],
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
        if file == 'swear.c':
            start, end = source.index('#define\tTAG_TURN'), source.index('/* INCITE: rumor flavor')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/swear.c"\n')
            out.write(source[start:end])
            start, end = source.index('#define\t\tSWITCH'), source.index('int\nd_bribe(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/swear.c"\n')
            out.write(source[start:end])
        if file == 'stack.c':
            for pattern in [r'^static ilist gives_today = NULL;\n',
                            r'^static int gives_turn = -1;\n', r'^static int gives_day = -1;\n']:
                out.write(re.search(pattern, source, re.M).group(0))
        if file == 'c2.c':
            start, end = source.index('#define DASH_LINE'), source.index('char *months[]')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/c2.c"\n')
            out.write(source[start:end])
            out.write(re.search(r'^char \*months\[\] = \{.*?^\};\n', source, re.M | re.S).group(0))
            out.write('int test_use_const_report_date = 1;\n')
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
            out.write(re.search(r'^#define\s+noble_cost\(pl\).*\n', source, re.M).group(0))
            start, end = source.index('struct flag_ent {'), source.index('static int\nflag_raised(')
            out.write(source[start:end])
            start, end = source.index('static char *wait_tags[]'), source.index('void\nclear_wait_parse(')
            out.write(source[start:end])
        if file == 'u.c':
            out.write(re.search(r'^static int dot_count = 0;\n', source, re.M).group(0))
            out.write(re.search(r'^int weight_display_flag = 0;\n', source, re.M).group(0))
            # Whole original RNG prefix: shared persistent entity stream, tags,
            # turn guard and all keyed leaves, not replacement helper draws.
            start = source.index('#define\tETAG_TURN')
            out.write('#include "rng.h"\n')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/u.c"\n')
            out.write(source[start:source.index('void\nkill_stack_ocean(int who)')])
            start, end = source.index('static char *num_s[]'), source.index('char *\nnice_num(')
            out.write(source[start:end])
        if file == 'build.c':
            for name in ['ship_loc_okay', 'temple_loc_okay', 'tower_loc_okay',
                         'mine_loc_okay', 'inn_loc_okay', 'castle_loc_okay']:
                out.write(f'static int {name}(struct command *, int);\n')
            start, end = source.index('struct build_ent\n'), source.index('static struct build_ent *\nfind_build(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/build.c"\n')
            out.write(source[start:end])
            start, end = source.index('static int improve_stone[]'), source.index('int\nv_improve(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/build.c"\n')
            out.write(source[start:end])
        if file == 'make.c':
            start, end = source.index('#define\tWHERE_SHIP'), source.index('static struct make *\nfind_make(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/make.c"\n')
            out.write(source[start:end])
        if file == 'produce.c':
            # Original complete production tables/macros, not a driver copy.
            start, end = source.index('#define\t\tMOUNTAIN_STONE'), source.index('static void\nreplenish(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/produce.c"\n')
            out.write(source[start:end])
            out.write(re.search(r'^static ilist collectors = NULL;\n', source, re.M).group(0))
            start, end = source.index('struct harvest\n'), source.index('static struct harvest *\nfind_harv(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/produce.c"\n')
            out.write(source[start:end])
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
        if file == 'basic.c':
            # Entire unchanged magic RNG prefix, including tags and turn guard.
            start, end = source.index('#define\tTAG_TURN'), source.index('int\nv_meditate(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/basic.c"\n')
            out.write(source[start:end])
        if file == 'use.c':
            # Complete original struct, callback declarations and table, unchanged.
            # Missing callbacks are abort-only references, never invoked by helpers.
            start, end = source.index('struct use_tbl_ent {'), source.index('\n\nint\nv_implicit(')
            out.write(f'#line {source.count(chr(10), 0, start) + 1} "olympia/use.c"\n')
            out.write(source[start:end])
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
            if name == 'v_flag':
                pattern = r'^int v_flag\(.*?^\}\n'
            if name == 'i_sail':
                pattern = r'^int i_sail\(.*?^\}\n'
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
void drv_social_rolls(int who, int target) {
 printf("social rolls %d %d %d %d %d\\n",who,target,soc_gift(who,target),soc_bribe(who,target,0,1,2),soc_bribe(who,target,1,1,100));
}
void drv_terror_roll(int who, int target) {
 printf("terror roll %d %d %d\\n",who,target,soc_terror(who,target));
}
int drv_use_callback(struct command *c, int sk, int finish) {
 struct use_tbl_ent *e = &use_tbl[find_use_entry(sk)];
 return finish ? e->finish(c) : e->start(c);
}
void drv_magic_rolls(int who, int other) {
 printf("magic rolls %d %d %d\\n",magc_omen(who,other),magc_med(who,0),magc_scry(who,other));
}
int drv_potion_roll(int who, int item, int sub) { return magc_potion(who,item,sub,1,2); }
void drv_use_table(void) {
 int i;
 for(i=0;i<(int)(sizeof(use_tbl)/sizeof(*use_tbl));i++) {
  struct use_tbl_ent *e=&use_tbl[i];
  printf("use row %d %s %d %d %d %d %d %d lookup %d skill %d\\n",i,e->allow?e->allow:"NULL",e->skill,
   e->start!=NULL,e->finish!=NULL,e->interrupt!=NULL,e->time,e->poll,find_use_entry(e->skill),use_tbl_skill(i));
 }
}
int drv_use_lookup(int sk) { return find_use_entry(sk); }
int drv_use_may(int who,int sk) { return may_use_skill(who,sk); }
int drv_use_correct(struct command *c) { return correct_use_item(c); }
int drv_use_meets(int who,int sk) { return meets_requirements(who,sk); }
void drv_use_consume(int who,int sk) { consume_requirements(who,sk); }
void drv_use_scroll(int who,int basis) { consume_scroll(who,basis); }
void drv_use_magic(struct command *c,int sk) { magically_speed_casting(c,sk); }
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
 else if (!strcmp(part,"wealth")) gm_faction_wealth(pl);
 else if (!strcmp(part,"nobles")) gm_nobles_list(pl);
 else if (!strcmp(part,"details")) gm_player_details(pl);
 else if (!strcmp(part,"notices")) list_all_notices(pl);
 else if (!strcmp(part,"scratch")) {
  int n, k, kinds[]={T_player,T_skill};
  for(k=0;k<2;k++) {
   loop_kind(kinds[k],n) { printf("gm temp %d %d\\n",n,bx[n]->temp); } next_kind;
  }
 }
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
int drv_promote_after(int a, int b) { return promote_after(a,b); }
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
void drv_list_player(int pl) { int flag = 0; list_a_player(stdout, pl, &flag); }
int drv_load(struct command *c) { return load_command(c); }
void drv_add_contact(int a, int b) { add_contact(a,b); }
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
static struct make drv_raw_make;
void drv_make_raw(int item, int worker, int inp1, int inp2, int aura, int factor, int pub, int where, int skill) {
 drv_raw_make = *find_make(item);
 drv_raw_make.worker=worker; drv_raw_make.inp1=inp1; drv_raw_make.inp2=inp2;
 drv_raw_make.aura=aura; drv_raw_make.factor=factor; drv_raw_make.public=pub;
 drv_raw_make.where=where; drv_raw_make.req_skill=skill;
}
void drv_build_table(void) {
 for (int i=0; build_tbl[i].what; i++) {
  struct build_ent *t=&build_tbl[i];
  int cb=t->loc_ok==ship_loc_okay ? 0 : t->loc_ok==temple_loc_okay ? 1 : t->loc_ok==inn_loc_okay ? 2 : t->loc_ok==castle_loc_okay ? 3 : t->loc_ok==mine_loc_okay ? 4 : t->loc_ok==tower_loc_okay ? 5 : -1;
  printf("build row %d %s %d %d %d %d %d %d %d %d %d %d %d %d [%s]\\n",i,t->what,t->skill_req,t->skill2,t->kind,cb,t->unfinished_subkind,t->finished_subkind,t->min_workers,t->worker_days,t->min_days,t->req_item,t->req_qty,t->capacity,t->default_name);
 }
}
int drv_build(struct command *c, char *action, char *s, int n) {
 if(!strcmp(action,"repair_points")) return repair_points(n);
 if(!strcmp(action,"raze_v")) return v_raze(c);
 if(!strcmp(action,"raze_d")) return d_raze(c);
 if(!strcmp(action,"repair_v")) return v_repair(c);
 if(!strcmp(action,"repair_d")) return d_repair(c);
 if(!strcmp(action,"repair_i")) return i_repair(c);
 if(!strcmp(action,"improve_v")) return v_improve(c);
 if(!strcmp(action,"improve_d")) return d_improve(c);
 if(!strcmp(action,"defense")) return fort_default_defense(n);
 if(!strcmp(action,"roll")) return ent_build(n);
 if(!strcmp(action,"v")) return v_build(c);
 if(!strcmp(action,"d")) return d_build(c);
 if(!strcmp(action,"unfinished")) { char *p=unfinished_building(c->who); printf("unfinished [%s]\\n",p?p:"NULL"); return p!=NULL; }
 struct build_ent *t=find_build(s);
 if(!strcmp(action,"find")) { printf("build find [%s] %d\\n",t?t->what:"NULL",fuzzy_build_match); return t!=NULL; }
 if(!t) abort();
 if(!strcmp(action,"loc")) return t->loc_ok(c,n);
 if(!strcmp(action,"materials")) return build_materials_check(c,t);
 if(!strcmp(action,"start")) return start_build(c,t,n);
 if(!strcmp(action,"structure")) return build_structure(c,t,n);
 if(!strcmp(action,"create")) { create_new_building(c,t,n); return 1; }
 abort();
}
void drv_make_table(void) {
 for (int i=0; i<sizeof(make_tbl)/sizeof(*make_tbl); i++) {
  struct make *t=&make_tbl[i];
  printf("make row %d %d %d %d %d %d [%s] %d %d %d %d %d\\n", i,t->item,t->inp1,t->inp2,t->req_skill,t->worker,t->got_em ? t->got_em : "NULL",t->public,t->where,t->aura,t->factor,t->days);
 }
}
int drv_make_find(int item) { return find_make(item)!=NULL; }
int drv_make_call(struct command *c, int action, int item, int number) {
 struct make *t=item == -1 ? &drv_raw_make : find_make(item);
 if(action==0) return v_generic_make(c,number,t);
 if(action==1) return d_generic_make(c,t);
 if(action==3) return v_second_make(c,number,t);
 if(action==4) return d_second_make(c,t);
 if(action==5) return d_make(c);
 if(action==6) return i_make(c);
 return i_generic_make(c,t);
}
ilist drv_collectors(void) { return collectors; }
void drv_collect_add(int who) { ilist_append(&collectors, who); }
void drv_harvest_table(void) {
 for (int i=0; i<sizeof(harv_tbl)/sizeof(*harv_tbl); i++) {
  struct harvest *t=&harv_tbl[i];
  printf("harvest row %d %d %d %d %d %d %d %d\\n", i,t->item,t->vis_item,t->mult,t->skill,t->worker,t->chance,t->public);
  printf("harvest text [%s] [%s] [%s] [%s]\\n", t->got_em ? t->got_em : "NULL",t->none_now ? t->none_now : "NULL",t->none_ever ? t->none_ever : "NULL",t->task_desc ? t->task_desc : "NULL");
 }
}
int drv_harvest_find(int item) { return find_harv(item)!=NULL; }
int drv_harvest_gen(int terr,int item) { return item_gen_here(terr,item); }
int drv_harvest_start(struct command *c,int item,int number,int days) { return v_generic_harvest(c,number,days,find_harv(item)); }
int drv_harvest_summary(struct command *c,int item) { return i_generic_harvest(c,find_harv(item)); }
int drv_harvest_daily(struct command *c,int item) { return d_generic_harvest(c,find_harv(item)); }
void drv_harvest_bump(int where,int item) { bump_other_collectors(where,find_harv(item)); }
void drv_menial(int who) {
 int roll=ent_menial(who);
 printf("menial %d %d [%s]\\n",who,roll,mage_menial_how(who));
}
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
int drv_sick_roll(int who, int health) { return ent_sick(who, health); }
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
int drv_prisoner_count(int pris) { return prisoner_gives_today(pris); }
void drv_prisoner_state(void) {
 printf("gives allocated %d turn %d day %d:",gives_today != NULL,gives_turn,gives_day);
 for(int i=0;i<ilist_len(gives_today);i++) printf(" %d",gives_today[i]);
 printf("\\n");
}
void drv_signal(int who, char *flag) {
 struct flag_ent *p = my_malloc(sizeof(*p));
 p->who = who; p->flag = str_save(flag); flag_ents_append(&flags, p);
}
void drv_flags(void) {
 printf("flags allocated %d count %d\\n",flags!=NULL,flag_ents_len(flags));
 for(int i=0;i<flag_ents_len(flags);i++) printf("flag %d [%s]\\n",flags[i]->who,flags[i]->flag);
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
    start = source.index('#define\tPROV_OR_CITY')
    end = source.index('\n\n\nint\ncreate_peasant_mob(', start)
    line = source.count('\n', 0, start) + 1
    out.write(f'#line {line} "olympia/npc.c"\n{source[start:end]}\n')
    out.write('''
void drv_cookie_table(void) {
 int i;
 for(i=0;i<(int)(sizeof(cookie_monster)/sizeof(cookie_monster[0]));i++) {
  struct cookie_monster_tbl *t=&cookie_monster[i];
  printf("cookie row %d %d %d %d %d %d %d %d %d [%s] [%s]\\n",i,
   t->cookie,t->kind,t->sk,t->ni,t->terrain,t->man_kind,t->low,t->high,
   t->not_here ? t->not_here : "NULL",t->no_cookies ? t->no_cookies : "NULL");
 }
}
int drv_cookie_find(int k) { return find_cookie(k)!=NULL; }
''')
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
routines = {'v_use_heal', 'v_use_slave', 'v_use_death', 'v_use_palantir',
            'd_use_palantir', 'v_use_proj_cast', 'v_use_quick_cast', 'v_use_drum',
            'v_use_faery_stone', 'v_use_orb', 'v_suffuse_ring', 'v_use_bta_skull'}
for m in re.finditer(r'^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),', table, re.M):
    routines.update(n for n in m.groups() if n != 'NULL')
source = open(f'{src}/use.c').read()
table = source[source.index('struct use_tbl_ent use_tbl[]'):source.index('\n\nint\nv_implicit(')]
for m in re.finditer(r'^\{[^,]*,[^,]*,\s*(\w+),\s*(\w+),\s*(\w+),', table, re.M):
    routines.update(n for n in m.groups() if n not in ('NULL', '0'))
with open(f'{dst}/stubs.c', 'w') as out:
    out.write('''#include <stdio.h>
#include <stdlib.h>
#include "z.h"
#include "oly.h"
#define STUB(n) { fprintf(stderr, "out-of-scope C dependency: %s\\n", n); abort(); }
''')
    for name in sorted(routines - real):
        args = 'struct command *c, int item' if name == 'v_suffuse_ring' else 'struct command *c'
        out.write(f'int {name}({args}) STUB("{name}")\n')
    out.write('''
int immediate, flush_always = 1, pretty_data_files, time_self;
int win_flag;
ilist trades_to_check;
''')
# Complete unchanged coverage.c, including all three families and static owners.
source = open(f'{src}/coverage.c').read()
with open(f'{dst}/coverage-command.c', 'w') as out:
    out.write('#line 1 "olympia/coverage.c"\n')
    out.write(source)
    out.write('''
void drv_ph_seed(int row, int ran, int acted) { ph_ran[row] = ran; ph_acted[row] = acted; }
void drv_coverage_state(int cmd) {
 printf("coverage allocated %d ncmds %d table %d\\n", cov_counts != NULL, cov_ncmds, cov_num_cmds());
 if (cov_counts && cmd >= 0 && cmd < cov_ncmds)
  printf("coverage row %d %d %d %d %d\\n", cmd, cov_counts[cmd][0], cov_counts[cmd][1], cov_counts[cmd][2], cov_counts[cmd][3]);
}
void drv_coverage_seed(int cmd, int outcome, int value) {
 assert(cov_counts && cmd >= 0 && cmd < cov_ncmds && outcome >= 0 && outcome < COV_MAX);
 cov_counts[cmd][outcome] = value;
}
void drv_ucov_state(void) {
 printf("ucov allocated %d nents %d table %d\\n", ucov_skills != NULL, ucov_nents, ucov_num_ents());
}
void drv_ucov_seed(int ent, int outcome, int value) { ucov_skills[ent][outcome] = value; }
''')
    # Names and evaluated enum values come from the original C header, never Go.
    header = open(f'{src}/oly.h').read()
    enum = header[header.index('PH_SAVAGE_SPAWN'):header.index('PH_MAX')]
    names = re.findall(r'\bPH_[A-Z_]+\b', enum)
    if len(names) != 74 or len(set(names)) != 74:
        sys.exit('unexpected phase enum')
    out.write('void drv_ph_metadata(void) {\n')
    for name in names:
        out.write(f'printf("{name} %d\\n", {name});\n')
    out.write('printf("PH_MAX %d\\n", PH_MAX);\n}\n')
PY

# Expose only the existing static reader; all parsing and command loading
# remains in the complete unchanged io.c translation unit.
cat >"${WORK}/io-wrapper.c" <<'C'
#include "io.c"
extern ilist load_q[5], run_q;
void drv_read_boxes(char *path) { read_boxes(path); }
/* Restart only the command record through the original writer/reader, with
 * process-local scheduler queues discarded as on process restart. Orders and
 * all gameplay records remain owned by their normal implementations. */
void drv_persist_command(int who) {
 FILE *fp = fopen(sout("%s/command-reload",libdir),"w");
 int i;
 assert(fp);
 fprintf(fp,"%d %s 0\n",who,kind_s[kind(who)]);
 print_command(fp,rp_command(who));
 fprintf(fp,"\n");
 fclose(fp);
 bx[who]->cmd = NULL;
 for(i=0;i<5;i++) ilist_clear(&load_q[i]);
 ilist_clear(&run_q);
 read_boxes("command-reload");
}
C
# Keep the entire original translation unit, including export_location.
# Expose static entry points without replacing bodies or dependencies.
cat >"${WORK}/state-export-wrapper.c" <<'C'
#include "state_export.c"
void drv_json(FILE *fp, char *s) { json_str(fp, s); }
void drv_export_unit(FILE *fp, int who) { export_unit(fp, who); }
void drv_export_location(FILE *fp, int where, int viewer) { export_location(fp, where, viewer); }
void drv_export_faction(FILE *fp, int pl) { export_faction(fp, pl); }
void drv_sort_ids(ilist l) { sort_ids(l); }
C
case "$(uname -s)" in
 Darwin) ldflags=(-Wl,-dead_strip) ;;
 *) ldflags=(-Wl,--gc-sections) ;;
esac
(cd "${WORK}/g3" && "${CC}" -std=c11 -w -ffunction-sections -fdata-sections "${ldflags[@]}" \
 -I olympia -I lib -I "${WORK}" \
 "${HERE}/driver.c" "${WORK}/helpers.c" "${WORK}/stubs.c" "${WORK}/hades-rng.c" "${WORK}/faery-rng.c" "${WORK}/npc-rng.c" \
 "${WORK}/coverage-command.c" "${WORK}/state-export-wrapper.c" \
 olympia/{code,glob,z,sout}.c "${WORK}/io-wrapper.c" \
 lib/{checked_alloc,ilist,exit_views,rng,rnd,cstrings,item_ents,skill_ents,orders,wait_args,flag_ents,admits,trades,req_ents,accept_ents}.c \
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
  marker_status=0
  grep -Ei 'out-of-scope|not (yet )?implemented|missing.*dependency:|NYI|stub|Sanitizer|runtime error:' \
   "${WORK}/${run}/stdout" "${WORK}/${run}/stderr" >&2 || marker_status=$?
  [ "$marker_status" = 1 ] || die "section $i invalid C reference or unreadable streams"
  if [[ "${section_names[$i]}" == 'section N22b '* ]]; then
   expected_status=0
   if [[ "${section_names[$i]}" == 'section N22b negative C32 manpower reaches shared damage assertion' ]]; then
    expected_status=134
    grep -Fxq 'assertion failure: olympia/u.c (1148): damage >= 0' "${WORK}/${run}/stderr" || die "wrong N22b assertion in section $i"
   fi
   [ "$status" = "$expected_status" ] || die "unexpected N22b status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N22a '* ]]; then
   [ "$status" = 0 ] || die "unexpected N22a status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N21 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N21 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N20 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N20 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N19 '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N19 negative boarding amount reaches direct charge assertion')
     expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N19 status $status in section $i"
   if [ "$expected_status" = 134 ]; then
    grep -Fxq 'assertion failure: olympia/u.c (645): qty >= 0' "${WORK}/${run}/stderr" || die "wrong N19 assertion in section $i"
   fi
  fi
  if [[ "${section_names[$i]}" == 'section N18 '* ]]; then
   expected_status=0
   expected_assert=''
   case "${section_names[$i]}" in
    'section N18 invalid completion changed ship asserts still ship at1612')
     expected_status=134; expected_assert='assertion failure: olympia/move.c (1612): is_ship(ship)' ;;
    'section N18 invalid completion changed captain asserts owner at1613')
     expected_status=134; expected_assert='assertion failure: olympia/move.c (1613): building_owner(ship) == c->who' ;;
    'section N18 invalid interruption changed ship asserts still ship at1657')
     expected_status=134; expected_assert='assertion failure: olympia/move.c (1657): is_ship(ship)' ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N18 status $status in section $i"
   if [ -n "$expected_assert" ]; then
    grep -Fxq "$expected_assert" "${WORK}/${run}/stderr" || die "wrong N18 assertion in section $i"
   fi
  fi
  if [[ "${section_names[$i]}" == 'section N17 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N17 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N16 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N16 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N15 '* ]]; then
   expected_status=0
   expected_assert=''
   case "${section_names[$i]}" in
    'section N15 finish mutated NPC loyalty preserves swear614 assertion after real charge')
     expected_status=134; expected_assert='assertion failure: olympia/swear.c (614): FALSE' ;;
    'section N15 negative amount preserves sub_item assertion at completion not start')
     expected_status=134; expected_assert='assertion failure: olympia/u.c (645): qty >= 0' ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N15 status $status in section $i"
   if [ -n "$expected_assert" ]; then
    grep -Fxq "$expected_assert" "${WORK}/${run}/stderr" || die "wrong N15 assertion in section $i"
   fi
  fi
  if [[ "${section_names[$i]}" == 'section FOLLOW288 '* || "${section_names[$i]}" == 'section FOLLOW290 '* ]]; then
   [ "$status" = 0 ] || die "unexpected follow-C status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N14b '* ]]; then
   [ "$status" = 0 ] || die "unexpected N14b status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N14a '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N14a invalid '*|'section N14a exhausted '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N14a status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N13b '* ]]; then
   [ "$status" = 0 ] || die "unexpected N13b status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N13a '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N13a raw overtarget remainder asserts original location') expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N13a status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N12b '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N12b invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N12b status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N12a '* ]]; then
   [ "$status" = 0 ] || die "unexpected N12a status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N11 '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N11 invalid CREDIT nonzero invalid item '*|\
    'section N11 invalid EMOTE absent target '*) expected_status=139 ;;
    'section N11 invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N11 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N10 '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N10 invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N10 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N09 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N09 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N08 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N08 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N07 '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N07 invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N07 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N06b '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N06b invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N06b status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N06a '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N06a invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N06a status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N05 '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N05 invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N05 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N04 '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section N04 invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected N04 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N03 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N03 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N02 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N02 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section N01 '* ]]; then
   [ "$status" = 0 ] || die "unexpected N01 status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section los '* ]]; then
   [ "$status" = 0 ] || die "unexpected LOS status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section skill grant '* ]]; then
   [ "$status" = 0 ] || die "unexpected skill grant status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section potion '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section potion exhausted '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected potion status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section aura cap '* ]]; then
   [ "$status" = 0 ] || die "unexpected aura cap status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section drum '* ]]; then
   [ "$status" = 0 ] || die "unexpected drum status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section hinder meditation '* ]]; then
   [ "$status" = 0 ] || die "unexpected hinder-meditation status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section guard stack '* ]]; then
   [ "$status" = 0 ] || die "unexpected GUARD/STACK status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section state export '* ]]; then
   [ "$status" = 0 ] || die "unexpected state export status $status in section $i"
  fi
  if [[ "${section_names[$i]}" == 'section phase coverage '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section phase coverage invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected phase coverage status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section use coverage '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section use coverage invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected USE coverage status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section use prerequisites '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section use prerequisites invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected USE prerequisite status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section coverage '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section coverage invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected command coverage status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section include assembly '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section include assembly invalid '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected include status $status in section $i; expected $expected_status"
  fi
  if [[ "${section_names[$i]}" == 'section GM '* ]]; then
   expected_status=0
   case "${section_names[$i]}" in
    'section GM economy invalid '*|'section GM stats invalid '*|\
    'section GM notices invalid '*|'section GM root invalid '*|\
    'section GM root output failure '*|'section GM root loyalty assertion '*) expected_status=134 ;;
   esac
   [ "$status" = "$expected_status" ] || die "unexpected GM report status $status in section $i; expected $expected_status"
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
 marker_status=0
 grep -Ei 'out-of-scope|not (yet )?implemented|missing.*dependency:|NYI|stub|Sanitizer|runtime error:' \
  "${native}/${run}.stdout" "${native}/${run}.stderr" >&2 || marker_status=$?
 [ "$marker_status" = 1 ] || die "invalid native C reference or unreadable streams"
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
