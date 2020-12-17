package com.atlas.cos.database.provider;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.IntStream;
import javax.persistence.EntityManager;

import com.app.database.provider.NamedQueryClient;
import com.atlas.cos.database.transformer.EquipmentDataTransformer;
import com.atlas.cos.entity.Equipment;
import com.atlas.cos.model.EquipmentData;

public class EquipmentProvider {
   private EquipmentProvider() {
   }

   public static List<EquipmentData> getForCharacter(EntityManager entityManager, int characterId) {
      return new NamedQueryClient<>(entityManager, Equipment.GET_FOR_CHARACTER, Equipment.class)
            .setParameter(Equipment.CHARACTER_ID, characterId)
            .list(new EquipmentDataTransformer());
   }

   public static Optional<EquipmentData> getByEquipmentId(EntityManager entityManager, int equipmentId) {
      return new NamedQueryClient<>(entityManager, Equipment.GET_FOR_EQUIPMENT, Equipment.class)
            .setParameter(Equipment.EQUIPMENT_ID, equipmentId)
            .element(new EquipmentDataTransformer());
   }

   public static Optional<Integer> findEquipmentInSlot(EntityManager entityManager, int characterId, short destinationSlot) {
      return new NamedQueryClient<>(entityManager, Equipment.GET_FOR_CHARACTER_BY_SLOT, Equipment.class)
            .setParameter(Equipment.CHARACTER_ID, characterId)
            .setParameter(Equipment.SLOT, destinationSlot)
            .element(new EquipmentDataTransformer())
            .map(EquipmentData::id);
   }

   public static Optional<Short> getNextFreeEquipmentSlot(EntityManager entityManager, int characterId) {

      Set<Short> usedSlots = EquipmentProvider.getForCharacter(entityManager, characterId)
            .stream()
            .map(EquipmentData::slot)
            .filter(slot -> slot >= 0)
            .collect(Collectors.toSet());
      if (usedSlots.isEmpty()) {
         return Optional.of((short) 0);
      }

      Set<Short> slots = IntStream.rangeClosed(Collections.min(usedSlots), Collections.max(usedSlots))
            .boxed()
            .map(Integer::shortValue)
            .collect(Collectors.toSet());

      slots.removeAll(usedSlots);

      return Optional.of(slots.stream()
            .findFirst()
            .orElse((short) usedSlots.size()));
   }
}
