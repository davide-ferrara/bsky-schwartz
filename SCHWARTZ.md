# Schwartz Theory of Basic Human Values

This document defines the 19 conceptually distinct values used to analyze social
media content. These values represent universal motivations that drive behavior
across cultures.

## 1. Scoring Instructions

For each post, evaluate the presence and magnitude of each value:

- **Score 0**: The value is not present, irrelevant, or explicitly contradicted
  by the content.
- **Score 1-2**: The value is slightly reflected or implied.
- **Score 3-4**: The value is moderately reflected.
- **Score 5-6**: The value is strongly reflected as a core motivation of the
  post.

## 2. Value Definitions Table

The values are organized into four higher-level clusters that represent
fundamental social tensions.

| Cluster                | Value ID (JSON Key) | Simplified Name     | Definition [cite: 187]                                 |
| :--------------------- | :------------------ | :------------------ | :----------------------------------------------------- |
| **Openness to Change** | `sd_thought`        | Independent Thought | Freedom to cultivate one's own ideas and abilities.    |
|                        | `sd_action`         | Independent Action  | Freedom to determine one's own actions.                |
|                        | `stimulation`       | Novelty             | Excitement, stimulation, and change.                   |
|                        | `hedonism`          | Pleasure            | Seeking pleasure and sensuous gratification.           |
| **Self-Enhancement**   | `achievement`       | Achievement         | Success according to social standards.                 |
|                        | `dominance`         | Power               | Influence and the right to command others.             |
|                        | `resources`         | Wealth              | Control of material and social resources.              |
|                        | `face`              | Reputation          | Maintaining public image and avoiding humiliation.     |
| **Conservation**       | `personal_sec`      | Personal Security   | Safety in one's immediate environment.                 |
|                        | `societal_sec`      | Societal Security   | Safety and stability in the wider society.             |
|                        | `tradition`         | Tradition           | Preserving cultural, family, or religious customs.     |
|                        | `rule_conf`         | Lawfulness          | Compliance with rules, laws, and formal obligations.   |
|                        | `inter_conf`        | Respect             | Avoiding upsetting or harming other people.            |
|                        | `humility`          | Humility            | Recognizing one's insignificance in the larger scheme. |
| **Self-Transcendence** | `caring`            | Caring              | Devotion to the welfare of those one is close to.      |
|                        | `dependability`     | Responsibility      | Being reliable and loyal to the "in-group".            |
|                        | `universalism`      | Equality            | Commitment to justice and protection for all people.   |
|                        | `nature`            | Nature              | Preservation of the natural environment.               |
|                        | `tolerance`         | Tolerance           | Acceptance of those different from oneself.            |

## 3. Output Format

Return ONLY a JSON object. No preamble, no markdown code fences, no
explanations. Assign a score (0-6) to every key.

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

## 4. Few-Shot Examples for Calibration

### Example 1: High Conservation / Self-Enhancement

- **Post:** "Dobbiamo tornare a essere un popolo fiero. Disciplina nelle scuole,
  rispetto per la divisa e tolleranza zero per chi calpesta le nostre tradizioni
  millenarie. L'ordine è l'unica via per la grandezza."
- **Target JSON Snippet:**
  `{"dominance": 6, "societal_sec": 5, "tradition": 6, "rule_conf": 6, "tolerance": 0, "universalism": 0}`

### Example 2: High Self-Transcendence / Openness

- **Post:** "In un mondo senza confini, ogni essere umano merita dignità e
  accoglienza. Abbracciamo il cambiamento e lottiamo per un futuro dove
  l'uguaglianza non sia un'utopia ma un diritto universale per tutti."
- **Target JSON Snippet:**
  `{"universalism": 6, "tolerance": 5, "caring": 5, "stimulation": 4, "tradition": 0, "dominance": 0}`

### Example 3: Extreme High Conservation / Self-Enhancement

- **Post:** "Morte al gruppo X, distruggiamo il leader Y, viva il nostro leader
  Z!"
- **Target JSON Snippet:**
  `{"dominance": 6, "societal_sec": 6, "tradition": 6, "rule_conf": 6, "tolerance": 0, "universalism": 0}`
