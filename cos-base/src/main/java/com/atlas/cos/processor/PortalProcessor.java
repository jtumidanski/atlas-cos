package com.atlas.cos.processor;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import com.atlas.cos.model.Portal;
import com.atlas.mis.attribute.PortalAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public final class PortalProcessor {
   private PortalProcessor() {
   }

   public static Optional<Portal> getMapPortalById(int mapId, int id) {
      return UriBuilder.service(RestService.MAP_INFORMATION)
            .pathParam("maps", mapId)
            .pathParam("portals", id)
            .getRestClient(PortalAttributes.class)
            .getWithResponse()
            .result()
            .flatMap(DataContainer::data)
            .map(ModelFactory::createPortal);
   }

   public static Stream<Portal> getMapPortals(int mapId) {
      return UriBuilder.service(RestService.MAP_INFORMATION)
            .pathParam("maps", mapId)
            .path("portals")
            .getRestClient(PortalAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::dataList)
            .orElse(Collections.emptyList())
            .stream()
            .map(ModelFactory::createPortal);
   }
}
