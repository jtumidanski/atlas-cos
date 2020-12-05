package com.atlas.cos.processor;

import java.awt.*;
import java.util.Optional;

import com.atlas.cos.model.Portal;

public final class MapProcessor {
   private MapProcessor() {
   }

   public static Optional<Integer> findClosestSpawnPoint(int mapId, int x, int y) {
      Point from = new Point(x, y);
      return PortalProcessor.getMapPortals(mapId).stream()
            .filter(MapProcessor::isSpawnPoint)
            .min((o1, o2) -> compareDistanceFromPoint(from, o1, o2))
            .map(Portal::id);
   }

   protected static boolean isSpawnPoint(Portal portal) {
      return (portal.type() == 0 || portal.type() == 1) && portal.targetMap() == 999999999;
   }

   protected static int compareDistanceFromPoint(Point from, Portal o1, Portal o2) {
      double o1Distance = new Point(o1.x(), o1.y()).distanceSq(from);
      double o2Distance = new Point(o2.x(), o2.y()).distanceSq(from);
      return Double.compare(o1Distance, o2Distance);
   }
}
