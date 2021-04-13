package com.atlas.cos.event.consumer;

import com.atlas.cos.command.ChangeMapCommand;
import com.atlas.cos.constant.CommandConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class ChangeMapCommandConsumer extends AbstractEventConsumer<ChangeMapCommand> {
   @Override
   public void handle(Long key, ChangeMapCommand command) {
      CharacterProcessor
            .updateMap(command.worldId(), command.channelId(), command.characterId(), command.mapId(), command.portalId());
   }

   @Override
   public Class<ChangeMapCommand> getEventClass() {
      return ChangeMapCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(CommandConstants.TOPIC_CHANGE_MAP_COMMAND);
   }
}
