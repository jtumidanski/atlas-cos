package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterExpressionChangedEvent;

public final class FacialExpressionProducer {
   private FacialExpressionProducer() {
   }

   public static void emit(int characterId, int mapId, int emote) {
      EventProducerRegistry.getInstance().send(EventConstants.EXPRESSION_CHANGED, characterId,
            new CharacterExpressionChangedEvent(characterId, mapId, emote));
   }
}
