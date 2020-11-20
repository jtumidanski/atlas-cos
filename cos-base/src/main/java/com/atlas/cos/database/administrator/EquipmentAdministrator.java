package com.atlas.cos.database.administrator;

import java.util.Optional;
import javax.persistence.EntityManager;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.transformer.EquipmentDataTransformer;
import com.atlas.cos.entity.Equipment;
import com.atlas.cos.model.EquipmentData;

public class EquipmentAdministrator {
   private EquipmentAdministrator() {
   }

   public static Optional<EquipmentData> create(EntityManager entityManager, int characterId, EquipmentData equipmentData) {
      Equipment equipment = new Equipment();
      equipment.setCharacterId(characterId);
      equipment.setSlot(equipmentData.slot());
      equipment.setItemId(equipmentData.itemId());
      equipment.setStrength(equipmentData.strength());
      equipment.setDexterity(equipmentData.dexterity());
      equipment.setIntelligence(equipmentData.intelligence());
      equipment.setLuck(equipmentData.luck());
      equipment.setHp(equipmentData.hp());
      equipment.setMp(equipmentData.mp());
      equipment.setWeaponAttack(equipmentData.weaponAttack());
      equipment.setWeaponDefense(equipmentData.weaponDefense());
      equipment.setMagicAttack(equipmentData.magicAttack());
      equipment.setMagicDefense(equipmentData.magicDefense());
      equipment.setAccuracy(equipmentData.accuracy());
      equipment.setAvoidability(equipmentData.avoidability());
      equipment.setHands(equipmentData.hands());
      equipment.setSpeed(equipmentData.speed());
      equipment.setJump(equipmentData.jump());
      equipment.setSlots(equipmentData.slots());
      return Optional.of(QueryAdministratorUtil.insertAndReturn(entityManager, equipment, new EquipmentDataTransformer()));
   }

   public static void updateSlot(EntityManager entityManager, int id, short newSlot) {
      QueryAdministratorUtil.update(entityManager, Equipment.class, id, equipment -> equipment.setSlot(newSlot));
   }
}
