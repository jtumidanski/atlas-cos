package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterExperienceEvent;

public final class CharacterExperienceGainProducer {
   private CharacterExperienceGainProducer() {
   }

   public static void gainExperience(int worldId, int channelId, int mapId, int characterId, int personalGain,
                                     int partyGain, boolean show, boolean chat, boolean white) {
      EventProducerRegistry.getInstance()
            .send(CharacterExperienceEvent.class, EventConstants.TOPIC_CHARACTER_EXPERIENCE_EVENT, worldId, channelId,
                  new CharacterExperienceEvent(worldId, channelId, mapId, characterId, personalGain, partyGain, show, chat,
                        white));
   }
}
