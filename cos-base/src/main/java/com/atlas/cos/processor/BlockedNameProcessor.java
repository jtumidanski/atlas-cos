package com.atlas.cos.processor;

import java.util.List;

import com.app.rest.util.stream.Collectors;
import com.atlas.cos.database.administrator.BlockedNameAdministrator;
import com.atlas.cos.database.provider.BlockedNameProvider;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;
import database.Connection;

public class BlockedNameProcessor {
   private static final Object lock = new Object();

   private static volatile BlockedNameProcessor instance;

   public static BlockedNameProcessor getInstance() {
      BlockedNameProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new BlockedNameProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder getNames() {
      return Connection.instance()
            .list(BlockedNameProvider::getNames)
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public ResultBuilder getName(String name) {
      return Connection.instance()
            .list(BlockedNameProvider::getNames)
            .stream()
            .filter(blockedName -> blockedName.name().equalsIgnoreCase(name))
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public void bulkAddBlockedNames(List<String> names) {
      Connection.instance()
            .with(entityManager -> BlockedNameAdministrator.create(entityManager, names));
   }
}
