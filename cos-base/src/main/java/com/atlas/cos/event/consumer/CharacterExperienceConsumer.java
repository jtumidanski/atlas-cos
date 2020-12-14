package com.atlas.cos.event.consumer;

import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.event.CharacterExperienceEvent;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.util.ExpTable;
import com.atlas.kafka.consumer.SimpleEventHandler;

public class CharacterExperienceConsumer implements SimpleEventHandler<CharacterExperienceEvent> {
   @Override
   public void handle(Long key, CharacterExperienceEvent event) {
      CharacterProcessor.getById(event.characterId())
            .ifPresent(character -> gainExperience(event.personalGain(), event.partyGain(), character));
   }

   private void gainExperience(long personalGain, long partyGain, CharacterData character) {
      CharacterData runningCharacter = character;

      long total = Math.max(personalGain + partyGain, -runningCharacter.experience());

      if (runningCharacter.level() < runningCharacter.maxClassLevel()) {
         long leftover = 0;
         long nextExp = runningCharacter.experience() + total;

         if (nextExp > (long) Integer.MAX_VALUE) {
            total = Integer.MAX_VALUE - runningCharacter.experience();
            leftover = nextExp - Integer.MAX_VALUE;
         }

         runningCharacter = CharacterProcessor.increaseExperience(runningCharacter.id(), (int) total);

         while (runningCharacter.experience() >= ExpTable.getExpNeededForLevel(runningCharacter.level())) {
            runningCharacter = CharacterProcessor.increaseLevel(runningCharacter.id(), true);

            if (runningCharacter.level() == runningCharacter.maxClassLevel()) {
               runningCharacter = CharacterProcessor.setExperience(runningCharacter.id(), 0);
               break;
            }
         }

         if (leftover > 0) {
            gainExperience(leftover, partyGain, runningCharacter);
         }
      }
   }

   @Override
   public Class<CharacterExperienceEvent> getEventClass() {
      return CharacterExperienceEvent.class;
   }

   @Override
   public String getConsumerId() {
      return "Character Service";
   }

   @Override
   public String getBootstrapServers() {
      return System.getenv("BOOTSTRAP_SERVERS");
   }

   @Override
   public String getTopic() {
      return System.getenv(EventConstants.TOPIC_CHARACTER_EXPERIENCE_EVENT);
   }
}
