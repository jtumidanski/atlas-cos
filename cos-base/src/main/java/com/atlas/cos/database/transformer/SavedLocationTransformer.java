package com.atlas.cos.database.transformer;

import com.atlas.cos.entity.SavedLocation;
import com.atlas.cos.model.SavedLocationData;

import transformer.SqlTransformer;

public class SavedLocationTransformer implements SqlTransformer<SavedLocationData, SavedLocation> {
   @Override
   public SavedLocationData transform(SavedLocation location) {
      return new SavedLocationData(location.getId(), location.getCharacterId(), location.getLocationType().name(),
            location.getMap(), location.getPortal());
   }
}
