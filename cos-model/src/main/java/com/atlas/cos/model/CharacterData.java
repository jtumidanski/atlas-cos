package com.atlas.cos.model;

import java.util.Arrays;

public record CharacterData(int id, int accountId, int worldId, String name, int level, int experience, int gachaponExperience,
                            int strength, int dexterity, int luck, int intelligence, int hp, int mp, int maxHp, int maxMp, int meso,
                            int hpMpUsed, int jobId, int skinColor, byte gender, int fame, int hair, int face, int ap, String sp,
                            int mapId, int spawnPoint, int gm) {
   public int maxClassLevel() {
      return isCygnus() ? 120 : 200;
   }

   public boolean isCygnus() {
      return getJobType() == 1;
   }

   public int getJobType() {
      return jobId / 1000;
   }

   public boolean isBeginnerJob() {
      return (jobId == 0 || jobId == 1000 || jobId == 2000);
   }

   public int[] sps() {
      return Arrays.stream(sp.split(","))
            .mapToInt(Integer::parseInt)
            .toArray();
   }
}
