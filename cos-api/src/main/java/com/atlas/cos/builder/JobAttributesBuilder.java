package com.atlas.cos.builder;

import builder.AttributeResultBuilder;
import builder.RecordBuilder;
import com.atlas.cos.attribute.JobAttributes;

public class JobAttributesBuilder extends RecordBuilder<JobAttributes, JobAttributesBuilder> implements AttributeResultBuilder {
   private static final String NAME = "NAME";

   private static final String CREATE_INDEX = "CREATE_INDEX";

   @Override
   public JobAttributes construct() {
      return new JobAttributes(get(NAME), get(CREATE_INDEX));
   }

   @Override
   public JobAttributesBuilder getThis() {
      return this;
   }

   public JobAttributesBuilder setName(String name) {
      return set(NAME, name);
   }

   public JobAttributesBuilder setCreateIndex(int createIndex) {
      return set(CREATE_INDEX, createIndex);
   }
}
