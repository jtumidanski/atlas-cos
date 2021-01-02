package com.atlas.cos.event.consumer;

import com.atlas.cos.command.AssignSpCommand;
import com.atlas.cos.constant.EventConstants;
import com.atlas.cos.database.administrator.SkillAdministrator;
import com.atlas.cos.database.provider.SkillProvider;
import com.atlas.cos.event.StatUpdateType;
import com.atlas.cos.event.producer.CharacterEnableActionsProducer;
import com.atlas.cos.event.producer.CharacterSkillUpdateProducer;
import com.atlas.cos.event.producer.CharacterStatUpdateProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.SkillData;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.SkillInformationProcessor;
import com.atlas.cos.processor.SkillProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;
import com.atlas.cos.util.SkillUtil;
import database.Connection;

import java.util.Collections;

public class AssignSpConsumer extends AbstractEventConsumer<AssignSpCommand> {
   @Override
   public void handle(Long key, AssignSpCommand command) {

      CharacterProcessor
            .getById(command.characterId())
            .ifPresent(character -> {
               SkillProcessor.getSkill(command.characterId(), command.skillId())
                     .ifPresent(skillData -> {
                        int skillBookId = SkillUtil.getSkillBook(command.skillId() / 10000);
                        int remainingSp = character.sps()[skillBookId];
                        int currentLevel = skillData.skillLevel();

                        boolean beginnerSkill = false;
                        if (command.skillId() % 10000000 > 999 && command.skillId() % 10000000 < 1003) {
                           int total = 0;
                           for (int i = 0; i < 3; i++) {
                              total += SkillProcessor.getSkill(character.id(), command.skillId())
                                    .map(SkillData::skillLevel)
                                    .orElse(0);
                           }
                           remainingSp = Math.min((character.level() - 1), 6) - total;
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
                              CharacterEnableActionsProducer.enableActions(character.id());
                           }

                           //TODO special handling for aran full swing and over swing.
                           changeSkillLevel(character.id(), skillData.skillId(), (byte) (currentLevel + 1), masterLevel,
                                 skillData.expiration());
                        }
                     });
            });
   }

   protected void adjustSp(int characterId, int amount, int skillBookId) {
      int newValue = CharacterProcessor.getById(characterId)
            .map(CharacterData::sps)
            .map(sps -> sps[skillBookId])
            .map(sp -> sp + amount)
            .map(sp -> Math.max(0, sp))
            .orElse(0);
      CharacterProcessor.updateSp(characterId, newValue, skillBookId);
      CharacterStatUpdateProducer.statsUpdated(characterId, Collections.singleton(StatUpdateType.AVAILABLE_SP));
   }

   protected void changeSkillLevel(int characterId, int skillId, byte level, int masterLevel, long expiration) {
      SkillProcessor.updateSkill(characterId, skillId, level, masterLevel, expiration);
      CharacterSkillUpdateProducer.updateSkill(characterId, skillId, level, masterLevel, expiration);
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
