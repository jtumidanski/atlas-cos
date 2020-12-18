package com.atlas.cos;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

import com.atlas.cos.event.CharacterCreatedEvent;
import com.atlas.cos.event.CharacterExperienceEvent;
import com.atlas.cos.event.CharacterInventoryModifyEvent;
import com.atlas.cos.event.CharacterLevelEvent;
import com.atlas.cos.event.CharacterStatUpdateEvent;
import com.atlas.cos.event.MapChangedEvent;
import com.atlas.cos.event.PickedUpItemEvent;
import com.atlas.cos.event.PickedUpMesoEvent;
import com.atlas.cos.event.PickedUpNxEvent;
import com.atlas.cos.processor.TopicDiscoveryProcessor;
import com.atlas.drg.command.CancelDropReservationCommand;
import com.atlas.drg.command.PickupDropCommand;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class EventProducerRegistry {
   private static final Object lock = new Object();

   private static volatile EventProducerRegistry instance;

   private final Map<Class<?>, Producer<Long, ?>> producerMap;

   private final Map<String, String> topicMap;

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
      producerMap.put(CharacterLevelEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(CharacterStatUpdateEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(PickedUpItemEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(PickedUpMesoEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(PickedUpNxEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(PickupDropCommand.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(CancelDropReservationCommand.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      producerMap.put(CharacterInventoryModifyEvent.class,
            KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS")));
      topicMap = new HashMap<>();
   }

   public <T> void send(Class<T> clazz, String topic, int worldId, int channelId, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(getTopic(topic), produceKey(worldId, channelId), event);
      getProducer(clazz).ifPresent(producer -> producer.send(record));
   }

   public <T> void send(Class<T> clazz, String topic, long key, T event) {
      ProducerRecord<Long, T> record = new ProducerRecord<>(getTopic(topic), key, event);
      getProducer(clazz).ifPresent(producer -> producer.send(record));
   }

   protected String getTopic(String id) {
      if (!topicMap.containsKey(id)) {
         topicMap.put(id, TopicDiscoveryProcessor.getTopic(id));
      }
      return topicMap.get(id);
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
