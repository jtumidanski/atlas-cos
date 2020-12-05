package com.atlas.cos.processor;

import com.atlas.cos.model.MapleJob;

import java.util.Optional;

public final class JobProcessor {
   private JobProcessor() {
   }

   public static Optional<MapleJob> getJobFromIndex(int createIndex) {
      if (createIndex == 0) {
         return Optional.of(MapleJob.NOBLESSE);
      } else if (createIndex == 1) {
         return Optional.of(MapleJob.BEGINNER);
      } else if (createIndex == 2) {
         return Optional.of(MapleJob.LEGEND);
      }
      return Optional.empty();
   }
}
