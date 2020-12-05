package com.atlas.cos.event.consumer;

import com.atlas.cos.model.Portal;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.processor.PortalProcessor;
import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.CharacterStatusEvent;
import com.atlas.csrv.event.CharacterStatusEventType;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class CharacterStatusConsumer implements SimpleEventHandler<CharacterStatusEvent> {
   @Override
   public void handle(Long aLong, CharacterStatusEvent event) {
      if (event.type().equals(CharacterStatusEventType.LOGIN)) {
         CharacterProcessor.getById(event.characterId())
               .ifPresent(characterData -> {
                  Portal portal = PortalProcessor.getMapPortalById(characterData.mapId(), characterData.spawnPoint())
                        .orElse(PortalProcessor.getMapPortalById(characterData.mapId(), 0).orElseThrow());
                  CharacterTemporalRegistry.getInstance().updatePosition(event.characterId(), portal.x(), portal.y());
               });
      }
   }

   @Override
   public Class<CharacterStatusEvent> getEventClass() {
      return CharacterStatusEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "Character Service";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return System.getenv(EventConstants.TOPIC_CHARACTER_STATUS);
   }
}
