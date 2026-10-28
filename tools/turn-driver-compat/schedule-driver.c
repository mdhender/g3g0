/* Copyright (c) 2026 Michael D Henderson.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Schedule-only instrumentation. No gameplay callback is implemented here.
 * schedule.c contains the original C stream/static definitions and daily body.
 */
#include <string.h>
#include "schedule.c"

static void print_stream(const char *name, const rng_stream *s, int turn)
{
    printf("%s: turn=%d seed=%08x %08x %08x %08x counter=%u\n",
           name, turn, s->seed[0], s->seed[1], s->seed[2], s->seed[3], s->counter);
}

int main(int argc, char **argv)
{
    char line[1024], op[32];
    int section = -1, want;
    FILE *fp;
    if (argc != 3 || !(fp = fopen(argv[1], "r"))) return 2;
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
            int turn, day;
            if (sscanf(line, "%*s %d %d", &turn, &day) != 2) return 2;
            sysclock.turn = (short)turn;
            sysclock.day = (short)day;
        } else if (!strcmp(op, "seed")) {
            unsigned char bytes[16];
            int first;
            if (sscanf(line, "%*s %d", &first) != 1) return 2;
            for (int i = 0; i < 16; i++) bytes[i] = (unsigned char)(first + i);
            FILE *seed = fopen("schedule-seed", "wb");
            if (!seed || fwrite(bytes, 1, 16, seed) != 16 || fclose(seed)) return 2;
            load_seed("schedule-seed");
        } else if (!strcmp(op, "cal")) {
            int pick, low, high;
            if (sscanf(line, "%*s %d %d %d", &pick, &low, &high) != 3) return 2;
            observed_cal_day(pick, low, high);
        } else if (!strcmp(op, "shuffle")) {
            int n;
            ilist l = NULL;
            if (sscanf(line, "%*s %d", &n) != 1 || n < 0 || n > 30) return 2;
            for (int i = 1; i <= n; i++) ilist_append(&l, i);
            observed_shuffle(l);
            ilist_reclaim(&l);
        } else if (!strcmp(op, "rnd")) {
            int low, high;
            if (sscanf(line, "%*s %d %d", &low, &high) != 2) return 2;
            printf("rnd -> %d\n", rnd(low, high));
        } else if (!strcmp(op, "streams")) {
            uint32_t m[4];
            rng_master_seed(m);
            printf("master: %08x %08x %08x %08x\n", m[0], m[1], m[2], m[3]);
            print_stream("calendar", &calendar_rng, calendar_rng_turn);
            print_stream("weather", &weather_rng, weather_rng_turn);
        } else if (!strcmp(op, "daily_events")) daily_events();
        else if (!strcmp(op, "days")) {
            int first, last;
            if (sscanf(line, "%*s %d %d", &first, &last) != 2 || last - first > 31) return 2;
            for (int d = first; d <= last; d++) {
                sysclock.day = (short)d;
                daily_events();
            }
        } else return 2;
    }
    fclose(fp);
    remove("schedule-seed");
    return 0;
}
