package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.drg.command.PickupDropCommand;
import com.atlas.drg.constant.EventConstants;

public final class PickupDropCommandProducer {
   private PickupDropCommandProducer() {
   }

   public static void emit(int characterId, int dropId) {
      EventProducerRegistry.getInstance()
            .send(PickupDropCommand.class, EventConstants.TOPIC_PICKUP_DROP_COMMAND, dropId,
                  new PickupDropCommand(dropId, characterId));
   }
}
