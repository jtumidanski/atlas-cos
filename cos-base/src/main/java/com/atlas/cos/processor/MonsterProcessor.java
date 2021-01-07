package com.atlas.cos.processor;

import java.util.HashMap;
import java.util.LinkedList;
import java.util.List;
import java.util.Map;
import java.util.Optional;

import com.atlas.cos.ConfigurationRegistry;
import com.atlas.cos.event.producer.CharacterExperienceGainProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.Monster;
import com.atlas.cos.model.MonsterExperienceDistribution;
import com.atlas.mis.attribute.MonsterDataAttributes;
import com.atlas.mis.constant.RestConstants;
import com.atlas.morg.rest.event.DamageEntry;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public final class MonsterProcessor {
   private MonsterProcessor() {
   }

   public static Optional<Monster> getMonster(int monsterId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("monsters", monsterId)
            .getRestClient(MonsterDataAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(ModelFactory::createMonster);
   }

   public static void distributeExperience(int worldId, int channelId, int mapId, Monster monster,
                                           List<DamageEntry> damageEntries) {
      MonsterExperienceDistribution distribution = MonsterProcessor.produceDistribution(mapId, monster,
            damageEntries);
      for (Map.Entry<Integer, Long> chrParticipation : distribution.solo().entrySet()) {
         float exp = chrParticipation.getValue() * distribution.experiencePerDamage();
         int characterId = chrParticipation.getKey();

         int level = CharacterProcessor.getById(characterId)
               .map(CharacterData::level)
               .orElse(1);

         distributePlayerExperience(characterId, level, exp, 0.0f, level, true,
               isWhiteExpGain(characterId, distribution.personalRatio(), distribution.standardDeviationRatio()), false);
      }

      //                              Set<Integer> underLeveled = new HashSet<>();
      //                              for (Map<Integer, Long> partyParticipation : distribution.party().values()) {
      //                                 distributePartyExperience(partyParticipation, experiencePerDamage, underLeveled, personalRatio, standardDeviationRatio);
      //                              }
      //
      //                              for (Integer mc : underLeveled) {
      //                                 mc.showUnderLeveledInfo(this);
      //                              }
   }

   protected static void distributePlayerExperience(int characterId, int level,
                                                    float experience, float partyBonusMod, int totalPartyLevel,
                                                    boolean highestPartyDamage, boolean whiteExpGain, boolean hasPartySharers) {
      float expSplitCommonMod = ConfigurationRegistry.getInstance().getConfiguration().expSplitCommonMod;
      float playerExp = (expSplitCommonMod * level) / totalPartyLevel;
      if (highestPartyDamage) {
         playerExp += ConfigurationRegistry.getInstance().getConfiguration().expSplitMvpMod;
      }

      playerExp *= experience;
      float bonusExp = partyBonusMod * playerExp;

      giveExpToCharacter(characterId, playerExp, bonusExp, whiteExpGain, hasPartySharers);
      //      giveFamilyRep(chr.getFamilyEntry());
   }

   protected static void giveExpToCharacter(int characterId, Float personalExp,
                                            Float partyExp, boolean white, boolean hasPartySharers) {
      //      if (attacker.isAlive()) {
      //         if (personalExp != null) {
      //            personalExp *= getStatusExpMultiplier(attacker, hasPartySharers);
      //            personalExp *= attacker.getExpRate();
      //         } else {
      //            personalExp = 0.0f;
      //         }

      //         Integer expBonus = attacker.getBuffedValue(MapleBuffStat.EXP_INCREASE);
      //         if (expBonus != null) {
      //            personalExp += expBonus;
      //         }

      int correctedPersonal = expValueToInteger(personalExp); // assuming no negative xp here

      //         if (partyExp != null) {
      //            partyExp *= getStatusExpMultiplier(attacker, hasPartySharers);
      //            partyExp *= attacker.getExpRate();
      //            partyExp *= YamlConfig.config.server.PARTY_BONUS_EXP_RATE;
      //         } else {
      //            partyExp = 0.0f;
      //         }

      int correctedParty = expValueToInteger(partyExp);

      CharacterExperienceGainProducer.gainExperience(characterId, correctedPersonal, correctedParty,
            true, false, white);
      //         attacker.increaseEquipExp(_personalExp);
      //      }
   }

   protected static int expValueToInteger(double exp) {
      if (exp > Integer.MAX_VALUE) {
         exp = Integer.MAX_VALUE;
      } else if (exp < Integer.MIN_VALUE) {
         exp = Integer.MIN_VALUE;
      }

      return (int) Math.round(exp);
   }

   protected static boolean isWhiteExpGain(int characterId, Map<Integer, Float> personalRatio, double standardDeviationRatio) {
      Float pr = personalRatio.get(characterId);
      if (pr == null) {
         return false;
      }

      return pr >= standardDeviationRatio;
   }

   public static MonsterExperienceDistribution produceDistribution(int mapId, Monster monster, List<DamageEntry> damageEntries) {
      int totalEntries = 0;

      //TODO incorporate party distribution.
      Map<Integer, Map<Integer, Long>> partyDistribution = new HashMap<>();
      Map<Integer, Long> soloDistribution = new HashMap<>();

      for (DamageEntry damageEntry : damageEntries) {
         if (CharacterProcessor.inMap(damageEntry.character(), mapId)) {
            soloDistribution.put(damageEntry.character(), damageEntry.damage());
         }
         totalEntries += 1;
      }

      //TODO account for healing.
      long totalDamage = monster.hp();
      float experiencePerDamage = ((float) monster.experience()) / totalDamage;

      Map<Integer, Float> personalRatio = new HashMap<>();
      List<Float> entryExperienceRatio = new LinkedList<>();
      for (Map.Entry<Integer, Long> e : soloDistribution.entrySet()) {
         float ratio = ((float) e.getValue()) / totalDamage;

         personalRatio.put(e.getKey(), ratio);
         entryExperienceRatio.add(ratio);
      }

      for (Map<Integer, Long> m : partyDistribution.values()) {
         float ratio = 0.0f;
         for (Map.Entry<Integer, Long> e : m.entrySet()) {
            float chrRatio = ((float) e.getValue()) / totalDamage;

            personalRatio.put(e.getKey(), chrRatio);
            ratio += chrRatio;
         }
         entryExperienceRatio.add(ratio);
      }

      double standardDeviationRatio = calcExperienceStandDevThreshold(entryExperienceRatio, totalEntries);
      return new MonsterExperienceDistribution(soloDistribution, partyDistribution, personalRatio, experiencePerDamage,
            standardDeviationRatio);
   }

   protected static double calcExperienceStandDevThreshold(List<Float> entryExpRatio, int totalEntries) {
      float avgExpReward = entryExpRatio.stream()
            .reduce(0.0f, Float::sum) / totalEntries;
      float varExpReward = entryExpRatio.stream()
            .reduce(0.0f, (result, next) -> (float) (result + Math.pow(next - avgExpReward, 2))) / entryExpRatio.size();
      return avgExpReward + Math.sqrt(varExpReward);
   }
}
