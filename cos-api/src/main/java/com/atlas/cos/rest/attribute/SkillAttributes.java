package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record SkillAttributes(Integer level, Integer masterLevel, Long expiration)
      implements AttributeResult {
}
