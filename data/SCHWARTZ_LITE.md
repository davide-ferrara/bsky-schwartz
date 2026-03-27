# Schwartz Values - Lite

Score each value 0-6 (0=absent, 6=core motivation).

## Value Definitions

**Openness to Change:**
- sd_thought: Freedom to cultivate ideas
- sd_action: Freedom to determine actions
- stimulation: Excitement and novelty
- hedonism: Pleasure and gratification

**Self-Enhancement:**
- achievement: Success by social standards
- dominance: Power over others
- resources: Control of material resources
- face: Maintaining public image

**Conservation:**
- personal_sec: Personal safety
- societal_sec: Societal stability
- tradition: Preserving cultural customs
- rule_conf: Compliance with rules
- inter_conf: Respect for others
- humility: Recognizing insignificance

**Self-Transcendence:**
- caring: Welfare of close others
- dependability: Reliability and loyalty
- universalism: Justice for all people
- nature: Preservation of environment
- tolerance: Acceptance of differences

## Output Format

Return ONLY a JSON object with all 19 values (score 0-6):

```json
{
  "sd_thought": 0,
  "sd_action": 0,
  "stimulation": 0,
  "hedonism": 0,
  "achievement": 0,
  "dominance": 0,
  "resources": 0,
  "face": 0,
  "personal_sec": 0,
  "societal_sec": 0,
  "tradition": 0,
  "rule_conf": 0,
  "inter_conf": 0,
  "humility": 0,
  "caring": 0,
  "dependability": 0,
  "universalism": 0,
  "nature": 0,
  "tolerance": 0
}
```