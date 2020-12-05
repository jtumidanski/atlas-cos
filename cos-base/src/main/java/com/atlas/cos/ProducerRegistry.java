package com.atlas.cos;

import com.atlas.cos.event.CharacterCreatedEvent;
import com.atlas.cos.event.MapChangedEvent;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;

public class ProducerRegistry {
   private static final Object lock = new Object();

   private static volatile ProducerRegistry instance;

   private final String producerId;

   private final String bootstrapServers;

   private Producer<Long, CharacterCreatedEvent> characterCreatedProducer;

   private Producer<Long, MapChangedEvent> mapChangedProducer;

   public static ProducerRegistry getInstance() {
      ProducerRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ProducerRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private ProducerRegistry() {
      producerId = "Character Service";
      bootstrapServers = System.getenv("BOOTSTRAP_SERVERS");
   }

   public Producer<Long, CharacterCreatedEvent> getCharacterCreatedProducer() {
      if (characterCreatedProducer == null) {
         synchronized (lock) {
            characterCreatedProducer = KafkaProducerFactory.createProducer(producerId, bootstrapServers);
         }
      }
      return characterCreatedProducer;
   }

   public Producer<Long, MapChangedEvent> getMapChangedProducer() {
      if (mapChangedProducer == null) {
         synchronized (lock) {
            mapChangedProducer = KafkaProducerFactory.createProducer(producerId, bootstrapServers);
         }
      }
      return mapChangedProducer;
   }
}
