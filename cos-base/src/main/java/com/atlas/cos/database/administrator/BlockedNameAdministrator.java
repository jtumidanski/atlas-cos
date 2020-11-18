package com.atlas.cos.database.administrator;

import java.util.List;
import java.util.stream.Collectors;
import javax.persistence.EntityManager;

import com.atlas.cos.database.transformer.BlockedNameDataTransformer;
import com.atlas.cos.entity.BlockedName;
import com.atlas.cos.model.BlockedNameData;

import accessor.AbstractQueryExecutor;

public class BlockedNameAdministrator extends AbstractQueryExecutor {
   private static final Object lock = new Object();

   private static volatile BlockedNameAdministrator instance;

   public static BlockedNameAdministrator getInstance() {
      BlockedNameAdministrator result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new BlockedNameAdministrator();
               instance = result;
            }
         }
      }
      return result;
   }

   public BlockedNameData create(EntityManager entityManager, String name) {
      BlockedName blockedName = new BlockedName();
      blockedName.setName(name);
      return insertAndReturn(entityManager, blockedName, new BlockedNameDataTransformer());
   }

   public List<BlockedNameData> create(EntityManager entityManager, List<String> names) {
      return names.stream()
            .map(name -> create(entityManager, name))
            .collect(Collectors.toList());
   }
}
