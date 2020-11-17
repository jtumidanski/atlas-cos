package com.atlas.cos.builder;

import com.atlas.cos.attribute.JobAttributes;

import builder.AttributeResultBuilder;
import builder.Builder;

public class JobAttributesBuilder extends Builder<JobAttributes, JobAttributesBuilder> implements AttributeResultBuilder {
   @Override
   public JobAttributes construct() {
      return new JobAttributes();
   }

   @Override
   public JobAttributesBuilder getThis() {
      return this;
   }

   public JobAttributesBuilder setName(String name) {
      return add(attr -> attr.setName(name));
   }

   public JobAttributesBuilder setCreateIndex(Integer createIndex) {
      return add(attr -> attr.setCreateIndex(createIndex));
   }
}
