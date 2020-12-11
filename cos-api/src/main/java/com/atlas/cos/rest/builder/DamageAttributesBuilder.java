package com.atlas.cos.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.rest.attribute.DamageAttributes;
import com.atlas.cos.rest.attribute.DamageType;

import builder.AttributeResultBuilder;

public class DamageAttributesBuilder extends RecordBuilder<DamageAttributes, DamageAttributesBuilder>
      implements AttributeResultBuilder {
   private DamageType type;

   private Integer maximum;

   @Override
   public DamageAttributes construct() {
      return new DamageAttributes(type, maximum);
   }

   @Override
   public DamageAttributesBuilder getThis() {
      return this;
   }

   public DamageAttributesBuilder setType(DamageType type) {
      this.type = type;
      return getThis();
   }

   public DamageAttributesBuilder setMaximum(Integer maximum) {
      this.maximum = maximum;
      return getThis();
   }
}
