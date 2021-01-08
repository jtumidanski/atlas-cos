package com.atlas.cos.event.consumer;

import com.atlas.cos.command.AdjustManaCommand;
import com.atlas.cos.constant.CommandConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class AdjustManaConsumer extends AbstractEventConsumer<AdjustManaCommand> {
   @Override
   public void handle(Long key, AdjustManaCommand command) {
      CharacterProcessor.adjustMana(command.characterId(), command.amount());
   }

   @Override
   public Class<AdjustManaCommand> getEventClass() {
      return AdjustManaCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(CommandConstants.TOPIC_ADJUST_MANA);
   }
}
