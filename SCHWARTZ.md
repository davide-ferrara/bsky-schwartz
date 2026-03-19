# Schwartz Values definition

This file contains the 19 Schwartz values. Schwartz values are a broad set of
human values that articulate complementary and opposing values, forming the
building blocks of many cultures.

The modern Schwartz theory articulates 19 conceptually distinct values, each of
which is placed into a circumplex model to articulate values that are
complementary or in tension with each other.

## The score range

To identify values in tweets, consider both the expression of the value and the
magnitude. We score each value’s expression on each post as an integer value. A
score of 0 indicates that the value does not exist or is not supported in the
tweet (i.e., the tweet contains content that contradicts a value). If a tweet
does contain content that supports the value, the scores range from 1 (i.e., the
value is slightly reflected in the content) to 6 (i.e., the value is strongly
reflected)

## Division of the spectrum

The spectrum is divided into 4 important basic values: Self-Transcendence,
Openness to Change, Conservation, Self-Enhancement

1. Openness to Change: Self-directed thoughts, Self-directed actions,
   Stimulation, Hedonism.

2. Self-Enhancement: Achievement, Dominance, Resources, Face.

3. Conservation: Personal Security, Societal Security, Tradition, Rule
   Conformity, Interpersonal Conformity, Humility.

4. Self-Transcendence: Caring, Dependability, Universal Concern, Preservation of
   Nature, Tolerance.

# Schwartz Value Definition Table

| Value                    | Definition                                                                         |
| ------------------------ | ---------------------------------------------------------------------------------- |
| Self-directed thoughts   | The freedom to cultivate one's own ideas and abilities                             |
| Self-directed actions    | The freedom to determine one's own actions                                         |
| Stimulation              | Excitement, stimulation, and change                                                |
| Hedonism                 | Hedonism                                                                           |
| Achievement              | Success according to social standards                                              |
| Dominance                | Influence and the right to command                                                 |
| Resources                | Control of material and social resources                                           |
| Face                     | Security and power through maintaining one's public image and avoiding humiliation |
| Personal Security        | Safety in one's immediate environment                                              |
| Societal Security        | Safety and stability in the wider society                                          |
| Tradition                | Maintaining and preserving cultural, family, or religious traditions               |
| Rule Conformity          | Compliance with rules, laws, and formal obligations                                |
| Interpersonal Conformity | Avoiding upsetting or harming other people                                         |
| Humility                 | Being humble                                                                       |
| Caring                   | Devotion to close others                                                           |
| Dependability            | Being responsible and loyal to others                                              |
| Universal Concern        | Commitment to equality, justice, and protection for all people                     |
| Preservation of Nature   | Preservation of the natural environment                                            |
| Tolerance                | Acceptance and understanding of those different from oneself                       |

## Return Format

Return an array of 19 integers (scores 0-6), one for each value in the table order:

```
[1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19]
```

Index mapping (0-based):
0. Self-directed thoughts
1. Self-directed actions
2. Stimulation
3. Hedonism
4. Achievement
5. Dominance
6. Resources
7. Face
8. Personal Security
9. Societal Security
10. Tradition
11. Rule Conformity
12. Interpersonal Conformity
13. Humility
14. Caring
15. Dependability
16. Universal Concern
17. Preservation of Nature
18. Tolerance
