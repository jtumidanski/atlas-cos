package com.atlas.cos.processor;

import java.util.stream.Stream;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Mappers;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.database.provider.EquipmentProvider;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;
import database.Connection;

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
      if (characterCreation) {
         boolean valid = validCharacterCreationItem(attributes.itemId());
         if (!valid) {
            return new ResultBuilder(Response.Status.FORBIDDEN);
         }
      }

      short nextOpenSlot = Connection.instance()
            .element(entityManager -> EquipmentProvider.getNextFreeEquipmentSlot(entityManager, characterId))
            .orElse((short) 0);
      return ItemProcessor.getInstance()
            .createEquipment(characterId, attributes.itemId(), nextOpenSlot, attributes.strength(), attributes.dexterity(),
                  attributes.intelligence(),
                  attributes.luck(), attributes.weaponAttack(), attributes.weaponDefense(), attributes.magicAttack(),
                  attributes.magicDefense(), attributes.accuracy(), attributes.avoidability(), attributes.speed(),
                  attributes.jump(), attributes.hp(), attributes.mp(), attributes.slots())
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.INTERNAL_SERVER_ERROR));
   }

   protected boolean validCharacterCreationItem(int itemId) {
      return Stream.of(
            1302000, 1312004, 1322005, 1442079,// weapons
            1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167,// bottom
            1060002, 1060006, 1061002, 1061008, 1062115, // top
            1072001, 1072005, 1072037, 1072038, 1072383,// shoes
            30000, 30010, 30020, 30030, 31000, 31040, 31050,// hair
            20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100  //face
      ).anyMatch(id -> id == itemId);
   }

   public ResultBuilder getEquipmentForCharacter(int characterId, int equipmentId) {
      return ItemProcessor.getInstance().getEquipmentForCharacter(characterId, equipmentId)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }
}
