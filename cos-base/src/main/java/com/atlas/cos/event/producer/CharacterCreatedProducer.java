package com.atlas.cos.event.producer;

import com.atlas.cos.ProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterCreatedEvent;
import com.atlas.cos.model.CharacterData;
import org.apache.kafka.clients.producer.ProducerRecord;

public final class CharacterCreatedProducer {
   private CharacterCreatedProducer() {
   }

   public static void notifyCharacterCreated(CharacterData data) {
      String topic = System.getenv(EventConstants.TOPIC_CHARACTER_CREATED_EVENT);
      long key = data.id();
      ProducerRegistry.getInstance()
            .getCharacterCreatedProducer()
            .send(new ProducerRecord<>(topic, key, new CharacterCreatedEvent(data.worldId(), data.id(), data.name())));
   }
}
