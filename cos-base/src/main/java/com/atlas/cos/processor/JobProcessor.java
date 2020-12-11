package com.atlas.cos.processor;

import com.atlas.cos.model.Job;

import java.util.Optional;

public final class JobProcessor {
   private JobProcessor() {
   }

   public static Optional<Job> getJobFromIndex(int createIndex) {
      if (createIndex == 0) {
         return Optional.of(Job.NOBLESSE);
      } else if (createIndex == 1) {
         return Optional.of(Job.BEGINNER);
      } else if (createIndex == 2) {
         return Optional.of(Job.LEGEND);
      }
      return Optional.empty();
   }
}
