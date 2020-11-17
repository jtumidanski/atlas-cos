package com.atlas.cos.builder;

import com.atlas.cos.attribute.BlockedNameAttributes;

import builder.AttributeResultBuilder;
import builder.Builder;

public class BlockedNameAttributesBuilder extends Builder<BlockedNameAttributes, BlockedNameAttributesBuilder>
      implements AttributeResultBuilder {
   @Override
   public BlockedNameAttributes construct() {
      return new BlockedNameAttributes();
   }

   @Override
   public BlockedNameAttributesBuilder getThis() {
      return this;
   }

   public BlockedNameAttributesBuilder setName(String name) {
      return add(attr -> attr.setName(name));
   }
}
