package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record ItemAttributes(Integer itemId, Integer quantity, Short slot) implements AttributeResult {
}
