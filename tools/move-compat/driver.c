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
