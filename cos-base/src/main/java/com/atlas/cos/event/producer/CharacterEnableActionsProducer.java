package com.atlas.cos.event.producer;

import com.atlas.cos.EventProducerRegistry;
import com.atlas.csrv.command.EnableActionsCommand;
import com.atlas.csrv.constant.CommandConstants;

public final class CharacterEnableActionsProducer {
   public static void enableActions(int characterId) {
      EventProducerRegistry.getInstance().send(CommandConstants.TOPIC_ENABLE_ACTIONS, characterId,
            new EnableActionsCommand(characterId));
   }
}
