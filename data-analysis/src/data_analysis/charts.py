import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
import numpy as np
from typing import Optional


def plot_value_comparison(
    results: dict[str, list],
    output_path: str = "comparison_chart.png",
    title: str = "GPT vs Minimax: Value Comparison",
    figsize: tuple[int, int] = (14, 8),
):
    """
    Plot side-by-side bars showing absolute values for each model.
    results: {"gpt": [AnalysisResult, ...], "minimax": [AnalysisResult, ...]}
    """
    values_list = [
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

    model_a = "gpt"
    model_b = "minimax"

    results_a = results[model_a]
    results_b = results[model_b]

    avg_a = []
    avg_b = []
    for val in values_list:
        scores_a = [r.values.get(val, 0) for r in results_a]
        scores_b = [r.values.get(val, 0) for r in results_b]
        avg_a.append(sum(scores_a) / len(scores_a) if scores_a else 0)
        avg_b.append(sum(scores_b) / len(scores_b) if scores_b else 0)

    x = np.arange(len(values_list))
    width = 0.35

    fig, ax = plt.subplots(figsize=figsize)

    bars_a = ax.bar(
        x - width / 2,
        avg_a,
        width,
        label="GPT",
        color="#3B82F6",
        edgecolor="black",
        linewidth=0.5,
    )
    bars_b = ax.bar(
        x + width / 2,
        avg_b,
        width,
        label="Minimax",
        color="#F97316",
        edgecolor="black",
        linewidth=0.5,
    )

    ax.set_xlabel("Schwartz Values", fontsize=10)
    ax.set_ylabel("Average Score (0-6)", fontsize=10)
    ax.set_title(title, fontsize=12, fontweight="bold")
    ax.set_xticks(x)
    ax.set_xticklabels(values_list, rotation=45, ha="right", fontsize=8)
    ax.legend()
    ax.set_ylim(0, 6)

    plt.tight_layout()
    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    print(f"Chart saved to {output_path}")
    plt.close()


def plot_value_differences(
    diffs: dict[str, float],
    output_path: str = "diff_chart.png",
    title: str = "GPT vs Minimax: Value Differences",
    figsize: tuple[int, int] = (14, 6),
):
    values = list(diffs.keys())
    differences = list(diffs.values())

    fig, ax = plt.subplots(figsize=figsize)

    colors = ["green" if d >= 0 else "red" for d in differences]

    bars = ax.bar(values, differences, color=colors, edgecolor="black", linewidth=0.5)

    ax.axhline(y=0, color="black", linestyle="-", linewidth=0.8)

    ax.set_xlabel("Schwartz Values", fontsize=10)
    ax.set_ylabel("Difference (GPT - Minimax)", fontsize=10)
    ax.set_title(title, fontsize=12, fontweight="bold")

    plt.xticks(rotation=45, ha="right", fontsize=8)
    plt.tight_layout()

    green_patch = mpatches.Patch(color="green", label="GPT > Minimax")
    red_patch = mpatches.Patch(color="red", label="Minimax > GPT")
    ax.legend(handles=[green_patch, red_patch], loc="upper right")

    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    print(f"Chart saved to {output_path}")
    plt.close()


def plot_cluster_differences(
    cluster_diffs: dict[str, float],
    output_path: str = "cluster_diff_chart.png",
    title: str = "GPT vs Minimax: Cluster Differences",
    figsize: tuple[int, int] = (10, 6),
):
    clusters = list(cluster_diffs.keys())
    differences = list(cluster_diffs.values())

    fig, ax = plt.subplots(figsize=figsize)

    colors = ["green" if d >= 0 else "red" for d in differences]

    bars = ax.bar(clusters, differences, color=colors, edgecolor="black", linewidth=0.5)

    ax.axhline(y=0, color="black", linestyle="-", linewidth=0.8)

    ax.set_xlabel("Value Clusters", fontsize=10)
    ax.set_ylabel("Difference (GPT - Minimax)", fontsize=10)
    ax.set_title(title, fontsize=12, fontweight="bold")

    plt.xticks(rotation=30, ha="right", fontsize=9)
    plt.tight_layout()

    green_patch = mpatches.Patch(color="green", label="GPT > Minimax")
    red_patch = mpatches.Patch(color="red", label="Minimax > GPT")
    ax.legend(handles=[green_patch, red_patch], loc="upper right")

    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    print(f"Chart saved to {output_path}")
    plt.close()
