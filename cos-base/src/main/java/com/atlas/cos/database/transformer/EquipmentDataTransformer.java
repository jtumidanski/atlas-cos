package com.atlas.cos.database.transformer;

import com.atlas.cos.builder.EquipmentDataBuilder;
import com.atlas.cos.entity.Equipment;
import com.atlas.cos.model.EquipmentData;

import transformer.SqlTransformer;

public class EquipmentDataTransformer implements SqlTransformer<EquipmentData, Equipment> {
   @Override
   public EquipmentData transform(Equipment equipment) {
      return new EquipmentDataBuilder()
            .setSlot(equipment.getSlot())
            .setItemId(equipment.getItemId())
            .setStrength(equipment.getStrength())
            .setDexterity(equipment.getDexterity())
            .setIntelligence(equipment.getIntelligence())
            .setLuck(equipment.getLuck())
            .setHp(equipment.getHp())
            .setMp(equipment.getMp())
            .setWeaponAttack(equipment.getWeaponAttack())
            .setWeaponDefense(equipment.getWeaponDefense())
            .setMagicAttack(equipment.getMagicAttack())
            .setMagicDefense(equipment.getMagicDefense())
            .setAccuracy(equipment.getAccuracy())
            .setAvoidability(equipment.getAvoidability())
            .setHands(equipment.getHands())
            .setSpeed(equipment.getSpeed())
            .setJump(equipment.getJump())
            .setSlots(equipment.getSlots())
            .build();
   }
}
