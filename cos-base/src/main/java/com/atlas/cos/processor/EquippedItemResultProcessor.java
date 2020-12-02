package com.atlas.cos.processor;

import builder.ResultBuilder;

public class EquippedItemResultProcessor {
   private static final Object lock = new Object();

   private static volatile EquippedItemResultProcessor instance;

   public static EquippedItemResultProcessor getInstance() {
      EquippedItemResultProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new EquippedItemResultProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public ResultBuilder equipForCharacter(int characterId, int id) {
      ItemProcessor.getInstance().equipItemForCharacter(characterId, id);
      return new ResultBuilder();
   }
}
