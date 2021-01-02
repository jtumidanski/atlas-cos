package com.atlas.cos.rest.processor;

import java.util.Optional;
import java.util.function.Function;
import javax.ws.rs.core.Response;

import com.app.rest.util.stream.Mappers;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.InventoryType;
import com.atlas.cos.model.Job;
import com.atlas.cos.model.skills.Beginner;
import com.atlas.cos.model.skills.Legend;
import com.atlas.cos.model.skills.Noblesse;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.EquipmentProcessor;
import com.atlas.cos.processor.ItemProcessor;
import com.atlas.cos.processor.JobProcessor;
import com.atlas.cos.processor.SkillProcessor;
import com.atlas.cos.rest.ResultObjectFactory;

import builder.ResultBuilder;

public final class CharacterSeedRequestProcessor {
   private CharacterSeedRequestProcessor() {
   }

   protected static Optional<Function<CharacterBuilder, Optional<CharacterData>>> getCreator(Job job) {
      if (Job.BEGINNER.equals(job)) {
         return Optional.of(CharacterProcessor::createBeginner);
      } else if (Job.NOBLESSE.equals(job)) {
         return Optional.of(CharacterProcessor::createNoblesse);
      } else if (Job.LEGEND.equals(job)) {
         return Optional.of(CharacterProcessor::createLegend);
      }
      return Optional.empty();
   }

   protected static CharacterData addEquippedItems(CharacterData characterData, int top, int bottom, int shoes, int weapon) {
      EquipmentProcessor.createAndEquip(characterData.id(), top, bottom, shoes, weapon);
      return characterData;
   }

   protected static CharacterData addOtherItems(CharacterData characterData) {
      Job.getById(characterData.jobId()).ifPresent(job -> addJobItems(characterData.id(), job));
      return characterData;
   }

   protected static void addJobItems(int characterId, Job job) {
      if (Job.BEGINNER.equals(job)) {
         ItemProcessor.createItemForCharacter(characterId, InventoryType.ETC, 4161001, 1);
      } else if (Job.NOBLESSE.equals(job)) {
         ItemProcessor.createItemForCharacter(characterId, InventoryType.ETC, 4161047, 1);
      } else if (Job.LEGEND.equals(job)) {
         ItemProcessor.createItemForCharacter(characterId, InventoryType.ETC, 4161048, 1);
      }
   }

   public static ResultBuilder create(int accountId, int worldId, String name, int jobIndex, int face, int hair, int hairColor,
                                      int skin, byte gender, int top, int bottom, int shoes, int weapon) {
      return JobProcessor.getJobFromIndex(jobIndex)
            .flatMap(CharacterSeedRequestProcessor::getCreator)
            .flatMap(creator -> creator.apply(new CharacterBuilder(accountId, worldId, name, skin, gender, hair + hairColor, face)))
            .map(character -> addEquippedItems(character, top, bottom, shoes, weapon))
            .map(CharacterSeedRequestProcessor::addOtherItems)
            .map(CharacterSeedRequestProcessor::addSkills)
            .map(ResultObjectFactory::create)
            .map(Mappers::singleCreatedResult)
            .orElseGet(ResultBuilder::forbidden);
   }

   protected static CharacterData addSkills(CharacterData characterData) {
      Job.getById(characterData.jobId()).ifPresent(job -> addSkills(characterData.id(), job));
      return characterData;
   }

   protected static void addSkills(int characterId, Job job) {
      if (job.isA(Job.BEGINNER)) {
         SkillProcessor.awardSkill(characterId, Beginner.RECOVERY);
         SkillProcessor.awardSkill(characterId, Beginner.NIMBLE_FEET);
         SkillProcessor.awardSkill(characterId, Beginner.THREE_SNAILS);
      } else if (job.isA(Job.LEGEND)) {

      } else if (job.isA(Job.NOBLESSE)) {
         SkillProcessor.awardSkill(characterId, Noblesse.RECOVERY);
         SkillProcessor.awardSkill(characterId, Noblesse.NIMBLE_FEET);
         SkillProcessor.awardSkill(characterId, Noblesse.THREE_SNAILS);
      }
   }
}
