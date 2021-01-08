package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.cos.command.AdjustMesoCommand;
import com.atlas.cos.constant.CommandConstants;

public final class AdjustMesoProducer {
   private AdjustMesoProducer() {
   }

   public static void command(int characterId, int amount) {
      EventProducerRegistry.getInstance().send(CommandConstants.TOPIC_ADJUST_MESO, characterId,
            new AdjustMesoCommand(characterId, amount, true));
   }
}
