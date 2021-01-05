package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.command.GainMesoCommand;
import com.atlas.cos.constant.CommandConstants;

public final class GainMesoProducer {
   private GainMesoProducer() {
   }

   public static void command(int characterId, int meso) {
      EventProducerRegistry.getInstance()
            .send(GainMesoCommand.class, CommandConstants.TOPIC_GAIN_MESO, characterId,
                  new GainMesoCommand(characterId, meso));
   }
}
