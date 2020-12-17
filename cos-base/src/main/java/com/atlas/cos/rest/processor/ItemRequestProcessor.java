package com.atlas.cos.rest.processor;

import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.cos.processor.ItemProcessor;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;

public final class ItemRequestProcessor {
   private ItemRequestProcessor() {
   }

   public static ResultBuilder getEquipmentForCharacter(int characterId, int equipmentId) {
      return ItemProcessor.getEquipmentForCharacter(characterId, equipmentId)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }

   public static ResultBuilder getEquippedItemsForCharacter(int characterId) {
      return ItemProcessor.getEquipmentForCharacter(characterId).stream()
            .filter(equipment -> equipment.slot() < 0)
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public static ResultBuilder getEquipsForCharacter(int characterId) {
      return ItemProcessor.getEquipmentForCharacter(characterId).stream()
            .filter(equipment -> equipment.slot() >= 0)
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
