package com.atlas.cos.event;

public record CharacterSkillUpdateEvent(int characterId, int skillId, int level, int masterLevel, long expiration) {
}
