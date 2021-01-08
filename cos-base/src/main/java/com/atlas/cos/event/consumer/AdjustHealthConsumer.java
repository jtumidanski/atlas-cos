package com.atlas.cos.event.consumer;

import com.atlas.cos.command.AdjustHealthCommand;
import com.atlas.cos.constant.CommandConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class AdjustHealthConsumer extends AbstractEventConsumer<AdjustHealthCommand> {
   @Override
   public void handle(Long key, AdjustHealthCommand command) {
      CharacterProcessor.adjustHealth(command.characterId(), command.amount());
   }

   @Override
   public Class<AdjustHealthCommand> getEventClass() {
      return AdjustHealthCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(CommandConstants.TOPIC_ADJUST_HEALTH);
   }
}
