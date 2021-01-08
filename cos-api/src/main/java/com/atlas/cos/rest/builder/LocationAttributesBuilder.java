package com.atlas.cos.rest.builder;

import com.app.common.builder.RecordBuilder;
import com.atlas.cos.rest.attribute.LocationAttributes;
import com.atlas.cos.rest.attribute.LocationType;

import builder.AttributeResultBuilder;

public class LocationAttributesBuilder extends RecordBuilder<LocationAttributes, LocationAttributesBuilder>
      implements AttributeResultBuilder {
   private LocationType type;

   private Integer mapId;

   private Integer portalId;

   @Override
   public LocationAttributes construct() {
      return new LocationAttributes(type, mapId, portalId);
   }

   @Override
   public LocationAttributesBuilder getThis() {
      return this;
   }

   public LocationAttributesBuilder setType(LocationType type) {
      this.type = type;
      return getThis();
   }

   public LocationAttributesBuilder setMapId(Integer mapId) {
      this.mapId = mapId;
      return getThis();
   }

   public LocationAttributesBuilder setPortalId(Integer portalId) {
      this.portalId = portalId;
      return getThis();
   }
}
