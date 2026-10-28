/* Copyright (c) 2026 Michael D Henderson.
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Scripted state construction only; all compared routines are real C.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>
#include "z.h"
#include "oly.h"
#include "rng.h"

void glob_init(void);
void init_lower(void);
int drv_load(struct command *c);
int find_command(char *s);
void drv_state(struct command *c, int s);
struct order_list *drv_orders(int pl, int who);
extern int cur_pri, cmd_wait;
extern ilist load_q[5], run_q;
extern int evening, auto_attack_flag;
extern ilist second_wait_list, wait_list;
ilist drv_collectors(void);
void drv_evening(void), drv_daily(void), drv_players(void), drv_sort_run(void);
void drv_sort_load(int pri);
int drv_fort_covers(int n);
char *drv_wait_check(struct command *c);
void drv_signal(int who, char *flag);
void drv_interrupted(void);
int drv_sub_item(int who, int item, int qty);
void drv_add_item(int who, int item, int qty);
int drv_bandit(int faery, int who, int where);
int drv_market(int who);
void drv_ocean(int who);
ilist drv_oceans(void);
trades_list drv_offers(int seller, int where, int except);
void drv_entity_rolls(int who, int sub, int where);
void drv_nasty(int slot), drv_hades_unit(int who);
int drv_nasty_roll(int slot, int sub, int low, int high);
struct exit_view *choose_npc_direction(int who, int where, int dir);
int drv_savage_count(void), drv_savage_hates(int where), drv_savage_here(int where);
void drv_unsworn(int who), drv_mob(int who);
void drv_auto_bandit(int who, int prog);
void drv_garrison_pay(void);
void drv_sink(int ship), drv_collapse(int fort);
int drv_coast(int where);
void drv_coastal(void), drv_mine(int mine), drv_inn(int inn), drv_loc_damage(void);
void drv_heal(void);
int drv_heal_roll(int who, int sub, int low, int high);
void drv_dogs(void), drv_ring(void), drv_curse(void);
void drv_aura(void), drv_ability(void), drv_shroud(void), drv_hinder(void);
void drv_barrier(void), drv_loyalty(void), drv_pillage(void), drv_relic(void), drv_hide_mage(void);
int drv_loyal_roll(int who, int sub);
void drv_income_inn(void), drv_income_temple(void);
int drv_income_roll(int inn, int sub, int low, int high);
void drv_maint(void), drv_starve(int who, int have), drv_animal(void);
int drv_starve_roll(int who, int item), drv_animal_roll(int who, int item, int idx);
void drv_decay_mines(void), drv_decay_ghosts(void), drv_decay_corpses(void), drv_decay_bodies(void);
int drv_corpse_roll(int who);
void drv_auto_drop(void), drv_link_decay(void), drv_quest_decay(void), drv_kings(void);
void drv_summary_provinces(void);
int drv_rank_threshold(int n), drv_rank_lower(int n), drv_rank_sup(int n);
extern int gold_inn, gold_temple;
void drv_city_gold(void), drv_replenish(int where, int item, int qty);
void drv_taxes(void), drv_claim(void), drv_np(void), drv_unformed(void);
void drv_special_open(void), drv_special_close(void), drv_clear_sent(void);
extern int gold_taxes;
int drv_bark_roll(int who, int hound);
void drv_event_days(void);
extern int gold_garrison;
extern int gold_provinces;
extern int gold_combat, gold_combat_indep;
extern int gold_pot_basket, gold_trade, gold_opium;
int drv_reduce(struct trade *t, int cost);
void drv_attempt(struct trade *b, struct trade *s);
void drv_inherit_rolls(int from, int item, int qty);
void drv_storms(int num, int sk);
void drv_storm_decay(void), drv_storm_move(void), drv_storm_touch(void);
void drv_move_storm(int storm, int dest);
void drv_hunt(int slot), drv_faery_unit(int who);
int drv_hunt_loc(int slot, int retry, int size), drv_hunt_qty(int slot);
char *drv_fetch(FILE *fp);
int drv_admission(char *acct, FILE *fp);
void drv_make_players(void);

static void wait_dump(struct command *c)
{
 int i;
 printf("wait args %d\n", wait_args_len(c->wait_parse));
 for (i = 0; i < wait_args_len(c->wait_parse); i++) {
  struct wait_arg *p = c->wait_parse[i];
  printf("wait arg %d %d %d [%s]\n", p->tag, p->a1, p->a2, p->flag ? p->flag : "NULL");
 }
}

static void list(char *label, ilist l)
{
 int i;
 printf("%s:", label);
 for (i = 0; i < ilist_len(l); i++) printf(" %d", l[i]);
 printf("\n");
}

/* Actual auto_attack has already emitted its output. Observe it, then abort:
 * this is deliberately NOT a successful combat replacement. */
int regular_combat(int who, int target, int seize, int ally)
{
 FILE *log;
 int ch;
 fflush(NULL);
 printf("combat log:\n");
 log = fopen(sout("%s/log/%d", libdir, player(who)), "r");
 if (log) { while ((ch = fgetc(log)) != EOF) putchar(ch); fclose(log); }
 printf("combat output %d\n", test_bit(p_player(player(who))->output, who));
 fprintf(stderr, "combat boundary: regular_combat\n");
 fflush(NULL); /* abort flushes on Darwin too; no successful return */
 abort();
}

int main(int argc, char **argv)
{
 FILE *fp, *log;
 char line[8192], op[32], text[8192], label[32];
 int section = -1, want, n, a, b, d, i, result, created = 0;
 struct command *c;
 if (argc != 4 || !(fp = fopen(argv[1], "r"))) return 2;
 want = atoi(argv[2]);
 libdir = argv[3];
 glob_init();
 init_lower();
 init_spaces(); /* main.c call_init_routines: print_att's continuation padding */
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
  if (!strcmp(op, "master")) {
   unsigned char seed[16]; FILE *s;
   sscanf(line,"%*s %d",&n);
   for (i=0;i<16;i++) seed[i]=n+i;
   s=fopen(sout("%s/randseed",libdir),"wb"); if(!s) return 2;
   fwrite(seed,1,16,s); fclose(s); load_seed(sout("%s/randseed",libdir));
  }
  else if (!strcmp(op, "weather_map")) {
   sscanf(line,"%*s %d",&a);
   alloc_box(8000,T_loc,sub_region);
   for(i=a-1;i>=0;i--) {
    n=9000+i; alloc_box(n,T_loc,1+i%4); set_where(n,8000);
    set_bit(&p_player(200)->locs,n);
   }
  }
  else if (!strcmp(op, "weather")) {
   n=a=b=d=0;
   sscanf(line,"%*s %31s %d %d %d %d",text,&n,&a,&b,&d);
   if(!strcmp(text,"natural")) natural_weather();
   else if(!strcmp(text,"some")) drv_storms(n,a);
   else if(!strcmp(text,"new")) printf("new storm %d\n",new_storm(n,a,b,d));
   else if(!strcmp(text,"start")) { c->a=n; printf("bind start %d\n",v_bind_storm(c)); }
   else if(!strcmp(text,"finish")) { c->a=n; printf("bind finish %d\n",d_bind_storm(c)); }
   else if(!strcmp(text,"count")) printf("provinces %d\n",nprovinces());
   else if(!strcmp(text,"decay")) drv_storm_decay();
   else if(!strcmp(text,"monthly")) drv_storm_move();
   else if(!strcmp(text,"touch")) drv_storm_touch();
   else if(!strcmp(text,"move")) drv_move_storm(n,a);
   else if(!strcmp(text,"bound")) move_bound_storms(n,a);
   else if(!strcmp(text,"dissipate")) dissipate_storm(n,a);
   else if(!strcmp(text,"pending")) {
    struct entity_misc *p=rp_misc(n);
    printf("pending %d %d %d\n",n,p ? p->npc_dir : 0,p ? p->storm_move : 0);
   }
   else if(!strcmp(text,"state")) {
    loop_storm(n) { printf("storm %d\n",n); save_box(stdout,n); } next_storm;
    loop_province(n) {
     printf("province %d\n",n); list("here",p_loc_info(n)->here_list);
    } next_province;
    printf("show_garrison %d\n",show_to_garrison);
    list("output",p_player(200)->output);
   }
   else if(!strcmp(text,"probe")) {
    ilist l=NULL;
    for(i=0;i<11;i++) ilist_append(&l,300+i*7);
    wthr_shuffle(l); list("weather shuffle",l); ilist_reclaim(&l);
    printf("weather strength %d\n",wthr_storm(n,a,2,3));
   }
   else return 2;
  }
  else if (!strcmp(op, "garrison")) {
   sscanf(line,"%*s %31s %d",text,&n);
   if(!strcmp(text,"ping")) ping_garrisons();
   else if(!strcmp(text,"pay")) drv_garrison_pay();
   else if(!strcmp(text,"tax")) garrison_gold();
   else if(!strcmp(text,"metric")) gold_provinces=n;
   else if(!strcmp(text,"subloc")) { sscanf(line,"%*s %*s %d %d",&n,&a); printf("province subloc %d %d %d\n",n,a,province_subloc(n,a)); }
   else if(!strcmp(text,"taxstate")) {
    printf("gold provinces %d\n",gold_provinces);
    loop_garrison(n) {
     struct entity_misc *p=rp_misc(n);
     printf("tax garrison %d misc %d tax %d forward %d strength %d\n",n,p!=NULL,p?p->garr_tax:0,p?p->garr_forward:0,garrison_strength(n));
    } next_garrison;
    loop_loc(n) { printf("tax loc %d temp %d\n",n,bx[n]->temp); } next_loc;
   }
   else if(!strcmp(text,"strength")) printf("garrison strength %d\n",garrison_strength(n));
   else if(!strcmp(text,"maint")) printf("maintenance %d\n",maint_cost(n));
   else if(!strcmp(text,"state")) {
    printf("garrison pay %d gold %d show %d\n",garrison_pay,gold_garrison,show_to_garrison);
    list("vector",out_vector);
    loop_char(n) {
     printf("garrison char %d leader %d default %d prisoner %d npc %d guard %d gold %d\n",
      n,stack_leader(n),default_garrison(n),is_prisoner(n),is_npc(n),char_guard(n),has_item(n,item_gold));
    } next_char;
    loop_player(n) { printf("player %d\n",n); list("output",p_player(n)->output); } next_player;
   }
   else return 2;
  }
  else if (!strcmp(op, "npc")) {
   sscanf(line,"%*s %31s %d",text,&n);
   if(!strcmp(text,"unsworn")) drv_unsworn(n);
   else if(!strcmp(text,"dispatch")) queue_npc_orders();
   else if(!strcmp(text,"bandit")) { sscanf(line,"%*s %*s %d %d",&n,&a); drv_auto_bandit(n,a); }
   else if(!strcmp(text,"mob")) drv_mob(n);
   else if(!strcmp(text,"undead")) auto_undead(n);
   else if(!strcmp(text,"begin")) begin_npc(n);
   else if(!strcmp(text,"probe")) {
    printf("behavior zero %d\n",npc_behavior(n,0,1,2));
    printf("behavior one %d\n",npc_behavior(n,1,1,2));
    printf("behavior wait %d\n",npc_behavior(n,1,10,20));
    printf("behavior wide %d\n",npc_behavior(n,99,-100000,100000));
   }
   else return 2;
  }
  else if (!strcmp(op, "savage")) {
   sscanf(line,"%*s %31s %d",text,&n);
   if(!strcmp(text,"init")) { init_savage_attacks(); printf("initial savages %d\n",drv_savage_count()); }
   else if(!strcmp(text,"auto")) auto_savage(n);
   else if(!strcmp(text,"hates")) printf("hates %d\n",drv_savage_hates(n));
   else if(!strcmp(text,"here")) printf("hate here %d\n",drv_savage_here(n));
   else if(!strcmp(text,"humans")) printf("humans %d\n",controlled_humans_here(n));
   else return 2;
  }
  else if (!strcmp(op, "savage_population")) {
   sscanf(line,"%*s %d",&n);
   for(i=0;i<n;i++) { alloc_box(1100+i,T_char,sub_ni); p_char(1100+i)->unit_item=item_savage; }
  }
  else if (!strcmp(op, "savage_rolls")) {
   sscanf(line,"%*s %d %d",&n,&a); begin_npc(n);
   printf("savage spawn %d\n",npc_spawn(n,0,1,100));
   printf("savage qty %d\n",npc_qty(n,a,3,25));
   printf("savage wait %d\n",npc_behavior(a,0,35,50));
  }
  else if (!strcmp(op, "savage_state")) {
   loop_units(indep_player,n) {
    printf("savage unit %d\n",n); save_box(stdout,n);
    printf("savage output %d\n",test_bit(p_player(player(n))->output,n));
   } next_unit;
   for(i=0;i<5;i++) list(sout("load_q %d",i),load_q[i]);
   list("run_q",run_q);
  }
  else if (!strcmp(op, "hades")) {
   sscanf(line,"%*s %31s %d",text,&n);
   if(!strcmp(text,"auto")) auto_hades();
   else if(!strcmp(text,"unit")) drv_hades_unit(n);
   else if(!strcmp(text,"spawn")) drv_nasty(n);
   else return 2;
  }
  else if (!strcmp(op, "hades_state")) {
   loop_units(hades_player,n) {
    printf("hades unit %d\n",n); save_box(stdout,n);
    printf("hades output %d\n",test_bit(p_player(player(n))->output,n));
   } next_unit;
   for(i=0;i<5;i++) list(sout("load_q %d",i),load_q[i]);
   list("run_q",run_q);
  }
  else if (!strcmp(op, "nasty_rolls")) {
   sscanf(line,"%*s %d %d",&n,&a);
   printf("nasty where %d\n",drv_nasty_roll(n,0,0,a-1));
   printf("nasty kind %d\n",drv_nasty_roll(n,1,1,4));
   printf("nasty spirits %d\n",drv_nasty_roll(n,2,25,75));
   printf("nasty lesser %d\n",drv_nasty_roll(n,2,50,150));
   printf("nasty greater %d\n",drv_nasty_roll(n,2,100,250));
  }
  else if (!strcmp(op, "faery")) {
   n=0; sscanf(line,"%*s %31s %d",text,&n);
   if(!strcmp(text,"auto")) auto_faery();
   else if(!strcmp(text,"spawn")) drv_hunt(n);
   else if(!strcmp(text,"unit")) drv_faery_unit(n);
   else return 2;
  }
  else if (!strcmp(op, "hunt_rolls")) {
   sscanf(line,"%*s %d %d",&n,&a);
   for(i=0;i<6;i++) printf("hunt loc %d %d\n",i,drv_hunt_loc(n,i,a));
   printf("hunt qty %d\n",drv_hunt_qty(n));
  }
  else if (!strcmp(op, "faery_state")) {
   loop_units(faery_player,n) {
    printf("faery unit %d\n",n); save_box(stdout,n);
   } next_unit;
   for(i=0;i<5;i++) list(sout("load_q %d",i),load_q[i]);
   list("run_q",run_q);
  }
  else if (!strcmp(op, "npc_move")) {
   sscanf(line,"%*s %d",&n); npc_move(n);
  }
  else if (!strcmp(op, "npc_choose")) {
   struct exit_view *e;
   sscanf(line,"%*s %d %d %d",&n,&a,&b); e=choose_npc_direction(n,a,b);
   if(e) printf("chosen %d %d\n",e->direction,e->destination);
   else printf("chosen none\n");
  }
  else if (!strcmp(op, "npc_rolls")) {
   sscanf(line,"%*s %d %d",&n,&a); begin_npc(a);
   for(i=1;i<=4;i++) printf("npc direction %d roll %d\n",i,npc_behavior(n,i,1,10));
  }
  else if (!strcmp(op, "arrival")) {
   sscanf(line,"%*s %31s %d %d",text,&n,&a);
   if (!strcmp(text,"move")) move_stack(n,a);
   else if (!strcmp(text,"known")) mark_loc_stack_known(n,a);
   else if (!strcmp(text,"touch")) touch_loc_after_move(n,a);
   else if (!strcmp(text,"weather")) update_weather_view_locs(n,a);
   else if (!strcmp(text,"clear")) clear_contacts(n);
   else if (!strcmp(text,"match")) match_trades(n);
   else if (!strcmp(text,"market")) printf("market %d\n",drv_market(n));
   else return 2;
  }
  else if (!strcmp(op, "escape")) {
   sscanf(line,"%*s %31s %d %d %d",text,&n,&a,&b);
   if(!strcmp(text,"check")) printf("escaped %d\n",check_prisoner_escape(n,a,b));
   else if(!strcmp(text,"run")) prisoner_escapes(n);
   else if(!strcmp(text,"move")) prisoner_movement_escape_check(n);
   else return 2;
  }
  else if (!strcmp(op, "land")) {
   sscanf(line,"%*s %d",&n); printf("land %d\n",find_nearest_land(n));
  }
  else if (!strcmp(op, "routes")) {
   int west, up, down, fields;
   fields=sscanf(line,"%*s %d %d %d %d %d %d %d",&n,&a,&b,&d,&west,&up,&down);
   ilist_clear(&p_loc(n)->prov_dest);
   ilist_append(&p_loc(n)->prov_dest,a); ilist_append(&p_loc(n)->prov_dest,b);
   ilist_append(&p_loc(n)->prov_dest,d); ilist_append(&p_loc(n)->prov_dest,west);
   if(fields>5) ilist_append(&p_loc(n)->prov_dest,up);
   if(fields>6) ilist_append(&p_loc(n)->prov_dest,down);
  }
  else if (!strcmp(op, "road")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b);
   p_gate(n)->to_loc=a; p_gate(n)->road_hidden=(schar)b;
  }
  else if (!strcmp(op, "delete_box")) {
   sscanf(line,"%*s %d",&n); change_box_kind(n,T_deleted);
  }
  else if (!strcmp(op, "shuffle")) {
   exit_views_list l;
   sscanf(line,"%*s %d %d",&n,&a);
   l=exits_from_loc_nsew_select(n,a,TRUE,TRUE);
   for(i=0;i<exit_views_len(l);i++) printf("exitdir %d %d\n",l[i]->direction,l[i]->destination);
  }
  else if (!strcmp(op, "move_state")) {
   list("vector",out_vector);
   loop_char(n) {
    printf("move char %d where %d moving %d guard %d\n",n,loc(n),char_moving(n),char_guard(n));
    list("here",p_loc_info(n)->here_list);
   } next_char;
   printf("show_garrison %d\n",show_to_garrison);
  }
  else if (!strcmp(op, "entity_rolls")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b); drv_entity_rolls(n,a,b);
  }
  else if (!strcmp(op, "escape_state")) {
   list("vector",out_vector);
   loop_char(n) {
    struct entity_char *p=rp_char(n); struct char_magic *m=rp_magic(n);
    printf("escape char %d where %d prisoner %d moving %d magic %d release %d\n",
     n,loc(n),p ? p->prisoner : 0,p ? p->moving : 0,m!=NULL,m ? m->swear_on_release : 0);
    list("here",p_loc_info(n)->here_list);
   } next_char;
  }
  else if (!strcmp(op, "ocean_add")) { sscanf(line,"%*s %d",&n); drv_ocean(n); }
  else if (!strcmp(op, "ocean_init")) init_ocean_chars();
  else if (!strcmp(op, "ocean_check")) check_ocean_chars();
  else if (!strcmp(op, "heal")) drv_heal();
  else if (!strcmp(op, "weekly")) weekly_prisoner_escape_check();
  else if (!strcmp(op, "dogs")) drv_dogs();
  else if (!strcmp(op, "bark")) { sscanf(line,"%*s %d",&n); bark_dogs(n); }
  else if (!strcmp(op, "ring")) drv_ring();
  else if (!strcmp(op, "curse")) drv_curse();
  else if (!strcmp(op, "closure")) {
   n=0; sscanf(line,"%*s %31s %d",text,&n);
   if (!strcmp(text,"auto")) drv_auto_drop();
   else if (!strcmp(text,"links")) drv_link_decay();
   else if (!strcmp(text,"quests")) drv_quest_decay();
   else if (!strcmp(text,"post")) post_month();
   else if (!strcmp(text,"ranks")) determine_noble_ranks();
   else if (!strcmp(text,"kings")) drv_kings();
   else if (!strcmp(text,"threshold")) printf("rank threshold %d %d\n",n,drv_rank_threshold(n));
   else if (!strcmp(text,"lower")) printf("rank lower %d %d\n",n,drv_rank_lower(n));
   else if (!strcmp(text,"sup")) printf("rank recursive %d %d\n",n,drv_rank_sup(n));
   else if (!strcmp(text,"ruler")) printf("top ruler %d %d\n",n,top_ruler(n));
   else if (!strcmp(text,"describe")) {
    printf("rank description %d [%s]\n",n,liner_desc(n)); show_chars_below(200,n);
   }
   else if (!strcmp(text,"state")) {
    printf("post flag %d month %d\n",post_has_been_run,oly_month(sysclock));
    loop_char(n) {
     printf("rank char %d allocated %d temp %d rank %d [%s]\n",n,rp_char(n)!=NULL,bx[n]->temp,char_rank(n),rank_s(n));
    } next_char;
   }
   else return 2;
  }
  /* Dense construction only. Every result still comes from original C.
   * One castle/owner and count provinces with ten-or-nine strength garrisons. */
  else if (!strcmp(op, "rankland")) {
   int reg, j;
   sscanf(line,"%*s %d %d %d %d %d",&n,&a,&b,&d,&reg);
   alloc_box(d,T_loc,sub_castle);
   if(loc(n)) ilist_rem_value(&p_loc_info(loc(n))->here_list,n);
   p_loc_info(n)->where=d; ilist_append(&p_loc_info(d)->here_list,n);
   for(j=1;j<=a;j++) {
    int prov=d+j, garr=d+100+j;
    alloc_box(prov,T_loc,sub_plain); alloc_box(garr,T_char,sub_garrison);
    p_char(garr)->unit_lord=201; p_misc(garr)->garr_castle=d;
    p_loc_info(prov)->where=reg; if(reg) ilist_append(&p_loc_info(reg)->here_list,prov);
    p_loc_info(garr)->where=prov; ilist_append(&p_loc_info(prov)->here_list,garr);
    gen_item(garr,item_soldier,b);
   }
   p_loc_info(d)->where=d+1; ilist_append(&p_loc_info(d+1)->here_list,d);
  }
  else if (!strcmp(op, "decay")) {
   sscanf(line,"%*s %31s %d",text,&n);
   if (!strcmp(text,"mines")) drv_decay_mines();
   else if (!strcmp(text,"ghosts")) drv_decay_ghosts();
   else if (!strcmp(text,"corpses")) drv_decay_corpses();
   else if (!strcmp(text,"bodies")) drv_decay_bodies();
   else if (!strcmp(text,"roll")) printf("corpse roll %d %d\n",n,drv_corpse_roll(n));
   else if (!strcmp(text,"chains")) {
    printf("mines:"); loop_collapsed_mine(n) { printf(" %d",n); } next_collapsed_mine; printf("\n");
    printf("bodies:"); loop_dead_body(n) { printf(" %d",n); } next_dead_body; printf("\n");
   }
   else return 2;
  }
  else if (!strcmp(op, "upkeep")) {
   n=a=b=0;
   sscanf(line,"%*s %31s %d %d %d",text,&n,&a,&b);
   if (!strcmp(text,"cost")) printf("unit maintenance %d\n",unit_maint_cost(n));
   else if (!strcmp(text,"charge")) printf("maintenance result %d\n",charge_maint_sup(n));
   else if (!strcmp(text,"all")) drv_maint();
   else if (!strcmp(text,"starve")) drv_starve(n,a);
   else if (!strcmp(text,"animals")) drv_animal();
   else if (!strcmp(text,"animal")) printf("animal roll %d %d %d %d\n",n,a,b,drv_animal_roll(n,a,b));
   else if (!strcmp(text,"flavor")) printf("starve roll %d %d %d\n",n,a,drv_starve_roll(n,a));
   else return 2;
  }
  else if (!strcmp(op, "income")) {
   sscanf(line,"%*s %31s %d %d %d %d",text,&n,&a,&b,&d);
   if (!strcmp(text,"inn")) drv_income_inn();
   else if (!strcmp(text,"temple")) drv_income_temple();
   else if (!strcmp(text,"roll")) printf("income roll %d %d %d\n",n,a,drv_income_roll(n,a,b,d));
   else if (!strcmp(text,"count")) printf("structures %d %d %d %d\n",n,a,b,count_loc_structures(n,a,b));
   else if (!strcmp(text,"metrics")) { gold_inn=n; gold_temple=a; }
   else if (!strcmp(text,"state")) printf("income gold inn %d temple %d\n",gold_inn,gold_temple);
   else return 2;
  }
  else if (!strcmp(op, "monthly")) {
   sscanf(line,"%*s %31s %d %d",text,&n,&a);
   if (!strcmp(text,"aura")) drv_aura();
   else if (!strcmp(text,"ability")) drv_ability();
   else if (!strcmp(text,"shroud")) drv_shroud();
   else if (!strcmp(text,"hinder")) drv_hinder();
   else if (!strcmp(text,"barrier")) drv_barrier();
   else if (!strcmp(text,"loyalty")) drv_loyalty();
   else if (!strcmp(text,"pillage")) drv_pillage();
   else if (!strcmp(text,"relic")) drv_relic();
   else if (!strcmp(text,"hide")) drv_hide_mage();
   else if (!strcmp(text,"roll")) printf("loyal roll %d %d %d\n",n,a,drv_loyal_roll(n,a));
   else if (!strcmp(text,"loyal")) printf("loyal %d %s\n",n,loyal_s(n));
   else if (!strcmp(text,"state")) {
    loop_char(n) {
     struct entity_char *p=rp_char(n); struct char_magic *m=rp_magic(n);
     printf("monthly char %d char %d magic %d fresh %d new %d\n",n,p!=NULL,m!=NULL,p?p->fresh_hire:0,p?p->new_lord:0);
     save_box(stdout,n);
    } next_char;
    loop_loc(n) {
     struct entity_subloc *p=rp_subloc(n);
     printf("monthly loc %d loc %d subloc %d recent %d\n",n,rp_loc(n)!=NULL,p!=NULL,p?p->recent_loot:0);
     save_box(stdout,n);
    } next_loc;
   }
   else return 2;
  }
  else if (!strcmp(op, "events")) daily_events();
  else if (!strcmp(op, "event_days")) drv_event_days();
  else if (!strcmp(op, "bark_roll")) {
   sscanf(line,"%*s %d %d",&n,&a); printf("bark roll %d\n",drv_bark_roll(n,a));
  }
  else if (!strcmp(op, "prison_roll")) {
   sscanf(line,"%*s %d",&n); printf("weekly roll %d\n",ent_prisoner(n,ESC_WEEKLY));
  }
  else if (!strcmp(op, "heal_roll")) {
   sscanf(line,"%*s %d %d %d %d",&n,&a,&b,&d);
   printf("heal roll %d\n",drv_heal_roll(n,a,b,d));
  }
  else if (!strcmp(op, "arrival_state")) {
   list("ocean",drv_oceans());
   loop_char(n) {
    struct entity_char *p=rp_char(n);
    printf("arrival char %d allocated %d time %d contacts_alloc %d\n",n,p!=NULL,p ? p->time_flying : 0,p && p->contact!=NULL);
    if(p) list("contacts",p->contact);
   } next_char;
   loop_player(n) {
    struct entity_player *p=rp_player(n);
    printf("arrival player %d allocated %d\n",n,p!=NULL);
    if(p) { list("known",p->known); list("touched",p->locs); list("weather",p->weather_seen); list("output",p->output); }
   } next_player;
  }
  else if (!strcmp(op, "accounts")) {
   n=0; sscanf(line,"%*s %s %d",text,&n);
   if (!strcmp(text,"taxes")) drv_taxes();
   else if (!strcmp(text,"claim")) drv_claim();
   else if (!strcmp(text,"np")) drv_np();
   else if (!strcmp(text,"unformed")) drv_unformed();
   else if (!strcmp(text,"sup")) add_unformed_sup(n);
   else if (!strcmp(text,"next")) printf("next np %d %d\n",n,next_np_turn(n));
   else if (!strcmp(text,"clear")) drv_clear_sent();
   else if (!strcmp(text,"open")) drv_special_open();
   else if (!strcmp(text,"close")) drv_special_close();
   else if (!strcmp(text,"state")) {
    printf("tax gold %d\n",gold_taxes);
    loop_player(n) {
     struct entity_player *p=rp_player(n);
     printf("account %d allocated %d\n",n,p!=NULL);
     if(p) {
      printf("np gained %d spent %d sent %d\n",p->np_gained,p->np_spent,p->sent_orders);
      list("unformed",p->unformed); list("output",p->output);
     }
     save_box(stdout,n);
    } next_player;
    loop_kind(T_unform,n) { save_box(stdout,n); } next_kind;
   }
   else return 2;
  }
  else if (!strcmp(op, "production")) {
   sscanf(line,"%*s %s %d %d %d",text,&n,&a,&b);
   if (!strcmp(text,"civ")) compute_civ_levels();
   else if (!strcmp(text,"locations")) location_production();
   else if (!strcmp(text,"mine")) mine_production(n);
   else if (!strcmp(text,"floor")) drv_replenish(n,a,b);
   else if (!strcmp(text,"taxes")) seed_taxes();
   else if (!strcmp(text,"gold")) drv_city_gold();
   else if (!strcmp(text,"post")) post_production();
   else if (!strcmp(text,"state")) {
    loop_loc(n) {
     printf("production loc %d temp %d double %d\n",n,bx[n]->temp,rp_misc(n) ? rp_misc(n)->opium_double : 0);
     save_box(stdout,n);
    } next_loc;
    printf("post %d\n",post_has_been_run);
    list("pending trades",trades_to_check);
   }
   else return 2;
  }
  else if (!strcmp(op, "market")) {
   int x=0;
   n=a=b=d=i=0;
   sscanf(line,"%*s %s %d %d %d %d %d %d",text,&n,&a,&b,&d,&i,&x);
   if (!strcmp(text,"add")) add_city_trade(n,a,b,d,i,x);
   else if (!strcmp(text,"options")) {
    struct trade *t=bx[n]->trades[a]; t->month_prod=b; t->expire=d; t->cloak=i;
   }
   else if (!strcmp(text,"opium")) opium_market_delta(n);
   else if (!strcmp(text,"expire")) expire_trades(n);
   else if (!strcmp(text,"restock")) loc_trade_sup(n,a);
   else if (!strcmp(text,"rings")) trade_suffuse_ring(n);
   else if (!strcmp(text,"all")) location_trades();
   else if (!strcmp(text,"rolls")) {
    begin_economy(n);
    printf("ring rolls %d %d %d\n",econ_stock(sub_suffuse_ring,n,1,3),econ_ring(sub_suffuse_ring,n,1,5),econ_cost(sub_suffuse_ring,n,0,12));
   }
   else if (!strcmp(text,"new_ring")) { begin_economy(n); created=new_suffuse_ring(n); printf("created %d\n",created); }
   else if (!strcmp(text,"state")) {
    for(i=0;i<trades_len(bx[n]->trades);i++) {
     struct trade *t=bx[n]->trades[i];
     printf("raw trade %d %d %d %d %d %d %d %d %d\n",i,t->who,t->kind,t->item,t->qty,t->cost,t->month_prod,t->expire,t->cloak);
    }
    save_box(stdout,n);
    loop_subkind(sub_suffuse_ring,a) { save_box(stdout,a); } next_subkind;
   }
   else if (!strcmp(text,"buy_ring")) {
    loop_subkind(sub_suffuse_ring,b) {
     if(item_unique(b)==n) { add_city_trade(a,BUY,b,1,2000,0); match_trades(a); break; }
    } next_subkind;
   }
   else return 2;
  }
  else if (!strcmp(op, "trade")) {
   struct trade *t=my_malloc(sizeof(*t));
   sscanf(line,"%*s %d %d %d %d %d %d %d",&n,&t->kind,&t->item,&t->qty,&t->cost,&t->who,&t->sort);
   trades_append(&bx[n]->trades,t);
  }
  else if (!strcmp(op, "economy")) {
   sscanf(line,"%*s %s %d %d %d %d",text,&n,&a,&b,&d);
   if (!strcmp(text,"all")) match_all_trades();
   else if (!strcmp(text,"check")) check_validated_trades();
   else if (!strcmp(text,"queue")) ilist_append(&trades_to_check,n);
   else if (!strcmp(text,"options")) { bx[n]->trades[a]->have_left=b; bx[n]->trades[a]->cloak=d; }
   else if (!strcmp(text,"reduce")) printf("reduced %d\n",drv_reduce(bx[n]->trades[a],b));
   else if (!strcmp(text,"attempt")) drv_attempt(bx[n]->trades[a],bx[b]->trades[d]);
   else if (!strcmp(text,"state")) {
    printf("trade gold %d pot basket %d opium %d\n",gold_trade,gold_pot_basket,gold_opium);
    list("pending trades",trades_to_check);
   }
   else { fprintf(stderr,"unknown economy %s\n",text); return 2; }
  }
  else if (!strcmp(op, "offers")) {
   trades_list l;
   sscanf(line,"%*s %d %d %d",&n,&a,&b); l=drv_offers(n,a,b);
   for(i=0;i<trades_len(l);i++) printf("offer %d %d %d %d %d\n",l[i]->who,l[i]->item,l[i]->qty,l[i]->cost,l[i]->sort);
  }
  else if (!strcmp(op, "ambush") || !strcmp(op, "bandit")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b);
   if (!strcmp(op,"bandit")) { created=drv_bandit(n,a,b); printf("created %d\n",created); }
   else if(n) faery_attack_check(a,b); else hades_attack_check(a,b);
  }
  else if (!strcmp(op, "region_rolls")) {
   sscanf(line,"%*s %d %d",&n,&a);
   /* Each statement is sequenced: never rely on C argument evaluation order. */
   printf("hads ambush %d\n",hads_ambush(n,a));
   printf("hads kind %d\n",hads_bandit_kind(n,a));
   printf("hads retal %d\n",hads_retal(n,a));
   printf("faer ambush %d\n",faer_ambush(n,a));
   printf("faer kind %d\n",faer_bandit_kind(n,a));
   printf("faer retal %d\n",faer_retal(n,a));
  }
  else if (!strcmp(op, "bandits")) {
   loop_char(n) {
    if(subkind(n)!=sub_ni) continue;
    printf("bandit %d\n",n); save_box(stdout,n);
    printf("bandit output %d\n",test_bit(p_player(player(n))->output,n));
   } next_char;
   for(i=0;i<5;i++) list(sout("load_q %d",i),load_q[i]);
   list("run_q",run_q);
  }
  else if (!strcmp(op, "select_bandit")) {
   loop_char(n) { if(subkind(n)==sub_ni) { c=rp_command(n); break; } } next_char;
  }
  else if (!strcmp(op, "learn")) {
   struct skill_ent *s=my_malloc(sizeof(*s));
   sscanf(line,"%*s %d %d %d %d",&n,&a,&b,&d);
   s->skill=a; s->know=(schar)b; s->experience=(short)d;
   skill_ents_append(&p_char(n)->skills,s);
  }
  else if (!strcmp(op, "payment")) {
   sscanf(line,"%*s %31s %d %d %d",text,&n,&a,&b);
   if(!strcmp(text,"can")) result=can_pay(n,a);
   else if(!strcmp(text,"charge")) result=charge(n,a);
   else if(!strcmp(text,"auto")) result=autocharge(n,a);
   else if(!strcmp(text,"stack")) result=stack_has_item(n,a);
   else if(!strcmp(text,"sub")) result=stack_sub_item(n,a,b);
   else if(!strcmp(text,"key")) result=has_use_key(n,a);
   else if(!strcmp(text,"stackkey")) result=stack_has_use_key(n,a);
   else return 2;
   printf("payment %d\n",result);
  }
  else if (!strcmp(op, "genitem")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b); gen_item(n,a,b);
  }
  else if (!strcmp(op, "output")) {
   sscanf(line,"%*s %d %d",&n,&a); printf("output %d\n",test_bit(p_player(n)->output,a));
  }
  else if (!strcmp(op, "orders_file")) {
   sscanf(line,"%*s %d",&n); save_player_orders(n);
   printf("orders file %d:\n",n); log=fopen(sout("%s/orders/%d",libdir,n),"r");
   if(log) { while((n=fgetc(log))!=EOF) putchar(n); fclose(log); }
  }
  else if (!strcmp(op, "emit")) {
   sscanf(line,"%*s %d %d",&n,&a); out(n,"routing %d  ",a);
  }
  else if (!strcmp(op, "garrison_view")) {
   sscanf(line,"%*s %d",&n); show_to_garrison = n;
  }
  else if (!strcmp(op, "lore_state")) {
   struct entity_player *p;
   sscanf(line,"%*s %d",&n); p = p_player(n);
   list("known",p->known); list("lore",p->deliver_lore);
   list("touched",p->locs); list("weather",p->weather_seen);
  }
  else if (!strcmp(op, "queue_lore")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b); queue_lore(n,a,b);
  }
  else if (!strcmp(op, "skill_lore")) scan_char_skill_lore();
  else if (!strcmp(op, "show_lore")) show_lore_sheets();
  else if (!strcmp(op, "gm_lore")) {
   sscanf(line,"%*s %d",&n); gm_show_all_skills(n);
  }
  /* Relative libdir keeps the actual perror paths identical in both engines
   * and both runs; no diagnostic filtering. Only fixture setup below. */
  else if (!strcmp(op, "lore_root")) {
   if(chdir(libdir)) return 2;
   libdir=".";
  }
  else if (!strcmp(op, "lore_dir")) {
   sscanf(line,"%*s %8191s",text);
   if(mkdir(text,0755)) return 2;
  }
  else if (!strcmp(op, "lore_file")) {
   char path[8192], *p;
   text[0]=0;
   sscanf(line,"%*s %8191s %8191[^\n]",path,text);
   log=fopen(path,"w"); if(!log) return 2;
   for(p=text;*p;p++) fputc(*p=='|' ? '\n' : *p,log);
   fclose(log);
  }
  else if (!strcmp(op, "post")) {
   char path[8192]; struct stat st;
   sscanf(line,"%*s %8191s",text);
   if(!strcmp(text,"players")) write_player_list();
   else if(!strcmp(text,"email")) write_email();
   else if(!strcmp(text,"totimes")) write_totimes();
   else if(!strcmp(text,"forward")) write_forwards();
   else if(!strcmp(text,"factions")) write_factions();
   else if(!strcmp(text,"all")) {
    write_player_list(); write_email(); write_totimes(); write_forwards(); write_factions();
   }
   else if(!strcmp(text,"file")) {
    sscanf(line,"%*s %*s %8191s",path);
    if(stat(path,&st)) printf("file %s missing\n",path);
    else printf("file %s %s\n",path,S_ISREG(st.st_mode) ? "regular" : "directory");
   }
   else if(!strcmp(text,"empty_player")) {
    sscanf(line,"%*s %*s %d",&n); bx[n]->x_player=NULL;
   }
   else if(!strcmp(text,"new")) {
    sscanf(line,"%*s %*s %d",&n); ilist_append(&new_players,n);
   }
   else if(!strcmp(text,"mail")) {
    sscanf(line,"%*s %*s %d %8191s",&n,text);
    if(!strcmp(text,"null")) p_player(n)->email=NULL;
    else if(!strcmp(text,"empty")) p_player(n)->email=str_save("");
    else if(!strcmp(text,"visible")) p_player(n)->vis_email=str_save("public@example.test");
    else if(!strcmp(text,"cis")) p_player(n)->compuserve=1;
    else return 2;
   }
   else return 2;
  }
  else if (!strcmp(op, "admission")) {
   char path[8192];
   sscanf(line,"%*s %8191s",text);
   if(!strcmp(text,"run")) add_new_players();
   else if(!strcmp(text,"probe")) drv_make_players();
   else if(!strcmp(text,"parse")) {
    sscanf(line,"%*s %*s %8191s %8191s",text,path);
    log=fopen(path,"r"); if(!log) return 2;
    printf("admission result %d\n",drv_admission(text,log));
    fclose(log);
   }
   else if(!strcmp(text,"fetch")) {
    char *s;
    sscanf(line,"%*s %*s %8191s %d",path,&a);
    log=fopen(path,"r"); if(!log) return 2;
    for(i=0;i<a;i++) {
     s=drv_fetch(log);
     if(s) printf("field %d [%s]\n",i,s);
     else printf("field %d NULL\n",i);
    }
    fclose(log);
   }
   else if(!strcmp(text,"hex")) {
    unsigned byte; char *p;
    sscanf(line,"%*s %*s %8191s %8191s",path,text);
    log=fopen(path,"w"); if(!log) return 2;
    for(p=text;*p;p+=2) {
     if(sscanf(p,"%2x",&byte)!=1) return 2;
     fputc(byte,log);
    }
    fclose(log);
   }
   else if(!strcmp(text,"lib")) {
    sscanf(line,"%*s %*s %8191s",path); libdir=str_save(path);
   }
   else if(!strcmp(text,"state")) {
    list("new players",new_players);
    list("vector",out_vector);
    loop_player(n) {
     printf("admission player %d record %d\n",n,rp_player(n)!=NULL);
     save_box(stdout,n);
    } next_player;
   }
   else return 2;
  }
  else if (!strcmp(op, "set_known")) { sscanf(line,"%*s %d %d",&n,&a); set_known(n,a); }
  else if (!strcmp(op, "additem")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b); drv_add_item(n,a,b);
  }
  else if (!strcmp(op, "scan_lore")) scan_char_item_lore();
  else if (!strcmp(op, "init_touched")) init_locs_touched();
  else if (!strcmp(op, "init_weather")) init_weather_views();
  else if (!strcmp(op, "newchar")) {
   int ni, where, health, pl, lk, ll;
   text[0] = 0;
   sscanf(line, "%*s %d %d %d %d %d %d %d %8191[^\n]", &n, &ni, &where, &health, &pl, &lk, &ll, text);
   created = new_char(n, ni, where, health, pl, lk, ll, text);
   printf("created %d\n", created);
  }
  else if (!strcmp(op, "newitem")) {
   sscanf(line, "%*s %d %d", &n, &a);
   created = create_unique_item(n, a); printf("created %d\n", created);
  }
  else if (!strcmp(op, "save_new")) save_box(stdout, created);
  else if (!strcmp(op, "delete_new")) { sscanf(line, "%*s %d", &n); destroy_unique_item(n, created); }
  else if (!strcmp(op, "drop_new")) { sscanf(line,"%*s %d %d",&n,&a); printf("drop new %d\n",drop_item(n,created,a)); }
  else if (!strcmp(op, "created_state"))
   printf("created kind %d sub %d unique %d\n",kind(created),subkind(created),item_unique(created));
  else if (!strcmp(op, "body_new")) {
   text[0] = 0;
   sscanf(line,"%*s %d %d %d %d %8191[^\n]",&n,&a,&b,&d,text);
   p_misc(created)->old_lord=n; p_misc(created)->save_name=*text ? str_save(text) : NULL;
   p_char(created)->prev_lord=a; p_char(created)->loy_kind=b; p_char(created)->loy_rate=d;
  }
  else if (!strcmp(op, "destroy")) { sscanf(line,"%*s %d %d",&n,&a); destroy_unique_item(n,a); }
  else if (!strcmp(op, "damage")) {
   sscanf(line,"%*s %31s %d %d %d %d",text,&n,&a,&b,&d);
   if(!strcmp(text,"coast")) printf("coast %d\n",drv_coast(n));
   else if(!strcmp(text,"ships")) drv_coastal();
   else if(!strcmp(text,"locations")) drv_loc_damage();
   else if(!strcmp(text,"mine")) drv_mine(n);
   else if(!strcmp(text,"inn")) drv_inn(n);
   else if(!strcmp(text,"count")) printf("workers %d depth %d\n",count_loc_char_item(n,item_worker),mine_depth(n));
   else if(!strcmp(text,"roll")) printf("wreck %d\n",wthr_wreck(n,a,b,d));
   else return 2;
  }
  else if (!strcmp(op, "structure")) {
   sscanf(line,"%*s %31s %d %d %d",text,&n,&a,&b);
   if(!strcmp(text,"damage")) printf("destroyed %d\n",add_structure_damage(n,a,b));
   else if(!strcmp(text,"sink")) drv_sink(n);
   else if(!strcmp(text,"collapse")) drv_collapse(n);
   else if(!strcmp(text,"remove")) get_rid_of_collapsed_mine(n);
   else if(!strcmp(text,"ocean")) kill_stack_ocean(n);
   else if(!strcmp(text,"state")) {
    struct entity_subloc *p=rp_subloc(n);
    struct entity_misc *m=rp_misc(n);
    printf("structure %d kind %d sub %d damage %d delay %d castle %d bind %d\n",
     n,kind(n),subkind(n),p ? p->damage : -1,m ? m->mine_delay : 0,
     m ? m->garr_castle : 0,m ? m->bind_storm : 0);
    list("bound",p ? p->bound_storms : NULL);
    list("here",rp_loc_info(n)->here_list);
   }
   else return 2;
  }
  else if (!strcmp(op, "death")) {
   sscanf(line,"%*s %31s %d %d",text,&n,&a);
   if(!strcmp(text,"kill")) kill_char(n,a);
   else if(!strcmp(text,"reclaim")) char_reclaim(n);
   else if(!strcmp(text,"subkind")) change_box_subkind(n,a);
   else if(!strcmp(text,"survive")) printf("survive %d\n",survive_fatal(n));
   else if(!strcmp(text,"forget")) printf("forget %d\n",forget_skill(n,a));
   else if(!strcmp(text,"location")) printf("location [%s]\n",char_rep_location(n));
   else if(!strcmp(text,"hades")) printf("hades location %d\n",random_hades_loc(n));
   else if(!strcmp(text,"state")) {
    int kinds[]={T_char,T_item,T_deadchar};
    for(i=0;i<3;i++) {
     printf("kind chain %d:",kinds[i]);
     loop_kind(kinds[i],n) { printf(" %d",n); } next_kind;
     printf("\n");
    }
   }
   else return 2;
  }
  else if (!strcmp(op, "inherit")) {
   sscanf(line,"%*s %31s %d %d %d %d %d",text,&n,&a,&b,&d,&i);
   if(!strcmp(text,"take")) take_unit_items(n,a,b);
   else if(!strcmp(text,"recipient")) printf("recipient %d\n",stackmate_inheritor(n));
   else if(!strcmp(text,"drop")) drop_stack(n,a);
   else if(!strcmp(text,"extract")) extract_stacked_unit(n);
   else if(!strcmp(text,"prisoner")) printf("move prisoner %d\n",move_prisoner(n,a,b));
   else if(!strcmp(text,"desert")) unit_deserts(n,a,b,d,i);
   else if(!strcmp(text,"cookie")) put_back_cookie(n);
   else if(!strcmp(text,"rolls")) drv_inherit_rolls(n,a,b);
   else if(!strcmp(text,"name")) {
    printf("plural box [%s]\n",plural_item_box(n,a));
    printf("box qty [%s]\n",box_name_qty(n,a));
   }
   else if(!strcmp(text,"state")) {
    printf("combat gold %d independent %d indent %d show_day %d\n",gold_combat,gold_combat_indep,indent,show_day);
    loop_player(n) { printf("inherit player %d\n",n); list("bits",p_player(n)->output); } next_player;
    for(i=0;i<5;i++) list(sout("load_q %d",i),load_q[i]);
    list("run_q",run_q);
   }
   else return 2;
  }
  else if (!strcmp(op, "moveitem")) {
   sscanf(line,"%*s %d %d %d %d",&n,&a,&b,&d);
   printf("moveitem %d\n",move_item(n,a,b,d));
  }
  else if (!strcmp(op, "dropitem") || !strcmp(op, "consume")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b);
   printf("inventory result %d\n",!strcmp(op,"consume") ? consume_item(n,a,b) : drop_item(n,a,b));
  }
  else if (!strcmp(op, "kindname")) {
   sscanf(line,"%*s %d",&n); printf("kind [%s] name [%s]\n",display_kind(n),box_name_kind(n));
  }
  else if (!strcmp(op, "reconcile")) check_token_units();
  else if (!strcmp(op, "process")) process_orders();
  else if (!strcmp(op, "token_owner")) {
   sscanf(line,"%*s %d",&n); printf("token player %d\n",token_player(n));
  }
  else if (!strcmp(op, "token_state")) {
   sscanf(line,"%*s %d",&n); save_box(stdout,n);
   loop_units(n,a) { printf("token row %d\n",a); save_box(stdout,a); } next_unit;
   list("pending trades",trades_to_check);
   for(i=0;i<5;i++) list(sout("load_q %d",i),load_q[i]);
   list("run_q",run_q);
   loop_player(a) { printf("token output %d\n",a); list("bits",p_player(a)->output); } next_player;
  }
  else if (!strcmp(op, "subitem")) {
   sscanf(line,"%*s %d %d %d",&n,&a,&b); printf("subitem %d\n",drv_sub_item(n,a,b));
  }
  else if (!strcmp(op, "report")) {
   sscanf(line,"%*s %31s %d",text,&n);
   if(!strcmp(text,"order")) determine_output_order();
   else if(!strcmp(text,"players")) player_ent_info();
   else if(!strcmp(text,"characters")) character_report();
   else if(!strcmp(text,"locations")) turn_end_loc_reports();
   else if(!strcmp(text,"banner")) player_banner();
   else if(!strcmp(text,"provinces")) drv_summary_provinces();
   else if(!strcmp(text,"templates")) list_order_templates();
   else if(!strcmp(text,"account")) report_account();
   else if(!strcmp(text,"summary")) summary_report();
   else if(!strcmp(text,"factions")) player_report();
   else if(!strcmp(text,"gm")) gm_report(n);
   else if(!strcmp(text,"include")) gen_include_section();
   else if(!strcmp(text,"include_sup")) gen_include_sup(n);
   else if(!strcmp(text,"bit")) {
    sscanf(line,"%*s %*s %d %d",&n,&a);
    set_bit(&p_player(n)->output,a);
   }
   else if(!strcmp(text,"block")) {
    if(mkdir(sout("%s/log/%d",libdir,n),0755)) return 2;
   }
   else if(!strcmp(text,"files")) {
    struct stat st;
    loop_player(n) {
     printf("report file %d %d\n",n,stat(sout("%s/log/%d",libdir,n),&st)==0);
    } next_player;
   }
   else if(!strcmp(text,"stored_orders")) {
    struct entity_player *p=rp_player(n);
    if(p) for(i=0;i<orders_len(p->orders);i++) {
     int j;
     printf("stored orders %d:",p->orders[i]->unit);
     for(j=0;j<cstrings_len(p->orders[i]->l);j++) printf(" [%s]",p->orders[i]->l[j]);
     printf("\n");
    }
   }
   else if(!strcmp(text,"known")) {
    struct entity_player *p=rp_player(n);
    if(p) list("raw known",p->known);
   }
   else if(!strcmp(text,"misc")) {
    printf("report misc %d %d\n",n,rp_misc(n)!=NULL);
    save_box(stdout,n);
   }
   else if(!strcmp(text,"style")) style(n);
   else if(!strcmp(text,"style_probe")) {
    out_path=MASTER; out_alt_who=OUT_BANNER;
    out(n,"style probe");
    out_path=0; out_alt_who=0;
   }
   else if(!strcmp(text,"skill")) {
    struct skill_ent *p=my_malloc(sizeof(*p));
    sscanf(line,"%*s %*s %d %d %d %d",&n,&a,&b,&d);
    p->skill=a; p->know=b; p->days_studied=d;
    skill_ents_append(&p_char(n)->skills,p);
   }
   else if(!strcmp(text,"backlinks")) {
    struct char_magic *p=rp_magic(n);
    printf("backlinks ready %d allocated %d\n",pledge_backlinks,p && p->pledged_to_us!=NULL);
    if(p) list("pledged",p->pledged_to_us);
    save_box(stdout,n);
   }
   else if(!strcmp(text,"backlink")) {
    sscanf(line,"%*s %*s %d %d",&n,&a);
    ilist_append(&p_magic(n)->pledged_to_us,a);
   }
   else if(!strcmp(text,"admit")) {
    struct admit *p=my_malloc(sizeof(*p));
    sscanf(line,"%*s %*s %d %d %d %d",&n,&a,&b,&d);
    p->targ=a; p->sense=b;
    if(d) ilist_append(&p->l,d);
    admits_append(&p_player(n)->admits,p);
   }
   else if(!strcmp(text,"member")) {
    sscanf(line,"%*s %*s %d %d",&n,&a);
    ilist_append(&p_player(n)->admits[admits_len(p_player(n)->admits)-1]->l,a);
   }
   else if(!strcmp(text,"state")) {
    struct entity_player *p=rp_player(n);
    struct att_ent *att=rp_disp(n);
    printf("report entity %d order %d player %d attitudes %d char %d magic %d\n",
     n,bx[n]->output_order,p!=NULL,att!=NULL,rp_char(n)!=NULL,rp_magic(n)!=NULL);
    for(i=0;i<item_ents_len(bx[n]->items);i++)
     printf("raw item %d %d\n",bx[n]->items[i]->item,bx[n]->items[i]->qty);
    if(p) {
     list("raw units",p->units);
     for(i=0;i<admits_len(p->admits);i++) {
      printf("raw admit %d %d\n",p->admits[i]->targ,p->admits[i]->sense);
      list("members",p->admits[i]->l);
     }
    }
    if(att) { list("raw hostile",att->hostile); list("raw neutral",att->neutral); list("raw defend",att->defend); }
    save_box(stdout,n);
   }
   else return 2;
  }
  else if (!strcmp(op, "join")) { sscanf(line,"%*s %d %d",&n,&a); join_stack(n,a); }
  else if (!strcmp(op, "leave")) { sscanf(line,"%*s %d",&n); leave_stack(n); }
  else if (!strcmp(op, "queue"))
   cstrings_append(&drv_orders(player(c->who), c->who)->l, str_save(line + 6));
  else if (!strcmp(op, "orders")) {
   sscanf(line, "%*s %d %d %8191[^\n]", &n, &a, text);
   cstrings_append(&drv_orders(n,a)->l,str_save(text));
  }
  else if (!strcmp(op, "transfer")) {
   sscanf(line, "%*s %d %d %d %d", &n, &a, &b, &d);
   set_lord(n,a,b,d);
  }
  else if (!strcmp(op, "generic_interrupt")) { sscanf(line,"%*s %d",&n); interrupt_order(n); }
  else if (!strcmp(op, "stops")) drv_interrupted();
  else if (!strcmp(op, "flush")) { sscanf(line,"%*s %d %d",&n,&a); flush_unit_orders(n,a); }
  else if (!strcmp(op, "stop_check")) {
   sscanf(line,"%*s %d %d",&n,&a); printf("stop %d\n",stop_order(n,a));
  }
  else if (!strcmp(op, "np_total")) { sscanf(line,"%*s %d",&n); printf("np %d\n",char_np_total(n)); }
  else if (!strcmp(op, "save")) { sscanf(line,"%*s %d",&n); save_box(stdout,n); }
  else if (!strcmp(op, "pstate")) {
   struct entity_player *p;
   sscanf(line,"%*s %d",&n); p = p_player(n);
   printf("player %d np %d gained %d spent %d backlinks %d\n",n,p->noble_points,p->np_gained,p->np_spent,pledge_backlinks);
   list("units",p->units);
   printf("output %d %d %d %d %d touched %d %d\n",test_bit(p->output,10),test_bit(p->output,200),test_bit(p->output,201),test_bit(p->output,1000),test_bit(p->output,1001),test_bit(p->locs,9000),test_bit(p->locs,9001));
   for (i=0;i<orders_len(p->orders);i++) {
    int j; struct order_list *q = p->orders[i];
    printf("orders %d:",q->unit);
    for(j=0;j<cstrings_len(q->l);j++) printf(" [%s]",q->l[j]);
    printf("\n");
   }
  }
  else if (!strcmp(op, "readlog")) {
   sscanf(line,"%*s %d",&n); fflush(NULL);
   printf("log %d:\n",n); log = fopen(sout("%s/log/%d",libdir,n),"r");
   if(log) { while((n=fgetc(log)) != EOF) putchar(n); fclose(log); }
  }
  else if (!strcmp(op, "load")) printf("load %d\n", drv_load(c));
  else if (!strcmp(op, "do")) do_command(c);
  else if (!strcmp(op, "finish")) printf("finish %d\n", finish_command(c));
  else if (!strcmp(op, "interrupt")) printf("interrupt %d\n", cmd_tbl[c->cmd].interrupt(c));
  else if (!strcmp(op, "init_wait")) init_wait_list();
  else if (!strcmp(op, "init_collect")) init_collect_list();
  else if (!strcmp(op, "poll_wait")) check_all_waits();
  else if (!strcmp(op, "clear_second")) clear_second_waits();
  else if (!strcmp(op, "wait_add")) { sscanf(line, "%*s %d", &n); ilist_append(&wait_list, n); }
  else if (!strcmp(op, "lists"))
  {
   list("wait_list", wait_list);
   list("collectors", drv_collectors());
   list("second_wait_list", second_wait_list);
   printf("second_wait %d\n", c->second_wait);
  }
  else if (!strcmp(op, "done")) command_done(c);
  else if (!strcmp(op, "shift")) cmd_shift(c);
  else if (!strcmp(op, "wait") || !strcmp(op, "wait_check"))
  {
   char *ret;
   if (!strcmp(op, "wait")) {
    clear_wait_parse(c);
    strcpy(text, line);
    if (!oly_parse(c, text)) return 2;
   }
   ret = drv_wait_check(c);
   printf("wait result [%s]\n", ret ? ret : "NULL");
  }
  else if (!strcmp(op, "wait_parse"))
  {
   char *ret = parse_wait_args(c);
   printf("wait parse [%s]\n", ret ? ret : "NULL");
   wait_dump(c);
  }
  else if (!strcmp(op, "wait_clear")) { clear_wait_parse(c); wait_dump(c); }
  else if (!strcmp(op, "wait_dump")) wait_dump(c);
  else if (!strcmp(op, "signal"))
  { sscanf(line,"%*s %d %8191[^\n]", &n, text); drv_signal(n, text); }
  else if (!strcmp(op, "name") || !strcmp(op, "plural"))
  {
   sscanf(line,"%*s %d %8191[^\n]", &n, text);
   if (!strcmp(op, "name")) set_name(n, text);
   else p_item(n)->plural_name = str_save(text);
  }
  else if (!strcmp(op, "identity")) {
   sscanf(line,"%*s %d %31s %8191s",&n,label,text);
   p_player(n)->full_name=str_save(label); p_player(n)->email=str_save(text);
  }
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
   else if (!strcmp(text,"last_order")) p_player(n)->last_order_turn = a;
   else if (!strcmp(text,"quest")) p_subloc(n)->quest_late = (schar)a;
   else if (!strcmp(text,"link_to")) ilist_append(&p_subloc(n)->link_to,a);
   else if (!strcmp(text,"link_open")) p_subloc(n)->link_open = (schar)a;
   else if (!strcmp(text,"link_when")) p_subloc(n)->link_when = (schar)a;
   else if (!strcmp(text,"rank")) p_char(n)->rank = (schar)a;
   else if (!strcmp(text,"mine_delay")) p_misc(n)->mine_delay = (schar)a;
   else if (!strcmp(text,"death_turn")) p_char(n)->death_time.turn = (short)a;
   else if (!strcmp(text,"created_death")) p_char(created)->death_time.turn = (short)a;
   else if (!strcmp(text,"melt")) p_char(n)->melt_me = (schar)a;
   else if (!strcmp(text,"sick")) p_char(n)->sick = (schar)a;
   else if (!strcmp(text,"magician")) p_magic(n)->magician = (schar)a;
   else if (!strcmp(text,"max_aura")) p_magic(n)->max_aura = a;
   else if (!strcmp(text,"cur_aura")) p_magic(n)->cur_aura = a;
   else if (!strcmp(text,"ability_shroud")) p_magic(n)->ability_shroud = (short)a;
   else if (!strcmp(text,"hinder")) p_magic(n)->hinder_meditation = (schar)a;
   else if (!strcmp(text,"hide_mage")) p_magic(n)->hide_mage = (schar)a;
   else if (!strcmp(text,"shroud")) p_loc(n)->shroud = (short)a;
   else if (!strcmp(text,"fresh")) p_char(n)->fresh_hire = (schar)a;
   else if (!strcmp(text,"recent")) p_subloc(n)->recent_loot = (schar)a;
   else if (!strcmp(text,"relic_decay")) p_item_magic(n)->relic_decay = (short)a;
   else if (!strcmp(text,"aura_bonus")) p_item_magic(n)->aura_bonus = (short)a;
   else if (!strcmp(text,"auraculum")) p_magic(n)->auraculum = a;
   else if (!strcmp(text,"item_aura")) p_item_magic(n)->aura = (short)a;
   else if (!strcmp(text,"creator")) p_item_magic(n)->creator = a;
   else if (!strcmp(text,"curse")) p_item_magic(n)->curse_loyalty = (schar)a;
   else if (!strcmp(text,"token")) p_magic(n)->token = a;
   else if (!strcmp(text,"required")) p_skill(n)->required_skill = a;
   else if (!strcmp(text,"noble_item")) p_char(n)->unit_item = a;
   else if (!strcmp(text,"prev")) p_char(n)->prev_lord = a;
   else if (!strcmp(text,"rate")) p_char(n)->loy_rate = a;
   else if (!strcmp(text,"pledge")) p_magic(n)->pledge = a;
   else if (!strcmp(text,"castle")) p_misc(n)->garr_castle = a;
   else if (!strcmp(text,"unit")) ilist_append(&p_player(n)->units,a);
   else if (!strcmp(text,"np")) p_player(n)->noble_points = (short)a;
   else if (!strcmp(text,"first_turn")) p_player(n)->first_turn = a;
   else if (!strcmp(text,"sent")) p_player(n)->sent_orders = (schar)a;
   else if (!strcmp(text,"unformed")) ilist_append(&p_player(n)->unformed,a);
   else if (!strcmp(text,"summer")) p_subloc(n)->summer_flag = (schar)a;
   else if (!strcmp(text,"uldim")) p_subloc(n)->uldim_flag = (schar)a;
   else if (!strcmp(text,"skill_np")) p_skill(n)->np_req = a;
   else if (!strcmp(text,"unique")) p_item(n)->who_has = a;
   else if (!strcmp(text,"token_num")) p_item_magic(n)->token_num = (schar)a;
   else if (!strcmp(text,"token_ni")) p_item_magic(n)->token_ni = a;
   else if (!strcmp(text,"lore")) p_item_magic(n)->lore = a;
   else if (!strcmp(text,"use_key")) p_item_magic(n)->use_key = a;
   else if (!strcmp(text,"man")) p_item(n)->is_man_item = (schar)a;
   else if (!strcmp(text,"animal")) p_item(n)->animal = (schar)a;
   else if (!strcmp(text,"attack")) p_item(n)->attack = (short)a;
   else if (!strcmp(text,"defense")) p_item(n)->defense = (short)a;
   else if (!strcmp(text,"default_garrison")) p_magic(n)->default_garr = (schar)a;
   else if (!strcmp(text,"weight")) p_item(n)->weight = a;
   else if (!strcmp(text,"land_cap")) p_item(n)->land_cap = a;
   else if (!strcmp(text,"ride_cap")) p_item(n)->ride_cap = a;
   else if (!strcmp(text,"guard")) p_char(n)->guard = (schar)a;
   else if (!strcmp(text,"queued_lore")) ilist_append(&p_player(n)->deliver_lore,a);
   else if (!strcmp(text,"touched")) set_bit(&p_player(n)->locs,a);
   else if (!strcmp(text,"weather_seen")) set_bit(&p_player(n)->weather_seen,a);
   else if (!strcmp(text,"weather_mage")) p_magic(n)->knows_weather = (schar)a;
   else if (!strcmp(text,"summoner")) p_misc(n)->summoned_by = a;
   else if (!strcmp(text,"npc_home")) p_misc(n)->npc_home = a;
   else if (!strcmp(text,"npc_cookie")) p_misc(n)->npc_cookie = a;
   else if (!strcmp(text,"npc_created")) p_misc(n)->npc_created = a;
   else if (!strcmp(text,"capturable")) p_item(n)->capturable = (schar)a;
   else if (!strcmp(text,"release")) p_magic(n)->swear_on_release = (schar)a;
   else if (!strcmp(text,"health")) p_char(n)->health = (schar)a;
   else if (!strcmp(text,"prisoner")) p_char(n)->prisoner = (schar)a;
   else if (!strcmp(text,"npc")) p_char(n)->npc_prog = (schar)a;
   else if (!strcmp(text,"npc_dir")) p_misc(n)->npc_dir = (schar)a;
   else if (!strcmp(text,"storm_move")) p_misc(n)->storm_move = a;
   else if (!strcmp(text,"bind_storm")) p_misc(n)->bind_storm = a;
   else if (!strcmp(text,"bound_storm")) ilist_append(&p_subloc(n)->bound_storms,a);
   else if (!strcmp(text,"loyal")) p_char(n)->loy_kind = (schar)a;
   else if (!strcmp(text,"moving")) p_char(n)->moving = a;
   else if (!strcmp(text,"time_flying")) p_char(n)->time_flying = (schar)a;
   else if (!strcmp(text,"hidden")) p_magic(n)->hide_self = (schar)a;
   else if (!strcmp(text,"horn")) p_magic(n)->ferry_flag = (schar)a;
   else if (!strcmp(text,"subkind")) bx[n]->skind = a;
   else if (!strcmp(text,"vulnerable")) p_misc(n)->only_vulnerable = a;
   else if (!strcmp(text,"strength")) p_misc(n)->storm_str = (short)a;
   else if (!strcmp(text,"damage")) p_subloc(n)->damage = (unsigned char)a;
   else if (!strcmp(text,"shaft_depth")) p_subloc(n)->shaft_depth = (short)a;
   else if (!strcmp(text,"castle_level")) p_subloc(n)->castle_lev = (schar)a;
   else if (!strcmp(text,"major")) p_subloc(n)->major = (schar)a;
   else if (!strcmp(text,"civ")) p_loc(n)->civ = (schar)a;
   else if (!strcmp(text,"pillage")) p_subloc(n)->loot = (schar)a;
   else if (!strcmp(text,"opium")) p_subloc(n)->opium_econ = a;
   else if (!strcmp(text,"double")) p_misc(n)->opium_double = (schar)a;
   else if (!strcmp(text,"temp")) bx[n]->temp = a;
   else if (!strcmp(text,"output_order")) bx[n]->output_order = a;
   else if (!strcmp(text,"garr_tax")) p_misc(n)->garr_tax = a;
   else if (!strcmp(text,"garr_forward")) p_misc(n)->garr_forward = a;
   else if (!strcmp(text,"safe")) p_subloc(n)->safe = (schar)a;
   else if (!strcmp(text,"loc_hidden")) p_loc(n)->hidden = (schar)a;
   else if (!strcmp(text,"barrier")) p_loc(n)->barrier = (short)a;
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
   else if (!strcmp(text,"fort_covers")) result = drv_fort_covers(n);
   else if (!strcmp(text,"depth")) result = loc_depth(n);
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
   else if (!strcmp(text, "auto_quit")) auto_quit_turns = n;
   else if (!strcmp(text, "post")) post_has_been_run = n;
   else if (!strcmp(text, "indep")) indep_player = n;
   else if (!strcmp(text, "gm")) gm_player = n;
   else if (!strcmp(text, "skill_player")) skill_player = n;
   else if (!strcmp(text, "faery_region")) faery_region = n;
   else if (!strcmp(text, "hades_region")) hades_region = n;
   else if (!strcmp(text, "cloud_region")) cloud_region = n;
   else if (!strcmp(text, "under_region")) under_region = n;
   else if (!strcmp(text, "tunnel_region")) tunnel_region = n;
   else if (!strcmp(text, "nowhere_region")) nowhere_region = n;
   else if (!strcmp(text, "nowhere_loc")) nowhere_loc = n;
   else if (!strcmp(text, "hades_player")) hades_player = n;
   else if (!strcmp(text, "faery_player")) faery_player = n;
   else if (!strcmp(text, "backlinks")) pledge_backlinks = n;
   else if (!strcmp(text, "garrison_pay")) garrison_pay = n;
   else if (!strcmp(text, "gold_garrison")) gold_garrison = n;
   else if (!strcmp(text, "show_day")) show_day = n;
   else if (!strcmp(text, "day")) sysclock.day = (short)n;
   else if (!strcmp(text, "epoch")) sysclock.days_since_epoch = n;
   else if (!strcmp(text, "turn")) sysclock.turn = (short)n;
   else if (!strcmp(text, "priority")) cur_pri = n;
   else if (!strcmp(text, "status")) c->status = (schar)n;
   else if (!strcmp(text, "wait")) c->wait = n;
   else if (!strcmp(text, "poll")) c->poll = (schar)n;
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
