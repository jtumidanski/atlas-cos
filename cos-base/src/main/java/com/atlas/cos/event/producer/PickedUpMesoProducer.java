package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.PickedUpMesoEvent;

public final class PickedUpMesoProducer {
   private PickedUpMesoProducer() {
   }

   public static void emit(int characterId, int meso) {
      EventProducerRegistry.getInstance()
            .send(PickedUpMesoEvent.class, EventConstants.TOPIC_PICKED_UP_MESO, characterId,
                  new PickedUpMesoEvent(characterId, meso));
   }
}
