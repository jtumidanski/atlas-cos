package com.atlas.cos.processor;

import com.atlas.cos.attribute.JobAttributes;
import com.atlas.cos.builder.JobAttributesBuilder;
import com.atlas.cos.model.MapleJob;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;

public class JobProcessor {
   private static final Object lock = new Object();

   private static volatile JobProcessor instance;

   public static JobProcessor getInstance() {
      JobProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new JobProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder getByCreateIndex(int createIndex) {
      int jobId = -1;
      if (createIndex == 0) {
         jobId = 1000;
      } else if (createIndex == 1) {
         jobId = 0;
      } else if (createIndex == 2) {
         jobId = 2000;
      }

      return new ResultBuilder().addData(getJob(jobId));
   }

   protected ResultObjectBuilder getJob(int jobId) {
      MapleJob mapleJob = MapleJob.getById(jobId);
      if (mapleJob != null) {
         int createIndex = -1;
         if (mapleJob.equals(MapleJob.BEGINNER)) {
            createIndex = 1;
         } else if (mapleJob.equals(MapleJob.NOBLESSE)) {
            createIndex = 0;
         } else if (mapleJob.equals(MapleJob.LEGEND)) {
            createIndex = 2;
         }

         return new ResultObjectBuilder(JobAttributes.class, jobId)
               .setAttribute(new JobAttributesBuilder()
                     .setName(mapleJob.name())
                     .setCreateIndex(createIndex)
               );
      }
      return null;
   }
}
