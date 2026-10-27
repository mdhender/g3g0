/* Copyright (c) 2026 Michael D Henderson.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Instrumentation only: real process_orders/daily_events/post_month bodies
 * call tracing phase callbacks. This is not a gameplay implementation.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "z.h"
#include "oly.h"

int time_self, cmd_wait = -1, auto_attack_flag, month_done;
int post_has_been_run, auto_quit_turns;
olytime sysclock;
static struct box *boxes[203];
struct box **bx = boxes;
int box_head[T_MAX];

int main(int argc, char **argv)
{
    char line[1024], op[32];
    int section = -1, want;
    FILE *fp;
    if (argc != 3 || !(fp = fopen(argv[1], "r"))) return 2;
    /* Minimal fixture: two player records and a NULL player record between.
     * The real clear_orders_sent body and real loop macros mutate these. */
    box_head[T_player] = 200;
    for (int n = 200; n <= 202; n++) {
        boxes[n] = calloc(1, sizeof(struct box));
        boxes[n]->x_next_kind = n == 202 ? 0 : n + 1;
        if (n != 201) {
            boxes[n]->x_player = calloc(1, sizeof(struct entity_player));
            boxes[n]->x_player->sent_orders = n == 200 ? 7 : -5;
        }
    }
    want = atoi(argv[2]);
    while (fgets(line, sizeof(line), fp)) {
        line[strcspn(line, "\n")] = 0;
        if (!line[0] || line[0] == '#') continue;
        if (!strncmp(line, "section ", 8)) {
            if (++section == want) printf("%s\n", line);
            continue;
        }
        if (section != want) continue;
        printf("> %s\n", line);
        if (sscanf(line, "%31s", op) != 1) return 2;
        if (!strcmp(op, "clock")) {
            int turn, day, epoch;
            if (sscanf(line, "%*s %d %d %d", &turn, &day, &epoch) != 3) return 2;
            sysclock.turn = (short)turn;
            sysclock.day = (short)day;
            sysclock.days_since_epoch = epoch;
        } else if (!strcmp(op, "process_orders")) process_orders();
        else if (!strcmp(op, "daily_events")) daily_events();
        else if (!strcmp(op, "post_month")) post_month();
        else if (!strcmp(op, "auto_quit")) {
            if (sscanf(line, "%*s %d", &auto_quit_turns) != 1) return 2;
        } else if (!strcmp(op, "state")) {
            printf("state: %d %d %d wait=%d auto=%d month=%d post=%d\n",
                   sysclock.turn, sysclock.day, sysclock.days_since_epoch,
                   cmd_wait, auto_attack_flag, month_done, post_has_been_run);
            printf("sent: %d %d\n", rp_player(200)->sent_orders, rp_player(202)->sent_orders);
        } else return 2;
    }
    fclose(fp);
    return 0;
}
