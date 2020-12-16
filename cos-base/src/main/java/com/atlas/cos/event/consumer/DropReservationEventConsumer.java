package com.atlas.cos.event.consumer;

import com.atlas.cos.processor.DropProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;
import com.atlas.drg.constant.EventConstants;
import com.atlas.drg.event.DropReservationEvent;
import com.atlas.drg.event.DropReservationType;

public class DropReservationEventConsumer extends AbstractEventConsumer<DropReservationEvent> {
   @Override
   public void handle(Long key, DropReservationEvent event) {
      if (event.type().equals(DropReservationType.SUCCESS)) {
         DropProcessor.attemptPickup(event.characterId(), event.dropId());
      }
   }

   @Override
   public Class<DropReservationEvent> getEventClass() {
      return DropReservationEvent.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_DROP_RESERVATION_EVENT);
   }
}
