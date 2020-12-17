package com.atlas.cos.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.model.EquipmentData;

public class EquipmentDataBuilder extends RecordBuilder<EquipmentData, EquipmentDataBuilder> {
   private int id;

   private int characterId;

   private int equipmentId;

   private short slot;

   @Override
   public EquipmentData construct() {
      return new EquipmentData(id, characterId, equipmentId, slot);
   }

   @Override
   public EquipmentDataBuilder getThis() {
      return this;
   }

   public EquipmentDataBuilder setId(int id) {
      this.id = id;
      return getThis();
   }

   public EquipmentDataBuilder setCharacterId(int characterId) {
      this.characterId = characterId;
      return getThis();
   }

   public EquipmentDataBuilder setEquipmentId(int equipmentId) {
      this.equipmentId = equipmentId;
      return getThis();
   }

   public EquipmentDataBuilder setSlot(short slot) {
      this.slot = slot;
      return getThis();
   }
}
