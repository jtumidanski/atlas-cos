package com.atlas.cos.model;

public enum MapleJob {
   BEGINNER(0),

   WARRIOR(100),
   FIGHTER(110), CRUSADER(111), HERO(112),
   PAGE(120), WHITE_KNIGHT(121), PALADIN(122),
   SPEARMAN(130), DRAGON_KNIGHT(131), DARK_KNIGHT(132),

   MAGICIAN(200),
   FP_WIZARD(210), FIRE_POISON_MAGICIAN(211), FIRE_POISON_ARCH_MAGICIAN(212),
   IL_WIZARD(220), ICE_LIGHTENING_MAGICIAN(221), ICE_LIGHTENING_ARCH_MAGICIAN(222),
   CLERIC(230), PRIEST(231), BISHOP(232),

   BOWMAN(300),
   HUNTER(310), RANGER(311), BOW_MASTER(312),
   CROSSBOWMAN(320), SNIPER(321), MARKSMAN(322),

   THIEF(400),
   ASSASSIN(410), HERMIT(411), NIGHT_LORD(412),
   BANDIT(420), CHIEF_BANDIT(421), SHADOWER(422),

   PIRATE(500),
   BRAWLER(510), MARAUDER(511), BUCCANEER(512),
   GUNSLINGER(520), OUTLAW(521), CORSAIR(522),

   MAPLE_LEAF_BRIGADIER(800),
   GM(900), SUPER_GM(910),

   NOBLESSE(1000),
   DAWN_WARRIOR_1(1100), DAWN_WARRIOR_2(1110), DAWN_WARRIOR_3(1111), DAWN_WARRIOR_4(1112),
   BLAZE_WIZARD_1(1200), BLAZE_WIZARD_2(1210), BLAZE_WIZARD_3(1211), BLAZE_WIZARD_4(1212),
   WIND_ARCHER_1(1300), WIND_ARCHER_2(1310), WIND_ARCHER_3(1311), WIND_ARCHER_4(1312),
   NIGHT_WALKER_1(1400), NIGHT_WALKER_2(1410), NIGHT_WALKER_3(1411), NIGHT_WALKER_4(1412),
   THUNDER_BREAKER_1(1500), THUNDER_BREAKER_2(1510), THUNDER_BREAKER_3(1511), THUNDER_BREAKER_4(1512),

   LEGEND(2000), EVAN(2001),
   ARAN1(2100), ARAN2(2110), ARAN3(2111), ARAN4(2112),

   EVAN1(2200), EVAN2(2210), EVAN3(2211), EVAN4(2212), EVAN5(2213), EVAN6(2214),
   EVAN7(2215), EVAN8(2216), EVAN9(2217), EVAN10(2218);

   final static int maxId = 22;    // maxId = (EVAN / 100);
   final int jobId;

   MapleJob(int id) {
      jobId = id;
   }

   public static int getMax() {
      return maxId;
   }

   public static MapleJob getById(int id) {
      for (MapleJob l : MapleJob.values()) {
         if (l.getId() == id) {
            return l;
         }
      }
      return null;
   }

   public static MapleJob getBy5ByteEncoding(int encoded) {
      return switch (encoded) {
         case 2 -> WARRIOR;
         case 4 -> MAGICIAN;
         case 8 -> BOWMAN;
         case 16 -> THIEF;
         case 32 -> PIRATE;
         case 1024 -> NOBLESSE;
         case 2048 -> DAWN_WARRIOR_1;
         case 4096 -> BLAZE_WIZARD_1;
         case 8192 -> WIND_ARCHER_1;
         case 16384 -> NIGHT_WALKER_1;
         case 32768 -> THUNDER_BREAKER_1;
         default -> BEGINNER;
      };
   }

   public int getId() {
      return jobId;
   }

   public boolean isA(MapleJob baseJob) {
      int baseBranch = baseJob.getId() / 10;
      return (getId() / 10 == baseBranch && getId() >= baseJob.getId()) || (baseBranch % 10 == 0 && getId() / 100 == baseJob.getId() / 100);
   }

   public int getJobNiche() {
      return (jobId / 100) % 10;
        
        /*
        case 0: BEGINNER;
        case 1: WARRIOR;
        case 2: MAGICIAN;
        case 3: BOWMAN;  
        case 4: THIEF;
        case 5: PIRATE;
        */
   }
}
