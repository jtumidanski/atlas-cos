package com.atlas.cos.processor;

import java.util.Optional;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Mappers;
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
      return getJobFromIndex(createIndex)
            .map(job -> getJob(createIndex, job))
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }

   public Optional<MapleJob> getJobFromIndex(int createIndex) {
      if (createIndex == 0) {
         return Optional.of(MapleJob.NOBLESSE);
      } else if (createIndex == 1) {
         return Optional.of(MapleJob.BEGINNER);
      } else if (createIndex == 2) {
         return Optional.of(MapleJob.LEGEND);
      }
      return Optional.empty();
   }

   protected ResultObjectBuilder getJob(int createIndex, MapleJob job) {
      return new ResultObjectBuilder(JobAttributes.class, job.getId())
            .setAttribute(new JobAttributesBuilder()
                  .setName(job.name())
                  .setCreateIndex(createIndex)
            );
   }
}
