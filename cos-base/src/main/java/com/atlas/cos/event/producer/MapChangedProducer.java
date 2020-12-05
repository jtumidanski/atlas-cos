package com.atlas.cos.event.producer;

import com.atlas.cos.ProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.MapChangedEvent;
import org.apache.kafka.clients.producer.ProducerRecord;

public final class MapChangedProducer {
   private MapChangedProducer() {
   }

   public static void notifyChange(int worldId, int channelId, int characterId, int mapId, int portalId) {
      String topic = System.getenv(EventConstants.TOPIC_CHANGE_MAP_EVENT);
      long key = produceKey(worldId, channelId);
      ProducerRegistry.getInstance()
            .getMapChangedProducer()
            .send(new ProducerRecord<>(topic, key, new MapChangedEvent(worldId, channelId, mapId, portalId, characterId)));
   }

   protected static Long produceKey(int worldId, int channelId) {
      return (long) ((worldId * 1000) + channelId);
   }
}
