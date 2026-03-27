# Social Media Post Analysis - Schwartz Basic Human Values

You are analyzing a social media post to identify Schwartz Basic Human Values.

## Scoring Guidelines

Rate each value on a scale of 0-6:

- **Score 0:** Value is absent, irrelevant, or explicitly contradicted
- **Score 1:** Value is barely hinted at, very implicit, slight implication
- **Score 2:** Value is mentioned indirectly or in passing
- **Score 3:** Value is mentioned but NOT central to the post's main message
- **Score 4:** Value is clearly present and moderately important to the post
- **Score 5:** Value is strongly emphasized but not the sole focus
- **Score 6:** Value is the core message or central theme of the post

**Key distinctions:**

- **Explicit vs Implicit:** Direct statements score higher than indirect
  references
- **Centrality:** Values central to the main argument score higher than
  peripheral mentions
- **Tone:** Passionate/emphatic language suggests higher scores than neutral
  statements
- **Context:** Consider the post's topic and emotional weight

## Analysis Rules

1. **Focus on what is said,** not what might be implied
2. **Score conservatively:** When uncertain, choose the lower score
3. **A value can only be central if it's directly addressed** in the post
4. **Multiple values can be present:** Rate each independently
5. **Emotional tone matters:** Stronger emotions → higher confidence in scoring
6. **Ambiguity:** If unclear, err on the side of lower scores

## Calibration Examples

### Example 1: High Conservation / Self-Enhancement

**Post (Italian):** "Dobbiamo difendere i nostri confini! La sicurezza nazionale
è oltre tutto. Le nostre tradizioni e la nostra identità sono sotto attacco.
Basta con questa debolezza!"

**Analysis:**

- Main topic: National security, borders, tradition preservation
- Emotional tone: Angry, emphatic
- Explicit values: societal_sec (national security), tradition (cultural
  identity), dominance (strength vs weakness)
- Central theme: Protection and preservation

**Correct scores:**

```json
{
  "societal_sec": 6,
  "tradition": 5,
  "dominance": 5,
  "inter_conf": 0,
  "caring": 0,
  "tolerance": 0,
  "universalism": 0
}
```

**Why these scores:**

- `societal_sec`: 6 (core message, explicitly stated, emphatic tone)
- `tradition`: 5 (strong emphasis but secondary to security)
- `dominance`: 4 (implicit call for strength, clear but not central)
- Others: 0 (not mentioned or contradicted)

### Example 2: High Self-Transcendence

**Post (Italian):** "Ogni essere umano merita dignità e rispetto. Dobbiamo
proteggere il pianeta per le future generazioni. L'uguaglianza non è un'utopia,
è un diritto."

**Analysis:**

- Main topic: Human rights, environmental protection, equality
- Emotional tone: Hopeful, principled
- Explicit values: universalism (equality), nature (environment), caring (future
  generations)
- Central theme: Justice and protection

**Correct scores:**

```json
{
  "universalism": 6,
  "nature": 5,
  "caring": 5,
  "tolerance": 2,
  "tradition": 0,
  "dominance": 0,
  "societal_sec": 0
}
```

**Why these scores:**

- `universalism`: 6 (explicit statement of equality as right, core message)
- `nature`: 5 (strong emphasis, second main theme)
- `caring`: 4 (mentioned in context of future generations)
- `tolerance`: 3 (implicit in "ogni essere umano", not central)
- Others: 0 (contradicted or absent)

### Example 3: Low/Mixed Values

**Post (Italian):** "Oggi sono andato al parco con il cane. Bella giornata."

**Analysis:**

- Main topic: Casual daily life
- Emotional tone: Neutral, descriptive
- Explicit values: None
- Central theme: None (descriptive post)

**Correct scores:**

```json
{
  "sd_thought": 0,
  "sd_action": 2,
  "stimulation": 2,
  "hedonism": 2,
  "achievement": 0,
  "dominance": 0,
  "tradition": 0,
  "caring": 0,
  "universalism": 0
}
```

**Why these scores:**

- `sd_action`: 2 (slight implication of personal choice activity)
- `stimulation`: 2 (possible enjoyment of outing, very implicit)
- `hedonism`: 2 (mentioned some enjoyment)
- Others: 0 (completely absent)

### Example 4: Moderate Openness to Change

**Post (Italian):** "Penso che ognuno dovrebbe poter scegliere il proprio
percorso nella vita. Niente dovrebbe imporci come vivere."

**Analysis:**

- Main topic: Personal freedom, autonomy
- Emotional tone: Assertive but not extreme
- Explicit values: sd_thought (ideas), sd_action (action)
- Central theme: Freedom from constraints

**Correct scores:**

```json
{
  "sd_thought": 4,
  "sd_action": 5,
  "stimulation": 0,
  "tradition": 2,
  "societal_sec": 0,
  "dominance": 0,
  "caring": 0
}
```

**Why these scores:**

- `sd_action`: 5 (explicit statement about choosing own path, central theme)
- `sd_thought`: 4 (implicit freedom of ideas, supporting theme)
- `tradition`: 2 (implicit rejection "niente dovrebbe imporci", very indirect)
- Others: 0 (not relevant)

## Output Format

Return **ONLY** a JSON object with all 19 Schwartz values. No explanations, no
markdown code fences, no preamble.

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

---

## Your Task

Based on the Schwartz Value Definitions provided, analyze the following social
media post and return a JSON object with scores 0-6 for all 19 values.

**POST:**

