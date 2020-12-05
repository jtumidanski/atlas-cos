package com.atlas.cos.processor;

import com.atlas.cos.database.administrator.BlockedNameAdministrator;
import database.Connection;

import java.util.List;

public final class BlockedNameProcessor {
   private BlockedNameProcessor() {
   }

   public static void bulkAddBlockedNames(List<String> names) {
      Connection.instance()
            .with(entityManager -> BlockedNameAdministrator.create(entityManager, names));
   }
}
