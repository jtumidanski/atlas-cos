package com.atlas.cos.util;

import com.atlas.cos.model.Job;
import com.atlas.cos.model.skills.Evan;

public final class SkillUtil {
   private SkillUtil() {
   }

   public static int getSkillBook(int jobId) {
      if (jobId >= 2210 && jobId <= 2218) {
         return jobId - 2209;
      }
      return 0;
   }

   public static boolean isFourthJob(int jobId, int skillId) {
      if (jobId == Job.EVAN4.getId()) {
         return false;
      }
      if (skillId == Evan.MAGIC_MASTERY || skillId == Evan.FLAME_WHEEL || skillId == Evan.HEROS_WILL || skillId == Evan.DARK_FOG
            || skillId == Evan.SOUL_STONE) {
         return true;
      }
      return jobId % 10 == 2;
   }
}
