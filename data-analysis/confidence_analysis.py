import argparse
import json
import sys
import time
from datetime import datetime
from pathlib import Path

import matplotlib.pyplot as plt
import numpy as np
import requests
from scipy import stats

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

CLUSTERS = [
    ("Openness to Change", ["sd_thought", "sd_action", "stimulation", "hedonism"]),
    ("Self-Enhancement", ["achievement", "dominance", "resources", "face"]),
    (
        "Conservation",
        [
            "personal_sec",
            "societal_sec",
            "tradition",
            "rule_conf",
            "inter_conf",
            "humility",
        ],
    ),
    (
        "Self-Transcendence",
        ["caring", "dependability", "universalism", "nature", "tolerance"],
    ),
]

MODEL_COLORS = {
    "gpt4o": "#3B82F6",
    "gemini3": "#F97316",
}

MODELS = {
    "gpt4o": "openai/gpt-4o-mini",
    "gemini3": "google/gemini-3.1-flash-lite-preview",
}

DEFAULT_MODELS = ["gpt4o", "gemini3"]

API_URL = "http://localhost:8080/api/analysis/by-url"


def call_analysis_api(feed_data: dict) -> list:
    """Call the API once and return results."""
    response = requests.post(API_URL, json=feed_data, timeout=120)
    response.raise_for_status()
    return response.json()


def run_multiple_analyses(feed_file: str, model_key: str, num_runs: int = 5) -> list:
    """Run analysis N times with same model, return all results."""
    with open(feed_file) as f:
        feed_data = json.load(f)

    feed_data["model"] = model_key

    all_runs = []

    for run_idx in range(num_runs):
        print(f"  Run {run_idx + 1}/{num_runs}...", end=" ", file=sys.stderr)
        try:
            results = call_analysis_api(feed_data)
            if not results:
                raise ValueError("Empty response from API")
            all_runs.append(results)
            print("✓", file=sys.stderr)
        except Exception as e:
            print(f"FAILED", file=sys.stderr)
            print(f"ERROR on run {run_idx + 1}: {e}", file=sys.stderr)
            sys.exit(1)

        if run_idx < num_runs - 1:
            time.sleep(5)

    return all_runs


def calculate_statistics(all_runs: list) -> dict:
    """Calculate mean, std, CI for each Schwartz value and scores."""
    stats_result = {}

    for value_name in VALUES:
        all_values = []
        for run in all_runs:
            for post in run:
                if "values" in post and value_name in post["values"]:
                    all_values.append(float(post["values"][value_name]))

        if len(all_values) > 0:
            mean = np.mean(all_values)
            std_dev = np.std(all_values, ddof=1) if len(all_values) > 1 else 0.0
            n = len(all_values)
            if n > 1 and std_dev > 0:
                ci_size = stats.t.ppf(0.975, df=n - 1) * std_dev / np.sqrt(n)
            else:
                ci_size = 0.0

            stats_result[value_name] = {
                "mean": round(float(mean), 3),
                "std": round(float(std_dev), 3),
                "ci_low": round(float(mean - ci_size), 3),
                "ci_high": round(float(mean + ci_size), 3),
            }
        else:
            stats_result[value_name] = {
                "mean": 0.0,
                "std": 0.0,
                "ci_low": 0.0,
                "ci_high": 0.0,
            }

    for score_name in ["score", "score_penalized"]:
        all_scores = []
        for run in all_runs:
            for post in run:
                if score_name in post:
                    all_scores.append(float(post[score_name]))

        if len(all_scores) > 0:
            mean = np.mean(all_scores)
            std_dev = np.std(all_scores, ddof=1) if len(all_scores) > 1 else 0.0
            n = len(all_scores)
            if n > 1 and std_dev > 0:
                ci_size = stats.t.ppf(0.975, df=n - 1) * std_dev / np.sqrt(n)
            else:
                ci_size = 0.0

            stats_result[score_name] = {
                "mean": round(float(mean), 3),
                "std": round(float(std_dev), 3),
                "ci_low": round(float(mean - ci_size), 3),
                "ci_high": round(float(mean + ci_size), 3),
            }
        else:
            stats_result[score_name] = {
                "mean": 0.0,
                "std": 0.0,
                "ci_low": 0.0,
                "ci_high": 0.0,
            }

    return stats_result


def plot_comparison_bar_chart(
    model_stats: dict, feed_name: str, output_path: str = "comparison_chart.png"
):
    """Plot horizontal bar chart comparing models with error bars."""
    num_models = len(model_stats)
    num_values = len(VALUES)

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(16, max(10, num_values * 0.5)))

    model_names = list(model_stats.keys())

    y_pos = np.arange(num_values)

    bar_height = 0.8 / num_models

    for i, model_name in enumerate(model_names):
        stats = model_stats[model_name]
        means = [stats[v]["mean"] for v in VALUES]
        errors_low = [stats[v]["mean"] - stats[v]["ci_low"] for v in VALUES]
        errors_high = [stats[v]["ci_high"] - stats[v]["mean"] for v in VALUES]
        errors = [errors_low, errors_high]

        color = MODEL_COLORS.get(model_name, f"C{i}")

        ax1.barh(
            y_pos + i * bar_height,
            means,
            bar_height,
            xerr=errors,
            label=model_name,
            color=color,
            alpha=0.8,
            capsize=3,
        )

    ax1.set_yticks(y_pos + bar_height * (num_models - 1) / 2)
    ax1.set_yticklabels(VALUES, fontsize=9)
    ax1.set_xlabel("Value (0-6)", fontsize=11)
    ax1.set_title(
        f"Schwartz Values Comparison\n{feed_name}", fontsize=13, weight="bold"
    )
    ax1.legend(loc="lower right", fontsize=10)
    ax1.grid(axis="x", alpha=0.3)
    ax1.set_xlim(0, 6)

    scores = ["score", "score_penalized"]
    y_pos_scores = np.arange(len(scores))

    for i, model_name in enumerate(model_names):
        stats = model_stats[model_name]
        means = [stats[s]["mean"] for s in scores]
        errors_low = [stats[s]["mean"] - stats[s]["ci_low"] for s in scores]
        errors_high = [stats[s]["ci_high"] - stats[s]["mean"] for s in scores]
        errors = [errors_low, errors_high]

        color = MODEL_COLORS.get(model_name, f"C{i}")

        ax2.barh(
            y_pos_scores + i * bar_height,
            means,
            bar_height,
            xerr=errors,
            label=model_name,
            color=color,
            alpha=0.8,
            capsize=3,
        )

    ax2.set_yticks(y_pos_scores + bar_height * (num_models - 1) / 2)
    ax2.set_yticklabels(["Score", "Score (Penalized)"], fontsize=10)
    ax2.set_xlabel("Score", fontsize=11)
    ax2.set_title(f"Overall Scores Comparison\n{feed_name}", fontsize=13, weight="bold")
    ax2.legend(loc="lower right", fontsize=10)
    ax2.grid(axis="x", alpha=0.3)

    plt.tight_layout()
    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    print(f"Comparison chart saved to {output_path}", file=sys.stderr)


def save_confidence_results(
    model_stats: dict, feed_name: str, num_runs: int, output_file: str
):
    """Save confidence results to JSON."""
    output_path = Path(output_file)
    output_path.parent.mkdir(parents=True, exist_ok=True)

    output = {
        "feed": feed_name,
        "runs": num_runs,
        "models": model_stats,
        "timestamp": datetime.now().isoformat(),
    }

    with open(output_file, "w") as f:
        json.dump(output, f, indent=2)

    print(f"Results saved to {output_file}", file=sys.stderr)


def print_summary(model_stats: dict, feed_name: str):
    """Print summary to stderr."""
    print(f"\n{'=' * 60}", file=sys.stderr)
    print(f"Summary for {feed_name}", file=sys.stderr)
    print(f"{'=' * 60}", file=sys.stderr)

    for cluster_name, cluster_vals in CLUSTERS:
        print(f"\n{cluster_name}:", file=sys.stderr)
        for v in cluster_vals:
            print(f"  {v:20s}:", file=sys.stderr)
            for model_name, stats in model_stats.items():
                s = stats[v]
                print(
                    f"    {model_name:15s}: {s['mean']:5.2f} ± {s['std']:4.2f} "
                    f"(CI: {s['ci_low']:5.2f} - {s['ci_high']:5.2f})",
                    file=sys.stderr,
                )

    print(f"\nOverall Scores:", file=sys.stderr)
    for score_name in ["score", "score_penalized"]:
        print(f"  {score_name:20s}:", file=sys.stderr)
        for model_name, stats in model_stats.items():
            s = stats[score_name]
            print(
                f"    {model_name:15s}: {s['mean']:7.2f} ± {s['std']:5.2f} "
                f"(CI: {s['ci_low']:7.2f} - {s['ci_high']:7.2f})",
                file=sys.stderr,
            )
    print(f"{'=' * 60}\n", file=sys.stderr)


def main():
    parser = argparse.ArgumentParser(
        description="Run confidence analysis on Bluesky feeds (single or multi-model)"
    )
    parser.add_argument(
        "--feed",
        required=True,
        help="Path to feed JSON file (e.g., feeds/conservation.json)",
    )
    parser.add_argument(
        "--model",
        help="Single model to use (e.g., gpt4o). Mutually exclusive with --models",
    )
    parser.add_argument(
        "--models",
        help="Comma-separated list of models to compare (e.g., gpt4o,gemini3). "
        "Default: gpt4o,gemini3",
    )
    parser.add_argument(
        "--runs",
        type=int,
        default=5,
        help="Number of runs per model for confidence interval (default: 5)",
    )
    parser.add_argument(
        "--output",
        help="Output JSON file path (default: results/{feed_name}_confidence.json)",
    )
    parser.add_argument(
        "--chart",
        help="Output chart file path (default: results/{feed_name}_chart.png)",
    )
    parser.add_argument(
        "--no-chart",
        action="store_true",
        help="Skip chart generation",
    )

    args = parser.parse_args()

    if args.model and args.models:
        print("ERROR: --model and --models are mutually exclusive", file=sys.stderr)
        sys.exit(1)

    feed_path = Path(args.feed)
    if not feed_path.exists():
        print(f"ERROR: Feed file not found: {args.feed}", file=sys.stderr)
        sys.exit(1)

    feed_name = feed_path.stem

    if args.models:
        model_list = [m.strip() for m in args.models.split(",")]
    elif args.model:
        model_list = [args.model]
    else:
        model_list = DEFAULT_MODELS

    for model_name in model_list:
        if model_name not in MODELS:
            print(
                f"ERROR: Unknown model '{model_name}'. Available models: {', '.join(MODELS.keys())}",
                file=sys.stderr,
            )
            sys.exit(1)

    if args.output:
        output_file = args.output
    else:
        output_file = f"results/{feed_name}_confidence.json"

    if args.chart:
        chart_file = args.chart
    elif not args.no_chart:
        chart_file = f"results/{feed_name}_chart.png"
    else:
        chart_file = None

    print(f"\n{'=' * 60}", file=sys.stderr)
    print(f"Confidence Analysis: {feed_name}", file=sys.stderr)
    print(f"Models: {', '.join(model_list)} | Runs: {args.runs}", file=sys.stderr)
    print(f"{'=' * 60}\n", file=sys.stderr)

    model_stats = {}

    for model_key in model_list:
        print(f"Processing model: {model_key}", file=sys.stderr)

        all_runs = run_multiple_analyses(args.feed, model_key, args.runs)

        print(f"  Calculating statistics...", file=sys.stderr)
        statistics = calculate_statistics(all_runs)

        model_stats[model_key] = statistics

        print(f"  Done!\n", file=sys.stderr)

    print(f"\nSaving results...", file=sys.stderr)
    save_confidence_results(model_stats, feed_name, args.runs, output_file)

    if chart_file:
        print(f"\nGenerating comparison chart...", file=sys.stderr)
        plot_comparison_bar_chart(model_stats, feed_name, chart_file)

    print_summary(model_stats, feed_name)


if __name__ == "__main__":
    main()
