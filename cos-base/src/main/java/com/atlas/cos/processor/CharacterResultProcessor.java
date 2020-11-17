package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Stream;
import javax.ws.rs.core.Response;

import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterAttributesBuilder;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.MapleJob;

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
      CharacterProcessor.getInstance().getByName(name)
            .ifPresent(characterData -> resultBuilder.addData(produceResultObjectForCharacter(characterData)));
      return resultBuilder;
   }

   public ResultBuilder createCharacter(CharacterAttributes attributes) {
      if (!validFace(attributes.getFace()) || !validHair(attributes.getHair())) {
         System.out.println("Owner from account '" + attributes.getAccountId() + "' tried to packet edit in char creation.");
         return new ResultBuilder(Response.Status.UNAUTHORIZED);
      }

      CharacterData character = null;
      MapleJob job = MapleJob.getById(attributes.getJobId());
      if (MapleJob.BEGINNER.equals(job)) {
         character = CharacterProcessor.getInstance().createBeginner(attributes);
      } else if (MapleJob.NOBLESSE.equals(job)) {
         character = CharacterProcessor.getInstance().createNoblesse(attributes);
      } else if (MapleJob.LEGEND.equals(job)) {
         character = CharacterProcessor.getInstance().createLegend(attributes);
      }

      ResultBuilder resultBuilder = new ResultBuilder(Response.Status.FORBIDDEN);
      if (character != null) {
         resultBuilder.setStatus(Response.Status.CREATED);
         resultBuilder.addData(produceResultObjectForCharacter(character));
      }
      return resultBuilder;
   }

   private final Set<Integer> IDs = new HashSet<>(
         Arrays.asList(
               1302000, 1312004, 1322005, 1442079,// weapons
               1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167,// bottom
               1060002, 1060006, 1061002, 1061008, 1062115, // top
               1072001, 1072005, 1072037, 1072038, 1072383,// shoes
               30000, 30010, 30020, 30030, 31000, 31040, 31050,// hair
               20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100  //face
         ));

   protected boolean validFace(int face) {
      return Arrays.asList(20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100).contains(face);
   }

   protected boolean validHair(int hair) {
      return Stream.of(30000, 30010, 30020, 30030, 31000, 31040, 31050)
            .anyMatch(id -> Math.floor(id / 10.0) == Math.floor(hair / 10.0));
   }
}
