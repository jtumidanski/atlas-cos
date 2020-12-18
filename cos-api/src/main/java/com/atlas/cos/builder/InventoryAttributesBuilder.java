package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.InventoryAttributes;

import builder.AttributeResultBuilder;

public class InventoryAttributesBuilder extends RecordBuilder<InventoryAttributes, InventoryAttributesBuilder>
      implements AttributeResultBuilder {
   private String type;

   private Integer capacity;

   @Override
   public InventoryAttributes construct() {
      return new InventoryAttributes(type, capacity);
   }

   @Override
   public InventoryAttributesBuilder getThis() {
      return this;
   }

   public InventoryAttributesBuilder setType(String type) {
      this.type = type;
      return getThis();
   }

   public InventoryAttributesBuilder setCapacity(Integer capacity) {
      this.capacity = capacity;
      return getThis();
   }
}
