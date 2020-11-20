package com.atlas.cos.processor;

import com.atlas.cos.database.provider.EquipmentProvider;
import com.atlas.cos.model.EquipmentData;

import builder.ResultBuilder;
import database.Connection;

public class EquippedItemResultProcessor {
   private static final Object lock = new Object();

   private static volatile EquippedItemResultProcessor instance;

   public static EquippedItemResultProcessor getInstance() {
      EquippedItemResultProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new EquippedItemResultProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder equipForCharacter(int characterId, int id) {
      Connection.instance()
            .element(entityManager -> EquipmentProvider.getById(entityManager, id))
            .map(EquipmentData::itemId)
            .flatMap(itemId -> ItemProcessor.getInstance().getEquipmentSlotDestination(itemId).findFirst())
            .ifPresent(destinationSlot -> ItemProcessor.getInstance().equipItemForCharacter(characterId, id, destinationSlot));

      return new ResultBuilder();
   }

   public ResultBuilder getEquippedItemsForCharacter(int characterId) {
      return new ResultBuilder();
   }
}
