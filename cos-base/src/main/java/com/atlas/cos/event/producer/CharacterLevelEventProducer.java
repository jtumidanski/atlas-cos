package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterLevelEvent;

public final class CharacterLevelEventProducer {
   private CharacterLevelEventProducer() {
   }

   public static void gainLevel(int worldId, int channelId, int mapId, int characterId) {
      EventProducerRegistry.getInstance()
            .send(CharacterLevelEvent.class, EventConstants.TOPIC_CHARACTER_LEVEL_EVENT, worldId, channelId,
                  new CharacterLevelEvent(worldId, channelId, mapId, characterId));
   }
}
