package com.atlas.cos.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.rest.attribute.EquipmentAttributes;

import builder.AttributeResultBuilder;

public class EquipmentAttributesBuilder extends RecordBuilder<EquipmentAttributes, EquipmentAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer equipmentId;

   private Short slot;

   @Override
   public EquipmentAttributes construct() {
      return new EquipmentAttributes(equipmentId, slot);
   }

   @Override
   public EquipmentAttributesBuilder getThis() {
      return this;
   }

   public EquipmentAttributesBuilder setEquipmentId(Integer equipmentId) {
      this.equipmentId = equipmentId;
      return getThis();
   }

   public EquipmentAttributesBuilder setSlot(Short slot) {
      this.slot = slot;
      return getThis();
   }
}
