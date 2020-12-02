package com.atlas.cos.processor;

import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;

public class ItemResultProcessor {
   private static final Object lock = new Object();

   private static volatile ItemResultProcessor instance;

   public static ItemResultProcessor getInstance() {
      ItemResultProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ItemResultProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder createEquipmentForCharacter(int characterId, EquipmentAttributes attributes, boolean characterCreation) {
      return ItemProcessor.getInstance().createEquipmentForCharacter(characterId, attributes.itemId(), attributes.strength(),
            attributes.dexterity(), attributes.intelligence(), attributes.luck(), attributes.weaponAttack(),
            attributes.weaponDefense(), attributes.magicAttack(), attributes.magicDefense(), attributes.accuracy(),
            attributes.avoidability(), attributes.speed(), attributes.jump(), attributes.hp(), attributes.mp(), attributes.slots(),
            characterCreation)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }

   public ResultBuilder getEquipmentForCharacter(int characterId, int equipmentId) {
      return ItemProcessor.getInstance().getEquipmentForCharacter(characterId, equipmentId)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }

   public ResultBuilder getEquippedItemsForCharacter(Integer characterId) {
      return ItemProcessor.getInstance().getEquipmentForCharacter(characterId).stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
