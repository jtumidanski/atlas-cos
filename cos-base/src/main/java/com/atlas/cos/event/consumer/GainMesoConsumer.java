package com.atlas.cos.event.consumer;

import com.atlas.cos.command.GainMesoCommand;
import com.atlas.cos.constant.CommandConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

@Deprecated
public class GainMesoConsumer extends AbstractEventConsumer<GainMesoCommand> {
   @Override
   public void handle(Long key, GainMesoCommand command) {
      CharacterProcessor.adjustMeso(command.characterId(), command.gain(), true);
   }

   @Override
   public Class<GainMesoCommand> getEventClass() {
      return GainMesoCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(CommandConstants.TOPIC_GAIN_MESO);
   }
}