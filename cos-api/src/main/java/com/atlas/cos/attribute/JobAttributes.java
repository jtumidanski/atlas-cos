package com.atlas.cos.attribute;

import rest.AttributeResult;

public class JobAttributes implements AttributeResult {
   private String name;

   private Integer createIndex;

   public String getName() {
      return name;
   }

   public void setName(String name) {
      this.name = name;
   }

   public Integer getCreateIndex() {
      return createIndex;
   }

   public void setCreateIndex(Integer createIndex) {
      this.createIndex = createIndex;
   }
}
