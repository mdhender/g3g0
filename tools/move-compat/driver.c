/* Copyright (c) 2026 Michael D Henderson.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Fixture construction and field dumps only; every compared helper is real C.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "z.h"
#include "oly.h"
#include "rnd.h"
#include "combat-structs.h"

int find_wield(struct wield *, int, struct fight *);

static void fight_dump(struct fight *f)
{
    printf("fight %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n",
           f->unit, f->kind, f->sav_num, f->num, f->new_health,
           f->attack, f->defense, f->behind, f->missile, f->protects,
           f->nprot, f->ally, f->prisoner, f->inside, f->seize_slot,
           f->survived_fatal, f->attack_item, f->defense_item, f->missile_item);
}

void glob_init(void);
void init_lower(void);
int hades_region, month_done, evening;
int immediate, flush_always = 1;
int immed_see_all;
extern int second_indent;
int drv_land(struct command *, struct exit_view *, int);
int drv_fly(struct command *, struct exit_view *, int);
int drv_water(struct command *, struct exit_view *, int, int);
int drv_check(struct command *, struct exit_view *, int);
int drv_can(int, struct command *);
int drv_outer(int, struct command *);
void drv_save(struct command *, struct exit_view *);
void drv_restore(struct command *, struct exit_view *);
void drv_suspend(int);
void drv_unstack(int);
int drv_sub_item(int, int, int);
int find_command(char *);

static void weights_dump(struct weights *w)
{
    printf("weights %d %d %d %d %d %d %d %d\n", w->animals,
           w->total_weight, w->land_cap, w->land_weight,
           w->ride_cap, w->ride_weight, w->fly_cap, w->fly_weight);
}

/* All fixture fields read by the travel closure, including inventory order,
 * ownership, location/stack order, and sout's report-path output mutations. */
static void state_dump(void)
{
    int n, i;
    printf("state factor %d\n", army_slow_factor);
    loop_boxes(n) {
        printf("box %d %d %d %d here", n, kind(n), subkind(n), loc(n));
        for (i = 0; i < ilist_len(rp_loc_info(n)->here_list); i++)
            printf(" %d", rp_loc_info(n)->here_list[i]);
        printf("\n");
        if (rp_char(n)) printf("unit %d %d %d %d\n", noble_item(n), is_prisoner(n), loyal_kind(n), lord(n));
        if (rp_item(n)) printf("item %d %d %d %d %d %d %d %d %d\n", item_weight(n), item_land_cap(n), item_ride_cap(n), item_fly_cap(n), item_animal(n), man_item(n), item_attack(n), item_defense(n), item_missile(n));
        if (rp_misc(n)) printf("storm %d\n", storm_strength(n));
        for (i = 0; i < item_ents_len(bx[n]->items); i++)
            printf("inv %d %d\n", bx[n]->items[i]->item, bx[n]->items[i]->qty);
        if (rp_player(n)) {
            printf("output:");
            for (i = 1; i < MAX_BOXES; i++) if (test_bit(rp_player(n)->output, i)) printf(" %d", i);
            printf("\n");
        }
    } next_box;
}

static void dump(exit_views_list l)
{
    int i;
    printf("exits %d\n", exit_views_len(l));
    for (i = 0; i < exit_views_len(l); i++) {
        struct exit_view *v = l[i];
        printf("exit %d %d %d %d %d %d %d %d %d %d %d %d %d %d\n",
               v->direction, v->destination, v->orig, v->distance,
               v->impassable, v->dest_hidden, v->orig_hidden, v->hidden,
               v->inside, v->road, v->water, v->in_transit,
               v->magic_barrier, v->hades_cost);
    }
}

int main(int argc, char **argv)
{
    char line[1024], op[32];
    struct command cmd = {0};
    struct exit_view v = {0};
    int target, section = -1;
    FILE *fp;
    if (argc != 4) return 2;
    target = atoi(argv[2]);
    fp = fopen(argv[1], "r");
    if (!fp) return 2;
    glob_init();
    init_lower();
    libdir = str_save(argv[3]);
    *strrchr(libdir, '/') = 0;
    while (fgets(line, sizeof line, fp)) {
        int a = 0, b = 0, c = 0, d = 0;
        line[strcspn(line, "\r\n")] = 0;
        if (!line[0] || line[0] == '#') continue;
        if (strncmp(line, "section ", 8) == 0) {
            if (++section > target) break;
            if (section == target) puts(line);
            continue;
        }
        if (section != target) continue;
        printf("> %s\n", line);
        sscanf(line, "%31s %d %d %d %d", op, &a, &b, &c, &d);
        if (!strcmp(op, "argument")) {
            cstrings_clear(&cmd.parse);
            cstrings_append(&cmd.parse, str_save("move"));
            cstrings_append(&cmd.parse, str_save(line + 9));
        } else if (!strcmp(op, "command")) { cmd.who = a; cmd.a = b; cmd.wait = c; cmd.cmd = d; }
        else if (!strcmp(op, "args")) { cmd.b = a; cmd.c = b; cmd.d = c; cmd.e = d; }
        else if (!strcmp(op, "tail")) { cmd.f = a; cmd.g = b; cmd.h = c; cmd.days_executing = d; }
        else if (!strcmp(op, "exit1")) { v.direction = a; v.destination = b; v.orig = c; v.distance = d; }
        else if (!strcmp(op, "exit2")) { v.impassable = a; v.dest_hidden = b; v.orig_hidden = c; v.hidden = d; }
        else if (!strcmp(op, "exit3")) { v.inside = a; v.road = b; v.water = c; v.in_transit = d; }
        else if (!strcmp(op, "exit4")) { v.magic_barrier = a; v.hades_cost = b; }
        else if (!strcmp(op, "parse")) {
            struct exit_view *ret = parse_exit_dir(&cmd, a, b ? "move" : NULL);
            if (ret) { v = *ret; printf("parsed\n"); }
            else printf("no exit\n");
        }
        else if (!strcmp(op, "check")) printf("check %d\n", drv_check(&cmd, &v, a));
        else if (!strcmp(op, "can")) printf("can %d\n", drv_can(a, &cmd));
        else if (!strcmp(op, "outer")) printf("outer %d\n", drv_outer(a, &cmd));
        else if (!strcmp(op, "save")) drv_save(&cmd, &v);
        else if (!strcmp(op, "restore")) drv_restore(&cmd, &v);
        else if (!strcmp(op, "dump_exit")) {
            exit_views_list l = NULL;
            exit_views_append(&l, &v); dump(l); exit_views_reclaim(&l);
        }
        else if (!strcmp(op, "dump_command")) printf("command %d %d %d %d %d %d %d %d %d %d %d %d\n", cmd.who, cmd.a, cmd.b, cmd.c, cmd.d, cmd.e, cmd.f, cmd.g, cmd.h, cmd.wait, cmd.cmd, cmd.days_executing);
        else if (!strcmp(op, "voyage")) {
            struct command *p = p_command(a);
            *p = cmd;
            p->who = a; p->state = b; p->wait = c;
            p->cmd = find_command(d ? "sail" : "move");
        }
        else if (!strcmp(op, "remove")) printf("remove %d\n", drv_sub_item(a, b, c));
        else if (!strcmp(op, "lose")) check_captain_loses_sailors(a, b, c);
        else if (!strcmp(op, "voyage_water")) {
            struct exit_view v = {0};
            v.destination = b; v.distance = c;
            printf("voyage_water %d\n", drv_water(rp_command(a), &v, b, d));
        }
        else if (!strcmp(op, "vector")) vector_add(a);
        else if (!strcmp(op, "voyage_state")) {
            int n, i;
            extern ilist out_vector;
            printf("vector:");
            for (i = 0; i < ilist_len(out_vector); i++) printf(" %d", out_vector[i]);
            printf("\n");
            loop_boxes(n) {
                struct command *p = rp_command(n);
                if (p) printf("voyage %d %d %d %d %d %d %d %d %d %d %d %d %d\n", p->who, p->state, p->a, p->b, p->c, p->d, p->e, p->f, p->g, p->h, p->wait, p->cmd, p->days_executing);
                if (rp_subloc(n)) printf("ship_moving %d %d\n", n, ship_moving(n));
            } next_box;
        }
        else if (!strcmp(op, "show_day")) show_day = a;
        else if (!strcmp(op, "frame")) { indent = a; second_indent = b; }
        else if (!strcmp(op, "garrison_view")) show_to_garrison = a;
        else if (!strcmp(op, "touch")) touch_loc_pl(a, b);
        else if (!strcmp(op, "touched")) set_bit(&p_player(a)->locs, b);
        else if (!strcmp(op, "pledge")) p_magic(a)->pledge = b;
        else if (!strcmp(op, "castle")) p_misc(a)->garr_castle = b;
        else if (!strcmp(op, "emit")) out(a, "routing %d  ", b);
        else if (!strcmp(op, "rulers")) {
            ilist l = players_who_rule_here(a); int i;
            printf("rulers:");
            for (i = 0; i < ilist_len(l); i++) printf(" %d", l[i]);
            printf("\n");
        }
        else if (!strcmp(op, "routing_state")) {
            int n, i;
            loop_boxes(n) {
                if (rp_misc(n)) printf("castle %d %d\n", n, garrison_castle(n));
                if (rp_magic(n)) printf("pledge %d %d\n", n, char_pledge(n));
                printf("player_data %d %d\n", n, rp_player(n) != NULL);
                if (rp_player(n)) {
                    printf("touched %d:", n);
                    for (i = 0; i < ilist_len(rp_player(n)->locs); i++) printf(" %d", rp_player(n)->locs[i]);
                    printf("\nknown %d:", n);
                    for (i = 0; i < ilist_len(rp_player(n)->known); i++) printf(" %d", rp_player(n)->known[i]);
                    printf("\noutput_order %d:", n);
                    for (i = 0; i < ilist_len(rp_player(n)->output); i++) printf(" %d", rp_player(n)->output[i]);
                    printf("\n");
                }
            } next_box;
        }
        else if (!strcmp(op, "see_all")) immed_see_all = a;
        else if (!strcmp(op, "admission")) {
            struct admit *p = my_malloc(sizeof(*p));
            p->targ = b; p->sense = c; admits_append(&p_player(a)->admits, p);
        }
        else if (!strcmp(op, "entry")) ilist_append(&p_player(a)->admits[b]->l, c);
        else if (!strcmp(op, "admit")) printf("admit %d\n", will_admit(a, b, c));
        else if (!strcmp(op, "default_garr")) p_magic(a)->default_garr = b;
        else if (!strcmp(op, "guard")) p_char(a)->guard = b;
        else if (!strcmp(op, "char_moving")) p_char(a)->moving = b;
        else if (!strcmp(op, "swear_release")) p_magic(a)->swear_on_release = b;
        else if (!strcmp(op, "suspend")) drv_suspend(a);
        else if (!strcmp(op, "resume")) restore_stack_actions(a);
        else if (!strcmp(op, "clear_guard")) clear_guard_flag(a);
        else if (!strcmp(op, "leave")) leave_stack(a);
        else if (!strcmp(op, "unstack")) drv_unstack(a);
        else if (!strcmp(op, "movement_state")) {
            int n, i, j;
            loop_boxes(n) {
                if (rp_char(n)) printf("action %d %d %d %d %d\n", n, char_moving(n), char_guard(n), loyal_kind(n), loyal_rate(n));
                if (rp_magic(n)) printf("magic %d %d %d\n", n, default_garrison(n), rp_magic(n)->swear_on_release);
                if (rp_player(n)) {
                    printf("player %d admits %d\n", n, admits_len(rp_player(n)->admits));
                    for (i = 0; i < admits_len(rp_player(n)->admits); i++) {
                        struct admit *p = rp_player(n)->admits[i];
                        printf("permission %d %d %d", p->targ, p->sense, p->flag);
                        for (j = 0; j < ilist_len(p->l); j++) printf(" %d", p->l[j]);
                        printf("\n");
                    }
                }
            } next_box;
        }
        else if (!strcmp(op, "box")) {
            alloc_box(a, b, c);
            rp_loc_info(a)->where = d;
            if (d) ilist_append(&rp_loc_info(d)->here_list, a);
        } else if (!strcmp(op, "dest")) ilist_append(&p_loc(a)->prov_dest, b);
        else if (!strcmp(op, "state")) state_dump();
        else if (!strcmp(op, "output")) immediate = a;
        else if (!strcmp(op, "output_bit")) printf("output_bit %d\n", test_bit(p_player(a)->output, b));
        else if (!strcmp(op, "report")) {
            FILE *log; int ch;
            printf("report %d:\n", a);
            log = fopen(sout("%s/log/%d", libdir, a), "r");
            if (log) { while ((ch = fgetc(log)) != EOF) putchar(ch); fclose(log); }
        }
        else if (!strcmp(op, "item")) {
            p_item(a)->weight = b; p_item(a)->land_cap = c; p_item(a)->ride_cap = d;
        }
        else if (!strcmp(op, "air")) { p_item(a)->fly_cap = b; p_item(a)->animal = c; p_item(a)->is_man_item = d; }
        else if (!strcmp(op, "fighter")) { p_item(a)->attack = b; p_item(a)->defense = c; p_item(a)->missile = d; }
        else if (!strcmp(op, "rank")) p_char(a)->rank = b;
        else if (!strcmp(op, "rank_s")) printf("rank [%s]\n", rank_s(a));
        else if (!strcmp(op, "aura")) { p_magic(a)->max_aura = b; p_magic(a)->auraculum = c; }
        else if (!strcmp(op, "aura_item")) { p_item_magic(a)->aura = b; p_item_magic(a)->aura_bonus = c; }
        else if (!strcmp(op, "bonuses")) { p_item_magic(a)->attack_bonus = b; p_item_magic(a)->defense_bonus = c; p_item_magic(a)->missile_bonus = d; }
        else if (!strcmp(op, "aura_result")) printf("aura %d %d\n", has_auraculum(a), max_eff_aura(a));
        else if (!strcmp(op, "wield")) {
            struct wield w = {91, -83, 77};
            struct fight f = {11, -12, 13, -14, 15, 101, -203, 17, 307, -19, 23, -29, 31, -37, 41, -43, 47, -53, 59};
            int result = find_wield(b ? &w : NULL, a, c ? &f : NULL);
            printf("wield %d %d %d %d\n", result, w.attack, w.defense, w.missile);
            fight_dump(&f);
            printf("wield_s [%s]\n", wield_s(a));
        }
        else if (!strcmp(op, "description_state")) {
            int n;
            state_dump();
            loop_boxes(n) {
                printf("description %d %d %d %d %d\n", n, rp_char(n) != NULL, char_rank(n), rp_magic(n) != NULL, rp_item_magic(n) != NULL);
                if (rp_magic(n)) printf("aura %d %d\n", char_max_aura(n), char_auraculum(n));
                if (rp_item_magic(n)) printf("bonuses %d %d %d %d %d\n", item_aura(n), item_aura_bonus(n), item_attack_bonus(n), item_defense_bonus(n), item_missile_bonus(n));
            } next_box;
        }
        else if (!strcmp(op, "unit")) { p_char(a)->unit_item = b; p_char(a)->prisoner = c; p_char(a)->loy_kind = d; }
        else if (!strcmp(op, "inv")) {
            struct item_ent *e = NULL; int i;
            for (i = 0; i < item_ents_len(bx[a]->items); i++)
                if (bx[a]->items[i]->item == b) { e = bx[a]->items[i]; break; }
            if (!e) { e = my_malloc(sizeof(*e)); e->item = b; item_ents_append(&bx[a]->items, e); }
            e->qty = c;
        }
        else if (!strcmp(op, "factor")) army_slow_factor = a;
        else if (!strcmp(op, "strength")) p_misc(a)->storm_str = b;
        else if (!strcmp(op, "unit_weights") || !strcmp(op, "stack_weights")) {
            struct weights w;
            if (!strcmp(op, "unit_weights")) determine_unit_weights(a, &w);
            else determine_stack_weights(a, &w);
            weights_dump(&w);
        }
        else if (!strcmp(op, "counts")) printf("counts %d %d\n", count_stack_move_nobles(a), count_stack_figures(a));
        else if (!strcmp(op, "land") || !strcmp(op, "fly") || !strcmp(op, "water")) {
            struct command cmd = {0}; struct exit_view v = {0}; int result;
            cmd.who = a; v.destination = b; v.distance = c;
            if (!strcmp(op, "land")) result = drv_land(&cmd, &v, d);
            else if (!strcmp(op, "fly")) result = drv_fly(&cmd, &v, d);
            else result = drv_water(&cmd, &v, b, d);
            printf("%s %d\n", op, result);
            printf("unchanged %d %d %d %d\n", cmd.who, v.destination, v.distance, cmd.wait);
        }
        else if (!strcmp(op, "hidden")) p_loc(a)->hidden = b;
        else if (!strcmp(op, "barrier")) p_loc(a)->barrier = b;
        else if (!strcmp(op, "lane")) p_loc(a)->sea_lane = b;
        else if (!strcmp(op, "uldim")) p_subloc(a)->uldim_flag = b;
        else if (!strcmp(op, "summer")) p_subloc(a)->summer_flag = b;
        else if (!strcmp(op, "open")) p_subloc(a)->link_open = b;
        else if (!strcmp(op, "from")) ilist_append(&p_subloc(a)->link_from, b);
        else if (!strcmp(op, "to")) ilist_append(&p_subloc(a)->link_to, b);
        else if (!strcmp(op, "road")) { p_gate(a)->to_loc = b; p_gate(a)->road_hidden = c; }
        else if (!strcmp(op, "lord")) p_char(a)->unit_lord = b;
        else if (!strcmp(op, "known")) set_bit(&p_player(a)->known, b);
        else if (!strcmp(op, "hades")) hades_region = a;
        else if (!strcmp(op, "clock")) { sysclock.turn = a; sysclock.days_since_epoch = b; evening = c; month_done = d; }
        else if (!strcmp(op, "day")) sysclock.day = a;
        else if (!strcmp(op, "moving")) p_subloc(a)->moving = b;
        else if (!strcmp(op, "master")) {
            int i; unsigned char seed[16]; FILE *s;
            for (i = 0; i < 16; i++) seed[i] = a + i;
            s = fopen(argv[3], "wb"); if (!s) return 2;
            fwrite(seed, 1, 16, s); fclose(s); load_seed(argv[3]);
        } else if (!strcmp(op, "legacy")) printf("legacy %d\n", rnd(a, b));
        else if (!strcmp(op, "direction")) printf("direction %d\n", location_direction(a, b));
        else if (!strcmp(op, "distance")) printf("distance %d\n", exit_distance(a, b));
        else if (!strcmp(op, "port")) printf("port %d\n", is_port_city(a));
        else if (!strcmp(op, "province_port")) printf("province_port %d\n", province_has_port_city(a));
        else if (!strcmp(op, "ocean")) printf("ocean %d\n", has_ocean_access(a));
        else if (!strcmp(op, "exits")) dump(exits_from_loc(a, b));
        else if (!strcmp(op, "nsew")) dump(exits_from_loc_nsew(a, b));
        else if (!strcmp(op, "select")) dump(exits_from_loc_nsew_select(a, b, c, d));
        else { fprintf(stderr, "unknown operation: %s\n", op); return 2; }
    }
    fclose(fp);
    return 0;
}
