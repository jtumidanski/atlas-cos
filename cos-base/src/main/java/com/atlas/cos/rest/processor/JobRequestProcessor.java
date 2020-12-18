package com.atlas.cos.rest.processor;

import com.app.rest.util.stream.Mappers;
import com.atlas.cos.attribute.JobAttributes;
import com.atlas.cos.builder.JobAttributesBuilder;
import com.atlas.cos.model.Job;
import com.atlas.cos.processor.JobProcessor;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;

public final class JobRequestProcessor {
   private JobRequestProcessor() {
   }

   public static ResultBuilder getByCreateIndex(int createIndex) {
      return JobProcessor.getJobFromIndex(createIndex)
            .map(job -> getJob(createIndex, job))
            .map(Mappers::singleOkResult)
            .orElseGet(ResultBuilder::notFound);
   }

   protected static ResultObjectBuilder getJob(int createIndex, Job job) {
      return new ResultObjectBuilder(JobAttributes.class, job.getId())
            .setAttribute(new JobAttributesBuilder()
                  .setName(job.name())
                  .setCreateIndex(createIndex)
            );
   }
}
