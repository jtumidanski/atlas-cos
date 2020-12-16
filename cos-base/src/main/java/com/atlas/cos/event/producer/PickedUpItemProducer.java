package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.PickedUpItemEvent;

public final class PickedUpItemProducer {
   private PickedUpItemProducer() {
   }

   public static void emit(int characterId, int itemId, int quantity) {
      EventProducerRegistry.getInstance()
            .send(PickedUpItemEvent.class, EventConstants.TOPIC_PICKED_UP_ITEM, characterId,
                  new PickedUpItemEvent(characterId, itemId, quantity));
   }
}
