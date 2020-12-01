package com.atlas.cos.event.consumer;

import com.atlas.cos.command.ChangeMapCommand;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.kafka.consumer.ConsumerRecordHandler;

public class ChangeMapCommandConsumer implements ConsumerRecordHandler<Long, ChangeMapCommand> {
   @Override
   public void handle(Long key, ChangeMapCommand changeMapCommand) {
      CharacterProcessor.getInstance().updateMap(changeMapCommand.worldId(), changeMapCommand.channelId(),
            changeMapCommand.characterId(), changeMapCommand.mapId(), changeMapCommand.portalId());
   }
}
