package com.atlas.cos.database.administrator;

import java.util.Optional;
import javax.persistence.EntityManager;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.provider.EquipmentProvider;
import com.atlas.cos.database.transformer.EquipmentDataTransformer;
import com.atlas.cos.entity.Equipment;
import com.atlas.cos.model.EquipmentData;

public class EquipmentAdministrator {
   private EquipmentAdministrator() {
   }

   public static void updateSlot(EntityManager entityManager, int equipmentId, short newSlot) {
      EquipmentProvider.getByEquipmentId(entityManager, equipmentId)
            .map(EquipmentData::id)
            .ifPresent(id -> QueryAdministratorUtil
                  .update(entityManager, Equipment.class, id, equipment -> equipment.setSlot(newSlot)));
   }

   public static Optional<EquipmentData> create(EntityManager entityManager, int characterId, Integer equipmentId, short slot) {
      Equipment equipment = new Equipment();
      equipment.setCharacterId(characterId);
      equipment.setEquipmentId(equipmentId);
      equipment.setSlot(slot);
      return Optional.of(QueryAdministratorUtil.insertAndReturn(entityManager, equipment, new EquipmentDataTransformer()));
   }
}
