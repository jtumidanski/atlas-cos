package com.atlas.cos.processor;

import java.util.Arrays;
import java.util.Optional;
import java.util.function.Function;
import java.util.stream.Stream;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Collectors;
import com.app.rest.util.stream.Mappers;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.MapleJob;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;
import database.Connection;

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
      return Connection.instance()
            .list(entityManager -> CharacterProvider.getForAccountAndWorld(entityManager, accountId, worldId))
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public ResultBuilder getByName(String name) {
      return Connection.instance()
            .list(entityManager -> CharacterProvider.getForName(entityManager, name))
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public ResultBuilder createCharacter(CharacterAttributes attributes) {
      if (!validFace(attributes.face()) || !validHair(attributes.hair())) {
         System.out.println("Owner from account '" + attributes.accountId() + "' tried to packet edit in char creation.");
         return new ResultBuilder(Response.Status.UNAUTHORIZED);
      }

      Function<CharacterAttributes, Optional<CharacterData>> creator;
      MapleJob job = MapleJob.getById(attributes.jobId());
      if (MapleJob.BEGINNER.equals(job)) {
         creator = CharacterProcessor.getInstance()::createBeginner;
      } else if (MapleJob.NOBLESSE.equals(job)) {
         creator = CharacterProcessor.getInstance()::createNoblesse;
      } else if (MapleJob.LEGEND.equals(job)) {
         creator = CharacterProcessor.getInstance()::createLegend;
      } else {
         return new ResultBuilder(Response.Status.NOT_IMPLEMENTED);
      }

      return creator.apply(attributes)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElse(new ResultBuilder(Response.Status.FORBIDDEN));
   }

   protected boolean validFace(int face) {
      return Arrays.asList(20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100).contains(face);
   }

   protected boolean validHair(int hair) {
      return Stream.of(30000, 30010, 30020, 30030, 31000, 31040, 31050)
            .anyMatch(id -> Math.floor(id / 10.0) == Math.floor(hair / 10.0));
   }

   public ResultBuilder getById(int characterId) {
      return CharacterProcessor.getInstance()
            .getById(characterId)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleOkResult)
            .orElse(new ResultBuilder(Response.Status.NOT_FOUND));
   }

   public ResultBuilder getForWorldInMap(int worldId, int mapId) {
      return Connection.instance()
            .list(entityManager -> CharacterProvider.getForWorldInMap(entityManager, worldId, mapId))
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
