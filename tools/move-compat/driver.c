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

void glob_init(void);
int hades_region, month_done, evening;
int immediate, flush_always = 1;
int drv_land(struct command *, struct exit_view *, int);
int drv_fly(struct command *, struct exit_view *, int);
int drv_water(struct command *, struct exit_view *, int, int);

/* Unreachable report routes must abort, never decide gameplay. */
ilist players_who_rule_here(int n)
{
    fprintf(stderr, "out-of-scope C dependency: players_who_rule_here\n");
    abort();
}

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
    int target, section = -1;
    FILE *fp;
    if (argc != 4) return 2;
    target = atoi(argv[2]);
    fp = fopen(argv[1], "r");
    if (!fp) return 2;
    glob_init();
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
        if (!strcmp(op, "box")) {
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
