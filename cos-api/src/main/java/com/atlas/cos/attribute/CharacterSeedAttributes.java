package com.atlas.cos.attribute;

import rest.AttributeResult;

public record CharacterSeedAttributes(Integer accountId, Integer worldId, String name, Integer jobIndex, Integer face, Integer hair,
                                      Integer hairColor, Integer skin, Byte gender, Integer top, Integer bottom, Integer shoes,
                                      Integer weapon) implements AttributeResult {
}
