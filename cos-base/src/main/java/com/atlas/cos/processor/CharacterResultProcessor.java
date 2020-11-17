package com.atlas.cos.processor;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterAttributesBuilder;
import com.atlas.cos.database.provider.CharacterProvider;
import com.sun.xml.bind.v2.runtime.JAXBContextImpl;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;
import database.DatabaseConnection;

public class CharacterResultProcessor {
   private static final Object lock = new Object();

   private static volatile CharacterResultProcessor instance;

   public static CharacterResultProcessor getInstance() {
      CharacterResultProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new CharacterResultProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder getForAccountAndWorld(int accountId, int worldId) {
      ResultBuilder resultBuilder = new ResultBuilder();
      DatabaseConnection.getInstance().withConnection(entityManager ->
            CharacterProvider.getInstance().getForAccountAndWorld(entityManager, accountId, worldId)
                  .stream()
                  .map(this::produceResultObjectForCharacter)
                  .forEach(resultBuilder::addData));
      return resultBuilder;
   }

   private ResultObjectBuilder produceResultObjectForCharacter(com.atlas.cos.model.CharacterData data) {
      return new ResultObjectBuilder(CharacterAttributes.class, data.id())
            .setAttribute(new CharacterAttributesBuilder()
                  .setAccountId(data.accountId())
                  .setWorldId(data.worldId())
                  .setName(data.name())
                  .setLevel(data.level())
                  .setExperience(data.experience())
                  .setGachaponExperience(data.gachaponExperience())
                  .setStrength(data.strength())
                  .setDexterity(data.dexterity())
                  .setLuck(data.luck())
                  .setIntelligence(data.intelligence())
                  .setHp(data.hp())
                  .setMp(data.mp())
                  .setMaxHp(data.maxHp())
                  .setMaxMp(data.maxMp())
                  .setMeso(data.meso())
                  .setHpMpUsed(data.hpMpUsed())
                  .setJobId(data.jobId())
                  .setSkinColor(data.skinColor())
                  .setGender(data.gender())
                  .setFame(data.fame())
                  .setHair(data.hair())
                  .setFace(data.face())
                  .setAp(data.ap())
                  .setSp(data.sp())
                  .setMapId(data.mapId())
                  .setSpawnPoint(data.spawnPoint())
                  .setGm(data.gm())
            );
   }

   public ResultBuilder getByName(String name) {
      ResultBuilder resultBuilder = new ResultBuilder();
      DatabaseConnection.getInstance().withConnection(entityManager ->
            CharacterProvider.getInstance().getForName(entityManager, name)
                  .stream()
                  .map(this::produceResultObjectForCharacter)
                  .forEach(resultBuilder::addData));
      return resultBuilder;
   }
}
