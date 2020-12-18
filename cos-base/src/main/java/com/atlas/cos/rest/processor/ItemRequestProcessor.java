package com.atlas.cos.rest.processor;

import com.atlas.cos.model.InventoryType;
import com.atlas.cos.processor.ItemProcessor;

import builder.ResultBuilder;

public final class ItemRequestProcessor {
   private ItemRequestProcessor() {
   }

   //   public static ResultBuilder getEquipmentForCharacter(int characterId, int equipmentId) {
   //      return ItemProcessor.getEquipmentForCharacter(characterId, equipmentId)
   //            .map(ResultObjectFactory::create)
   //            .map(Mappers::singleOkResult)
   //            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   //   }
   //
   //   public static ResultBuilder getEquippedItemsForCharacter(int characterId) {
   //      return ItemProcessor.getEquipmentForCharacter(characterId).stream()
   //            .filter(equipment -> equipment.slot() < 0)
   //            .map(ResultObjectFactory::create)
   //            .collect(Collectors.toResultBuilder());
   //   }
   //
   //   public static ResultBuilder getEquipsForCharacter(int characterId) {
   //      return ItemProcessor.getEquipmentForCharacter(characterId).stream()
   //            .filter(equipment -> equipment.slot() >= 0)
   //            .map(ResultObjectFactory::create)
   //            .collect(Collectors.toResultBuilder());
   //   }

//   public static ResultBuilder getItemsForCharacterByInventoryType(int characterId, String type) {
//      return InventoryType.getByName(type)
//            .map(inventoryType -> {
//               if (inventoryType.equals(InventoryType.EQUIP)) {
//
//               }
//               ItemProcessor.getItemsForCharacter(characterId, type)
//            })
//   }
}
