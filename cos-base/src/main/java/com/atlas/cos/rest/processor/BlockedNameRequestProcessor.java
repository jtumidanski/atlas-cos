package com.atlas.cos.rest.processor;

import builder.ResultBuilder;
import com.app.rest.util.stream.Collectors;
import com.atlas.cos.database.administrator.BlockedNameAdministrator;
import com.atlas.cos.database.provider.BlockedNameProvider;
import com.atlas.cos.rest.ResultObjectFactory;
import database.Connection;

import java.util.List;

public final class BlockedNameRequestProcessor {
   private BlockedNameRequestProcessor() {
   }

   public static ResultBuilder getNames() {
      return Connection.instance()
            .list(BlockedNameProvider::getNames)
            .stream()
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }

   public static ResultBuilder getName(String name) {
      return Connection.instance()
            .list(BlockedNameProvider::getNames)
            .stream()
            .filter(blockedName -> blockedName.name().equalsIgnoreCase(name))
            .map(ResultObjectFactory::create)
            .collect(Collectors.toResultBuilder());
   }
}
