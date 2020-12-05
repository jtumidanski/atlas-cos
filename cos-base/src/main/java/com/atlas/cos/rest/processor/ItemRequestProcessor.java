package com.atlas.cos.rest.processor;

import builder.ResultBuilder;
import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.processor.ItemProcessor;
import com.atlas.cos.rest.ResultObjectFactory;

import javax.ws.rs.core.Response;

public final class ItemRequestProcessor {
   private ItemRequestProcessor() {
   }

   public static ResultBuilder createEquipmentForCharacter(int characterId, EquipmentAttributes attributes, boolean characterCreation) {
      return ItemProcessor.createEquipmentForCharacter(characterId, attributes.itemId(), attributes.strength(),
            attributes.dexterity(), attributes.intelligence(), attributes.luck(), attributes.weaponAttack(),
            attributes.weaponDefense(), attributes.magicAttack(), attributes.magicDefense(), attributes.accuracy(),
            attributes.avoidability(), attributes.speed(), attributes.jump(), attributes.hp(), attributes.mp(), attributes.slots(),
            characterCreation)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }

   public static ResultBuilder getEquipmentForCharacter(int characterId, int equipmentId) {
      return ItemProcessor.getEquipmentForCharacter(characterId, equipmentId)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }

   public static ResultBuilder getEquippedItemsForCharacter(Integer characterId) {
      return ItemProcessor.getEquipmentForCharacter(characterId).stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
