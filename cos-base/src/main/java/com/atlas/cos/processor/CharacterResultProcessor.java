package com.atlas.cos.processor;

import java.util.Collections;

import com.app.rest.util.stream.Collectors;
import com.atlas.cos.database.provider.CharacterProvider;
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
      return Connection.getInstance()
            .withResult(entityManager -> CharacterProvider.getInstance().getForAccountAndWorld(entityManager,
                  accountId, worldId))
            .orElse(Collections.emptyList())
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public ResultBuilder getByName(String name) {
      return Connection.getInstance().withResult(entityManager -> CharacterProvider.getInstance().getForName(entityManager, name))
            .orElse(Collections.emptyList())
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
