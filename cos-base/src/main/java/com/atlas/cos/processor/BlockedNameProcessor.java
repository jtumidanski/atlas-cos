package com.atlas.cos.processor;

import java.util.List;

import com.atlas.cos.attribute.BlockedNameAttributes;
import com.atlas.cos.builder.BlockedNameAttributesBuilder;
import com.atlas.cos.database.administrator.BlockedNameAdministrator;
import com.atlas.cos.database.provider.BlockedNameProvider;

import builder.ResultBuilder;
import builder.ResultObjectBuilder;
import database.DatabaseConnection;

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
      ResultBuilder resultBuilder = new ResultBuilder();
      DatabaseConnection.getInstance().withConnection(entityManager ->
            BlockedNameProvider.getInstance().getNames(entityManager).forEach(blockedNameData ->
                  resultBuilder.addData(new ResultObjectBuilder(BlockedNameAttributes.class, blockedNameData.name())
                        .setAttribute(new BlockedNameAttributesBuilder().setName(blockedNameData.name()))
                  )));
      return resultBuilder;
   }

   public void bulkAddBlockedNames(List<String> names) {
      DatabaseConnection.getInstance().withConnection(entityManager ->
            names.forEach(name -> BlockedNameAdministrator.getInstance().createBlockedName(entityManager, name)));
   }

   public ResultBuilder getName(String name) {
      ResultBuilder resultBuilder = new ResultBuilder();
      DatabaseConnection.getInstance().withConnection(entityManager ->
            BlockedNameProvider.getInstance().getNames(entityManager).stream()
                  .filter(blockedName -> blockedName.name().equalsIgnoreCase(name))
                  .forEach(blockedNameData ->
                        resultBuilder.addData(new ResultObjectBuilder(BlockedNameAttributes.class, blockedNameData.name())
                              .setAttribute(new BlockedNameAttributesBuilder().setName(blockedNameData.name()))
                        )));
      return resultBuilder;
   }
}
