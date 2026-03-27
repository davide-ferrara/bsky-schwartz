import matplotlib.pyplot as plt
import numpy as np


def plot_comparison_chart(avg_data: dict, output_path: str = "chart.png"):
    """Grafico a barre standard (NON polare)"""
    model_names = list(avg_data.keys())
    values_a = [avg_data[model_names[0]].get(v, 0) for v in VALUES]
    values_b = [avg_data[model_names[1]].get(v, 0) for v in VALUES]

    x = np.arange(len(VALUES))
    width = 0.35

    fig, ax = plt.subplots(figsize=(14, 8))

    ax.bar(x - width / 2, values_a, width, label=model_names[0], color="#3B82F6")
    ax.bar(x + width / 2, values_b, width, label=model_names[1], color="#F97316")

    ax.set_ylabel("Average Score (0-6)")
    ax.set_title(f"Schwartz Values Comparison")
    ax.set_xticks(x)
    ax.set_xticklabels(VALUES, rotation=45, ha="right")
    ax.legend()
    ax.set_ylim(0, 6)
    plt.tight_layout()
    plt.savefig(output_path, dpi=150)


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

VALUES = [val for _, sublist in CLUSTERS for val in sublist]


def plot_radar_chart(avg_data: dict, output_path: str = "radar.png"):
    model_names = list(avg_data.keys())
    cluster_colors = {
        "Openness to Change": "#E879A0",  # Rosa
        "Self-Enhancement": "#F97316",  # Arancione
        "Conservation": "#7DD3FC",  # Azzurro
        "Self-Transcendence": "#22C55E",  # Verde
    }

    num_vars = len(VALUES)
    angles = np.linspace(0, 2 * np.pi, num_vars, endpoint=False).tolist()
    angles += angles[:1]

    fig, ax = plt.subplots(figsize=(12, 12), subplot_kw=dict(polar=True))

    ax.set_theta_offset(np.pi / 2)
    ax.set_theta_direction(-1)

    model_colors = ["#3B82F6", "#F97316"]
    for i, model in enumerate(model_names):
        values = [avg_data[model].get(v, 0) for v in VALUES]
        values += values[:1]
        ax.plot(angles, values, "o-", linewidth=2, label=model, color=model_colors[i])
        ax.fill(angles, values, alpha=0.15, color=model_colors[i])

    ax.set_xticks(angles[:-1])
    ax.set_xticklabels(VALUES, size=9)
    ax.set_ylim(0, 6)

    label_colors = []
    for v in VALUES:
        for cluster_name, cluster_vals in CLUSTERS:
            if v in cluster_vals:
                label_colors.append(cluster_colors[cluster_name])
                break

    for label, color in zip(ax.get_xticklabels(), label_colors):
        label.set_color(color)
        label.set_fontweight("bold")

    current_idx = 0
    for _, cluster_vals in CLUSTERS:
        current_idx += len(cluster_vals)
        ax.axvline(
            angles[current_idx % len(VALUES)], color="grey", linestyle="--", alpha=0.3
        )

    plt.tight_layout()
    plt.savefig(output_path, dpi=150)
    print(f"Radar chart saved to {output_path}")
