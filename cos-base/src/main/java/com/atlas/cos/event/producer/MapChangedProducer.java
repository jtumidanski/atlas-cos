package com.atlas.cos.event.producer;

import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.MapChangedEvent;
import com.atlas.kafka.KafkaProducerFactory;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class MapChangedProducer {
   private static final Object lock = new Object();

   private static volatile MapChangedProducer instance;

   private final Producer<Long, MapChangedEvent> producer;

   public static MapChangedProducer getInstance() {
      MapChangedProducer result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new MapChangedProducer();
               instance = result;
            }
         }
      }
      return result;
   }

   private MapChangedProducer() {
      producer = KafkaProducerFactory.createProducer("Character Service", System.getenv("BOOTSTRAP_SERVERS"));
   }

   public void notifyChange(int worldId, int channelId, int characterId, int mapId, int portalId) {
      String topic = System.getenv(EventConstants.TOPIC_CHANGE_MAP_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key, new MapChangedEvent(worldId, channelId, mapId, portalId, characterId)));
   }

   protected Long produceKey(int worldId, int channelId) {
      return (long) ((worldId * 1000) + channelId);
   }
}
