package com.atlas.cos.event.consumer;

import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterLevelEvent;
import com.atlas.cos.processor.CharacterProcessor;

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
      return System.getenv(EventConstants.TOPIC_CHARACTER_LEVEL_EVENT);
   }
}
