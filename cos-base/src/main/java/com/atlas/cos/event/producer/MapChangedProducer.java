package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.MapChangedEvent;

public final class MapChangedProducer {
   private MapChangedProducer() {
   }

   public static void notifyChange(int worldId, int channelId, int characterId, int mapId, int portalId) {
      EventProducerRegistry.getInstance()
            .send(MapChangedEvent.class, EventConstants.TOPIC_CHANGE_MAP_EVENT, worldId, channelId, new MapChangedEvent(worldId,
                  channelId, mapId, portalId, characterId));
   }
}
