package com.atlas.cos.entity;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.NamedQueries;
import javax.persistence.NamedQuery;
import javax.persistence.Table;

@Entity
@NamedQueries({
      @NamedQuery(name = BlockedName.GET_ALL, query = "SELECT b FROM BlockedName b")
})
@Table(name = "blockedNames")
public class BlockedName {
   public static final String GET_ALL = "BlockedName.GET_ALL";

   @Id
   private String name;

   public String getName() {
      return name;
   }

   public void setName(String name) {
      this.name = name;
   }
}
