package com.atlas.cos.event.consumer;

import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterLevelEvent;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class CharacterLevelConsumer extends AbstractEventConsumer<CharacterLevelEvent> {
   @Override
   public void handle(Long key, CharacterLevelEvent event) {
      CharacterProcessor.increaseLevel(event.characterId());
   }

   @Override
   public Class<CharacterLevelEvent> getEventClass() {
      return CharacterLevelEvent.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_CHARACTER_LEVEL_EVENT);
   }
}
