package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.PickedUpNxEvent;

public final class PickedUpNxProducer {
   private PickedUpNxProducer() {
   }

   public static void emit(int characterId, int gain) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_PICKED_UP_NX, characterId,
            new PickedUpNxEvent(characterId, gain));
   }
}
