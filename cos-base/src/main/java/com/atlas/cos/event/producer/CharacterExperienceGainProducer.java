package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterExperienceEvent;

public final class CharacterExperienceGainProducer {
   private CharacterExperienceGainProducer() {
   }

   public static void gainExperience(int characterId, int personalGain, int partyGain, boolean show, boolean chat, boolean white) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_CHARACTER_EXPERIENCE_EVENT, characterId,
            new CharacterExperienceEvent(characterId, personalGain, partyGain, show, chat, white));
   }
}
