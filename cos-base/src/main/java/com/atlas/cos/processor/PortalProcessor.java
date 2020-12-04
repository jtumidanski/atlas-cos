package com.atlas.cos.processor;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

import com.atlas.cos.model.Portal;
import com.atlas.mis.attribute.PortalAttributes;
import com.atlas.shared.rest.RestService;
import com.atlas.shared.rest.UriBuilder;

import rest.DataContainer;

public class PortalProcessor {
   private static final Object lock = new Object();

   private static volatile PortalProcessor instance;

   public static PortalProcessor getInstance() {
      PortalProcessor result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new PortalProcessor();
               instance = result;
            }
         }
      }
      return result;
   }

   public Optional<Portal> getMapPortalById(int mapId, int id) {
      return UriBuilder.service(RestService.MAP_INFORMATION)
            .pathParam("maps", mapId)
            .pathParam("portals", id)
            .getRestClient(PortalAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getData)
            .map(ModelFactory::createPortal);
   }

   public List<Portal> getMapPortals(int mapId) {
      return UriBuilder.service(RestService.MAP_INFORMATION)
            .pathParam("maps", mapId)
            .path("portals")
            .getRestClient(PortalAttributes.class)
            .getWithResponse()
            .result()
            .map(DataContainer::getDataAsList)
            .orElse(Collections.emptyList())
            .stream()
            .map(ModelFactory::createPortal)
            .collect(Collectors.toList());
   }
}
