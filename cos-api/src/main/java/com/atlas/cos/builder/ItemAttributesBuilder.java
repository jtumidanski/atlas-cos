package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.attribute.ItemAttributes;

import builder.AttributeResultBuilder;

public class ItemAttributesBuilder extends RecordBuilder<ItemAttributes, ItemAttributesBuilder> implements AttributeResultBuilder {
   private Integer itemId;

   private Integer quantity;

   private Short slot;

   @Override
   public ItemAttributes construct() {
      return new ItemAttributes(itemId, quantity, slot);
   }

   @Override
   public ItemAttributesBuilder getThis() {
      return this;
   }

   public ItemAttributesBuilder setItemId(Integer itemId) {
      this.itemId = itemId;
      return getThis();
   }

   public ItemAttributesBuilder setQuantity(Integer quantity) {
      this.quantity = quantity;
      return getThis();
   }

   public ItemAttributesBuilder setSlot(Short slot) {
      this.slot = slot;
      return getThis();
   }
}
