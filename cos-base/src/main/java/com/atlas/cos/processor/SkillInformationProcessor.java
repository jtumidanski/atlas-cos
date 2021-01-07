package com.atlas.cos.processor;

import java.util.Optional;
import java.util.concurrent.CompletableFuture;

import com.app.rest.util.RestResponseUtil;
import com.atlas.cos.model.SkillInformation;
import com.atlas.shared.rest.UriBuilder;
import com.atlas.sis.constant.RestConstants;
import com.atlas.sis.rest.attribute.SkillAttributes;

import rest.DataContainer;

public final class SkillInformationProcessor {
   private SkillInformationProcessor() {
   }

   public static CompletableFuture<SkillInformation> getSkillInformation(int skillId) {
      return UriBuilder.service(RestConstants.SERVICE)
            .pathParam("skills", skillId)
            .getAsyncRestClient(SkillAttributes.class)
            .get()
            .thenApply(RestResponseUtil::result)
            .thenApply(DataContainer::data)
            .thenApply(Optional::get)
            .thenApply(ModelFactory::createSkillInformation);
   }
}
