import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
from typing import Optional


def plot_value_differences(
    diffs: dict[str, float],
    output_path: str = "diff_chart.png",
    title: str = "GPT vs Qwen: Value Differences",
    figsize: tuple[int, int] = (14, 6),
):
    values = list(diffs.keys())
    differences = list(diffs.values())

    fig, ax = plt.subplots(figsize=figsize)

    colors = ["green" if d >= 0 else "red" for d in differences]

    bars = ax.bar(values, differences, color=colors, edgecolor="black", linewidth=0.5)

    ax.axhline(y=0, color="black", linestyle="-", linewidth=0.8)

    ax.set_xlabel("Schwartz Values", fontsize=10)
    ax.set_ylabel("Difference (GPT - Qwen)", fontsize=10)
    ax.set_title(title, fontsize=12, fontweight="bold")

    plt.xticks(rotation=45, ha="right", fontsize=8)
    plt.tight_layout()

    green_patch = mpatches.Patch(color="green", label="GPT > Qwen")
    red_patch = mpatches.Patch(color="red", label="Qwen > GPT")
    ax.legend(handles=[green_patch, red_patch], loc="upper right")

    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    print(f"Chart saved to {output_path}")
    plt.close()


def plot_cluster_differences(
    cluster_diffs: dict[str, float],
    output_path: str = "cluster_diff_chart.png",
    title: str = "GPT vs Qwen: Cluster Differences",
    figsize: tuple[int, int] = (10, 6),
):
    clusters = list(cluster_diffs.keys())
    differences = list(cluster_diffs.values())

    fig, ax = plt.subplots(figsize=figsize)

    colors = ["green" if d >= 0 else "red" for d in differences]

    bars = ax.bar(clusters, differences, color=colors, edgecolor="black", linewidth=0.5)

    ax.axhline(y=0, color="black", linestyle="-", linewidth=0.8)

    ax.set_xlabel("Value Clusters", fontsize=10)
    ax.set_ylabel("Difference (GPT - Qwen)", fontsize=10)
    ax.set_title(title, fontsize=12, fontweight="bold")

    plt.xticks(rotation=30, ha="right", fontsize=9)
    plt.tight_layout()

    green_patch = mpatches.Patch(color="green", label="GPT > Qwen")
    red_patch = mpatches.Patch(color="red", label="Qwen > GPT")
    ax.legend(handles=[green_patch, red_patch], loc="upper right")

    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    print(f"Chart saved to {output_path}")
    plt.close()
