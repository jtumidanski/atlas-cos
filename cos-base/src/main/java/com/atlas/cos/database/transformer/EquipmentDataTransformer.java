package com.atlas.cos.database.transformer;

import com.atlas.cos.builder.EquipmentDataBuilder;
import com.atlas.cos.entity.Equipment;
import com.atlas.cos.model.EquipmentData;

import transformer.SqlTransformer;

public class EquipmentDataTransformer implements SqlTransformer<EquipmentData, Equipment> {
   @Override
   public EquipmentData transform(Equipment equipment) {
      return new EquipmentDataBuilder()
            .setId(equipment.getId())
            .setCharacterId(equipment.getCharacterId())
            .setEquipmentId(equipment.getEquipmentId())
            .setSlot(equipment.getSlot())
            .build();
   }
}
