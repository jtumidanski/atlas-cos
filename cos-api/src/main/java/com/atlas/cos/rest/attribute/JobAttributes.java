package com.atlas.cos.rest.attribute;

import rest.AttributeResult;

public record JobAttributes(String name, Integer createIndex) implements AttributeResult {
}
