package com.atlas.cos.event.consumer;

import com.atlas.cos.command.AssignApCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class AssignApConsumer extends AbstractEventConsumer<AssignApCommand> {
   @Override
   public void handle(Long key, AssignApCommand command) {
      CharacterProcessor
            .getById(command.characterId())
            .ifPresent(character -> {
               switch (command.type()) {
                  case STRENGTH -> CharacterProcessor.assignStrDexIntLuk(character, 1, 0, 0, 0);
                  case DEXTERITY -> CharacterProcessor.assignStrDexIntLuk(character, 0, 1, 0, 0);
                  case INTELLIGENCE -> CharacterProcessor.assignStrDexIntLuk(character, 0, 0, 1, 0);
                  case LUCK -> CharacterProcessor.assignStrDexIntLuk(character, 0, 0, 0, 1);
                  case HP -> CharacterProcessor.assignHpMp(character, 1, 0);
                  case MP -> CharacterProcessor.assignHpMp(character, 0, 1);
               }
            });
   }

   @Override
   public Class<AssignApCommand> getEventClass() {
      return AssignApCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_ASSIGN_AP_COMMAND);
   }
}
