package com.atlas.cos.model;

import java.util.Map;

public record MonsterExperienceDistribution(Map<Integer, Long> solo, Map<Integer, Map<Integer, Long>> party,
                                            Map<Integer, Float> personalRatio, float experiencePerDamage,
                                            double standardDeviationRatio) {
}
