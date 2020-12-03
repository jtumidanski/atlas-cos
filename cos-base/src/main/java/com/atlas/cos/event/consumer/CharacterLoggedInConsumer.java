package com.atlas.cos.event.consumer;

import com.atlas.cos.model.Portal;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.CharacterTemporalRegistry;
import com.atlas.cos.processor.PortalProcessor;
import com.atlas.csrv.event.CharacterLoggedInEvent;
import com.atlas.kafka.consumer.ConsumerRecordHandler;

public class CharacterLoggedInConsumer implements ConsumerRecordHandler<Long, CharacterLoggedInEvent> {
   @Override
   public void handle(Long aLong, CharacterLoggedInEvent event) {
      CharacterProcessor.getInstance().getById(event.characterId())
            .ifPresent(characterData -> {
               Portal portal = PortalProcessor.getInstance()
                     .getMapPortalById(characterData.mapId(), characterData.spawnPoint())
                     .orElse(PortalProcessor.getInstance().getMapPortalById(characterData.mapId(), 0).orElseThrow());
               CharacterTemporalRegistry.getInstance().updatePosition(event.characterId(), portal.x(), portal.y());
            });
   }
}
