package com.atlas.cos.rest.processor;

import builder.ResultBuilder;
import com.atlas.cos.processor.ItemProcessor;

public final class EquippedItemRequestProcessor {
   private EquippedItemRequestProcessor() {
   }

   public static ResultBuilder equipForCharacter(int characterId, int id) {
      ItemProcessor.equipItemForCharacter(characterId, id);
      return new ResultBuilder();
   }
}
