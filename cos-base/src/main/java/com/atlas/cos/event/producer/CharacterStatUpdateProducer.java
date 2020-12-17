package com.atlas.cos.event.producer;

import java.util.Collection;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterStatUpdateEvent;
import com.atlas.cos.event.StatUpdateType;

public final class CharacterStatUpdateProducer {
   private CharacterStatUpdateProducer() {
   }

   public static void statsUpdated(int characterId, Collection<StatUpdateType> updates) {
      EventProducerRegistry.getInstance()
            .send(CharacterStatUpdateEvent.class, EventConstants.TOPIC_CHARACTER_STAT_EVENT, characterId,
                  new CharacterStatUpdateEvent(characterId, updates));
   }
}
