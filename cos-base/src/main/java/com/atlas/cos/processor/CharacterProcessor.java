package com.atlas.cos.processor;

import java.util.Optional;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.builder.CharacterBuilder;
import com.atlas.cos.database.administrator.CharacterAdministrator;
import com.atlas.cos.database.provider.CharacterProvider;
import com.atlas.cos.event.producer.CharacterCreatedProducer;
import com.atlas.cos.event.producer.MapChangedProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.util.ExpTable;

import database.Connection;

public final class CharacterProcessor {
   private CharacterProcessor() {
   }

   public static Optional<CharacterData> getByName(String name) {
      return Connection.instance()
            .list(entityManager -> CharacterProvider.getForName(entityManager, name))
            .stream().findFirst();
   }

   public static Optional<CharacterData> getById(int id) {
      return Connection.instance()
            .element(entityManager -> CharacterProvider.getById(entityManager, id));
   }

   protected static Optional<CharacterData> create(CharacterBuilder builder) {
      CharacterData original = builder.build();

      Optional<CharacterData> result = Connection.instance().element(entityManager ->
            CharacterAdministrator.create(entityManager,
                  original.accountId(), original.worldId(), original.name(), original.level(),
                  original.strength(), original.dexterity(), original.luck(), original.intelligence(),
                  original.maxHp(), original.maxMp(), original.jobId(), original.gender(), original.hair(),
                  original.face(), original.mapId())
      );

      result.ifPresent(CharacterCreatedProducer::notifyCharacterCreated);
      return result;
   }

   public static Optional<CharacterData> createBeginner(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 10000);
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public static Optional<CharacterData> createNoblesse(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 130030000);
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public static Optional<CharacterData> createLegend(CharacterAttributes attributes) {
      CharacterBuilder builder = new CharacterBuilder(attributes, 1, 914000000);
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      return create(builder);
   }

   public static Optional<CharacterData> createBeginner(CharacterBuilder builder) {
      //giveItem(recipe, 4161001, 1, MapleInventoryType.ETC);
      builder.setMapId(10000);
      return create(builder);
   }

   public static Optional<CharacterData> createNoblesse(CharacterBuilder builder) {
      //giveItem(recipe, 4161047, 1, MapleInventoryType.ETC);
      builder.setMapId(130030000);
      return create(builder);
   }

   public static Optional<CharacterData> createLegend(CharacterBuilder builder) {
      //giveItem(recipe, 4161048, 1, MapleInventoryType.ETC);
      builder.setMapId(914000000);
      return create(builder);
   }

   public static void updateMap(int worldId, int channelId, int characterId, int mapId, int portalId) {
      Connection.instance().with(entityManager -> CharacterAdministrator.updateMap(entityManager, characterId, mapId));

      PortalProcessor.getMapPortalById(mapId, portalId)
            .ifPresent(portal -> CharacterTemporalRegistry.getInstance()
                  .updatePosition(characterId, portal.x(), portal.y()));

      MapChangedProducer.notifyChange(worldId, channelId, characterId, mapId, portalId);
   }

   public static void updateSpawnPoint(int characterId, int newSpawnPoint) {
      Connection.instance()
            .with(entityManager -> CharacterAdministrator.updateSpawnPoint(entityManager, characterId, newSpawnPoint));
   }

   public static int getWeaponAttack(CharacterData character) {
      int weaponAttack = 0;

      weaponAttack += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::weaponAttack)
            .sum();

      //TODO
      // apply Aran Combo
      // apply ThunderBreaker Marauder energy charge
      // apply Marksman Boost or Bowmaster Expert
      // apply weapon attack buffs
      // apply blessing
      // apply active projectile

      return weaponAttack;
   }

   public static int getStrength(CharacterData character) {
      int strength = character.strength();

      //TODO
      // apply Maple Warrior

      strength += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::strength)
            .sum();

      return strength;
   }

   public static int getDexterity(CharacterData character) {
      int dexterity = character.dexterity();

      //TODO
      // apply Maple Warrior

      dexterity += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::dexterity)
            .sum();

      return dexterity;
   }

   public static int getLuck(CharacterData character) {
      int luck = character.luck();

      //TODO
      // apply Maple Warrior

      luck += ItemProcessor.getEquipmentForCharacter(character.id()).stream()
            .mapToInt(EquipmentData::luck)
            .sum();

      return luck;
   }

   public static boolean inMap(int characterId, int mapId) {
      return getById(characterId)
            .map(CharacterData::mapId)
            .filter(id -> id == mapId)
            .isPresent();
   }

   public static CharacterData increaseExperience(int characterId, int amount) {
      return Connection.instance()
            .element(entityManager -> {
               CharacterAdministrator.increaseExperience(entityManager, characterId, amount);
               return CharacterProvider.getById(entityManager, characterId);
            }).orElseThrow();
   }

   public static CharacterData increaseLevel(int characterId, boolean takeExp) {
      //      boolean isBeginner = isBeginnerJob();
      //      if (YamlConfig.config.server.USE_AUTOASSIGN_STARTERS_AP && isBeginner && level < 11) {
      //         effLock.lock();
      //         statWriteLock.lock();
      //         try {
      //            gainAp(5, true);
      //
      //            int str = 0, dex = 0;
      //            if (level < 6) {
      //               str += 5;
      //            } else {
      //               str += 4;
      //               dex += 1;
      //            }
      //
      //            assignStrDexIntLuk(str, dex, 0, 0);
      //         } finally {
      //            statWriteLock.unlock();
      //            effLock.unlock();
      //         }
      //      } else {
      //         int remainingAp = 5;
      //
      //         if (isCygnus()) {
      //            if (level > 10) {
      //               if (level <= 17) {
      //                  remainingAp += 2;
      //               } else if (level < 77) {
      //                  remainingAp++;
      //               }
      //            }
      //         }
      //
      //         gainAp(remainingAp, true);
      //      }
      //
      //      levelUpHealthAndManaPoints(isBeginner);

      CharacterData characterData = getById(characterId).orElseThrow();

      if (takeExp) {
         int newExperience = characterData.experience() - ExpTable.getExpNeededForLevel(characterData.level());
         if (newExperience < 0) {
            newExperience = 0;
         }
         setExperience(characterId, newExperience);
      }

      if (characterData.level() + 1 >= characterData.maxClassLevel()) {
         setExperience(characterId, 0);
         Connection.instance().with(entityManager -> CharacterAdministrator.setLevel(entityManager, characterId,
               characterData.maxClassLevel()));

         //         if (level == maxClassLevel) {
         ////            if (!this.isGM()) {
         ////               if (YamlConfig.config.server.PLAYERNPC_AUTODEPLOY) {
         ////                  ThreadManager.getInstance()
         ////                        .newTask(() -> MaplePlayerNPC.spawnPlayerNPC(GameConstants.getHallOfFameMapId(job), MapleCharacter.this));
         ////               }
         ////
         ////               final String names = (getMedalText() + name);
         ////               MessageBroadcaster.getInstance().sendWorldServerNotice(world, ServerNoticeType.LIGHT_BLUE,
         ////                     I18nMessage.from("LEVEL_200").with(names, maxClassLevel, names));
         ////            }
         //         }

         //To prevent levels past the maximum
      } else {
         Connection.instance().with(entityManager -> CharacterAdministrator.setLevel(entityManager, characterId,
               characterData.level() + 1));
      }

      //      levelUpGainSp();
      //
      //      effLock.lock();
      //      statWriteLock.lock();
      //      try {
      //         recalculateLocalStats();
      //         changeHpMp(localMaxHp, localMaxMp, true);
      //
      //         List<Pair<MapleStat, Integer>> statIncreases = new ArrayList<>(10);
      //         statIncreases.add(new Pair<>(MapleStat.AVAILABLE_AP, remainingAp));
      //         statIncreases.add(new Pair<>(MapleStat.AVAILABLE_SP, remainingSp[GameConstants.getSkillBook(job.getId())]));
      //         statIncreases.add(new Pair<>(MapleStat.HP, hp));
      //         statIncreases.add(new Pair<>(MapleStat.MP, mp));
      //         statIncreases.add(new Pair<>(MapleStat.EXP, exp.get()));
      //         statIncreases.add(new Pair<>(MapleStat.LEVEL, level));
      //         statIncreases.add(new Pair<>(MapleStat.MAX_HP, clientMaxHp));
      //         statIncreases.add(new Pair<>(MapleStat.MAX_MP, clientMaxMp));
      //         statIncreases.add(new Pair<>(MapleStat.STR, str));
      //         statIncreases.add(new Pair<>(MapleStat.DEX, dex));
      //
      //         PacketCreator.announce(client, new UpdatePlayerStats(statIncreases, true, this));
      //      } finally {
      //         statWriteLock.unlock();
      //         effLock.unlock();
      //      }

      //      MasterBroadcaster.getInstance().sendToAllInMap(getMap(), new ShowForeignEffect(getId(), 0), false, this);
      //      setMPC(new MaplePartyCharacter(this));
      //      silentPartyUpdate();
      //
      //      if (this.guildId > 0) {
      //         getGuild().ifPresent(
      //               guild -> MasterBroadcaster.getInstance().sendToGuild(guild, new NotifyLevelUp(2, level, name), false, this.getId()));
      //      }

      //      if (level % 20 == 0) {
      //         if (YamlConfig.config.server.USE_ADD_SLOTS_BY_LEVEL) {
      //            if (!isGM()) {
      //               for (byte i = 1; i < 5; i++) {
      //                  gainSlots(i, 4, true);
      //               }
      //               MessageBroadcaster.getInstance().yellowMessage(this, I18nMessage.from("INVENTORY_EXPANSION_ON_LEVEL").with(level));
      //            }
      //         }
      //         if (YamlConfig.config.server.USE_ADD_RATES_BY_LEVEL) { //For the rate upgrade
      //            revertLastPlayerRates();
      //            setPlayerRates();
      //            MessageBroadcaster.getInstance().yellowMessage(this, I18nMessage.from("USER_INCREASE_RATES_ON_LEVEL").with(level));
      //         }
      //      }

      //      if (YamlConfig.config.server.USE_PERFECT_PITCH && level >= 30) {
      //         //milestones?
      //         if (MapleInventoryManipulator.checkSpace(client, 4310000, (short) 1, "")) {
      //            MapleInventoryManipulator.addById(client, 4310000, (short) 1, "", -1);
      //         }
      //      } else if (level == 10) {
      //         Runnable r = () -> {
      //            if (leaveParty()) {
      //               showHint("You have reached #blevel 10#k, therefore you must leave your #rstarter party#k.");
      //            }
      //         };
      //
      //         ThreadManager.getInstance().newTask(r);
      //      }
      //
      //      levelUpMessages();
      //      guildUpdate();
      //
      //      MapleFamilyProcessor.getInstance().giveReputationToCharactersSenior(getFamilyEntry(), level, getName());
      return getById(characterId).orElseThrow();
   }

   public static CharacterData setExperience(int characterId, int experience) {
      return Connection.instance()
            .element(entityManager -> {
               CharacterAdministrator.setExperience(entityManager, characterId, experience);
               return CharacterProvider.getById(entityManager, characterId);
            }).orElseThrow();
   }
}
