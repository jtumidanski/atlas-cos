package com.atlas.cos.processor;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.function.Consumer;
import java.util.stream.Stream;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.builder.EquipmentDataBuilder;
import com.atlas.cos.database.administrator.EquipmentAdministrator;
import com.atlas.cos.database.provider.EquipmentProvider;
import com.atlas.cos.model.EquipmentData;
import com.atlas.iis.attribute.EquipmentAttributes;
import com.atlas.iis.attribute.EquipmentSlotAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import database.Connection;
import rest.DataBody;
import rest.DataContainer;

public final class ItemProcessor {
   private ItemProcessor() {
   }

   protected static <T> void setIfNotNull(Consumer<T> setter, T value) {
      if (value != null) {
         setter.accept(value);
      }
   }

   public static Optional<EquipmentData> createEquipmentForCharacter(int characterId, int itemId, boolean characterCreation) {
      return createEquipmentForCharacter(characterId, itemId, null, null, null, null,
            null, null, null, null, null, null, null,
            null, null, null, null, characterCreation);
   }

   public static Optional<EquipmentData> createEquipmentForCharacter(int characterId, int itemId, Integer strength,
                                                                     Integer dexterity,
                                                                     Integer intelligence, Integer luck, Integer weaponAttack,
                                                                     Integer weaponDefense, Integer magicAttack,
                                                                     Integer magicDefense,
                                                                     Integer accuracy, Integer avoidability, Integer speed,
                                                                     Integer jump,
                                                                     Integer hp, Integer mp, Integer slots,
                                                                     boolean characterCreation) {
      if (characterCreation) {
         boolean valid = validCharacterCreationItem(itemId);
         if (!valid) {
            return Optional.empty();
         }
      }

      short nextOpenSlot = Connection.instance()
            .element(entityManager -> EquipmentProvider.getNextFreeEquipmentSlot(entityManager, characterId))
            .orElse((short) 0);
      return ItemProcessor.createEquipment(characterId, itemId, nextOpenSlot, strength, dexterity, intelligence, luck,
            weaponAttack, weaponDefense, magicAttack, magicDefense, accuracy, avoidability, speed, jump, hp, mp, slots);
   }

   protected static Optional<EquipmentData> createEquipment(int characterId, int itemId, Short slot, Integer strength,
                                                            Integer dexterity,
                                                            Integer intelligence, Integer luck, Integer weaponAttack,
                                                            Integer weaponDefense,
                                                            Integer magicAttack, Integer magicDefense, Integer accuracy,
                                                            Integer avoidability,
                                                            Integer speed, Integer jump, Integer hp, Integer mp, Integer slots) {
      EquipmentDataBuilder equipmentBuilder = new EquipmentDataBuilder()
            .setItemId(itemId)
            .setSlot(slot);

      UriBuilder.service(RestService.ITEM_INFORMATION)
            .pathParam("equipment", itemId)
            .getRestClient(EquipmentAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .ifPresent(body -> {
               equipmentBuilder.setStrength(body.getAttributes().strength());
               equipmentBuilder.setDexterity(body.getAttributes().dexterity());
               equipmentBuilder.setIntelligence(body.getAttributes().intelligence());
               equipmentBuilder.setLuck(body.getAttributes().luck());
               equipmentBuilder.setWeaponAttack(body.getAttributes().weaponAttack());
               equipmentBuilder.setWeaponDefense(body.getAttributes().weaponDefense());
               equipmentBuilder.setMagicAttack(body.getAttributes().magicAttack());
               equipmentBuilder.setMagicDefense(body.getAttributes().magicDefense());
               equipmentBuilder.setAccuracy(body.getAttributes().accuracy());
               equipmentBuilder.setAvoidability(body.getAttributes().avoidability());
               equipmentBuilder.setSpeed(body.getAttributes().speed());
               equipmentBuilder.setJump(body.getAttributes().jump());
               equipmentBuilder.setHp(body.getAttributes().hp());
               equipmentBuilder.setMp(body.getAttributes().mp());
               equipmentBuilder.setSlots(body.getAttributes().slots());
            });

      setIfNotNull(equipmentBuilder::setStrength, strength);
      setIfNotNull(equipmentBuilder::setDexterity, dexterity);
      setIfNotNull(equipmentBuilder::setIntelligence, intelligence);
      setIfNotNull(equipmentBuilder::setLuck, luck);
      setIfNotNull(equipmentBuilder::setWeaponAttack, weaponAttack);
      setIfNotNull(equipmentBuilder::setWeaponDefense, weaponDefense);
      setIfNotNull(equipmentBuilder::setMagicAttack, magicAttack);
      setIfNotNull(equipmentBuilder::setMagicDefense, magicDefense);
      setIfNotNull(equipmentBuilder::setAccuracy, accuracy);
      setIfNotNull(equipmentBuilder::setAvoidability, avoidability);
      setIfNotNull(equipmentBuilder::setSpeed, speed);
      setIfNotNull(equipmentBuilder::setJump, jump);
      setIfNotNull(equipmentBuilder::setHp, hp);
      setIfNotNull(equipmentBuilder::setMp, mp);
      setIfNotNull(equipmentBuilder::setSlots, slots);

      EquipmentData equipment = equipmentBuilder.build();

      return Connection.instance()
            .element(entityManager -> EquipmentAdministrator.create(entityManager, characterId, equipment));
   }

   protected static Stream<Short> getEquipmentSlotDestination(int itemId) {
      return UriBuilder.service(RestService.ITEM_INFORMATION)
            .pathParam("equipment", itemId)
            .path("slots")
            .getRestClient(EquipmentSlotAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::dataList)
            .orElse(Collections.emptyList()).stream()
            .map(DataBody::getAttributes)
            .map(EquipmentSlotAttributes::slot);
   }

   public static void equipItemForCharacter(int characterId, int id) {
      Connection.instance()
            .element(entityManager -> EquipmentProvider.getById(entityManager, id))
            .map(EquipmentData::itemId)
            .flatMap(itemId -> getEquipmentSlotDestination(itemId).findFirst())
            .ifPresent(destinationSlot -> equipItemForCharacter(characterId, id, destinationSlot));
   }

   protected static void equipItemForCharacter(int characterId, int id, short destinationSlot) {
      Connection.instance().with(entityManager ->
            QueryAdministratorUtil.inTransaction(entityManager, transactionEm -> {
               short temporarySlot = Short.MIN_VALUE;

               EquipmentProvider.findEquipmentInSlot(transactionEm, characterId, destinationSlot)
                     .ifPresent(itemToMove -> EquipmentAdministrator.updateSlot(transactionEm, itemToMove, temporarySlot));

               short currentSlot = EquipmentProvider
                     .getById(transactionEm, id)
                     .map(EquipmentData::slot)
                     .orElse(EquipmentProvider.getNextFreeEquipmentSlot(transactionEm, characterId).orElse((short) 0));

               EquipmentAdministrator.updateSlot(transactionEm, id, destinationSlot);

               EquipmentProvider.findEquipmentInSlot(transactionEm, characterId, temporarySlot)
                     .ifPresent(itemToMove -> EquipmentAdministrator.updateSlot(transactionEm, itemToMove, currentSlot));
            }));
   }

   public static Optional<EquipmentData> getEquipmentForCharacter(int characterId, int equipmentId) {
      return Connection.instance().element(entityManager -> EquipmentProvider.getById(entityManager, equipmentId));
   }

   public static List<EquipmentData> getEquipmentForCharacter(int characterId) {
      return Connection.instance().list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId));
   }

   public static Optional<EquipmentData> getEquipedItemForCharacterBySlot(int characterId, short slotId) {
      return Connection.instance().list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId)).stream()
            .filter(equipmentData -> equipmentData.slot() == slotId)
            .findFirst();
   }

   protected static boolean validCharacterCreationItem(int itemId) {
      return Stream.of(
            1302000, 1312004, 1322005, 1442079,// weapons
            1040002, 1040006, 1040010, 1041002, 1041006, 1041010, 1041011, 1042167,// bottom
            1060002, 1060006, 1061002, 1061008, 1062115, // top
            1072001, 1072005, 1072037, 1072038, 1072383,// shoes
            30000, 30010, 30020, 30030, 31000, 31040, 31050,// hair
            20000, 20001, 20002, 21000, 21001, 21002, 21201, 20401, 20402, 21700, 20100  //face
      ).anyMatch(id -> id == itemId);
   }
}
