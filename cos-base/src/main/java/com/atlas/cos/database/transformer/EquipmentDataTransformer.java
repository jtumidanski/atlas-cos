package com.atlas.cos.database.transformer;

import com.atlas.cos.entity.Equipment;
import com.atlas.cos.model.EquipmentData;

import transformer.SqlTransformer;

public class EquipmentDataTransformer implements SqlTransformer<EquipmentData, Equipment> {
   @Override
   public EquipmentData transform(Equipment equipment) {
      return new EquipmentData(equipment.getId(), equipment.getCharacterId(), equipment.getEquipmentId(), equipment.getSlot());
   }
}
