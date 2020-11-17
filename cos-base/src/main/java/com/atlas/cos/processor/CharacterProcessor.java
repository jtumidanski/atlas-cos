package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.entity.Character;
import com.atlas.cos.model.CharacterData;

import database.DatabaseConnection;

public class CharacterProcessor {
   private static final Object lock = new Object();

   private static volatile CharacterProcessor instance;

   public static CharacterProcessor getInstance() {
      CharacterProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public Optional<CharacterData> getByName(String name) {
      return DatabaseConnection.getInstance().withConnectionResult(entityManager ->
            CharacterProvider.getInstance().getForName(entityManager, name)
                  .stream()
                  .findFirst()
                  .orElse(null)
      );
   }

   public CharacterData createBeginner(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 10000);
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   protected CharacterData create(CharacterBuilder builder) {
      CharacterData original = builder.build();

      CharacterData characterData = DatabaseConnection.getInstance().withConnectionResult(entityManager ->
            CharacterAdministrator.getInstance().create(entityManager,
                  original.accountId(), original.worldId(), original.name(), original.level(),
                  original.strength(), original.dexterity(), original.luck(), original.intelligence(),
                  original.maxHp(), original.maxMp(), original.jobId(), original.gender(), original.hair(),
                  original.face(), original.mapId())
      ).orElseThrow();
      
      return characterData;
   }

   public CharacterData createNoblesse(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 130030000);
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public CharacterData createLegend(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 914000000);
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      return create(builder);
   }
}