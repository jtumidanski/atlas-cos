package com.atlas.cos.event.consumer;

import com.atlas.cos.command.AssignSpCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.SkillInformationProcessor;
import com.atlas.cos.processor.SkillProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;
import com.atlas.cos.util.SkillUtil;

public class AssignSpConsumer extends AbstractEventConsumer<AssignSpCommand> {
   @Override
   public void handle(Long key, AssignSpCommand command) {

      CharacterProcessor
            .getById(command.characterId())
            .ifPresent(character -> {
               int skillBookId = SkillUtil.getSkillBook(command.skillId() / 10000);
               int remainingSp = character.sps()[skillBookId];

               SkillProcessor.getSkill(command.characterId(), command.skillId())
                     .ifPresent(skillData -> {
                        int currentLevel = skillData.skillLevel();

                        boolean beginnerSkill = false;
                        if (command.skillId() % 10000000 > 999 && command.skillId() % 10000000 < 1003) {
                           beginnerSkill = true;
                        }

                        //TODO retrieve maxLevel from SkillInformation
                        int skillMaxLevel = SkillInformationProcessor.getSkillInformation(command.skillId())
                              .join().effects().size();

                        int masterLevel = skillData.masterLevel();
                        int maxLevel = (SkillUtil.isFourthJob(character.jobId(), command.skillId()) ? masterLevel : skillMaxLevel);
                        if (remainingSp > 0 && currentLevel + 1 <= maxLevel) {
                           if (!beginnerSkill) {
                              adjustSp(character.id(), -1, skillBookId);
                           } else {
                              // enable actions.
                           }

                           //TODO special handling for aran full swing and over swing.
                           long expiration = 0;
                           changeSkillLevel(character.id(), skillData.skillId(), (byte) (currentLevel + 1), masterLevel,
                                 expiration);
                        }
                     });
            });
   }

   protected void adjustSp(int characterId, int amount, int skillBookId) {
   }

   protected void changeSkillLevel(int characterId, int skillId, byte level, int masterLevel, long expiration) {
   }

   @Override
   public Class<AssignSpCommand> getEventClass() {
      return AssignSpCommand.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(EventConstants.TOPIC_ASSIGN_SP_COMMAND);
   }
}
