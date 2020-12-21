package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.SkillAttributes;

import builder.AttributeResultBuilder;

public class SkillAttributesBuilder extends RecordBuilder<SkillAttributes, SkillAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer level;

   private Integer masterLevel;

   private Long expiration;

   @Override
   public SkillAttributes construct() {
      return new SkillAttributes(level, masterLevel, expiration);
   }

   @Override
   public SkillAttributesBuilder getThis() {
      return this;
   }

   public SkillAttributesBuilder setLevel(Integer level) {
      this.level = level;
      return getThis();
   }

   public SkillAttributesBuilder setMasterLevel(Integer masterLevel) {
      this.masterLevel = masterLevel;
      return getThis();
   }

   public SkillAttributesBuilder setExpiration(Long expiration) {
      this.expiration = expiration;
      return getThis();
   }
}
