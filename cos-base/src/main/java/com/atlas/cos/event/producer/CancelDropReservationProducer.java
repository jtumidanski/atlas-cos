package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.drg.command.CancelDropReservationCommand;
import com.atlas.drg.constant.EventConstants;

public final class CancelDropReservationProducer {
   private CancelDropReservationProducer() {
   }

   public static void emit(int dropId, int characterId) {
      EventProducerRegistry.getInstance()
            .send(CancelDropReservationCommand.class, EventConstants.TOPIC_CANCEL_DROP_RESERVATION_COMMAND, dropId,
                  new CancelDropReservationCommand(characterId, dropId));
   }
}
