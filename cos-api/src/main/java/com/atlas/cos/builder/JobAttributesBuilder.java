package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.JobAttributes;

import builder.AttributeResultBuilder;

public class JobAttributesBuilder extends RecordBuilder<JobAttributes, JobAttributesBuilder> implements AttributeResultBuilder {
   private String name;

   private Integer createIndex;

   @Override
   public JobAttributes construct() {
      return new JobAttributes(name, createIndex);
   }

   @Override
   public JobAttributesBuilder getThis() {
      return this;
   }

   public JobAttributesBuilder setName(String name) {
      this.name = name;
      return getThis();
   }

   public JobAttributesBuilder setCreateIndex(Integer createIndex) {
      this.createIndex = createIndex;
      return getThis();
   }
}
