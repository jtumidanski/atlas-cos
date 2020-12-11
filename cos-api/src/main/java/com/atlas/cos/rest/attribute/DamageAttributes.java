package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record DamageAttributes(DamageType type, Integer maximum) implements AttributeResult {
}
