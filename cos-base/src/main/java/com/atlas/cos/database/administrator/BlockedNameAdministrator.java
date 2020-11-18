package com.atlas.cos.database.administrator;

import java.util.List;
import java.util.stream.Collectors;
import javax.persistence.EntityManager;

import com.app.database.util.QueryAdministratorUtil;
import com.atlas.cos.database.transformer.BlockedNameDataTransformer;
import com.atlas.cos.entity.BlockedName;
import com.atlas.cos.model.BlockedNameData;

public class BlockedNameAdministrator {
   private BlockedNameAdministrator() {
   }

   public static BlockedNameData create(EntityManager entityManager, String name) {
      BlockedName blockedName = new BlockedName();
      blockedName.setName(name);
      return QueryAdministratorUtil.insertAndReturn(entityManager, blockedName, new BlockedNameDataTransformer());
   }

   public static List<BlockedNameData> create(EntityManager entityManager, List<String> names) {
      return names.stream()
            .map(name -> create(entityManager, name))
            .collect(Collectors.toList());
   }
}
