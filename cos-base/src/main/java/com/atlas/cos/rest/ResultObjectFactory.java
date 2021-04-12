package com.atlas.cos.rest;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.CharacterTemporalData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.EquipmentStatistics;
import com.atlas.cos.model.Inventory;
import com.atlas.cos.model.InventoryItem;
import com.atlas.cos.model.InventoryItemType;
import com.atlas.cos.model.ItemData;
import com.atlas.cos.model.SavedLocationData;
import com.atlas.cos.model.SkillData;
import com.atlas.cos.processor.EquipmentProcessor;
import com.atlas.cos.processor.ItemProcessor;
import com.atlas.cos.rest.attribute.CharacterAttributes;
import com.atlas.cos.rest.attribute.EquipmentAttributes;
import com.atlas.cos.rest.attribute.EquipmentStatisticsAttributes;
import com.atlas.cos.rest.attribute.InventoryAttributes;
import com.atlas.cos.rest.attribute.ItemAttributes;
import com.atlas.cos.rest.attribute.LocationAttributes;
import com.atlas.cos.rest.attribute.LocationType;
import com.atlas.cos.rest.attribute.SkillAttributes;
import com.atlas.cos.rest.builder.CharacterAttributesBuilder;
import com.atlas.cos.rest.builder.EquipmentAttributesBuilder;
import com.atlas.cos.rest.builder.EquipmentStatisticsAttributesBuilder;
import com.atlas.cos.rest.builder.InventoryAttributesBuilder;
import com.atlas.cos.rest.builder.ItemAttributesBuilder;
import com.atlas.cos.rest.builder.LocationAttributesBuilder;
import com.atlas.cos.rest.builder.SkillAttributesBuilder;

import builder.LinkDataBuilder;
import builder.RelationshipBuilder;
import builder.RelationshipsBuilder;
import builder.ResultObjectBuilder;
import rest.AttributeResult;

public class ResultObjectFactory {
   public static ResultObjectBuilder create(CharacterData data) {

      CharacterTemporalData temporalData = CharacterTemporalRegistry.getInstance()
            .getTemporalData(data.id());

      return new ResultObjectBuilder(CharacterAttributes.class, data.id())
            .setAttribute(new CharacterAttributesBuilder()
                  .setAccountId(data.accountId())
                  .setWorldId(data.worldId())
                  .setName(data.name())
                  .setLevel(data.level())
                  .setExperience(data.experience())
                  .setGachaponExperience(data.gachaponExperience())
                  .setStrength(data.strength())
                  .setDexterity(data.dexterity())
                  .setLuck(data.luck())
                  .setIntelligence(data.intelligence())
                  .setHp(data.hp())
                  .setMp(data.mp())
                  .setMaxHp(data.maxHp())
                  .setMaxMp(data.maxMp())
                  .setMeso(data.meso())
                  .setHpMpUsed(data.hpMpUsed())
                  .setJobId(data.jobId())
                  .setSkinColor(data.skinColor())
                  .setGender(data.gender())
                  .setFame(data.fame())
                  .setHair(data.hair())
                  .setFace(data.face())
                  .setAp(data.ap())
                  .setSp(data.sp())
                  .setMapId(data.mapId())
                  .setSpawnPoint(data.spawnPoint())
                  .setGm(data.gm())
                  .setX(temporalData.x())
                  .setY(temporalData.y())
                  .setStance(temporalData.stance())
            );
   }

   public static ResultObjectBuilder create(EquipmentData data) {
      return new ResultObjectBuilder(EquipmentAttributes.class, data.id())
            .setAttribute(new EquipmentAttributesBuilder()
                  .setEquipmentId(data.equipmentId())
                  .setSlot(data.slot())
            )
            .setRelationships(new RelationshipsBuilder()
                  .addRelationship("equipmentStatistics", new RelationshipBuilder()
                        .addData(new LinkDataBuilder()
                              .setId(String.valueOf(data.equipmentId()))
                              .setType(EquipmentStatisticsAttributes.class)
                        )
                  )
            );
   }

   protected static ResultObjectBuilder create(ItemData data) {
      return new ResultObjectBuilder(ItemAttributes.class, data.id())
            .setAttribute(new ItemAttributesBuilder()
                  .setItemId(data.itemId())
                  .setQuantity(data.quantity())
                  .setSlot(data.slot())
            );
   }

   public static ResultObjectBuilder create(SavedLocationData data) {
      return new ResultObjectBuilder(LocationAttributes.class, data.id())
            .setAttribute(new LocationAttributesBuilder()
                  .setType(LocationType.valueOf(data.type()))
                  .setMapId(data.mapId())
                  .setPortalId(data.portalId())
            );
   }

   protected static Class<? extends AttributeResult> getInventoryItemType(InventoryItemType type) {
      return type.equals(InventoryItemType.EQUIPMENT) ? EquipmentAttributes.class : ItemAttributes.class;
   }

   public static ResultObjectBuilder create(Inventory inventory) {
      RelationshipBuilder relationshipBuilder = new RelationshipBuilder();
      inventory.items().stream()
            .map(inventoryItem -> new LinkDataBuilder()
                  .setId(String.valueOf(inventoryItem.id()))
                  .setType(getInventoryItemType(inventoryItem.type()))
            )
            .forEach(relationshipBuilder::addData);

      return new ResultObjectBuilder(InventoryAttributes.class, (int) inventory.id())
            .setAttribute(new InventoryAttributesBuilder()
                  .setType(inventory.type())
                  .setCapacity(inventory.capacity())
            )
            .setRelationships(new RelationshipsBuilder()
                  .addRelationship("inventoryItems", relationshipBuilder)
            );
   }

   public static ResultObjectBuilder create(InventoryItem inventoryItem) {
      if (inventoryItem.type().equals(InventoryItemType.EQUIPMENT)) {
         return EquipmentProcessor.getEquipmentById(inventoryItem.id())
               .map(ResultObjectFactory::create)
               .orElseThrow();
      } else {
         return ItemProcessor.getItemById(inventoryItem.id())
               .map(ResultObjectFactory::create)
               .orElseThrow();
      }
   }

   public static ResultObjectBuilder create(EquipmentStatistics statistics) {
      return new ResultObjectBuilder(EquipmentStatisticsAttributes.class, statistics.id())
            .setAttribute(new EquipmentStatisticsAttributesBuilder()
                  .setItemId(statistics.itemId())
                  .setStrength(statistics.strength())
                  .setDexterity(statistics.dexterity())
                  .setIntelligence(statistics.intelligence())
                  .setLuck(statistics.luck())
                  .setHp(statistics.hp())
                  .setMp(statistics.mp())
                  .setWeaponAttack(statistics.weaponAttack())
                  .setMagicAttack(statistics.magicAttack())
                  .setWeaponDefense(statistics.weaponDefense())
                  .setMagicDefense(statistics.magicDefense())
                  .setAccuracy(statistics.accuracy())
                  .setAvoidability(statistics.avoidability())
                  .setHands(statistics.hands())
                  .setSpeed(statistics.speed())
                  .setJump(statistics.jump())
                  .setSlots(statistics.slots())
            );
   }

   public static ResultObjectBuilder create(SkillData data) {
      return new ResultObjectBuilder(SkillAttributes.class, data.skillId())
            .setAttribute(new SkillAttributesBuilder()
                  .setLevel(data.skillLevel())
                  .setMasterLevel(data.masterLevel())
                  .setExpiration(data.expiration())
            );
   }
}
