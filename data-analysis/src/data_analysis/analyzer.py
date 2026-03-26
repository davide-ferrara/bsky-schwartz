from dataclasses import dataclass
from typing import Optional

from .api_client import APIClient


VALUES = [
    "sd_thought",
    "sd_action",
    "stimulation",
    "hedonism",
    "achievement",
    "dominance",
    "resources",
    "face",
    "personal_sec",
    "societal_sec",
    "tradition",
    "rule_conf",
    "inter_conf",
    "humility",
    "caring",
    "dependability",
    "universalism",
    "nature",
    "tolerance",
]

CLUSTERS = {
    "Openness to Change": ["sd_thought", "sd_action", "stimulation", "hedonism"],
    "Self-Enhancement": ["achievement", "dominance", "resources", "face"],
    "Conservation": [
        "personal_sec",
        "societal_sec",
        "tradition",
        "rule_conf",
        "inter_conf",
        "humility",
    ],
    "Self-Transcendence": [
        "caring",
        "dependability",
        "universalism",
        "nature",
        "tolerance",
    ],
}


@dataclass
class AnalysisResult:
    uri: str
    model: str
    values: dict[str, int]
    score: float
    stats: dict


class Analyzer:
    def __init__(self, api_client: APIClient):
        self.api = api_client

    def analyze_uri(self, uri: str, model: str) -> Optional[AnalysisResult]:
        try:
            data = self.api.analyze_by_uri(uri, model)
            return AnalysisResult(
                uri=uri,
                model=model,
                values=data.get("values", {}),
                score=data.get("score", 0.0),
                stats=data.get("stats", {}),
            )
        except Exception as e:
            print(f"Error analyzing {uri} with {model}: {e}")
            return None

    def compare_models(self, uris: list[str], models: list[str] = None) -> dict:
        if models is None:
            models = ["gpt", "qwen"]

        results = {model: [] for model in models}

        for uri in uris:
            for model in models:
                result = self.analyze_uri(uri, model)
                if result:
                    results[model].append(result)

        return results

    def compute_differences(self, results: dict) -> dict[str, float]:
        if len(results) != 2:
            raise ValueError("Exactly 2 models required for comparison")

        model_names = list(results.keys())
        model_a, model_b = model_names[0], model_names[1]
        results_a, results_b = results[model_a], results[model_b]

        if len(results_a) == 0 or len(results_b) == 0:
            raise ValueError("No results to compare")

        diffs = {}
        for val in VALUES:
            scores_a = [r.values.get(val, 0) for r in results_a]
            scores_b = [r.values.get(val, 0) for r in results_b]
            diffs[val] = (sum(scores_a) / len(scores_a)) - (
                sum(scores_b) / len(scores_b)
            )

        return diffs

    def compute_cluster_diffs(self, diffs: dict) -> dict[str, float]:
        cluster_diffs = {}
        for cluster_name, values in CLUSTERS.items():
            cluster_diffs[cluster_name] = sum(diffs.get(v, 0) for v in values) / len(
                values
            )
        return cluster_diffs
