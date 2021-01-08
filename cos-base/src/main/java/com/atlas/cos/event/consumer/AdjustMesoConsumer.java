package com.atlas.cos.event.consumer;

import com.atlas.cos.command.AdjustMesoCommand;
import com.atlas.cos.constant.CommandConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class AdjustMesoConsumer extends AbstractEventConsumer<AdjustMesoCommand> {
   @Override
   public void handle(Long key, AdjustMesoCommand command) {
      CharacterProcessor.adjustMeso(command.characterId(), command.amount(), command.show());
   }

   @Override
   public Class<AdjustMesoCommand> getEventClass() {
      return AdjustMesoCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(CommandConstants.TOPIC_ADJUST_MESO);
   }
}
