package com.atlas.cos.model;

public record CharacterData(int id, int accountId, int worldId, String name, int level, int experience, int gachaponExperience,
                            int strength, int dexterity, int luck, int intelligence, int hp, int mp, int maxHp, int maxMp, int meso,
                            int hpMpUsed, int jobId, int skinColor, int gender, int fame, int hair, int face, int ap, String sp,
                            int mapId, int spawnPoint, int gm) {
}
