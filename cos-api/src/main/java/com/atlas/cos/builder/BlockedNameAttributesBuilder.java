package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.BlockedNameAttributes;

import builder.AttributeResultBuilder;

public class BlockedNameAttributesBuilder extends RecordBuilder<BlockedNameAttributes, BlockedNameAttributesBuilder>
      implements AttributeResultBuilder {
   private String name;

   @Override
   public BlockedNameAttributes construct() {
      return new BlockedNameAttributes(name);
   }

   @Override
   public BlockedNameAttributesBuilder getThis() {
      return this;
   }

   public BlockedNameAttributesBuilder setName(String name) {
      this.name = name;
      return getThis();
   }
}
