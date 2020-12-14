package com.atlas.cos;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

import com.atlas.cos.event.CharacterCreatedEvent;
import com.atlas.cos.event.CharacterExperienceEvent;
import com.atlas.cos.event.MapChangedEvent;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class EventProducerRegistry {
   private static final Object lock = new Object();

   private static volatile EventProducerRegistry instance;

   private final Map<Class<?>, Producer<Long, ?>> producerMap;

   public static EventProducerRegistry getInstance() {
      EventProducerRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new EventProducerRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private EventProducerRegistry() {
      producerMap = new HashMap<>();
      producerMap.put(CharacterCreatedEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(MapChangedEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(CharacterExperienceEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
   }

   public <T> void send(Class<T> clazz, String topic, int worldId, int channelId, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(System.getenv(topic), produceKey(worldId, channelId), event);
      getProducer(clazz).ifPresent(producer -> producer.send(record));
   }

   public <T> void send(Class<T> clazz, String topic, long key, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(System.getenv(topic), key, event);
      getProducer(clazz).ifPresent(producer -> producer.send(record));
   }

   protected <T> Optional<Producer<Long, T>> getProducer(Class<T> clazz) {
      Producer<Long, T> producer = null;
      if (producerMap.containsKey(clazz)) {
         producer = (Producer<Long, T>) producerMap.get(clazz);
      }
      return Optional.ofNullable(producer);
   }

   protected static Long produceKey(int worldId, int channelId) {
      return (long) ((worldId * 1000) + channelId);
   }
}
