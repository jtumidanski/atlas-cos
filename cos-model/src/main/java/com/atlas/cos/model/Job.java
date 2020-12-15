package com.atlas.cos.model;

import java.util.Arrays;
import java.util.Optional;

public enum Job {
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

   private static final int[] jobUpgradeBlob = {1, 20, 60, 110, 190};

   private static final int[] jobUpgradeSpUp = {0, 1, 2, 3, 6};

   private final int jobId;

   Job(int id) {
      jobId = id;
   }

   public static Optional<Job> getById(int id) {
      return Arrays.stream(values())
            .filter(possible -> possible.getId() == id)
            .findFirst();
   }

   public boolean isA(Job baseJob) {
      int baseBranch = baseJob.getId() / 10;
      return this.getId() / 10 == baseBranch && this.getId() >= baseJob.getId()
            || baseBranch % 10 == 0 && this.getId() / 100 == baseJob.getId() / 100;
   }

   public int getId() {
      return jobId;
   }

   public Job getJobStyle(int strength, int dexterity) {
      return Job.getJobStyle(jobId, strength, dexterity);
   }

   public static Job getJobStyle(int jobId, int strength, int dexterity) {
      return getJobStyleInternal(jobId, (byte) ((strength > dexterity) ? 0x80 : 0x40));
   }

   public static Job getJobStyleInternal(int jobId, byte opt) {
      int jobType = jobId / 100;

      if (jobType == Job.WARRIOR.getId() / 100 || jobType == Job.DAWN_WARRIOR_1.getId() / 100
            || jobType == Job.ARAN1.getId() / 100) {
         return Job.WARRIOR;
      } else if (jobType == Job.MAGICIAN.getId() / 100 || jobType == Job.BLAZE_WIZARD_1.getId() / 100
            || jobType == Job.EVAN1.getId() / 100) {
         return Job.MAGICIAN;
      } else if (jobType == Job.BOWMAN.getId() / 100 || jobType == Job.WIND_ARCHER_1.getId() / 100) {
         if (jobId / 10 == Job.CROSSBOWMAN.getId() / 10) {
            return Job.CROSSBOWMAN;
         } else {
            return Job.BOWMAN;
         }
      } else if (jobType == Job.THIEF.getId() / 100 || jobType == Job.NIGHT_WALKER_1.getId() / 100) {
         return (Job.THIEF);
      } else if (jobType == Job.PIRATE.getId() / 100 || jobType == Job.THUNDER_BREAKER_1.getId() / 100) {
         if (opt == (byte) 0x80) {
            return Job.BRAWLER;
         } else {
            return Job.GUNSLINGER;
         }
      }

      return Job.BEGINNER;
   }

   public int getJobBranch() {
      if (jobId % 1000 == 0) {
         return 0;
      } else if (jobId % 100 == 0) {
         return 1;
      } else {
         return 2 + (jobId % 10);
      }
   }

   public boolean hasSpTable() {
      return switch (this) {
         case EVAN, EVAN1, EVAN2, EVAN3, EVAN4, EVAN5, EVAN6, EVAN7, EVAN8, EVAN9, EVAN10 -> true;
         default -> false;
      };
   }
   public int getJobUpgradeLevelRange() {
      return jobUpgradeBlob[getJobBranch()];
   }

   public int getChangeJobSpUpgrade() {
      return jobUpgradeSpUp[getJobBranch()];
   }
}
