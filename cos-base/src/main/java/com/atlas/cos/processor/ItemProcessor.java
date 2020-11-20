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

public class ItemProcessor {
   private static final Object lock = new Object();

   private static volatile ItemProcessor instance;

   public static ItemProcessor getInstance() {
      ItemProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new ItemProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   protected <T> void setIfNotNull(Consumer<T> setter, T value) {
      if (value != null) {
         setter.accept(value);
      }
   }

   public Optional<EquipmentData> createEquipment(int characterId, int itemId, Short slot, Integer strength, Integer dexterity,
                                                  Integer intelligence, Integer luck, Integer weaponAttack, Integer weaponDefense,
                                                  Integer magicAttack, Integer magicDefense, Integer accuracy, Integer avoidability,
                                                  Integer speed, Integer jump, Integer hp, Integer mp, Integer slots) {
      EquipmentDataBuilder equipmentBuilder = new EquipmentDataBuilder()
            .setItemId(itemId)
            .setSlot(slot);

      UriBuilder.service(RestService.ITEM_INFORMATION)
            .pathParam("equipment", itemId)
            .getRestClient(EquipmentAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getData)
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

   public Stream<Short> getEquipmentSlotDestination(int itemId) {
      return UriBuilder.service(RestService.ITEM_INFORMATION)
            .pathParam("equipment", itemId)
            .path("slots")
            .getRestClient(EquipmentSlotAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getDataAsList)
            .orElse(Collections.emptyList()).stream()
            .map(DataBody::getAttributes)
            .map(EquipmentSlotAttributes::slot);
   }

   public void equipItemForCharacter(int characterId, int id, short destinationSlot) {
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

   public Optional<EquipmentData> getEquipmentForCharacter(int characterId, int equipmentId) {
      return Connection.instance().element(entityManager -> EquipmentProvider.getById(entityManager, equipmentId));
   }

   public List<EquipmentData> getEquipmentForCharacter(int characterId) {
      return Connection.instance().list(entityManager -> EquipmentProvider.getForCharacter(entityManager, characterId));
   }
}
