package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.MesoGainedEvent;

public final class MesoGainedProducer {
   private MesoGainedProducer() {
   }

   public static void emit(int characterId, int meso) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_MESO_GAINED, characterId,
            new MesoGainedEvent(characterId, meso));
   }
}
