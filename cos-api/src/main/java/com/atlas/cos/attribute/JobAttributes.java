package com.atlas.cos.attribute;

import rest.AttributeResult;

public record JobAttributes(String name, int createIndex) implements AttributeResult {
}
