/* Copyright (c) 2026 Michael D Henderson.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Scripted state construction only; all compared routines are real C.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "z.h"
#include "oly.h"

void glob_init(void);
void init_lower(void);
int drv_load(struct command *c);
int find_command(char *s);
void drv_state(struct command *c, int s);
struct order_list *drv_orders(int pl, int who);
extern int cur_pri, cmd_wait;
extern ilist load_q[5], run_q;
extern int evening, auto_attack_flag;
extern ilist second_wait_list;
void drv_evening(void), drv_daily(void), drv_players(void), drv_sort_run(void);
void drv_sort_load(int pri);

static void list(char *label, ilist l)
{
 int i;
 printf("%s:", label);
 for (i = 0; i < ilist_len(l); i++) printf(" %d", l[i]);
 printf("\n");
}

int main(int argc, char **argv)
{
 FILE *fp, *log;
 char line[8192], op[32], text[8192], label[32];
 int section = -1, want, n, a, b, d, i, result;
 struct command *c;
 if (argc != 4 || !(fp = fopen(argv[1], "r"))) return 2;
 want = atoi(argv[2]);
 libdir = argv[3];
 glob_init();
 init_lower();
 alloc_box(200, T_player, sub_pl_regular);
 alloc_box(201, T_player, sub_pl_regular);
 p_player(200);
 p_player(201);
 alloc_box(1000, T_char, 0);
 alloc_box(1001, T_char, 0);
 alloc_box(sk_hide_self, T_skill, 0);
 p_char(1000)->unit_lord = 200;
 p_char(1001)->unit_lord = 200;
 c = p_command(1000);
 c->who = 1000;
 c->state = STATE_DONE;
 while (fgets(line, sizeof(line), fp))
 {
  line[strcspn(line, "\n")] = 0;
  if (!line[0] || line[0] == '#') continue;
  if (!strncmp(line, "section ", 8))
  {
   section++;
   if (section == want) printf("%s\n", line);
   continue;
  }
  if (section != want) continue;
  printf("> %s\n", line);
  sscanf(line, "%31s", op);
  if (!strcmp(op, "queue"))
   cstrings_append(&drv_orders(player(c->who), c->who)->l, str_save(line + 6));
  else if (!strcmp(op, "load")) printf("load %d\n", drv_load(c));
  else if (!strcmp(op, "do")) do_command(c);
  else if (!strcmp(op, "finish")) printf("finish %d\n", finish_command(c));
  else if (!strcmp(op, "done")) command_done(c);
  else if (!strcmp(op, "shift")) cmd_shift(c);
  else if (!strcmp(op, "entity"))
  { sscanf(line,"%*s %d %d %d",&n,&a,&b); alloc_box(n,a,b); }
  else if (!strcmp(op, "link"))
  {
   sscanf(line,"%*s %d %d",&n,&a);
   if (loc(a)) ilist_rem_value(&p_loc_info(loc(a))->here_list,a);
   p_loc_info(a)->where = n;
   ilist_append(&p_loc_info(n)->here_list,a);
  }
  else if (!strcmp(op, "inventory"))
  {
   struct item_ent *e = my_malloc(sizeof(*e));
   sscanf(line,"%*s %d %d %d",&n,&a,&b);
   e->item = a; e->qty = b; item_ents_append(&bx[n]->items,e);
  }
  else if (!strcmp(op, "field"))
  {
   sscanf(line,"%*s %d %31s %d",&n,text,&a);
   if (!strcmp(text,"lord")) p_char(n)->unit_lord = a;
   else if (!strcmp(text,"npc")) p_char(n)->npc_prog = (schar)a;
   else if (!strcmp(text,"loyal")) p_char(n)->loy_kind = (schar)a;
   else if (!strcmp(text,"moving")) p_char(n)->moving = a;
   else if (!strcmp(text,"hidden")) p_magic(n)->hide_self = (schar)a;
   else if (!strcmp(text,"subkind")) bx[n]->skind = a;
   else if (!strcmp(text,"vulnerable")) p_misc(n)->only_vulnerable = a;
   else if (!strcmp(text,"strength")) p_misc(n)->storm_str = (short)a;
   else if (!strcmp(text,"safe")) p_subloc(n)->safe = (schar)a;
   else if (!strcmp(text,"loc_hidden")) p_loc(n)->hidden = (schar)a;
   else if (!strcmp(text,"hostile")) ilist_append(&p_disp(n)->hostile,a);
   else if (!strcmp(text,"neutral")) ilist_append(&p_disp(n)->neutral,a);
   else if (!strcmp(text,"defend")) ilist_append(&p_disp(n)->defend,a);
   else if (!strcmp(text,"contact")) ilist_append(&p_char(n)->contact,a);
   else if (!strcmp(text,"garr_host")) ilist_append(&p_misc(n)->garr_host,a);
   else if (!strcmp(text,"known")) set_bit(&p_player(n)->known,a);
   else return 2;
  }
  else if (!strcmp(op, "predicate"))
  {
   n = a = b = 0;
   sscanf(line,"%*s %31s %d %d %d",text,&n,&a,&b);
   if (!strcmp(text,"hostile")) result = is_hostile(n,a);
   else if (!strcmp(text,"contacted")) result = contacted(n,a);
   else if (!strcmp(text,"char_where")) result = char_where(n,a,b);
   else if (!strcmp(text,"char_here")) result = char_here(n,a);
   else if (!strcmp(text,"known")) result = test_known(n,a);
   else if (!strcmp(text,"loc_hidden")) result = loc_hidden(n);
   else if (!strcmp(text,"has_item")) result = has_item(n,a);
   else if (!strcmp(text,"safe")) result = in_safe_now(n);
   else if (!strcmp(text,"province")) result = province(n);
   else if (!strcmp(text,"owner")) result = building_owner(n);
   else if (!strcmp(text,"weather")) result = weather_here(n,a);
   else if (!strcmp(text,"really_hidden")) result = char_really_hidden(n);
   else if (!strcmp(text,"npc")) result = is_npc(n);
   else if (!strcmp(text,"vulnerable")) result = only_defeatable(n);
   else if (!strcmp(text,"char_alloc")) result = rp_char(n) != NULL;
   else return 2;
   printf("predicate %d\n",result);
  }
  else if (!strcmp(op, "initial")) initial_command_load();
  else if (!strcmp(op, "start")) start_phase();
  else if (!strcmp(op, "evening")) drv_evening();
  else if (!strcmp(op, "daily")) drv_daily();
  else if (!strcmp(op, "players")) drv_players();
  else if (!strcmp(op, "auto")) check_all_auto_attacks();
  else if (!strcmp(op, "empty_chars"))
  { change_box_kind(1000,T_deleted); change_box_kind(1001,T_deleted); }
  else if (!strcmp(op, "min")) printf("min %d\n", min_pri_ready());
  else if (!strcmp(op, "sort_load")) { sscanf(line,"%*s %d",&n); drv_sort_load(n); }
  else if (!strcmp(op, "sort_run")) drv_sort_run();
  else if (!strcmp(op, "phase"))
   printf("phase evening %d auto %d second %d\n",evening,auto_attack_flag,ilist_len(second_wait_list));
  else if (!strcmp(op, "second"))
  { c->second_wait = TRUE; ilist_append(&second_wait_list,c->who); }
  else if (!strcmp(op, "select")) { sscanf(line,"%*s %d",&n); c = rp_command(n); }
  else if (!strcmp(op, "place"))
  {
   sscanf(line,"%*s %d",&n);
   if (!valid_box(n)) alloc_box(n,T_loc,sub_plain);
   p_loc_info(c->who)->where = n;
   ilist_append(&p_loc_info(n)->here_list,c->who);
  }
  else if (!strcmp(op, "resolve_wait")) cmd_wait = find_command("wait");
  else if (!strcmp(op, "skill"))
  {
   struct skill_ent *s = my_malloc(sizeof(*s));
   sscanf(line, "%*s %d %d", &a, &b);
   s->skill = sk_hide_self; s->know = (schar)a; s->experience = (short)b;
   skill_ents_append(&p_char(c->who)->skills, s);
  }
  else if (!strcmp(op, "parse"))
  {
   strcpy(text, line + 6);
   result = oly_parse(c, text);
   printf("parse %d\n", result);
  }
  else if (!strcmp(op, "allow"))
  {
   strcpy(text, line + 6);
   printf("allow %d\n", check_allow(c, !strcmp(text, "NULL") ? NULL : !strcmp(text, "EMPTY") ? "" : text));
  }
  else if (!strcmp(op, "who"))
  {
   sscanf(line, "%*s %d", &n);
   c = p_command(n);
   c->who = n;
   c->state = STATE_DONE;
  }
  else if (!strcmp(op, "run"))
  {
   sscanf(line, "%*s %d %d %d %d", &a, &b, &d, &i);
   c->wait = a; c->poll = (schar)b; c->status = (schar)d; c->inhibit_finish = (schar)i;
   c->pri = (schar)cmd_tbl[c->cmd].pri;
   drv_state(c, STATE_RUN);
  }
  else if (!strcmp(op, "set"))
  {
   sscanf(line, "%*s %31s %d", text, &n);
   if (!strcmp(text, "immediate")) immediate = n;
   else if (!strcmp(text, "show_day")) show_day = n;
   else if (!strcmp(text, "day")) sysclock.day = (short)n;
   else if (!strcmp(text, "turn")) sysclock.turn = (short)n;
   else if (!strcmp(text, "priority")) cur_pri = n;
   else if (!strcmp(text, "status")) c->status = (schar)n;
   else if (!strcmp(text, "debug")) c->debug = (schar)n;
   else if (!strcmp(text, "days")) c->days_executing = n;
   else if (!strcmp(text, "restrict")) p_misc(c->who)->cmd_allow = (schar)n;
   else if (!strcmp(text, "lord")) p_char(c->who)->unit_lord = n;
   else if (!strcmp(text, "moving")) p_char(c->who)->moving = n;
   else if (!strcmp(text, "parent")) p_loc_info(c->who)->where = n;
   else if (!strcmp(text, "kind")) bx[c->who]->kind = n;
   else if (!strcmp(text, "waitcmd")) cmd_wait = n;
   else if (!strcmp(text, "count")) p_player(player(c->who))->cmd_count = (short)n;
   else if (!strcmp(text, "prisoner")) p_char(c->who)->prisoner = (schar)n;
   else if (!strcmp(text, "health")) p_char(c->who)->health = (schar)n;
   else if (!strcmp(text, "second")) c->second_wait = (schar)n;
   else if (!strcmp(text, "cmdpri")) c->pri = (schar)n;
   else if (!strcmp(text, "state")) c->state = (schar)n;
   else if (!strcmp(text, "auto")) auto_attack_flag = n;
   else return 2;
  }
  else if (!strcmp(op, "table"))
  {
   for (i = 0; cmd_tbl[i].name; i++)
    printf("handler %d %s %d %d %d\n", i, cmd_tbl[i].name,
           cmd_tbl[i].start != NULL, cmd_tbl[i].finish != NULL, cmd_tbl[i].interrupt != NULL);
  }
  else if (!strcmp(op, "dump"))
  {
   int pl = player(c->who);
   struct entity_player *p = rp_player(pl);
   struct order_list *q = drv_orders(pl, c->who);
   printf("cmd %d who %d state %d status %d wait %d poll %d pri %d cond %d fuzzy %d days %d debug %d inhibit %d cur %d\n",
     c->cmd, c->who, c->state, c->status, c->wait, c->poll, c->pri,
     c->conditional, c->fuzzy, c->days_executing, c->debug, c->inhibit_finish, cur_pri);
   printf("args %d %d %d %d %d %d %d %d; parse", c->a,c->b,c->c,c->d,c->e,c->f,c->g,c->h);
   for (i = 0; i < cstrings_len(c->parse); i++) printf(" [%s]", c->parse[i]);
   printf("\nplayer %d count %d last %d format %d notab %d output %d hidden %d orders %d\n",
     pl,p->cmd_count,p->last_order_turn,p->format,p->notab,test_bit(p->output,c->who),
     char_hidden(c->who),cstrings_len(q->l));
   for (i = 0; i < 5; i++) { sprintf(label,"load_q %d",i); list(label,load_q[i]); }
   list("run_q", run_q);
  }
  else if (!strcmp(op, "log"))
  {
   fflush(NULL);
   log = fopen(sout("%s/log/%d",libdir,player(c->who)), "r");
   printf("log:\n");
   if (log) { while ((n = fgetc(log)) != EOF) putchar(n); fclose(log); }
  }
  else if (!strcmp(op, "null_do")) do_command(NULL);
  else if (!strcmp(op, "null_finish")) finish_command(NULL);
  else if (!strcmp(op, "null_load")) drv_load(NULL);
  else if (!strcmp(op, "invalid_finish"))
  {
   struct command invalid = {0};
   invalid.who = 9999;
   finish_command(&invalid);
  }
  else return 2;
 }
 close_logfile();
 fclose(fp);
 return 0;
}
