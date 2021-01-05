package com.atlas.cos.event.consumer;

import com.atlas.cos.command.GainMesoCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class GainMesoConsumer extends AbstractEventConsumer<GainMesoCommand> {
   @Override
   public void handle(Long key, GainMesoCommand command) {
      CharacterProcessor.gainMeso(command.characterId(), command.gain());
   }

   @Override
   public Class<GainMesoCommand> getEventClass() {
      return GainMesoCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_GAIN_MESO);
   }
}
