package com.atlas.cos.builder;

import builder.AttributeResultBuilder;
import builder.RecordBuilder;
import com.atlas.cos.attribute.BlockedNameAttributes;

public class BlockedNameAttributesBuilder extends RecordBuilder<BlockedNameAttributes, BlockedNameAttributesBuilder> implements AttributeResultBuilder {
   private static final String NAME = "NAME";

   @Override
   public BlockedNameAttributes construct() {
      return new BlockedNameAttributes(get(NAME));
   }

   @Override
   public BlockedNameAttributesBuilder getThis() {
      return this;
   }

   public BlockedNameAttributesBuilder setName(String name) {
      return set(NAME, name);
   }
}
