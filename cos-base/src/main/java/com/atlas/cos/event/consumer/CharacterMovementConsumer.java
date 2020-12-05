package com.atlas.cos.event.consumer;

import com.atlas.cos.model.CharacterData;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.CharacterTemporalRegistry;
import com.atlas.cos.processor.MapProcessor;
import com.atlas.csrv.constant.EventConstants;
import com.atlas.csrv.event.CharacterMovementEvent;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class CharacterMovementConsumer implements SimpleEventHandler<CharacterMovementEvent> {
   @Override
   public void handle(Long aLong, CharacterMovementEvent event) {
      // Update tracked position and stance.
      if (event.x() != null && event.y() != null) {
         CharacterTemporalRegistry.getInstance().update(event.characterId(), event.x(), event.y(), event.stance());
      } else if (event.stance() != null) {
         CharacterTemporalRegistry.getInstance().updateStance(event.characterId(), event.stance());
      }

      // Update spawn point.
      if (event.x() != null && event.y() != null) {
         CharacterProcessor.getInstance().getById(event.characterId())
               .map(CharacterData::mapId)
               .flatMap(mapId -> MapProcessor.getInstance().findClosestSpawnPoint(mapId, event.x(), event.y()))
               .ifPresent(spawnPoint -> CharacterProcessor.getInstance().updateSpawnPoint(event.characterId(), spawnPoint));
      }
   }

   @Override
   public Class<CharacterMovementEvent> getEventClass() {
      return CharacterMovementEvent.class;
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
      return System.getenv(EventConstants.TOPIC_CHARACTER_MOVEMENT);
   }
}
