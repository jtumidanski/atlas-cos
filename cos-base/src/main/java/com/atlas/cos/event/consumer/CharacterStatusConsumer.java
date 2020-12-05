package com.atlas.cos.event.consumer;

import com.atlas.cos.model.Portal;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.CharacterTemporalRegistry;
import com.atlas.cos.processor.PortalProcessor;
import com.atlas.csrv.event.CharacterStatusEvent;
import com.atlas.csrv.event.CharacterStatusEventType;
import com.atlas.kafka.consumer.ConsumerRecordHandler;

public class CharacterStatusConsumer implements ConsumerRecordHandler<Long, CharacterStatusEvent> {
   @Override
   public void handle(Long aLong, CharacterStatusEvent event) {
      if (event.type().equals(CharacterStatusEventType.LOGIN)) {
         CharacterProcessor.getInstance().getById(event.characterId())
               .ifPresent(characterData -> {
                  Portal portal = PortalProcessor.getInstance()
                        .getMapPortalById(characterData.mapId(), characterData.spawnPoint())
                        .orElse(PortalProcessor.getInstance().getMapPortalById(characterData.mapId(), 0).orElseThrow());
                  CharacterTemporalRegistry.getInstance().updatePosition(event.characterId(), portal.x(), portal.y());
               });
      }
   }
}
