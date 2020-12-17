package com.atlas.cos.event.consumer;

import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.PickedUpMesoEvent;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class PickedUpMesoConsumer extends AbstractEventConsumer<PickedUpMesoEvent> {
   @Override
   public void handle(Long key, PickedUpMesoEvent event) {
      CharacterProcessor.gainMeso(event.characterId(), event.gain());
   }

   @Override
   public Class<PickedUpMesoEvent> getEventClass() {
      return PickedUpMesoEvent.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_PICKED_UP_MESO);
   }
}
