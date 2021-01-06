package com.atlas.cos.event.consumer;

import com.atlas.cos.FacialExpressionRegistry;
import com.atlas.cos.command.ChangeFacialExpression;
import com.atlas.cos.constant.CommandConstants;
import com.atlas.cos.event.producer.FacialExpressionProducer;
import com.atlas.cos.model.CharacterData;
import com.atlas.cos.processor.CharacterProcessor;
import com.atlas.cos.processor.TopicDiscoveryProcessor;

public class CharacterExpressionCommandConsumer extends AbstractEventConsumer<ChangeFacialExpression> {
   @Override
   public void handle(Long key, ChangeFacialExpression command) {
      CharacterProcessor
            .getById(command.characterId())
            .ifPresent(character -> changeFacialExpression(character, command.emote()));
   }

   protected void changeFacialExpression(CharacterData character, int expression) {
      FacialExpressionProducer.emit(character.id(), character.mapId(), expression);
      FacialExpressionRegistry.getInstance().registerChange(character.id(), getReturnToNormal(character.id()));
   }

   protected Runnable getReturnToNormal(int characterId) {
      return () ->
            CharacterProcessor.getById(characterId)
                  .ifPresent(character -> FacialExpressionProducer.emit(character.id(), character.mapId(), 0));
   }

   @Override
   public Class<ChangeFacialExpression> getEventClass() {
      return ChangeFacialExpression.class;
   }

   @Override
   public String getTopic() {
      return TopicDiscoveryProcessor.getTopic(CommandConstants.CHANGE_FACIAL_EXPRESSION);
   }
}
