/* Copyright (c) 2026 Michael D Henderson.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Small scripted calls into the real pinned routines extracted by run.sh.
 * stdout and stderr are captured separately; no stream normalization.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "z.h"
#include "oly.h"

int time_self = 0;

int
main(int argc, char **argv)
{
	char line[1024], op[32];
	int section = -1, want, n, i, c;
	olytime clock = {0};
	FILE *fp;
	if (argc != 3 || !(fp = fopen(argv[1], "r")))
		return 2;
	want = atoi(argv[2]);
	while (fgets(line, sizeof(line), fp))
	{
		line[strcspn(line, "\n")] = 0;
		if (!line[0] || line[0] == '#')
			continue;
		if (!strncmp(line, "section ", 8))
		{
			section++;
			if (section == want)
				printf("%s\n", line);
			continue;
		}
		if (section != want)
			continue;
		printf("> %s\n", line);
		if (sscanf(line, "%31s", op) != 1)
			return 2;
		if (!strcmp(op, "clock"))
		{
			int turn, day, epoch;
			if (sscanf(line, "%*s %d %d %d", &turn, &day, &epoch) != 3)
				return 2;
			clock.turn = (short)turn;
			clock.day = (short)day;
			clock.days_since_epoch = epoch;
		}
		else if (!strcmp(op, "increment"))
		{
			if (sscanf(line, "%*s %d", &n) != 1) return 2;
			for (i = 0; i < n; i++) olytime_increment(&clock);
		}
		else if (!strcmp(op, "turn_change"))
			olytime_turn_change(&clock);
		else if (!strcmp(op, "dump"))
			printf("clock: %d %d %d; month %d year %d\n", clock.turn,
			       clock.day, clock.days_since_epoch, oly_month(clock), oly_year(clock));
		else if (!strcmp(op, "dot"))
		{
			if (sscanf(line, "%*s %d %d", &n, &c) != 2) return 2;
			for (i = 0; i < n; i++) print_dot(c);
		}
		else if (!strcmp(op, "time_self"))
		{
			if (sscanf(line, "%*s %d", &time_self) != 1) return 2;
		}
		else if (!strcmp(op, "stage"))
			stage(line + 6);
		else if (!strcmp(op, "stage_null"))
			stage(NULL);
		else if (!strcmp(op, "stage_empty"))
			stage("");
		else return 2;
	}
	fclose(fp);
	return 0;
}
