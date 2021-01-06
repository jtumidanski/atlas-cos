package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterCreatedEvent;
import com.atlas.cos.model.CharacterData;

public final class CharacterCreatedProducer {
   private CharacterCreatedProducer() {
   }

   public static void notifyCharacterCreated(CharacterData data) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_CHARACTER_CREATED_EVENT, data.id(),
            new CharacterCreatedEvent(data.worldId(), data.id(), data.name()));
   }
}
