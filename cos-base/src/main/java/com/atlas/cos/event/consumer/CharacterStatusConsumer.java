package com.atlas.cos.event.consumer;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.model.Portal;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.PortalProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;
import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.CharacterStatusEvent;
import com.atlas.csrv.event.CharacterStatusEventType;

public class CharacterStatusConsumer extends AbstractEventConsumer<CharacterStatusEvent> {
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
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_CHARACTER_STATUS);
   }
}
