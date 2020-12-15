package com.atlas.cos.event.consumer;

import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterExperienceEvent;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class CharacterExperienceConsumer extends AbstractEventConsumer<CharacterExperienceEvent> {
   @Override
   public void handle(Long key, CharacterExperienceEvent event) {
      CharacterProcessor.gainExperience(event.worldId(), event.channelId(), event.mapId(), event.characterId(), event.personalGain() + event.partyGain());
   }

   @Override
   public Class<CharacterExperienceEvent> getEventClass() {
      return CharacterExperienceEvent.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_CHARACTER_EXPERIENCE_EVENT);
   }
}
