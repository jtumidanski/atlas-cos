package com.atlas.cos.event.consumer;

import com.atlas.cos.command.ChangeMapCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.processor.CharacterProcessor;

public class ChangeMapCommandConsumer extends AbstractEventConsumer<ChangeMapCommand> {
   @Override
   public void handle(Long key, ChangeMapCommand changeMapCommand) {
      CharacterProcessor.updateMap(changeMapCommand.worldId(), changeMapCommand.channelId(),
            changeMapCommand.characterId(), changeMapCommand.mapId(), changeMapCommand.portalId());
   }

   @Override
   public Class<ChangeMapCommand> getEventClass() {
      return ChangeMapCommand.class;
   }

   @Override
   public String getTopic() {
      return System.getenv(EventConstants.TOPIC_CHANGE_MAP_COMMAND);
   }
}
