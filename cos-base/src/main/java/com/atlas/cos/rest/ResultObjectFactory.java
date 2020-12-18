package com.atlas.cos.rest;

import com.atlas.cos.CharacterTemporalRegistry;
import com.atlas.cos.attribute.BlockedNameAttributes;
import com.atlas.cos.attribute.CharacterAttributes;
import com.atlas.cos.attribute.EquipmentAttributes;
import com.atlas.cos.attribute.EquipmentStatistics;
import com.atlas.cos.attribute.InventoryAttributes;
import com.atlas.cos.attribute.ItemAttributes;
import com.atlas.cos.attribute.LocationAttributes;
import com.atlas.cos.attribute.LocationType;
import com.atlas.cos.builder.BlockedNameAttributesBuilder;
import com.atlas.cos.builder.CharacterAttributesBuilder;
import com.atlas.cos.builder.EquipmentAttributesBuilder;
import com.atlas.cos.builder.EquipmentStatisticsAttributesBuilder;
import com.atlas.cos.builder.InventoryAttributesBuilder;
import com.atlas.cos.builder.ItemAttributesBuilder;
import com.atlas.cos.builder.LocationAttributesBuilder;
import com.atlas.cos.model.BlockedNameData;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.model.CharacterTemporalData;
import com.atlas.cos.model.EquipData;
import com.atlas.cos.model.EquipmentData;
import com.atlas.cos.model.Inventory;
import com.atlas.cos.model.InventoryItem;
import com.atlas.cos.model.InventoryItemType;
import com.atlas.cos.model.ItemData;
import com.atlas.cos.model.SavedLocationData;
import com.atlas.cos.processor.ItemProcessor;

import builder.LinkDataBuilder;
import builder.RelationshipBuilder;
import builder.RelationshipsBuilder;
import builder.ResultObjectBuilder;
import rest.AttributeResult;

public class ResultObjectFactory {
   public static ResultObjectBuilder create(BlockedNameData blockedNameData) {
      return new ResultObjectBuilder(BlockedNameAttributes.class, blockedNameData.name())
            .setAttribute(new BlockedNameAttributesBuilder().setName(blockedNameData.name()));
   }

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
                              .setType(EquipmentStatistics.class)
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
         return ItemProcessor.getEquipmentById(inventoryItem.id())
               .map(ResultObjectFactory::create)
               .orElseThrow();
      } else {
         return ItemProcessor.getItemById(inventoryItem.id())
               .map(ResultObjectFactory::create)
               .orElseThrow();
      }
   }

   public static ResultObjectBuilder create(EquipData equipData) {
      return new ResultObjectBuilder(EquipmentStatistics.class, equipData.id())
            .setAttribute(new EquipmentStatisticsAttributesBuilder()
                  .setItemId(equipData.itemId())
                  .setStrength(equipData.strength())
                  .setDexterity(equipData.dexterity())
                  .setIntelligence(equipData.intelligence())
                  .setLuck(equipData.luck())
                  .setHp(equipData.hp())
                  .setMp(equipData.mp())
                  .setWeaponAttack(equipData.weaponAttack())
                  .setMagicAttack(equipData.magicAttack())
                  .setWeaponDefense(equipData.weaponDefense())
                  .setMagicDefense(equipData.magicDefense())
                  .setAccuracy(equipData.accuracy())
                  .setAvoidability(equipData.avoidability())
                  .setHands(equipData.hands())
                  .setSpeed(equipData.speed())
                  .setJump(equipData.jump())
                  .setSlots(equipData.slots())
            );
   }
}
