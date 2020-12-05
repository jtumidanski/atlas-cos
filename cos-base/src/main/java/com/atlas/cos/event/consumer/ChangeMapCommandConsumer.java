package com.atlas.cos.event.consumer;

import com.atlas.cos.command.ChangeMapCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class ChangeMapCommandConsumer implements SimpleEventHandler<ChangeMapCommand> {
   @Override
   public void handle(Long key, ChangeMapCommand changeMapCommand) {
      CharacterProcessor.getInstance().updateMap(changeMapCommand.worldId(), changeMapCommand.channelId(),
            changeMapCommand.characterId(), changeMapCommand.mapId(), changeMapCommand.portalId());
   }

   @Override
   public Class<ChangeMapCommand> getEventClass() {
      return ChangeMapCommand.class;
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
      return System.getenv(EventConstants.TOPIC_CHANGE_MAP_COMMAND);
   }
}
